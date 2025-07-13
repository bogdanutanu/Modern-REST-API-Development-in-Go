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
	AddSession(username string) (*Session, error)
	GetSession(token string) (*Session, error)
	GetUserByUsername(username string) (*User, error)
	GetUser(userID uint) (*User, error)
	CreateShoppingList(parentSpan opentracing.Span, list *ShoppingList) error
	GetAllShoppingLists() ([]ShoppingList, error)
	GetShoppingList(id string) (*ShoppingList, error)
	UpdateShoppingList(id string, list *ShoppingList) error
	DeleteShoppingList(id string) error
	PatchShoppingList(id string, patch *ShoppingListPatch) error
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

func (r *Repository) GetUser(userID uint) (*User, error) {
	var user User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	var user User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repository) AddSession(username string) (*Session, error) {
	var user User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	token := strconv.Itoa(rand.Intn(100000000000))
	session := Session{Token: token, Expires: time.Now().Add(7 * 24 * time.Hour), UserID: user.ID}
	createResult := r.db.Create(&session)
	if createResult.Error != nil {
		return nil, createResult.Error
	}
	return &session, nil
}

func (r *Repository) GetSession(token string) (*Session, error) {
	var session Session
	result := r.db.Where("token = ? AND expires > ?", token, time.Now()).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r *Repository) CreateShoppingList(parentSpan opentracing.Span, list *ShoppingList) error {
	span := opentracing.StartSpan("AddShopingList",
		opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	span.LogKV("my-custom-data", "relevant data in the trace")
	result := r.db.Create(list)
	return result.Error
}

func (r *Repository) GetAllShoppingLists() ([]ShoppingList, error) {
	var results []ShoppingList
	result := r.db.Find(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (r *Repository) GetShoppingList(id string) (*ShoppingList, error) {
	var list ShoppingList
	result := r.db.Where("id = ?", id).First(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return &list, nil
}

func (r *Repository) UpdateShoppingList(id string, list *ShoppingList) error {
	result := r.db.Model(&ShoppingList{}).Where("id = ?", id).Updates(list)
	return result.Error
}

func (r *Repository) DeleteShoppingList(id string) error {
	result := r.db.Where("id = ?", id).Delete(&ShoppingList{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) PatchShoppingList(id string, patch *ShoppingListPatch) error {
	result := r.db.Model(&ShoppingList{}).Where("id = ?", id).Updates(patch)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
