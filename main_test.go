package main

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestAlbumsByArtist(t *testing.T) {
	// set up a database connection
	testDB, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/recordings")

	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer func(testDB *sql.DB) {
		err := testDB.Close()
		if err != nil {
			t.Fatalf("Failed to close test database: %v", err)
		}
	}(testDB)

	//Assign the test database to the global db variable
	db = testDB

	//insert the test data
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", "Test Album",
		"Test Artist", 9.99)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	//Retrieve the last inserted ID
	lastInsertedId, err := result.LastInsertId()

	defer func() {
		_, err := db.Exec("DELETE FROM album WHERE artist = ?", "Test Artist")
		if err != nil {
			t.Fatalf("Failed to delete test data: %v", err)
		}
	}()

	//Call the function to be tested
	albums, err := albumsByArtist("Test Artist")
	if err != nil {
		t.Fatalf("Failed to get albums by artist: %v", err)
	}

	//Check the result
	if len(albums) != 1 {
		t.Fatalf("Expected 1 album, got %d", len(albums))
	}

	if albums[0].Title != "Test Album" {
		t.Errorf("Expected album title Test Album, got %s", albums[0].Title)
	}

	if albums[0].Artist != "Test Artist" {
		t.Errorf("Expected album artist Test Artist, got %s", albums[0].Artist)
	}

	if albums[0].Price != 9.99 {
		t.Errorf("Expected album price 9.99, got %f", albums[0].Price)
	}

	//testing albumById
	album, err := albumById(lastInsertedId)
	if err != nil {
		t.Fatalf("Failed to get album by id: %v", err)
	}

	if album.Title != "Test Album" {
		t.Errorf("Expected album title Test Album, got %s", album.Title)
	}

	//testing deleteAlbum
	_, err = removeAlbumById(lastInsertedId)
	if err != nil {
		t.Fatalf("Failed to remove album by id: %v", err)
	}

	_, err = albumById(lastInsertedId)
	if err == nil {
		t.Errorf("Failed to get album by id: %v", err)
	} else if err.Error() != fmt.Sprintf("albumById %d: sql: no rows in result set", lastInsertedId) {
		t.Errorf("Unexpected error when getting deleted album: %v", err)
	}

}
