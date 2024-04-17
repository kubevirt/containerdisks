package build

import (
	"fmt"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/partial"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

const (
	LabelShaSum = "shasum"
	ImageOS     = "linux"
)

func ContainerDiskConfig(checksum string, envVariables map[string]string) v1.Config {
	labels := map[string]string{
		LabelShaSum: checksum,
	}

	var env []string
	for k, v := range envVariables {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// OCI runtimes like crun-vm [1] may be able to run a containerdisk container image directly, in which case it
	// would be unfortunate to require the user to pass in a placeholder entrypoint to docker-run and equivalents,
	// since engines like Docker require an entrypoint to be set, e.g.:
	//
	//   $ docker run --runtime crun-vm quay.io/containerdisks/fedora:39 unused-entrypoint
	//
	// Instead we set the entrypoint to a predetermined "no-entrypoint" value that runtimes like crun-vm can
	// interpret as there not being an entrypoint.
	//
	// [1] https://github.com/containers/crun-vm
	entrypoint := []string{"no-entrypoint"}

	return v1.Config{Labels: labels, Env: env, Entrypoint: entrypoint}
}

func ContainerDisk(imgPath, imgArch string, config v1.Config) (v1.Image, error) {
	layer, err := tarball.LayerFromOpener(StreamLayerOpener(imgPath))
	if err != nil {
		return nil, fmt.Errorf("error creating an image layer from disk: %v", err)
	}

	img := mutate.MediaType(empty.Image, types.DockerManifestSchema2)
	img, err = mutate.AppendLayers(img, layer)
	if err != nil {
		return nil, fmt.Errorf("error appending the image layer: %v", err)
	}

	cf, err := img.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error getting the image config file: %v", err)
	}

	// Modify the config file
	cf.Architecture = imgArch
	cf.OS = ImageOS
	cf.Config = config

	img, err = mutate.ConfigFile(img, cf)
	if err != nil {
		return nil, fmt.Errorf("error setting the image config file: %v", err)
	}

	return img, nil
}

func ContainerDiskIndex(images []v1.Image) (v1.ImageIndex, error) {
	var indexAddendum []mutate.IndexAddendum

	for _, image := range images {
		configFile, err := image.ConfigFile()
		if err != nil {
			return nil, err
		}

		descriptor, err := partial.Descriptor(image)
		if err != nil {
			return nil, err
		}
		descriptor.Platform = configFile.Platform()

		indexAddendum = append(indexAddendum, mutate.IndexAddendum{
			Add:        image,
			Descriptor: *descriptor,
		})
	}

	idx := mutate.IndexMediaType(empty.Index, types.DockerManifestList)
	return mutate.AppendManifests(idx, indexAddendum...), nil
}
