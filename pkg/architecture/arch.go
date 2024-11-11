package architecture

import "fmt"

func GetImageArchitecture(arch string) string {
	const archS390x = "s390x"

	switch arch {
	case "x86_64":
		return "amd64"
	case "aarch64":
		return "arm64"
	case archS390x:
		return archS390x
	default:
		panic(fmt.Sprintf("can't map unknown architecture %s to image architecture", arch))
	}
}
