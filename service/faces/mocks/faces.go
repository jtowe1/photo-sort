package mocks

import (
	"bytes"
	"github.com/jtowe1/photo-sort/service/faces"
)

type FaceService struct {
}

func (f *FaceService) CompareFaces(photo1 []byte, photo2 []byte) (*faces.FaceCompareResult, error) {
	if bytes.Compare(photo1, photo2) == 0 {
		return &faces.FaceCompareResult{
				DoFacesMatch: true,
				PercentMatch: 99,
			},
			nil
	} else {
		return &faces.FaceCompareResult{
				DoFacesMatch: false,
				PercentMatch: 0,
			},
			nil
	}
}
