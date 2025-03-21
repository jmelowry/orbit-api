package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Orbiter struct {
	Orbiter  string `json:"orbiter"`
	Category string `json:"category"`
	Scope    int    `json:"scope"`
	Priority int    `json:"priority"`
}

type OrbiterResponse struct {
	Orbiter  string `json:"orbiter"`
	Category string `json:"category"`
	Scope    int    `json:"scope"`
	Priority int    `json:"priority"`
	Angle    int    `json:"angle"`
	Distance int    `json:"distance"`
	Size     int    `json:"size"`
}

var orbiters = make(map[string][]Orbiter)
var users = make(map[string]User)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	// Expected URL: /api/users or /api/users/{username}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[1] != "api" || parts[2] != "users" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// If a username is provided, it'll be in parts[3]
	var username string
	if len(parts) >= 4 && parts[3] != "" {
		username = parts[3]
	}

	switch r.Method {
	case "POST":
		// Create a new user; expect the user JSON in the request body.
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, exists := users[user.Username]; exists {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}
		users[user.Username] = user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	case "PUT":
		// Update an existing user; the username must be in the URL.
		if username == "" {
			http.Error(w, "Username is required in URL", http.StatusBadRequest)
			return
		}
		var userUpdate User
		if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		existingUser, exists := users[username]
		if !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		// For simplicity, only update the email.
		if userUpdate.Email != "" {
			existingUser.Email = userUpdate.Email
		}
		users[username] = existingUser
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingUser)
	case "DELETE":
		// Delete a user; the username must be provided in the URL.
		if username == "" {
			http.Error(w, "Username is required in URL", http.StatusBadRequest)
			return
		}
		if _, exists := users[username]; !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		delete(users, username)
		w.WriteHeader(http.StatusNoContent)
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		// If no username is provided, return all users.
		if username == "" {
			var allUsers []User
			for _, u := range users {
				allUsers = append(allUsers, u)
			}
			json.NewEncoder(w).Encode(allUsers)
		} else {
			if user, exists := users[username]; exists {
				json.NewEncoder(w).Encode(user)
			} else {
				http.Error(w, "User not found", http.StatusNotFound)
			}
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func orbitersHandler(w http.ResponseWriter, r *http.Request) {
	// Expected URL: /{username}/api/orbiters
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[2] != "api" || parts[3] != "orbiters" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	user := parts[1]
	switch r.Method {
	case "POST":
		var orb Orbiter
		if err := json.NewDecoder(r.Body).Decode(&orb); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		orbiters[user] = append(orbiters[user], orb)
		var responseData []OrbiterResponse
		for _, o := range orbiters[user] {
			responseData = append(responseData, OrbiterResponse{
				Orbiter:  o.Orbiter,
				Category: o.Category,
				Scope:    o.Scope,
				Priority: o.Priority,
				Angle:    o.Priority * 15,
				Distance: o.Scope * 3,
				Size:     o.Scope,
			})
		}
		resp := map[string]interface{}{
			"starmap":      fmt.Sprintf("/%s/api/starmap", user),
			"orbital_data": responseData,
			"log":          fmt.Sprintf("New celestial object detected in %s's system.", user),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	case "GET":
		var responseData []OrbiterResponse
		for _, o := range orbiters[user] {
			responseData = append(responseData, OrbiterResponse{
				Orbiter:  o.Orbiter,
				Category: o.Category,
				Scope:    o.Scope,
				Priority: o.Priority,
				Angle:    o.Priority * 15,
				Distance: o.Scope * 3,
				Size:     o.Scope,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {

	http.HandleFunc("/api/users/", usersHandler)
	http.HandleFunc("/api/users", usersHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/orbiters") {
			orbitersHandler(w, r)
			return
		}
		fmt.Fprintf(w, "Orbit API 🚀")
	})
	port := "8080"
	log.Println("Server starting on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
