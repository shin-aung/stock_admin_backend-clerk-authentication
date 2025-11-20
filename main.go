package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/rs/cors"
)

func main() {
    _ = godotenv.Load()

    if err := initDB(); err != nil {
        log.Fatalf("db init: %v", err)
    }
    defer pool.Close()

    r := mux.NewRouter()

    api := r.PathPrefix("/api").Subrouter()
    api.Handle("/admins", ClerkAuthMiddleware(http.HandlerFunc(listAdminsHandler))).Methods("GET")
    api.Handle("/admins/sync", ClerkAuthMiddleware(http.HandlerFunc(upsertAdminHandler))).Methods("POST")

    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"},
        AllowCredentials: true,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        Debug:            true,
    })

    handler := c.Handler(r)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("listening on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}
