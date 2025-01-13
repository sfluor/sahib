package clients

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"sahib/model"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const ElixirURL = "https://quest.ms.mff.cuni.cz/cgi-bin/elixir/index.fcgi?mode=home"

func QueryElixir(word string) (*model.Translations, error) {
	result := &model.Translations{
		Link: ElixirURL,
	}

	start := time.Now()
	defer func() {
		result.Elapsed = elapsed(start)
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	formFields := []struct {
		key string
		val string
	}{
		{"text", word},
		{"code", "Unicode"},
		{"submit", "Resolve"},
		{"mode", "resolve"},
		{".cgifields", "code"},
		{".cgifields", "fuzzy"},
		{".cgifields", "quick"},
	}

	for _, f := range formFields {
		err := writer.WriteField(f.key, f.val)
		if err != nil {
			return result, fmt.Errorf("Couldn't write form field %s: %v\n", f.key, err)
		}
	}

	err := writer.Close()
	if err != nil {
		return result, fmt.Errorf("Couldn't write form multipart %v\n", err)
	}

	res, err := queryURL("POST", ElixirURL, body, map[string]string{ContentType: writer.FormDataContentType()}, false)
	if err != nil {
		return result, fmt.Errorf("Failed to query elixir: %w", err)
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the words
	doc.Find(".lexeme").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		tag := s.Find(".xtag").Text()
		orth := s.Find(".orth").Text()
		reflex := s.Find(".reflex").Text()

		result.List = append(result.List, model.Translation{
			Meta:        tag,
			Arabic:      orth,
			Translation: strings.ReplaceAll(reflex, "\"", ""),
		})
	})

	return result, nil
}
