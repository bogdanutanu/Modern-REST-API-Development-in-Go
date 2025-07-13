package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var (
	token     string
	randomGen *rand.Rand
)

func random(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"value": randomGen.Intn(100)})
}

func setSeed(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var data map[string]int
	json.NewDecoder(r.Body).Decode(&data)
	randomGen = rand.New(rand.NewSource(int64(data["seed"])))
}

func login(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)
	if data["user"] == "user" && data["password"] == "password" {
		token = strconv.Itoa(rand.Intn(100000000000))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func main() {
	token = strconv.Itoa(rand.Intn(100000000000))
	randomGen = rand.New(rand.NewSource(0))
	http.HandleFunc("/random", random)
	http.HandleFunc("/seed", setSeed)
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
