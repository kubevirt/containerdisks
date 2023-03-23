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
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
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
	ImageMetadata(imgRef string, insecure bool) (*ImageInfo, error)
	PushImage(ctx context.Context, img v1.Image, imgRef string) error
	CopyImage(ctx context.Context, srcRef, dstRef string, insecure bool) error
}

type RepositoryImpl struct {
}

func (r RepositoryImpl) ImageMetadata(imgRef string, insecure bool) (imageInfo *ImageInfo, retErr error) {
	sys := &types.SystemContext{
		OCIInsecureSkipTLSVerify: insecure,
	}
	if insecure {
		sys.DockerInsecureSkipTLSVerify = types.OptionalBoolTrue
	}
	ctx := context.Background()
	src, err := parseImageSource(ctx, sys, fmt.Sprintf("docker://%s", imgRef))
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
		return nil, errors.Wrapf(err, "Error inspecting image")
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

func (r RepositoryImpl) PushImage(ctx context.Context, img v1.Image, imgRef string) error {
	return crane.Push(img, imgRef, crane.WithContext(ctx))
}

func (r RepositoryImpl) CopyImage(ctx context.Context, srcRef, dstRef string, insecure bool) error {
	options := []crane.Option{
		crane.WithContext(ctx),
	}

	if insecure {
		options = append(options, crane.Insecure)
	}

	return crane.Copy(srcRef, dstRef, options...)
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
	ec := getErrorCode(err)
	if ec == nil {
		return false
	}

	if ec.ErrorCode().Error() == "unknown" {
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

	if ec, ok := err.(errcode.ErrorCoder); ok {
		return ec
	}

	return nil
}
