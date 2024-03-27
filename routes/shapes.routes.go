package routes

import (
	"encoding/json"
	"net/http"
	"geo-data/models"
	"gorm.io/gorm"
	"strconv"
)

func CreateShape(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		if userID == "" {
			http.Error(w, "UserID is required", http.StatusBadRequest)
			return
		}

		var updates map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var geoData models.GeoData
		result := db.Where("user_id = ?", userID).Order("created_at desc").First(&geoData)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusNotFound)
			return
		}

		updateResult := db.Model(&geoData).Updates(updates)
		if updateResult.Error != nil {
			http.Error(w, updateResult.Error.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(geoData)
	}
}

func UpdateShape(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		idStr := query.Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var geoData models.GeoData
		if err := db.First(&geoData, id).Error; err != nil {
			http.Error(w, "GeoData not found", http.StatusNotFound)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&geoData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		db.Save(&geoData)
		json.NewEncoder(w).Encode(geoData)
	}
}

func DeleteShape(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		idStr := query.Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var geoData models.GeoData
		if err := db.Delete(&geoData, id).Error; err != nil {
			http.Error(w, "GeoData not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
