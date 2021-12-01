package common

type Options struct {
	AllowInsecureRegistry bool
	Registry              string
	DryRun                bool
	PublishImagesOptions  PublishImageOptions
	PublishDocsOptions    PublishDocsOptions
}

type PublishImageOptions struct {
	ForceBuild bool
	Focus      string
	Workers    int
}

type PublishDocsOptions struct {
	TokenFile string
}
