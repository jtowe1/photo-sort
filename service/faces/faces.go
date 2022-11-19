package faces

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

type FaceServiceInterface interface {
	CompareFaces(photo1 []byte, photo2 []byte) (*FaceCompareResult, error)
}

type FaceService struct {
	awsConfig         *aws.Config
	rekognitionClient *rekognition.Client
}

type FaceCompareResult struct {
	DoFacesMatch bool
	PercentMatch float32
}

func NewFaceService() (*FaceService, error) {
	//logger := logging.NewStandardLogger(os.Stdout)
	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
		//config.WithLogger(logger),
		//config.WithClientLogMode(aws.LogRequest|aws.LogResponseWithBody),
		config.WithSharedConfigProfile("photo-sort"),
	)
	if err != nil {
		return nil, err
	}
	client := rekognition.NewFromConfig(awsConfig)

	faceService := FaceService{
		awsConfig:         &awsConfig,
		rekognitionClient: client,
	}

	return &faceService, nil
}

func (f *FaceService) CompareFaces(photo1 []byte, photo2 []byte) (*FaceCompareResult, error) {
	res, err := f.rekognitionClient.CompareFaces(context.TODO(), &rekognition.CompareFacesInput{
		SourceImage: &types.Image{
			Bytes:    photo1,
			S3Object: nil,
		},
		TargetImage: &types.Image{
			Bytes:    photo2,
			S3Object: nil,
		},
		QualityFilter:       "",
		SimilarityThreshold: nil,
	})

	if err != nil {
		return nil, err
	}

	faceCompareResult := &FaceCompareResult{
		DoFacesMatch: doFacesMatch(res),
		PercentMatch: percentMatch(res),
	}

	return faceCompareResult, nil
}

func doFacesMatch(output *rekognition.CompareFacesOutput) bool {
	if len(output.FaceMatches) > 0 {
		return true
	} else {
		return false
	}
}

func percentMatch(output *rekognition.CompareFacesOutput) float32 {
	var maxMatch float32

	for _, element := range output.FaceMatches {
		if *element.Similarity > maxMatch {
			maxMatch = *element.Similarity
		}
	}

	return maxMatch
}
