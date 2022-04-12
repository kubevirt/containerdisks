package common

type Options struct {
	AllowInsecureRegistry bool
	DryRun                bool
	Focus                 string
	Registry              string
	ImagesOptions         ImagesOptions
	PublishDocsOptions    PublishDocsOptions
	PublishImagesOptions  PublishImageOptions
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
