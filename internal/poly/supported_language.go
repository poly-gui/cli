package poly

type SupportedLanguage string

const SupportedLanguageTypeScript SupportedLanguage = "ts"

func IsLanguageSupported(lang string) bool {
	switch lang {
	case "ts":
		return true
	default:
		return false
	}
}
