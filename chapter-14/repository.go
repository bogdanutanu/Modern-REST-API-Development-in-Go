package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Init() error
	AddSession(parentSpan opentracing.Span, username string) (*Session, error)
	GetSession(parentSpan opentracing.Span, token string) (*Session, error)
	GetUserByUsername(parentSpan opentracing.Span, username string) (*User, error)
	GetUser(parentSpan opentracing.Span, userID uint) (*User, error)
	CreateShoppingList(parentSpan opentracing.Span, list *ShoppingList) error
	GetAllShoppingLists(parentSpan opentracing.Span) ([]ShoppingList, error)
	GetShoppingList(parentSpan opentracing.Span, id string) (*ShoppingList, error)
	UpdateShoppingList(parentSpan opentracing.Span, id string, list *ShoppingList) error
	DeleteShoppingList(parentSpan opentracing.Span, id string) error
	PatchShoppingList(parentSpan opentracing.Span, id string, patch *ShoppingListPatch) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(database string) (*Repository, error) {
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, &Session{}, &ShoppingList{})
	if err != nil {
		return nil, err
	}

	return &Repository{db}, nil
}

func (r *Repository) Init() error {
	adminUser := User{
		Username: "admin",
		Password: "password",
		Role:     "admin",
	}
	result := r.db.Where(User{Username: "admin"}).FirstOrCreate(&adminUser)
	return result.Error
}

func (r *Repository) GetUser(parentSpan opentracing.Span, userID uint) (*User, error) {
	span := childSpan(parentSpan, "GetUser")
	defer span.Finish()

	var user User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(parentSpan opentracing.Span, username string) (*User, error) {
	span := childSpan(parentSpan, "GetUserByUsername")
	defer span.Finish()

	var user User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repository) AddSession(parentSpan opentracing.Span, username string) (*Session, error) {
	span := childSpan(parentSpan, "AddSession")
	defer span.Finish()

	user, err := r.GetUserByUsername(span, username)
	if err != nil {
		return nil, err
	}
	token := strconv.Itoa(rand.Intn(100000000000))
	session := Session{Token: token, Expires: time.Now().Add(7 * 24 * time.Hour), UserID: user.ID}
	createResult := r.db.Create(&session)
	if createResult.Error != nil {
		return nil, createResult.Error
	}
	return &session, nil
}

func (r *Repository) GetSession(parentSpan opentracing.Span, token string) (*Session, error) {
	span := childSpan(parentSpan, "GetSession")
	defer span.Finish()

	var session Session
	result := r.db.Where("token = ? AND expires > ?", token, time.Now()).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r *Repository) CreateShoppingList(parentSpan opentracing.Span, list *ShoppingList) error {
	span := childSpan(parentSpan, "AddShopingList")
	defer span.Finish()
	span.LogKV("my-custom-data", "relevant data in the trace")
	result := r.db.Create(list)
	return result.Error
}

func (r *Repository) GetAllShoppingLists(parentSpan opentracing.Span) ([]ShoppingList, error) {
	span := childSpan(parentSpan, "GetAllShoppingLists")
	defer span.Finish()

	var results []ShoppingList
	result := r.db.Find(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (r *Repository) GetShoppingList(parentSpan opentracing.Span, id string) (*ShoppingList, error) {
	span := childSpan(parentSpan, "GetShoppingList")
	defer span.Finish()

	var list ShoppingList
	result := r.db.Where("id = ?", id).First(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return &list, nil
}

func (r *Repository) UpdateShoppingList(parentSpan opentracing.Span, id string, list *ShoppingList) error {
	span := childSpan(parentSpan, "UpdateShoppingList")
	defer span.Finish()

	result := r.db.Model(&ShoppingList{}).Where("id = ?", id).Updates(list)
	return result.Error
}

func (r *Repository) DeleteShoppingList(parentSpan opentracing.Span, id string) error {
	span := childSpan(parentSpan, "DeleteShoppingList")
	defer span.Finish()

	result := r.db.Where("id = ?", id).Delete(&ShoppingList{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) PatchShoppingList(parentSpan opentracing.Span, id string, patch *ShoppingListPatch) error {
	span := childSpan(parentSpan, "PatchShoppingList")
	defer span.Finish()

	result := r.db.Model(&ShoppingList{}).Where("id = ?", id).Updates(patch)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func childSpan(parentSpan opentracing.Span, name string) opentracing.Span {
	return opentracing.StartSpan(name, opentracing.ChildOf(parentSpan.Context()))
}
