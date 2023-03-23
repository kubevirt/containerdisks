package build

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

func TestStreamLayer(t *testing.T) {
	type args struct {
		imageContent string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should create the tar file in an expected format",
			args: args{imageContent: "hello"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			imageName := filepath.Join(t.TempDir(), "image")

			err := os.WriteFile(imageName, []byte(tt.args.imageContent), 0600)
			g.Expect(err).ToNot(HaveOccurred())

			imageStat, err := os.Stat(imageName)
			g.Expect(err).ToNot(HaveOccurred())

			reader, err := StreamLayerOpener(imageName)()
			g.Expect(err).ToNot(HaveOccurred())

			tarReader := tar.NewReader(reader)

			dir, err := tarReader.Next()
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(dir.Name).To(Equal("disk/"))
			g.Expect(int32(dir.Typeflag)).To(Equal(tar.TypeDir))
			g.Expect(dir.Uid).To(Equal(107))
			g.Expect(dir.Gid).To(Equal(107))

			image, err := tarReader.Next()
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(image.Name).To(Equal("disk/disk.img"))
			g.Expect(int32(image.Typeflag)).To(Equal(tar.TypeReg))
			g.Expect(image.Size).To(Equal(imageStat.Size()))
			g.Expect(image.Uid).To(Equal(107))
			g.Expect(image.Gid).To(Equal(107))
			data, err := io.ReadAll(tarReader)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(string(data)).To(Equal(tt.args.imageContent))
		})
	}
}
