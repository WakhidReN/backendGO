package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type asset struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Kategori string `json:"kategori"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_assets")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
}

func getAssets() ([]asset, error) {
	rows, err := db.Query("SELECT id, name, kategori FROM assets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []asset

	for rows.Next() {
		var a asset
		err := rows.Scan(&a.ID, &a.Name, &a.Kategori)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	return assets, nil
}

func getAssetDetails(id string) (asset, error) {
	var a asset
	err := db.QueryRow("SELECT id, name, kategori FROM assets WHERE id = ?", id).Scan(&a.ID, &a.Name, &a.Kategori)
	if err != nil {
		if err == sql.ErrNoRows {
			return a, err
		}
		return a, err
	}

	return a, nil
}

func assets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		assets, err := getAssets()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(assets)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func assetDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "GET" {
		id := r.FormValue("id")
		a, err := getAssetDetails(id)

		if err != nil {
			if err != sql.ErrNoRows {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			result, err := json.Marshal(a)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(result)
			return
		}
	}

	http.Error(w, "", http.StatusBadRequest)
}

func main() {
	initDB()

	http.HandleFunc("/assets", assets)
	http.HandleFunc("/asset", assetDetails)

	fmt.Println("Starting web server at http://localhost:8084/")
	http.ListenAndServe(":8084", nil)
}