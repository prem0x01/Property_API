package handlers

import (
	"Property_App/config"
	"Property_App/models"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func InitPropertyHandler(database *sql.DB) {
	db = database
}

func PropertyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		viewProperties(w, r)
	case "POST":
		addProperty(w, r)
	case "PUT":
		updateProperty(w, r)
	case "DELETE":
		deleteProperty(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

import (
    "context"
    "encoding/json"
    "net/http"
    "Property_App/config"
    "Property_App/models"
    "database/sql"
    "encoding/base64"
    "time"
)

func viewProperties(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()

	config.Logger.Info("Checking Redis cache for properties")


    cachedProperties, err := config.RedisClient.Get(ctx, "properties").Result()
    if err == nil {
		config.Logger.Info("Serving properties from Redis cache")
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cachedProperties)) 
        return
    }

	config.Logger.Warn("Cache miss, fectcing properties from database")

    rows, err := db.Query(`SELECT  
        p.property_id, p.type, p.p_address, p.prize, p.map_link, p.img_path, p.user_id,
        u.name, u.email
        FROM properties p
        JOIN users u ON p.user_id = u.user_id`)
    if err != nil {
		config.Logger.Error("Failed to fetch properties from database", logrus.Fields{"error": err})
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var properties []struct {
        Property  models.Property `json:"property"`
        UserName  string          `json:"user_name"`
        UserEmail string          `json:"user_email"`
    }

    for rows.Next() {
        var property models.Property
        var imageData []byte
        var userName, userEmail string

        if err := rows.Scan(&property.PropertyID, &property.Type, &property.PAddress, &property.Prize, &property.MapLink, &imageData, &property.UserID,
            &userName, &userEmail); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }


        property.Img = []byte(base64.StdEncoding.EncodeToString(imageData))

        properties = append(properties, struct {
            Property  models.Property `json:"property"`
            UserName  string          `json:"user_name"`
            UserEmail string          `json:"user_email"`
        }{
            Property:  property,
            UserName:  userName,
            UserEmail: userEmail,
        })
    }

    jsonData, _ := json.Marshal(properties)

	config.Logger.Info("Successfully fetched properties from database, caching in Redis")
    config.RedisClient.Set(ctx, "properties", string(jsonData), 10*time.Minute)

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonData)

	config.Logger.Info("Response sent successfully for properties")
}


func addProperty(w http.ResponseWriter, r *http.Request) {
	var p models.Property
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt,a err := db.Prepare("INSERT INTO properties(user_id, type, p_address, prize, map_link, img) VALUES(?, ?, ?, ?, ?, ?) RETURNING property_id")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(p.UserID, p.Type, p.PAddress, p.Prize, p.MapLink, p.Img).Scan(&p.PropertyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updateProperty(w http.ResponseWriter, r *http.Request) {
	var p models.Property
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if p.PropertyID == 0 {
		http.Error(w, "Property ID is required", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("UPDATE properties SET type=?, p_address=?, prize=?, map_link=?, img=? WHERE property_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.Type, p.PAddress, p.Prize, p.MapLink, p.Img, p.PropertyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update status", http.StatusInternalServerError)
		return
	}
	
	if rowsAffected == 0 {
		http.Error(w, "No property found with the given IP", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Property updated successfully"})
}

func deleteProperty(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("property_id")
	propertyID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("DELETE FROM properties WHERE property_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err = stmt.Exec(propertyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking delete status", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No property found with the given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
