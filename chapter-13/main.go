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
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	http.HandleFunc("GET /lists", addCacheHeaders(authRequired(handleListLists)))
	http.HandleFunc("POST /lists", MetricsMiddleware(adminRequired(handleCreateList)))
	http.HandleFunc("GET /lists/{id}", authRequired(handleGetList))
	http.HandleFunc("PUT /lists/{id}", adminRequired(handleUpdateList))
	http.HandleFunc("DELETE /lists/{id}", adminRequired(handleDeleteList))
	http.HandleFunc("PATCH /lists/{id}", adminRequired(handlePatchList))
	http.HandleFunc("POST /lists/{id}/push", adminRequired(handleListPush))
	http.HandleFunc("POST /login", handleLogin)

	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("listening on port :8888")
	http.ListenAndServe(":8888", nil)
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
func handleCreateList(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Creating new shopping list",
		slog.String("ip", r.RemoteAddr),
		slog.String("user", r.Header.Get("X-User")),
		slog.String("request_id", r.Header.Get("X-Request-ID")),
	)
	var list ShoppingList
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		slog.Info("Invalid request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list.ID = uint(rand.Int())

	parentSpan := opentracing.GlobalTracer().StartSpan("handleCreateList")
	defer parentSpan.Finish()
	ext.HTTPMethod.Set(parentSpan, r.Method)
	ext.HTTPUrl.Set(parentSpan, r.URL.Path)

	err = repository.CreateShoppingList(parentSpan, &list)
	if err != nil {
		slog.Error("Failed to create new shopping list", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
func handleListLists(w http.ResponseWriter, r *http.Request) {
	lists, err := repository.GetAllShoppingLists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(lists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
func handleDeleteList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := repository.DeleteShoppingList(id)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}
	listsCache.Remove(id)
	w.WriteHeader(http.StatusNoContent)
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
func handleUpdateList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updatedList ShoppingList
	err := json.NewDecoder(r.Body).Decode(&updatedList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repository.UpdateShoppingList(id, &updatedList)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}
	listsCache.Remove(id)

	if err := json.NewEncoder(w).Encode(updatedList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
func handlePatchList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var patch ShoppingListPatch
	err := json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = repository.PatchShoppingList(id, &patch)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}
	listsCache.Remove(id)

	list, err := repository.GetShoppingList(id)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
func handleGetList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	list, ok := listsCache.Get(id)
	if !ok {
		dbList, err := repository.GetShoppingList(id)
		if err != nil {
			http.Error(w, "List not found", http.StatusNotFound)
			return
		}
		listsCache.Add(id, *dbList)
	}

	data, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
func handleListPush(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var item ListPushAction
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list, err := repository.GetShoppingList(id)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	list.Items = append(list.Items, item.Item)
	err = repository.UpdateShoppingList(id, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	listsCache.Remove(id)

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var data LoginRequest
	json.NewDecoder(r.Body).Decode(&data)
	user, err := repository.GetUserByUsername(data.Username)
	if err == nil && user.Password == data.Password {
		session, err := repository.AddSession(user.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": session.Token})
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func authRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token = token[7:]
		_, err := repository.GetSession(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func adminRequired(next http.HandlerFunc) http.HandlerFunc {
	return authRequired(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = token[7:]
		session, err := repository.GetSession(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		user, err := repository.GetUser(session.UserID)
		if err != nil || user.Role != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

func addCacheHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=300")
		w.Header().Set("Expires", time.Now().Add(5*time.Minute).Format(http.TimeFormat))
		next(w, r)
	}
}
