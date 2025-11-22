package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Reponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/greet/", greetHandler)

	fmt.Println("Starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "Getting user data")
	case "POST":
		fmt.Fprintf(w, "Creating new user")
	case "PUT":
		fmt.Fprintf(w, "Updating user")
	case "DELETE":
		fmt.Fprintf(w, "Deleting user")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid JSON")
		return
	}

	fmt.Fprintf(w, "User created: %s (%s)", user.Name, user.Email)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my Go server!")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Reponse{
		Message: "Hello from the API!",
		Status:  200,
	}

	json.NewEncoder(w).Encode(response)
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/greet/"):]
	fmt.Fprintf(w, "Hello, %s!", name)
}
