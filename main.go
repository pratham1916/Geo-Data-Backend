package main

import (
	"fmt"
	"geo-data/config"
	"geo-data/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := config.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to the database: ", err)
		return
	}

	http.HandleFunc("/register", corsMiddleware(routes.Register(client)))
	http.HandleFunc("/login", corsMiddleware(routes.Login(client)))
	http.HandleFunc("/geodata", corsMiddleware(routes.CreateGeoData(client)))
	http.HandleFunc("/geodata/list", corsMiddleware(routes.ListGeoData(client)))
	http.HandleFunc("/geodata/user", corsMiddleware(routes.GetGeoDataByUser(client)))

	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}



// http.HandleFunc("/shape", corsMiddleware(routes.CreateShape(client, ctx)))
// 	http.HandleFunc("/shape/update", corsMiddleware(routes.UpdateShape(client, ctx)))
// 	http.HandleFunc("/shape/delete", corsMiddleware(routes.DeleteShape(client, ctx)))