package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/docker/distribution/registry/api/errcode"
	v2 "github.com/docker/distribution/registry/api/v2"
	"github.com/pkg/errors"
)

type ImageInfo struct {
	Tag           string `json:",omitempty"`
	Created       *time.Time
	DockerVersion string
	Labels        map[string]string
	Architecture  string
	Os            string
	Layers        []string
	Env           []string
}

type Repository interface {
	ImageMetadata(imageRef string) (*ImageInfo, error)
}

type RepositoryImpl struct {
}

func (r RepositoryImpl) ImageMetadata(imageRef string, insecure bool) (imageInfo *ImageInfo, retErr error) {
	sys := &types.SystemContext{
		OCIInsecureSkipTLSVerify: insecure,
	}
	if insecure {
		sys.DockerInsecureSkipTLSVerify = types.OptionalBoolTrue
	}
	ctx := context.Background()
	src, err := parseImageSource(ctx, sys, fmt.Sprintf("docker://%s", imageRef))
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing image")
	}

	defer func() {
		if err := src.Close(); err != nil {
			retErr = errors.Wrapf(retErr, fmt.Sprintf("(could not close image: %v) ", err))
		}
	}()

	img, err := image.FromUnparsedImage(ctx, sys, image.UnparsedInstance(src, nil))
	if err != nil {
		return nil, errors.Wrapf(err, "Error parsing manifest for image")
	}
	imgInspect, err := img.Inspect(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Error inpspecting image")
	}
	imageInfo = &ImageInfo{
		Tag: imgInspect.Tag,
		// Digest is set below.
		Created:       imgInspect.Created,
		DockerVersion: imgInspect.DockerVersion,
		Labels:        imgInspect.Labels,
		Architecture:  imgInspect.Architecture,
		Os:            imgInspect.Os,
		Layers:        imgInspect.Layers,
		Env:           imgInspect.Env,
	}

	return imageInfo, nil
}

func parseImageSource(ctx context.Context, sys *types.SystemContext, name string) (types.ImageSource, error) {
	ref, err := alltransports.ParseImageName(name)
	if err != nil {
		return nil, err
	}

	return ref.NewImageSource(ctx, sys)
}

func IsManifestUnknownError(err error) bool {
	ec := getErrorCode(err)
	if ec == nil {
		return false
	}
	switch ec.ErrorCode() {
	case v2.ErrorCodeManifestUnknown:
		return true
	default:
		return false
	}
}

func IsRepositoryUnknownError(err error) bool {
	ec := getErrorCode(err)
	if ec == nil {
		return false
	}

	switch ec.ErrorCode() {
	case v2.ErrorCodeNameUnknown:
		return true
	default:
		return false
	}
}

func IsTagUnknownError(err error) bool {
	e := getErrorCode(err)
	if e.ErrorCode().Error() == "unknown" {
		// errors like this have no explicit error handling: "unknown: Tag 5.2 was deleted or has expired. To pull, revive via time machine"
		if strings.Contains(err.Error(), "was deleted or has expired. To pull, revive via time machine") {
			return true
		}
	}
	return false
}

func getErrorCode(err error) errcode.ErrorCoder {
	for {
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			err = unwrapped
		} else {
			break
		}
	}

	errors, ok := err.(errcode.Errors)
	if !ok || len(errors) == 0 {
		return nil
	}
	err = errors[0]
	ec, ok := err.(errcode.ErrorCoder)
	if !ok {
		return nil
	}
	return ec
}
