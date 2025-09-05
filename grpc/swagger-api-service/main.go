package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Serve swagger.json as a static file
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger.json")
	})

	// Serve Swagger UI and tell it where swagger.json is located
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8085/swagger.json"),
	))

	// // scalar ui
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
	// 		SpecURL: "./swagger.json",
	// 		CustomOptions: scalar.CustomOptions{
	// 			PageTitle: "Simple API",
	// 		},
	// 		DarkMode: true,
	// 	})

	// 	if err != nil {
	// 		fmt.Printf("%v", err)
	// 	}
	// 	fmt.Fprintln(w, htmlContent)
	// })

	// Start server
	fmt.Println("Server started at http://localhost:8085")
	http.ListenAndServe(":8085", r)
}
