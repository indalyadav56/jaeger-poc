package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
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

	// Serve Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("./swagger.json"),
	))

	// scalar ui
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}
		fmt.Fprintln(w, htmlContent)
	})

	// Start server
	fmt.Println("Server started at http://localhost:8085")
	http.ListenAndServe(":8085", r)
}
