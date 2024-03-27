package routes

import (
	"context"
	"encoding/json"
	"geo-data/auth"
	"geo-data/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usersCollection := client.Database("geo-data").Collection("users")

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Error decoding registration information", http.StatusBadRequest)
			return
		}

		user.Password, err = auth.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		result, err := usersCollection.InsertOne(ctx, user)
		if err != nil {
			http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		user.ID = result.InsertedID.(primitive.ObjectID)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func Login(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usersCollection := client.Database("geo-data").Collection("users")

		var creds models.User
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Error decoding login information", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var user models.User
		err = usersCollection.FindOne(ctx, bson.M{"email": creds.Email}).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid login credentials", http.StatusBadRequest)
			return
		}

		if !auth.CheckPasswordHash(creds.Password, user.Password) {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(user)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Login Successful", "token": token, "userId": user.ID.Hex()})
	}
}
