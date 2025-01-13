package model

type TranslationsAndSource struct {
	Translations *Translations
	Source       string
}

type Translations struct {
	Elapsed string
	Link   string
	List   []Translation
	Error  string
}

type Translation struct {
	Arabic      string
	Translation string
	Meta        string
}
