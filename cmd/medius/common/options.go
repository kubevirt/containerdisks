package common

type Options struct {
	AllowInsecureRegistry bool
	DryRun                bool
	Focus                 string
	Registry              string
	ImagesOptions         ImagesOptions
	PublishDocsOptions    PublishDocsOptions
	PublishImagesOptions  PublishImageOptions
	VerifyImagesOptions   VerifyImageOptions
}

type ImagesOptions struct {
	ResultsFile string
	Workers     int
}

type PublishDocsOptions struct {
	TokenFile string
}

type PublishImageOptions struct {
	ForceBuild bool
}

type VerifyImageOptions struct {
	Namespace string
	Timeout   int
}
