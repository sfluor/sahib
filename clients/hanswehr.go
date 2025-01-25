package clients

import (
	"database/sql"
	"fmt"
	"sahib/model"
	"strings"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
)

type HansWehr struct {
	db    *sql.DB
	forms map[string]verbForm
}

func NewHansWehrClient(path string) (*HansWehr, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Failed to init hans wehr client: %w", err)
	}

	forms := []verbForm{
		{
			key:      "I",
			template: "فَعَل/فَعُل/فَعِل",
			desc:     "Basic root",
			example:  "ضَرَبَ - He hit",
		},
		{
			key:      "II",
			template: "فَعّل",
			desc:     "Doing something intensively/ repeatedly, doing or causing something to someone else",
			example:  "علّم - He taught",
		},
		{
			key:      "III",
			template: "فَاعَل",
			desc:     "To try to do something, to do something with someone else",
			example:  "قاتل - He fought",
		},
		{
			key:      "IV",
			template: "أَفْعَل",
			desc:     "Transitive, immediate, doing something to other/ someone else, causing something",
			example:  "اكْرَمَ - He honored",
		},
		{
			key:      "V",
			template: "تَفَعّل",
			desc:     "Doing something intensively/ repeatedly, doing or causing something to yourself",
			example:  "تَمَتَّعَ - He enjoyed",
		},
		{
			key:      "VI",
			template: "تَفَاعَل",
			desc:     "Doing something with each other, to pretend to do something, expressing a state",
			example:  "تَبادَلَ - He exchanged",
		},
		{
			key:      "VII",
			template: "اِنْفَعَل",
			desc:     "Intransitive, Passive meaning",
			example:  "اِنكَسَرَ - He broke",
		},
		{
			key:      "VIII",
			template: "اِفْتَعَل",
			desc:     "No consistent meaning pattern, being in a state of something",
			example:  "اِجتَنَبَ - He avoided",
		},
		{
			key:      "IX",
			template: "اِفْعَل",
			desc:     "Used for colors or defects",
			example:  "اِحمرّ - He became red",
		},
		{
			key:      "X",
			template: "اِسْتَفْعَل",
			desc:     "To seek or ask something, wanting, trying",
			example:  "اِسْتَغفر - He sought forgiveness",
		},
		{
			key:      "XI",
			template: "اِفْعالَّ",
			desc:     "Like Form IX used for colors or defects but more temporary or intense",
			example:  "اِحْمارَّ - He became temporarily or extremely red",
		},
		{
			key:      "XII",
			template: "اِفْعَوْعَلَ",
			desc:     "Like Form XI tend to refer to a colour or physical quality",
			example:  "اِخْشَوْشَنَ - He became rough, coarse",
		},
	}

	formsMap := make(map[string]verbForm, len(forms))
	for _, f := range forms {
		formsMap[f.key] = f
	}

	return &HansWehr{db: db, forms: formsMap}, nil
}

func removeDiacritics(input string) string {
	result := make([]rune, 0, len(input))
	for _, r := range input {
		if !unicode.Is(unicode.Mn, r) {
			result = append(result, r)
		}
	}
	return strings.TrimSpace(string(result))
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
    // Remove the diacritics since everything is stored without diacritics in the DB
	rows, err := h.db.Query(q, removeDiacritics(word))

	if err != nil {
		return nil, fmt.Errorf("failed to query word in sqlite db: %w", err)
	}

	entries := []model.Definition{}

	defer rows.Close()
	for rows.Next() {
		e := model.Definition{}
		err = rows.Scan(&e.ID, &e.Word, &e.Definition, &e.Root, &e.RootDef, &e.QuranCount)
		if err != nil {
			return nil, fmt.Errorf("error while scanning row: %w", err)
		}

		// Make the definition easier to read.
		e.Definition = patchForms(e.Definition, h.forms)
		e.RootDef.String = patchForms(e.RootDef.String, h.forms)

        // Hide the root if it's the same as the current word.
        if e.Definition == e.RootDef.String {
            e.RootDef.Valid = false
            e.RootDef.String = ""
            e.Root.Valid = false
            e.Root.String = ""
        }

		entries = append(entries, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error at scan end: %w", err)
	}

	return &model.Definitions{Definitions: entries}, nil
}

func patchForms(s string, forms map[string]verbForm) string {
	for key, form := range forms {
		inp := fmt.Sprintf("<b>%s</b>", key)
		out := fmt.Sprintf(
			`<hr /><b data-placement="right" data-tooltip="%s (%s)">%s (%s)</b>`,
			form.desc,
			form.example,
			key,
			form.template)
		s = strings.ReplaceAll(s, inp, out)
	}

	return s
}

func (h *HansWehr) Close() error {
	return h.db.Close()
}

type verbForm struct {
	key      string
	template string
	desc     string
	example  string
}
