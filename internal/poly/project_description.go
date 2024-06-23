package poly

type ProjectDescription struct {
	DebugWorkspacePath string
	DebugMode          bool
	FullPath           string
	AppName            string
	OrganizationName   string
	PackageName        string
	Language           SupportedLanguage
}
