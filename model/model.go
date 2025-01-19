package model

import (
	"database/sql"
	"fmt"
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
    RootDef sql.NullString
    QuranCount sql.NullInt64
}

func (d *Definition) TruncatedRootDef() string {
    if !d.Root.Valid {
        return ""
    }

    def := fmt.Sprintf("%s\n---\n%s", d.Root.String, d.RootDef.String)

    maxC := 128
    if len(def) > maxC {
        def = def[:maxC]
    }

    fmt.Printf("def: %s\n", def)

    return def
}

