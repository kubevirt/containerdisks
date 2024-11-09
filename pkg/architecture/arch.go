package architecture

import "fmt"

const ArchS390x = "s390x"

func GetImageArchitecture(arch string) string {
	switch arch {
	case "x86_64":
		return "amd64"
	case "aarch64":
		return "arm64"
	case ArchS390x:
		return ArchS390x
	default:
		panic(fmt.Sprintf("can't map unknown architecture %s to image architecture", arch))
	}
}
