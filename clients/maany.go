package clients

import (
	"fmt"
	"log"
	"sahib/model"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func QueryMaany(word string) (*model.Translations, error) {
	url := fmt.Sprintf("https://www.almaany.com/en/dict/ar-en/%s/?c=Tout", word)
	results := &model.Translations{
		Link: url,
	}
	start := time.Now()
	defer func() {
		results.Elapsed = elapsed(start)
	}()

	res, err := queryURL("GET", url, nil, nil, true)
	if err != nil {
		return results, fmt.Errorf("failed to query maany: %w", err)
	}
	defer res.Body.Close()

	log.Printf("Done request")

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return results, fmt.Errorf("failed to parse html body: %w", err)
	}

	log.Printf("Done parsing")

	toTashkil := []string{}
	doc.Find(".panel-lightyellow").Find(".row").Each(func(i int, s *goquery.Selection) {
		arabic := strings.Trim(s.Find(".text-left").Text(), " ..")
		translation := strings.Trim(s.Find(".text-right").Text(), " .")

		toTashkil = append(toTashkil, strings.TrimSpace(arabic))

		results.List = append(results.List, model.Translation{Arabic: arabic, Translation: translation})
	})

	// TODO: tashkil is very slow so do it in two paths
	// tashkil, err := tashkil(toTashkil)
	// if err != nil {
	// 	log.Printf("Couldn't add tashkil to result: %+v: %s", results, err)
	// 	return results, nil
	// }

	// for i := range results.List {
	// 	results.List[i].Arabic = tashkil[i]
	// }

	return results, nil
}
