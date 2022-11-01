package sort

import (
	"github.com/jtowe1/photo-sort/service/faces/mocks"
	"os"
	"strings"
	"testing"
)

func TestPhotosByLocalPathWithMatches(t *testing.T) {
	sort := ServiceSort{
		FaceService: &mocks.FaceService{},
	}

	tempDir := t.TempDir()
	err := os.WriteFile(tempDir+"/person1photo1.jpg", []byte("person1"), 777)
	err = os.WriteFile(tempDir+"/person1photo2.jpg", []byte("person1"), 777)
	err = os.WriteFile(tempDir+"/person2photo1.jpg", []byte("person2"), 777)
	err = os.WriteFile(tempDir+"/person2photo2.jpg", []byte("person2"), 777)
	if err != nil {
		t.Fatal(err)
	}

	err = sort.PhotosByLocalPath(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	dir, err := os.ReadDir(tempDir)

	hasSortedFolder := false
	for _, entry := range dir {
		if strings.Contains(entry.Name(), "sorted") {
			hasSortedFolder = true
		}
	}

	if hasSortedFolder {
		t.Log("found sorted folder")
	} else {
		t.Fatal("did not find sorted folder")
	}

	album0, err := os.ReadDir(tempDir + "/sorted/album0")
	album1, err := os.ReadDir(tempDir + "/sorted/album1")
	person1Sorted := false
	person2Sorted := false
	if album0[0].Name() == "person1photo1.jpg" && album0[1].Name() == "person1photo2.jpg" {
		person1Sorted = true
	}

	if album1[0].Name() == "person2photo1.jpg" && album1[1].Name() == "person2photo2.jpg" {
		person2Sorted = true
	}

	if person1Sorted && person2Sorted {
		t.Log("person 1 and 2 sorted")
	} else {
		t.Fatal("person 1 and 2 NOT sorted")
	}
}
