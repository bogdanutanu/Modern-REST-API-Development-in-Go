// @title Shopping List API
// @version 0.1
// @description An API for managing shopping lists
// @host my-shopping-lists.com/api
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/labstack/echo-contrib/echoprometheus"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

// ShoppingListPatch represents partial updates to a shopping list
// @Description Partial shopping list update structure
type ShoppingListPatch struct {
	Name  *string  `json:"name"`
	Items []string `json:"items"`
}

// ListPushAction represents an action to add an item to a list
// @Description Action to push an item to a shopping list
type ListPushAction struct {
	Item string `json:"item"`
}

// ShopingList represents a shopping list with items
// @Description Shopping list with items
type ShoppingList struct {
	gorm.Model
	Name   string   `gorm:"not null"`
	UserID uint     `gorm:"not null"`
	Items  []string `gorm:"serializer:json"`
}

type Session struct {
	gorm.Model
	Token   string    `gorm:"uniqueIndex;not null"`
	Expires time.Time `gorm:"not null"`
	UserID  uint      `gorm:"not null"`
}

type User struct {
	gorm.Model
	Role          string `gorm:"not null"`
	Username      string `gorm:"uniqueIndex;not null"`
	Password      string `gorm:"not null"`
	Sessions      []Session
	ShoppingLists []ShoppingList
}

// LoginRequest represents user login credentials
// @Description User login request structure
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var listsCache *lru.Cache[string, ShoppingList]

var repository RepositoryInterface

var metricsService *Metrics

func main() {
	e := echo.New()
	api := e.Group("/api")
	api.Use(authRequired)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://example.com"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
	e.Use(echojwt.JWT([]byte("my-secret")))
	e.Use(echoprometheus.NewMiddleware("myapp"))
	e.GET("/metrics", echoprometheus.NewHandler())

	metricsService = NewMetrics()

	var err error
	listsCache, err = lru.New[string, ShoppingList](128)
	if err != nil {
		fmt.Println("Unable to initialize the lists cache:", err.Error())
		os.Exit(1)
	}
	repository, err = NewRepository("./database.db")
	if err != nil {
		fmt.Println("Unable to open the database:", err.Error())
		os.Exit(1)
	}
	if err := repository.Init(); err != nil {
		fmt.Println("Unable to initialize the database:", err.Error())
		os.Exit(1)
	}

	closer, err := InitTracer("my-api")
	if err != nil {
		log.Fatalf("Could not initialize the tracer: %v", err)
	}
	defer closer.Close()

	e.POST("/api/login", handleLogin)
	api.GET("/lists", handleListLists)
	api.POST("/lists", adminRequired(handleCreateList))
	api.GET("/lists/:id", handleGetList)
	api.PUT("/lists/:id", adminRequired(handleUpdateList))
	api.PATCH("/lists/:id", adminRequired(handlePatchList))
	api.DELETE("/lists/:id", adminRequired(handleDeleteList))
	api.POST("/lists/:id/push", adminRequired(handleListPush))
	e.Logger.Fatal(e.Start(":8080"))
}

// handleCreateList creates a new shopping list
// @Summary Create a new shopping list
// @Description Create a new shopping list with the provided data
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param list body ShopingList true "Shopping list data"
// @Success 201 {object} ShopingList
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal Server Error"
// @Router /lists [post]
func handleCreateList(c echo.Context) error {
	var list ShoppingList
	if err := c.Bind(&list); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	err := repository.CreateShoppingList(nil, &list)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create list")
	}
	return c.JSON(http.StatusCreated, list)
}

// handleListLists get all the shopping lists
// @Summary List all shopping lists
// @Description Get all shopping lists for the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ShopingList
// @Failure 401 {string} string "Unauthorized"
// @Router /lists [get]
func handleListLists(c echo.Context) error {
	lists, err := repository.GetAllShoppingLists()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, lists)
}

// handleDeleteList deletes a shopping list by ID
// @Summary Delete a shopping list
// @Description Delete a shopping list by its ID
// @Tags lists
// @Security BearerAuth
// @Param id path string true \"Shopping list ID\"
// @Success 204 \"No Content\"
// @Failure 401 {string} string \"Unauthorized\"
// @Failure 403 {string} string \"Forbidden\"
// @Failure 404 {string} string \"List not found\"
// @Router /lists/{id} [delete]
func handleDeleteList(c echo.Context) error {
	id := c.Param("id")
	err := repository.DeleteShoppingList(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "List not found")
	}
	listsCache.Remove(id)
	return c.NoContent(http.StatusNoContent)
}

