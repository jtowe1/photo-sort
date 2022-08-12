package sort

import (
	"fmt"
	"github.com/jtowe1/photo-sort/service/faces"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"os"
	"path"
	"strconv"
)

func PhotosByLocalPath(pathToPhotos string) {
	faceService := faces.NewFaceService()
	photos, err := os.ReadDir(pathToPhotos)
	if err != nil {
		log.Fatal(err)
	}
	unsortedPhotos := make([]os.DirEntry, len(photos))
	copy(unsortedPhotos, photos)

	albumNameBase := "album"
	albumNameIndex := 0
	albums := make(map[string][]string)
	var matchedPhotos []string

	for _, photo := range photos {
		pathToPhoto := pathToPhotos + string(os.PathSeparator) + photo.Name()
		ext := path.Ext(pathToPhoto)
		if ext != ".jpg" && ext != ".jpeg" {
			log.Printf("non jpg/jpeg detected; skipping %s\n", pathToPhoto)
			continue
		}
		if slices.Contains(matchedPhotos, pathToPhoto) {
			continue
		}
		for _, unsortedPhoto := range unsortedPhotos {
			pathToUnsortedPhoto := pathToPhotos + string(os.PathSeparator) + unsortedPhoto.Name()
			ext := path.Ext(pathToUnsortedPhoto)
			if ext != ".jpg" && ext != ".jpeg" {
				log.Printf("non jpg/jpeg detected; skipping %s\n", pathToUnsortedPhoto)
				continue
			}
			if pathToUnsortedPhoto == pathToPhoto {
				continue
			}
			if slices.Contains(matchedPhotos, pathToUnsortedPhoto) {
				continue
			}

			photo1, err := os.ReadFile(pathToPhoto)
			if err != nil {
				log.Fatal(err)
			}

			photo2, err := os.ReadFile(pathToUnsortedPhoto)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("calling service on %s and %s\t", pathToPhoto, pathToUnsortedPhoto)
			faceMatchResult, err := faceService.CompareFaces(photo1, photo2)

			if faceMatchResult.DoFacesMatch {
				fmt.Println("matched!")

				albumName := albumNameBase + strconv.Itoa(albumNameIndex)
				albumPath := pathToPhotos + "/sorted/" + albumName

				err = os.MkdirAll(albumPath, 0777)
				if err != nil {
					log.Fatal(err)
				}

				matchedPhotos = append(matchedPhotos, pathToPhoto)
				matchedPhotos = append(matchedPhotos, pathToUnsortedPhoto)

				if !slices.Contains(albums[albumName], pathToPhoto) {
					albums[albumName] = append(albums[albumName], pathToPhoto)
					fileInfo, _ := os.Stat(pathToPhoto)
					copyFile(pathToPhoto, albumPath+string(os.PathSeparator)+fileInfo.Name())
				}

				if !slices.Contains(albums[albumName], pathToUnsortedPhoto) {
					albums[albumName] = append(albums[albumName], pathToUnsortedPhoto)
					fileInfo, _ := os.Stat(pathToUnsortedPhoto)
					copyFile(pathToUnsortedPhoto, albumPath+string(os.PathSeparator)+fileInfo.Name())
				}
			} else {
				fmt.Println("no match")
			}
		}
		albumNameIndex++
	}
}

func copyFile(sourcePath string, destinationPath string) {
	file, err := os.Create(destinationPath)
	if err != nil {
		log.Fatal(err)
	}

	source, err := os.Open(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, source)
	if err != nil {
		log.Fatal(err)
	}
}
