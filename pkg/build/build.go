package build

import (
	"fmt"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

const (
	LabelShaSum       = "shasum"
	ImageArchitecture = "amd64"
)

func ContainerDiskConfig(checksum string, additionalLabels map[string]string) v1.Config {
	labels := map[string]string{
		LabelShaSum: checksum,
	}
	for k, v := range additionalLabels {
		labels[k] = v
	}

	// Add all labels also as ENV variable for compatibility with crun-vm
	// Replace illegal characters with underscore
	var env []string
	r := strings.NewReplacer(".", "_", "/", "_", "-", "_")
	for k, v := range labels {
		env = append(env, fmt.Sprintf("%s=%s", strings.ToUpper(r.Replace(k)), v))
	}

	return v1.Config{Labels: labels, Env: env}
}

func ContainerDisk(imgPath string, config v1.Config) (v1.Image, error) {
	img := empty.Image
	layer, err := tarball.LayerFromOpener(StreamLayerOpener(imgPath))
	if err != nil {
		return nil, fmt.Errorf("error creating an image layer from disk: %v", err)
	}

	img, err = mutate.AppendLayers(img, layer)
	if err != nil {
		return nil, fmt.Errorf("error appending the image layer: %v", err)
	}

	cf, err := img.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error getting the image config file: %v", err)
	}

	// Modify the config file
	cf.Architecture = ImageArchitecture
	cf.Config = config

	img, err = mutate.ConfigFile(img, cf)
	if err != nil {
		return nil, fmt.Errorf("error setting the image config file: %v", err)
	}

	return img, nil
}
