package common

type Options struct {
	AllowInsecureRegistry bool
	Registry              string
	DryRun                bool
	PublishImagesOptions  PublishImageOptions
	PublishDocsOptions    PublishDocsOptions
	Focus                 string
	Workers               int
}

type PublishImageOptions struct {
	ForceBuild bool
}

type PublishDocsOptions struct {
	TokenFile string
}
