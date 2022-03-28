package build

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

const (
	LabelShaSum = "shasum"
)

func BuildContainerDisk(imgPath string, checksum string) (v1.Image, error) {
	img := empty.Image
	layerStream, errChan := StreamLayer(imgPath)
	layer, err := tarball.LayerFromReader(layerStream)
	if err != nil {
		return nil, fmt.Errorf("error creating an image layer from disk: %v", err)
	}

	img, err = mutate.AppendLayers(img, layer)
	if err != nil {
		return nil, fmt.Errorf("error appending the image layer: %v", err)
	}

	img, err = mutate.Config(img, v1.Config{Labels: map[string]string{LabelShaSum: checksum}})
	if err != nil {
		return nil, fmt.Errorf("error appending labels to the image: %v", err)
	}

	if err := <-errChan; err != nil {
		return nil, fmt.Errorf("error creating the tar file with the disk: %v", err)
	}
	return img, nil
}

func PushImage(ctx context.Context, img v1.Image, name string) error {
	if err := crane.Push(img, name, crane.WithContext(ctx)); err != nil {
		return fmt.Errorf("error pushing image %q: %v", img, err)
	}
	return nil
}