// handleUpdateList updates a shopping list completely
// @Summary Update a shopping list
// @Description Update a shopping list with new data (full replacement)
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true \"Shopping list ID\"
// @Param list body ShopingList true \"Updated shopping list data\"
// @Success 200 {object} ShopingList
// @Failure 400 {string} string \"Bad Request\"
// @Failure 401 {string} string \"Unauthorized\"
// @Failure 403 {string} string \"Forbidden\"
// @Failure 404 {string} string \"List not found\"
// @Failure 500 {string} string \"Internal Server Error\"
// @Router /lists/{id} [put]
func handleUpdateList(c echo.Context) error {
	id := c.Param("id")
	var updatedList ShoppingList
	if err := c.Bind(&updatedList); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := repository.UpdateShoppingList(id, &updatedList)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "List not found")
	}
	listsCache.Remove(id)

	return c.JSON(http.StatusOK, updatedList)
}

// handlePatchList partially updates a shopping list
// @Summary Partially update a shopping list
// @Description Update specific fields of a shopping list (partial update)
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true \"Shopping list ID\"
// @Param patch body ShoppingListPatch true \"Partial shopping list data\"
// @Success 200 {object} ShopingList
// @Failure 400 {string} string \"Bad Request\"
// @Failure 401 {string} string \"Unauthorized\"
// @Failure 403 {string} string \"Forbidden\"
// @Failure 404 {string} string \"List not found\"
// @Failure 500 {string} string \"Internal Server Error\"
// @Router /lists/{id} [patch]
func handlePatchList(c echo.Context) error {
	id := c.Param("id")
	var patch ShoppingListPatch
	if err := c.Bind(&patch); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := repository.PatchShoppingList(id, &patch)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "List not found")
	}
	listsCache.Remove(id)

	list, err := repository.GetShoppingList(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "List not found")
	}

	return c.JSON(http.StatusOK, list)
}

// handleGetList retrieves a specific shopping list by ID
// @Summary Get a shopping list by ID
// @Description Retrieve a shopping list by its ID with caching support
// @Tags lists
// @Produce json
// @Security BearerAuth
// @Param id path string true \"Shopping list ID\"
// @Success 200 {object} ShopingList
// @Failure 401 {string} string \"Unauthorized\"
// @Failure 404 {string} string \"List not found\"
// @Failure 500 {string} string \"Internal Server Error\"
// @Router /lists/{id} [get]
func handleGetList(c echo.Context) error {
	id := c.Param("id")

	list, err := repository.GetShoppingList(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "List not found")
	}

	return c.JSON(http.StatusOK, list)
}

// handleListPush adds an item to a shopping list
// @Summary Add an item to a shopping list
// @Description Add a new item to an existing shopping list
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true \"Shopping list ID\"
// @Param item body ListPushAction true \"Item to add to the list\"
// @Success 200 {object} ShopingList
// @Failure 400 {string} string \"Bad Request\"
// @Failure 401 {string} string \"Unauthorized\"
// @Failure 403 {string} string \"Forbidden\"
// @Failure 404 {string} string \"List not found\"
// @Failure 500 {string} string \"Internal Server Error\"
// @Router /lists/{id}/push [post]
func handleListPush(c echo.Context) error {
	id := c.Param("id")
	var item ListPushAction
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	list, err := repository.GetShoppingList(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "List not found")
	}

	list.Items = append(list.Items, item.Item)
	err = repository.UpdateShoppingList(id, list)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	listsCache.Remove(id)

	return c.JSON(http.StatusOK, list)
}

// handleLogin authenticates a user and returns a session token
// @Summary User login
// @Description Authenticate user credentials and return a session token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true \"User login credentials\"
// @Success 200 {object} map[string]string \"token\"
// @Failure 401 {string} string \"Unauthorized\"
// @Failure 500 {string} string \"Internal Server Error\"
// @Router /login [post]
func handleLogin(c echo.Context) error {
	var data LoginRequest
	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	
	user, err := repository.GetUserByUsername(data.Username)
	if err == nil && user.Password == data.Password {
		session, err := repository.AddSession(user.Username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"token": session.Token})
	}
	return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
}

func authRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
		}

		token := auth[7:]
		session, err := repository.GetSession(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		c.Set("session", session)

		return next(c)
	}
}

func adminRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session := c.Get("session").(*Session)

		user, err := repository.GetUser(session.UserID)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied")
		}

		if user.Role != "admin" {
			return echo.NewHTTPError(http.StatusForbidden, "Admin access required")
		}

		return next(c)
	}
}

func addCacheHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "public, max-age=300")
			c.Response().Header().Set("Expires", time.Now().Add(5*time.Minute).Format(http.TimeFormat))
			return next(c)
		}
	}
}
