package common

type Options struct {
	AllowInsecureRegistry bool
	Registry              string
	DryRun                bool
	PublishImagesOptions  PublishImageOptions
	PublishDocsOptions    PublishDocsOptions
	Focus                 string
}

type PublishImageOptions struct {
	ForceBuild bool
	Workers    int
}

type PublishDocsOptions struct {
	TokenFile string
}
