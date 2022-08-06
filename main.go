package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "sikanderkhan"
	password = "password"
	dbname = "recordings"
)

type Album struct {
	ID int64
	TITLE string
	ARTIST string
	PRICE float64
}

func main() {
	psqlconn := fmt.Sprintf("host= %s port = %d user = %s password =%s dbname = %s sslmode=disable", host, port,user, password, dbname)

	// Making db connection
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err, "conneting..")
	defer db.Close()
	
	// Pinging db
	pingErr := db.Ping()
	CheckError(pingErr, "pinging..")
	fmt.Println("Connected")

	// Getting albums by artist name
	albums, err := albumsByArtist("SIKANDER", db)
	CheckError(err, "albumByArtist Call")
	fmt.Printf("Albums found by Artist: %v \n", albums)

	// Getting album by id
	album, err := albumById(4, db)
	CheckError(err, "albumById Call")
	fmt.Printf("Album found by Id: %v \n", album)

	// Insert new album
	albId, err := insertAlbum(Album{
		ID: 9,
		TITLE: "MY SONG2",
		ARTIST: "SIKANDER",
		PRICE: 78.53,
	}, db)
	CheckError(err, "insertAlbum")
	fmt.Printf("Created album with id %d \n", albId)
}

func CheckError(err error, msg string){
	if err != nil {
		log.Fatalf("Error while %s with follwoing error %s", msg, err)
	}
}

func albumsByArtist(name string, db *sql.DB) ([]Album, error){
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	CheckError(err, "selecting..")

	for rows.Next(){
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.TITLE, &alb.ARTIST, &alb.PRICE); err != nil {
			return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	rowError := rows.Err()
	if rowError != nil{
		return nil, fmt.Errorf("albumByArtist %q: %v", name, err)
	}	

	return albums, nil
}

func albumById(id int64, db *sql.DB) (Album, error) {
	var alb Album
	query := `SELECT * FROM album WHERE id =$1;`
	rows := db.QueryRow(query, id)

	if err := rows.Scan(&alb.ID, &alb.TITLE, &alb.ARTIST, &alb.PRICE); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("no rows found err msg: %v", err)
		}
		return alb, fmt.Errorf("albumById %q: %v", id, err)
	}
	return alb, nil 
}

func insertAlbum(album Album, db *sql.DB) (int64, error){
	query := `INSERT INTO album VALUES ($1, $2, $3, $4) RETURNING id;`
	// , album.ID, album.TITLE, album.ARTIST, album.PRICE
	lastid := 0
	fmt.Println(query)
	err := db.QueryRow(query, album.ID, album.TITLE, album.ARTIST, album.PRICE).Scan(&lastid)
	if err != nil {
		return int64(lastid), err
	}

	return int64(lastid), nil
}