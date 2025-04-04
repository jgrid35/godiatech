package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joho/godotenv"
)

type MovieMetadata struct {
	Title     string `json:"Title"`
	Year      string `json:"Year"`
	Rated     string `json:"Rated"`
	Released  string `json:"Released"`
	Runtime   string `json:"Runtime"`
	Genre     string `json:"Genre"`
	Director  string `json:"Director"`
	Writer    string `json:"Writer"`
	Actors    string `json:"Actors"`
	Plot      string `json:"Plot"`
	Languages string `json:"Language"`
	Country   string `json:"Country"`
	Awards    string `json:"Awards"`
	Poster    string `json:"Poster"`
	Ratings   []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	DVD        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}

func getMovieMetadata(queryType string) func(movie string) MovieMetadata {
	return func(movie string) MovieMetadata {
		envs, err := godotenv.Read(".env")
		apiKey := envs["OMDB_API_KEY"]
		res, err := http.Get("http://www.omdbapi.com/?" + queryType + "=" + movie + "&apikey=" + apiKey)
		if err != nil {
			fmt.Println("Error fetching movie metadata:", err)
		}

		var movieMetadata MovieMetadata
		resBody, err := io.ReadAll(res.Body)
		err = json.Unmarshal(resBody, &movieMetadata)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		fmt.Println("Response Body:", movieMetadata)
		return movieMetadata
	}
}

var (
	getMovieMetadataByID    = getMovieMetadata("i")
	getMovieMetadataByTitle = getMovieMetadata("t")
)
