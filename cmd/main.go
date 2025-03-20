package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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

func orbitersHandler(w http.ResponseWriter, r *http.Request) {
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
