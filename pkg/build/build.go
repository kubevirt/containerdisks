package build

import (
	"fmt"

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
	layer, err := tarball.LayerFromOpener(StreamLayerOpener(imgPath))
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

	return img, nil
}
