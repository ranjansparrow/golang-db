package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

var db *sql.DB

// added json for web server
type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func main() {
	//load environment variables from .env file
	enVError := godotenv.Load()
	if enVError != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	//Get a database handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	http.HandleFunc("/albums", getAlbumsByArtist)
	http.HandleFunc("/album", getAlbumById)
	http.HandleFunc("/addAlbum", addNewAlbum)
	http.HandleFunc("/deleteById", deleteAlbumById)

	log.Fatal(http.ListenAndServe(":8080", nil))

	//albums, err := albumsByArtist("Pink Floyd")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Albums found: %v\n", albums)

	//alb, err := albumById(4)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Album found: %v\n", alb)
	//
	//albId, err := addAlbum(Album{
	//	Title:  "The Modern Sound of Betty Carter",
	//	Artist: "Betty Carter",
	//	Price:  49.99,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("ID of added album: %d\n", albId)
}

// albumsByArtist returns all albums for a given artist name
func albumsByArtist(name string) ([]Album, error) {
	//An albums slice to hold data from the returned rows
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %s: %v", name, err)

	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("albumsByArtist %s: could not close rows: %v", name, err)
		}
	}(rows)
	//Loop through rows, using Scan to assign column data to struct fields
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %s: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %s: %v", name, err)
	}
	return albums, nil
}

func albumById(id int64) (Album, error) {
	//An album to hold data from the returned row
	var alb Album

	ranjan := db.QueryRow("select * from album where id =?", id)
	if err := ranjan.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumById %d: %v", id, err)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("insert into album (title, artist, price) values (?,?,?)",
		alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

// web server functions
func getAlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	artist := r.URL.Query().Get("artist")
	albums, err := albumsByArtist(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(albums)
	if err != nil {
		return
	}
}

func getAlbumById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}
	album, err := albumById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(album)
	if err != nil {
		return
	}
}

func addNewAlbum(w http.ResponseWriter, r *http.Request) {
	var alb Album
	if err := json.NewDecoder(r.Body).Decode(&alb); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := addAlbum(alb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]int64{"id": id})
	if err != nil {
		return
	}
}

func removeAlbumById(id int64) (int64, error) {
	//An album to hold data from the returned row

	result, err := db.Exec("Delete from album where id = ?", id)

	if err != nil {
		log.Fatal("Error deleting album")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("removeAlbumById %d: no rows affected", id)
	}

	return rowsAffected, nil
}

func deleteAlbumById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	_, err = removeAlbumById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Album deleted successfully"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
