package clients

import (
	"database/sql"
	"fmt"
	"sahib/model"

	_ "github.com/mattn/go-sqlite3"
)


type HansWehr struct {
    db *sql.DB
}

func NewHansWehrClient(path string) (*HansWehr, error) {
	db, err := sql.Open("sqlite3", path)
    if err != nil  {
        return nil, fmt.Errorf("Failed to init hans wehr client: %w", err)
    }

    return &HansWehr{db : db}, nil
}

func (h *HansWehr) Query(word string) (*model.Definitions, error) {

    q := `
SELECT
    d1.id, d1.word, d1.definition, d2.word, d2.definition, d1.quran_occurrence
FROM
    DICTIONARY d1
    INNER JOIN
    DICTIONARY d2
    ON d2.id = d1.parent_id
    WHERE d1.word = ?
    LIMIT 10
    `
	rows, err := h.db.Query(q, word)

    if err != nil {
        return nil, fmt.Errorf("failed to query word in sqlite db: %w", err)
    }

    entries := []model.Definition{}

	defer rows.Close()
	for rows.Next() {
        e:= model.Definition{}
		err = rows.Scan(&e.ID, &e.Word, &e.Definition, &e.Root, &e.RootDef, &e.QuranCount)
        if err != nil {
            return nil, fmt.Errorf("error while scanning row: %w", err)
        }

        entries = append(entries, e)
	}
	err = rows.Err()
    if err != nil {
        return nil, fmt.Errorf("error at scan end: %w", err)
    }

    return &model.Definitions{Definitions: entries}, nil
}

func (h *HansWehr) Close() error {
    return h.db.Close()
}

