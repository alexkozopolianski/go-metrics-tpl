package main

import (
	"net/http"

	handler "github.com/alexkozopolianski/go-metrics-tpl/internal/handlers"
	"github.com/alexkozopolianski/go-metrics-tpl/internal/storage"
)

func main() {
	mux := http.NewServeMux()

	memStorage := storage.NewMemStorage()
	handler := handler.NewHandler(memStorage)

	mux.HandleFunc("/update/", handler.Update)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
