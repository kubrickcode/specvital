package entity

type Language string

const (
	LanguageAr Language = "ar" // Arabic (العربية)
	LanguageCs Language = "cs" // Czech (Čeština)
	LanguageDa Language = "da" // Danish (Dansk)
	LanguageDe Language = "de" // German (Deutsch)
	LanguageEl Language = "el" // Greek (Ελληνικά)
	LanguageEn Language = "en" // English
	LanguageEs Language = "es" // Spanish (Español)
	LanguageFi Language = "fi" // Finnish (Suomi)
	LanguageFr Language = "fr" // French (Français)
	LanguageHi Language = "hi" // Hindi (हिन्दी)
	LanguageId Language = "id" // Indonesian (Bahasa Indonesia)
	LanguageIt Language = "it" // Italian (Italiano)
	LanguageJa Language = "ja" // Japanese (日本語)
	LanguageKo Language = "ko" // Korean (한국어)
	LanguageNl Language = "nl" // Dutch (Nederlands)
	LanguagePl Language = "pl" // Polish (Polski)
	LanguagePt Language = "pt" // Portuguese (Português)
	LanguageRu Language = "ru" // Russian (Русский)
	LanguageSv Language = "sv" // Swedish (Svenska)
	LanguageTh Language = "th" // Thai (ไทย)
	LanguageTr Language = "tr" // Turkish (Türkçe)
	LanguageUk Language = "uk" // Ukrainian (Українська)
	LanguageVi Language = "vi" // Vietnamese (Tiếng Việt)
	LanguageZh Language = "zh" // Chinese (中文)
)

var validLanguages = map[Language]struct{}{
	LanguageAr: {},
	LanguageCs: {},
	LanguageDa: {},
	LanguageDe: {},
	LanguageEl: {},
	LanguageEn: {},
	LanguageEs: {},
	LanguageFi: {},
	LanguageFr: {},
	LanguageHi: {},
	LanguageId: {},
	LanguageIt: {},
	LanguageJa: {},
	LanguageKo: {},
	LanguageNl: {},
	LanguagePl: {},
	LanguagePt: {},
	LanguageRu: {},
	LanguageSv: {},
	LanguageTh: {},
	LanguageTr: {},
	LanguageUk: {},
	LanguageVi: {},
	LanguageZh: {},
}

func (l Language) IsValid() bool {
	_, ok := validLanguages[l]
	return ok
}

func (l Language) String() string {
	return string(l)
}

type ConversionResult struct {
	ConvertedName string
	IsFromCache   bool
	OriginalName  string
}

type TestInput struct {
	FilePath       string
	Framework      string
	Line           int
	Modifier       string
	OriginalName   string
	Status         string
	SuiteHierarchy string
}

type FileTestGroup struct {
	FilePath  string
	Framework string
	Suites    []SuiteTestGroup
}

type SuiteTestGroup struct {
	Hierarchy string
	Tests     []TestInput
}
