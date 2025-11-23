package gcp

import (
	"context"

	vision "cloud.google.com/go/vision/v2/apiv1"
)

func NewVisionClient(ctx context.Context) (*vision.ImageAnnotatorClient, error) {
	return vision.NewImageAnnotatorClient(ctx)
}
