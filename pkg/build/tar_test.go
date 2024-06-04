package build

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tar", func() {
	It("StreamLayer should create the tar file in an expected format", func() {
		const imageContent = "hello"

		imageName := filepath.Join(GinkgoT().TempDir(), "image")
		err := os.WriteFile(imageName, []byte(imageContent), 0o600)
		Expect(err).ToNot(HaveOccurred())

		imageStat, err := os.Stat(imageName)
		Expect(err).ToNot(HaveOccurred())

		reader, err := StreamLayerOpener(imageName)()
		Expect(err).ToNot(HaveOccurred())

		tarReader := tar.NewReader(reader)

		dir, err := tarReader.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(dir.Name).To(Equal("disk/"))
		Expect(int32(dir.Typeflag)).To(Equal(tar.TypeDir))
		Expect(dir.Uid).To(Equal(107))
		Expect(dir.Gid).To(Equal(107))

		image, err := tarReader.Next()
		Expect(err).ToNot(HaveOccurred())
		Expect(image.Name).To(Equal("disk/disk.img"))
		Expect(int32(image.Typeflag)).To(Equal(tar.TypeReg))
		Expect(image.Size).To(Equal(imageStat.Size()))
		Expect(image.Uid).To(Equal(107))
		Expect(image.Gid).To(Equal(107))
		data, err := io.ReadAll(tarReader)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(data)).To(Equal(imageContent))
	})
})

func TestTar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tar Suite")
}
