package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	result, err := http.Get("http://localhost:8888/")
	if err != nil {
		panic(err)
	}
	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response:", string(data))
}
