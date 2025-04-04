package main

import (
	"database/sql"
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Movie struct {
	Title     string `json:"title"`
	Subtitles string `json:"subtitles"`
}

func connectToDatabase() *sql.DB {
	// Load environment variables
	envs, err := godotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	// Get database connection parameters from environment variables
	dbUser := envs["DB_USER"]
	dbPassword := envs["DB_PASSWORD"]
	dbHost := envs["DB_HOST"]
	dbPort := envs["DB_PORT"]
	dbName := envs["DB_NAME"]

	// Create the connection string
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database successfully!")
	return db
}

func insertMovie(db *sql.DB, movieMetadata MovieMetadata) error {
	sqlStatement := `INSERT INTO movies (title,
		year, rated, released, runtime, genre, director,
		writer, actors, plot, languages, country, awards,
		poster, metascore, imdb_rating, imdb_votes,
		imdb_id, type, dvd, box_office, production,
		website, response)
	VALUES ($1, $2, $3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12, $13,
		$14, $15, $16, $17, $18,
		$19, $20, $21, $22, $23,
		$24)
	ON CONFLICT (imdb_id) DO NOTHING`

	_, err := db.Exec(sqlStatement,
		movieMetadata.Title,
		movieMetadata.Year,
		movieMetadata.Rated,
		movieMetadata.Released,
		movieMetadata.Runtime,
		movieMetadata.Genre,
		movieMetadata.Director,
		movieMetadata.Writer,
		movieMetadata.Actors,
		movieMetadata.Plot,
		movieMetadata.Languages,
		movieMetadata.Country,
		movieMetadata.Awards,
		movieMetadata.Poster,
		movieMetadata.Metascore,
		movieMetadata.ImdbRating,
		movieMetadata.ImdbVotes,
		movieMetadata.ImdbID,
		movieMetadata.Type,
		movieMetadata.DVD,
		movieMetadata.BoxOffice,
		movieMetadata.Production,
		movieMetadata.Website,
		movieMetadata.Response)
	if err != nil {
		return fmt.Errorf("failed to insert movie: %w", err)
	}

	return nil
}

func getAllMovies(db *sql.DB) ([]Movie, error) {
	sqlStatement := `SELECT title FROM movies`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %w", err)
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.Title); err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func createTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		year VARCHAR(4),
		rated VARCHAR(10),
		released VARCHAR(20),
		runtime VARCHAR(20),
		genre VARCHAR(100),
		director VARCHAR(100),
		writer VARCHAR(100),
		actors VARCHAR(100),
		plot TEXT,
		languages VARCHAR(100),
		country VARCHAR(100),
		awards VARCHAR(100),
		poster VARCHAR(255),
		metascore VARCHAR(10),
		imdb_rating VARCHAR(10),
		imdb_votes VARCHAR(20),
		imdb_id VARCHAR(20) UNIQUE,
		type VARCHAR(20),
		dvd VARCHAR(20),
		box_office VARCHAR(20),
		production VARCHAR(100),
		website VARCHAR(255),
		response VARCHAR(10)
	)`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}
