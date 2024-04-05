package architecture

import "fmt"

func GetImageArchitecture(arch string) string {
	switch arch {
	case "x86_64":
		return "amd64"
	case "aarch64":
		return "arm64"
	default:
		panic(fmt.Sprintf("can't map unknown architecture %s to image architecture", arch))
	}
}
