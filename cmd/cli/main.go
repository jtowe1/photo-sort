package main

import (
	"flag"
	"fmt"
	"github.com/jtowe1/photo-sort/service/faces"
	"github.com/jtowe1/photo-sort/service/sort"
	"log"
	"os"
)

func main() {
	var pathToPhotos string
	flag.StringVar(
		&pathToPhotos,
		"pathToPhotos",
		"",
		"The path to the photos to sort",
	)
	flag.Parse()

	_, err := os.Stat(pathToPhotos)
	if err != nil {
		log.Printf("problem with path: %s", pathToPhotos)
		log.Fatal(err)
	}

	faceService, err := faces.NewFaceService()
	if err != nil {
		log.Fatal(err)
	}

	serviceSort := sort.ServiceSort{
		FaceService: faceService,
	}

	err = serviceSort.PhotosByLocalPath(pathToPhotos)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"photos sorted and placed in %s",
		pathToPhotos+string(os.PathSeparator)+"sorted",
	)
}
