package common

type Options struct {
	AllowInsecureRegistry bool
	DryRun                bool
	Focus                 string
	Registry              string
	ImagesOptions         ImagesOptions
	PublishDocsOptions    PublishDocsOptions
	PublishImagesOptions  PublishImageOptions
	PromoteImageOptions   PromoteImageOptions
	VerifyImagesOptions   VerifyImageOptions
}

type ImagesOptions struct {
	ResultsFile string
	Workers     int
}

type PromoteImageOptions struct {
	SourceRegistry string
	TargetRegistry string
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
