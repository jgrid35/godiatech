package main

import (
	"fmt"
	iofs "io/fs"
	"net"
	"regexp"

	"github.com/hirochachacha/go-smb2"
	"github.com/joho/godotenv"
)

type MovieFile struct {
	file     string
	srtFile  string
	imdbFile string
}

func importMovies() {
	movieChan := make(chan MovieFile)
	quit := make(chan int)
	db := connectToDatabase()
	go readSambaFiles(movieChan, quit)
	for {
		select {
		case movie := <-movieChan:
			var movieMetadata MovieMetadata
			if movie.imdbFile != "" {
				movieMetadata = getMovieMetadataByID(movie.imdbFile)
			} else if movie.file != "" {
				movieMetadata = getMovieMetadataByTitle(movie.file)
			}
			fmt.Println("Movie Metadata:", movieMetadata)
			err := insertMovie(db, movieMetadata)
			if err != nil {
				panic(err)
			}
		case <-quit:
			return
		}
	}
}

func readSambaFiles(c chan MovieFile, quit chan int) {
	var envs map[string]string
	envs, err := godotenv.Read(".env")

	freeboxUser := envs["FREEBOX_USERNAME"]
	freeboxPassword := envs["FREEBOX_PASSWORD"]
	freeboxIP := envs["FREEBOX_IP"]
	freeboxSmbPort := envs["FREEBOX_SMB_PORT"]
	freeboxMediaFolder := envs["FREEBOX_MEDIA_FOLDER"]

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", freeboxIP, freeboxSmbPort))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     freeboxUser,
			Password: freeboxPassword,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer s.Logoff()

	fs, err := s.Mount("Freebox")
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	matches, err := iofs.Glob(fs.DirFS(freeboxMediaFolder), "*")
	if err != nil {
		panic(err)
	}

	imdbFileRegexp := "^tt[0-9]*"
	fileNameRegexp := "(\\.mkv|\\.avi|\\.mp4)$"
	srtFileNameRegexp := "(\\.srt)$"

	for _, match := range matches {
		matches, _ := iofs.Glob(fs.DirFS(fmt.Sprintf("%v/%v", freeboxMediaFolder, match)), "*")
		var movie MovieFile

		for _, match := range matches {
			imdbFileMatch, _ := regexp.MatchString(imdbFileRegexp, match)
			if imdbFileMatch {
				movie.imdbFile = match
			}

			fileMatch, _ := regexp.MatchString(fileNameRegexp, match)
			if fileMatch {
				movie.file = match
			}

			srtFileMatch, _ := regexp.MatchString(srtFileNameRegexp, match)
			if srtFileMatch {
				movie.srtFile = match
			}
		}

		c <- movie
	}
	fmt.Println("No more files to read")
	quit <- 0
	fmt.Println("Exiting readSmbFiles")
}
