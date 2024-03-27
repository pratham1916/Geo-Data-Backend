package routes

import (
    "context"
    "encoding/json"
    "geo-data/models"
    "io/ioutil"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func CreateGeoData(client *mongo.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
            return
        }
        if err := r.ParseMultipartForm(10 << 20); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        userID := r.FormValue("user_id")
        title := r.FormValue("title")
        file, _, err := r.FormFile("file_path")
        if err != nil {
            http.Error(w, "File upload error: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer file.Close()

        fileBytes, err := ioutil.ReadAll(file)
        if err != nil {
            http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
            return
        }

        var geoJSON map[string]interface{}
        if err := json.Unmarshal(fileBytes, &geoJSON); err != nil {
            http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusInternalServerError)
            return
        }

        geometry, err := json.Marshal(geoJSON["features"].([]interface{})[0].(map[string]interface{})["geometry"])
        if err != nil {
            http.Error(w, "Error processing geometry: "+err.Error(), http.StatusInternalServerError)
            return
        }

        geoData := models.GeoData{
            UserID:   userID,
            Geometry: string(geometry),
            Title:    title,
        }

        collection := client.Database("geo-data").Collection("geodata")
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        _, err = collection.InsertOne(ctx, geoData)
        if err != nil {
            http.Error(w, "Error saving geometry to database: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": true,
            "message": "Geometry saved successfully",
        })
    }
}

func ListGeoData(client *mongo.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var geodata []models.GeoData

        collection := client.Database("geo-data").Collection("geodata")
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        cursor, err := collection.Find(ctx, bson.M{})
        if err != nil {
            http.Error(w, "Error fetching geodata: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer cursor.Close(ctx)

        for cursor.Next(ctx) {
            var gd models.GeoData
            if err = cursor.Decode(&gd); err != nil {
                http.Error(w, "Error reading geodata: "+err.Error(), http.StatusInternalServerError)
                return
            }
            geodata = append(geodata, gd)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(geodata)
    }
}

func GetGeoDataByUser(client *mongo.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.URL.Query().Get("user_id")
        if userID == "" {
            http.Error(w, "User ID is required", http.StatusBadRequest)
            return
        }

        var geodata []models.GeoData

        collection := client.Database("geo-data").Collection("geodata")
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
        if err != nil {
            http.Error(w, "Error fetching user's geodata: "+err.Error(), http.StatusInternalServerError)
            return
        }
        defer cursor.Close(ctx)

        for cursor.Next(ctx) {
            var gd models.GeoData
            if err = cursor.Decode(&gd); err != nil {
                http.Error(w, "Error reading user's geodata: "+err.Error(), http.StatusInternalServerError)
                return
            }
            geodata = append(geodata, gd)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(geodata)
    }
}
