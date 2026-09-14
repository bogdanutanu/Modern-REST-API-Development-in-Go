package main

import (
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"github.com/opentracing/opentracing-go"
)

type RepositoryInterface interface {
	Init() error
	AddSession(parentSpan opentracing.Span, username string) (*Session, error)
	GetSession(parentSpan opentracing.Span, token string) (*Session, error)
	CreateShoppingList(parentSpan opentracing.Span, list *ShoppingList) error
	GetAllShoppingLists(parentSpan opentracing.Span) ([]ShoppingList, error)
	GetShoppingList(parentSpan opentracing.Span, id string) (*ShoppingList, error)
	UpdateShoppingList(parentSpan opentracing.Span, id string, list *ShoppingList) error
	DeleteShoppingList(parentSpan opentracing.Span, id string) error
	PatchShoppingList(parentSpan opentracing.Span, id string, patch *ShoppingListPatch) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(database string) (*Repository, error) {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

func (r *Repository) Init() error {
	if _, err := r.db.Exec("CREATE TABLE IF NOT EXISTS users (role VARCHAR, username VARCHAR PRIMARY KEY, password VARCHAR)"); err != nil {
		return err
	}
	if _, err := r.db.Exec("CREATE TABLE IF NOT EXISTS sessions (token VARCHAR PRIMARY KEY, expires TIMESTAMP, username VARCHAR)"); err != nil {
		return err
	}
	if _, err := r.db.Exec("CREATE TABLE IF NOT EXISTS shopping_lists (id VARCHAR PRIMARY KEY, name VARCHAR, items TEXT)"); err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddSession(parentSpan opentracing.Span, username string) (*Session, error) {
	span := opentracing.StartSpan("AddSessionRepository", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	token := strconv.Itoa(rand.Intn(100000000000))
	session := Session{Token: token, Expires: time.Now().Add(7 * 24 * time.Hour), Username: username}
	query := sq.Insert("sessions").Columns("token", "expires", "username").Values(session.Token, session.Expires, session.Username)
	_, err := query.RunWith(r.db).Exec()
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) GetSession(parentSpan opentracing.Span, token string) (*Session, error) {
	span := opentracing.StartSpan("GetSessionRepository", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	query := sq.Select("token", "expires", "username").From("sessions").Where(sq.Eq{"token": token}, sq.Gt{"expires": time.Now()})
	row := query.RunWith(r.db).QueryRow()
	session := Session{}
	if err := row.Scan(&session.Token, &session.Expires, &session.Username); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) CreateShoppingList(parentSpan opentracing.Span, list *ShoppingList) error {
	span := opentracing.StartSpan("CreateShoppingListRepository",
		opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	span.LogKV("my-custom-data", "relevant data in the trace")
	query := sq.Insert("shopping_lists").Columns("id", "name", "items").Values(strconv.Itoa(list.ID), list.Name, strings.Join(list.Items, ","))
	_, err := query.RunWith(r.db).Exec()
	return err
}

func (r *Repository) GetAllShoppingLists(parentSpan opentracing.Span) ([]ShoppingList, error) {
	span := opentracing.StartSpan("GetAllShoppingListsRepository", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	query := sq.Select("id", "name", "items").From("shopping_lists")
	rows, err := query.RunWith(r.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []ShoppingList
	for rows.Next() {
		var list ShoppingList
		var idStr, itemsStr string
		if err := rows.Scan(&idStr, &list.Name, &itemsStr); err != nil {
			return nil, err
		}
		list.ID, _ = strconv.Atoi(idStr)
		if itemsStr != "" {
			list.Items = strings.Split(itemsStr, ",")
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (r *Repository) GetShoppingList(parentSpan opentracing.Span, id string) (*ShoppingList, error) {
	span := opentracing.StartSpan("GetShoppingListRepository", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	query := sq.Select("id", "name", "items").From("shopping_lists").Where(sq.Eq{"id": id})
	row := query.RunWith(r.db).QueryRow()

	var list ShoppingList
	var idStr, itemsStr string
	if err := row.Scan(&idStr, &list.Name, &itemsStr); err != nil {
		return nil, err
	}
	list.ID, _ = strconv.Atoi(idStr)
	if itemsStr != "" {
		list.Items = strings.Split(itemsStr, ",")
	}
	return &list, nil
}

func (r *Repository) UpdateShoppingList(parentSpan opentracing.Span, id string, list *ShoppingList) error {
	span := opentracing.StartSpan("UpdateShoppingListRepository", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	query := sq.Update("shopping_lists").Where(sq.Eq{"id": id}).Set("name", list.Name).Set("items", strings.Join(list.Items, ","))
	_, err := query.RunWith(r.db).Exec()
	return err
}

func (r *Repository) DeleteShoppingList(parentSpan opentracing.Span, id string) error {
	span := opentracing.StartSpan("DeleteShoppingListRepository", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	query := sq.Delete("shopping_lists").Where(sq.Eq{"id": id})
	_, err := query.RunWith(r.db).Exec()
	return err
}

func (r *Repository) PatchShoppingList(parentSpan opentracing.Span, id string, patch *ShoppingListPatch) error {
	span := opentracing.StartSpan("PatchShoppingListRepository", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()
	query := sq.Update("shopping_lists").Where(sq.Eq{"id": id})
	if patch.Name != nil {
		query = query.Set("name", *patch.Name)
	}
	if patch.Items != nil {
		query = query.Set("items", strings.Join(patch.Items, ","))
	}
	_, err := query.RunWith(r.db).Exec()
	if err != nil {
		return err
	}
	return nil
}
