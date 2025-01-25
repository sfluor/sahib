package model

import (
	"database/sql"
)

const (
	SourceWehr     = "HansWehr"
	SourceElixir     = "Elixir"
	SourceMaany      = "Maany"
	SourcePerplexity = "Perplexity"

	ApiKey = "apiKey"
	Search = "search"
    Lang= "lang"
)

var AllSources = []string{
    SourceWehr,
    SourceElixir,
    SourceMaany,
    SourcePerplexity,
}

func SourceAndLangIds() []string {
    languages := Languages()
    r := make([]string, 0, len(AllSources) + len(languages))
    for _, src := range AllSources {
        r = append(r, "#" + src)
    }
    for _, lang := range languages {
        r = append(r, "#" + lang.Short)
    }

    return r
}

func Languages() []Language {
    return []Language{
        {
            Short: "lang_fr",
            Code: "fr",
            Name: "French",
            Logo:  "🇫🇷",
        },
        {
            Short: "lang_en",
            Code: "en",
            Name: "English",
            Logo:  "🇬🇧",
        },
    }
}

type Language struct {
    Short string
    Code string
    Name string
    Logo string
}

type TranslationsAndSource struct {
	Translations *Translations
	Source       string
}

type Translations struct {
	Elapsed string
	Link    string
	List    []Translation
	Error   string
}

type Translation struct {
	Arabic      string
	Translation string
	Meta        string
}

type Definitions struct {
	Definitions []Definition
}

type Definition struct {
	ID         int
	Word       string
	Definition string
	Root       sql.NullString
	RootDef    sql.NullString
	QuranCount sql.NullInt64
}
