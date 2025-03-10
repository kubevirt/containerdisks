package common

type Options struct {
	AllowInsecureRegistry bool
	DryRun                bool
	Focus                 string
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
	Registry  string
	TokenFile string
}

type PublishImageOptions struct {
	ForceBuild     bool
	NoFail         bool
	SourceRegistry string
	TargetRegistry string
}

type VerifyImageOptions struct {
	Registry           string
	Namespace          string
	NoFail             bool
	Timeout            int
	TargetArchitecture string
}
