package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/handlers" 
)

type User struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

var users = make(map[string]string) 

func registerHandler(w http.ResponseWriter, r *http.Request) {
    
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusNoContent)
        return
    }

    if r.Method == http.MethodPost {
        var user User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil || user.Email == "" || user.Password == "" {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        
        if _, exists := users[user.Email]; exists {
            http.Error(w, "User already exists", http.StatusConflict)
            return
        }

        users[user.Email] = user.Password
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Configuração de CORS
        w.WriteHeader(http.StatusCreated)
        fmt.Fprintln(w, "User registered successfully")
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusNoContent)
        return
    }

    if r.Method == http.MethodPost {
        var user User
        err := json.NewDecoder(r.Body).Decode(&user)
        if err != nil || user.Email == "" || user.Password == "" {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        
        storedPassword, exists := users[user.Email]
        if !exists || storedPassword != user.Password {
            http.Error(w, "Invalid email or password", http.StatusUnauthorized)
            return
        }

        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") 
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Login successful")
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func main() {
    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)

    
    corsObj := handlers.AllowedOrigins([]string{"http://localhost:3000"}) 
    corsHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
    corsMethods := handlers.AllowedMethods([]string{"POST", "OPTIONS"})

    
    fmt.Println("Servidor rodando na porta 8080...")
    http.ListenAndServe(":8080", handlers.CORS(corsHeaders, corsMethods, corsObj)(http.DefaultServeMux))
}
