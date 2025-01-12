package model

type TranslationsAndSource struct {
	Translations Translations
	Source       string
}

type Translations struct {
	List  []Translation
	Error string
}

type Translation struct {
	Arabic      string
	Translation string
	Meta        string
}
