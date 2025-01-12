package clients

import (
	"fmt"
	"sahib/model"

	"github.com/PuerkitoBio/goquery"
)

func QueryMaany(word string) ([]model.MaanyResp, error) {
	url := fmt.Sprintf("https://www.almaany.com/fr/dict/ar-fr/%s/?c=Tout", word)

	res, err := queryURL("GET", url, nil, nil, true)
	if err != nil {
		return nil, fmt.Errorf("failed to query maany: %w", err)
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html body: %w", err)
	}

	results := []model.MaanyResp{}

	doc.Find(".panel-lightyellow").Find(".row").Each(func(i int, s *goquery.Selection) {
		translation := s.Find(".text-left").Text()
		arabic := s.Find(".text-right").Text()

		results = append(results, model.MaanyResp{Arabic: arabic, Translation: translation})
	})

	return results, nil
}
