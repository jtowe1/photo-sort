package sort

import (
	"github.com/jtowe1/photo-sort/service/faces"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
)

const AlbumBaseName = "album"
const SortedFolderName = "sorted"

type ServiceSort struct {
	FaceService faces.FaceServiceInterface
	DebugLog    bool
}

func (ss *ServiceSort) PhotosByLocalPath(pathToPhotos string) error {
	var wg sync.WaitGroup
	photos, err := validatePath(pathToPhotos)
	if err != nil {
		return err
	}

	unsortedPhotos := make([]os.DirEntry, len(*photos))
	copy(unsortedPhotos, *photos)

	albumNameIndex := 0
	albums := make(map[string][]string)
	var matchedPhotos []string

	for _, photo := range *photos {
		pathToPhoto := pathToPhotos + string(os.PathSeparator) + photo.Name()
		ext := path.Ext(pathToPhoto)
		if ext != ".jpg" && ext != ".jpeg" {
			ss.debugLogF("non jpg/jpeg detected; skipping %s\n", pathToPhoto)
			continue
		}
		if slices.Contains(matchedPhotos, pathToPhoto) {
			continue
		}
		for _, unsortedPhoto := range unsortedPhotos {
			pathToUnsortedPhoto := pathToPhotos + string(os.PathSeparator) + unsortedPhoto.Name()
			ext := path.Ext(pathToUnsortedPhoto)
			if ext != ".jpg" && ext != ".jpeg" {
				ss.debugLogF("non jpg/jpeg detected; skipping %s\n", pathToUnsortedPhoto)
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
				return err
			}

			photo2, err := os.ReadFile(pathToUnsortedPhoto)
			if err != nil {
				return err
			}

			ss.debugLogF("calling service on %s and %s\t", pathToPhoto, pathToUnsortedPhoto)
			faceMatchResult, err := ss.FaceService.CompareFaces(photo1, photo2)
			if err != nil {
				return err
			}

			if faceMatchResult.DoFacesMatch {
				ss.debugLogF("matched!\n")

				albumName := buildAlbumName(albumNameIndex)
				albumPath := buildAlbumPath(pathToPhotos, albumName)

				err = os.MkdirAll(albumPath, 0777)
				if err != nil {
					return err
				}

				matchedPhotos = append(matchedPhotos, pathToPhoto)
				matchedPhotos = append(matchedPhotos, pathToUnsortedPhoto)

				if !slices.Contains(albums[albumName], pathToPhoto) {
					albums[albumName] = append(albums[albumName], pathToPhoto)
					fileInfo, _ := os.Stat(pathToPhoto)
					wg.Add(1)
					go func() {
						defer wg.Done()
						copyFile(pathToPhoto, albumPath+string(os.PathSeparator)+fileInfo.Name())
					}()

				}

				if !slices.Contains(albums[albumName], pathToUnsortedPhoto) {
					albums[albumName] = append(albums[albumName], pathToUnsortedPhoto)
					fileInfo, _ := os.Stat(pathToUnsortedPhoto)
					wg.Add(1)
					go func() {
						defer wg.Done()
						copyFile(pathToUnsortedPhoto, albumPath+string(os.PathSeparator)+fileInfo.Name())
					}()

				}
			} else {
				ss.debugLogF("no match\n")
			}
		}
		albumNameIndex++
	}

	wg.Wait()
	return nil
}

func (ss *ServiceSort) debugLogF(format string, v ...any) {
	if ss.DebugLog {
		log.Printf(format, v...)
	}
}

func buildAlbumName(albumNameIndex int) string {
	return AlbumBaseName + strconv.Itoa(albumNameIndex)
}

func buildAlbumPath(pathToPhotos string, albumName string) string {
	return pathToPhotos + string(os.PathSeparator) + SortedFolderName + string(os.PathSeparator) + albumName
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

func validatePath(path string) (*[]os.DirEntry, error) {
	photos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return &photos, nil
}
