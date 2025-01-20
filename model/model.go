package model

import (
	"database/sql"
)

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
    ID int
    Word string
    Definition string
    Root sql.NullString
    RootID int
    RootDef sql.NullString
    QuranCount sql.NullInt64
}
