package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"sahib/components"
	"sahib/model"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const ElixirURL = "https://quest.ms.mff.cuni.cz/cgi-bin/elixir/index.fcgi?mode=home"

func main() {
	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		component := components.Index()
		component.Render(r.Context(), w)
	}

	http.HandleFunc("GET /", mainHandler)

	http.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		search := r.FormValue("search")
		log.Printf("Searching for: %s", search)
		tok := "apikey"
		res, err := queryPerplexity(tok, search)

		if err != nil {
			errStr := fmt.Sprintf("An error occurred querying perplexity: %w", err)
			log.Printf(errStr)
			http.Error(w, errStr, 500)
			return
		}

		elixir, err := queryElixir(search)
		if err != nil {
			errStr := fmt.Sprintf("An error occurred querying elixir: %w", err)
			log.Printf(errStr)
			http.Error(w, errStr, 500)
			return
		}

		component := components.Result(res, elixir, ElixirURL)
		component.Render(r.Context(), w)
	})

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func queryElixir(word string) ([]model.ElixirResp, error) {

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
			return nil, fmt.Errorf("Couldn't write form field %s: %v\n", f.key, err)
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("Couldn't write form multipart %v\n", err)
	}

	req, err := http.NewRequest("POST", ElixirURL, body)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	result := []model.ElixirResp{}
	// Find the words
	doc.Find(".lexeme").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		tag := s.Find(".xtag").Text()
		orth := s.Find(".orth").Text()
		reflex := s.Find(".reflex").Text()

		result = append(result, model.ElixirResp{
			Tag:         tag,
			Arabic:      orth,
			Translation: reflex,
		})
	})

	return result, nil
}

func prompt(word string) string {
	count := 5
	target := "french"
	return fmt.Sprintf(`Give me %d examples of useful and relevant sentences from medias or stories with the proper arabic diacritics (harakats) on all words.

The word is: %s

The translation language should be %s

The output should be in JSON format like so:

{
    "translation": "The translation of the word",
    "examples": [
         {"sentence": "An example of sentence", "translation": "The translation of the setence in the target language"},
     ]
}
`, count, word, target)
}

type PerplexityAPIResp struct {
	Choices []PerplexityAPIChoice `json:"choices"`
}

type PerplexityAPIChoice struct {
	Message PerplexityAPIMessage `json:"message"`
}

type PerplexityAPIMessage struct {
	Content string `json:"content"`
}

func queryPerplexity(token string, word string) (model.PerplexityResp, error) {
	resp := model.PerplexityResp{}
	//    url := "https://api.perplexity.ai/chat/completions"
	//
	// // Create the request body using map[string]interface{}
	// requestBody := map[string]interface{}{
	// 	"model": "llama-3.1-sonar-small-128k-online",
	// 	"messages": []map[string]string{
	// 		{
	// 			"role":    "system",
	// 			"content": "Don't repeat yourself, be precise and concise.",
	// 		},
	// 		{
	// 			"role":    "user",
	// 			"content": prompt(word),
	// 		},
	// 	},
	// 	"max_tokens":               "1000",
	// 	"temperature":              0.2,
	// 	"top_p":                    0.9,
	// 	"search_domain_filter":     []string{"perplexity.ai"},
	// 	"return_images":            false,
	// 	"return_related_questions": false,
	// 	"search_recency_filter":    "month",
	// 	"top_k":                    0,
	// 	"stream":                   false,
	// 	"presence_penalty":         0,
	// 	"frequency_penalty":        1,
	// }
	//
	//
	// // Serialize the request body to JSON
	// jsonBody, err := json.Marshal(requestBody)
	// if err != nil {
	//        return resp, fmt.Errorf("Error serializing request body: %w\n", err)
	// }
	//
	// // Create a new HTTP request
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	// if err != nil {
	//        return resp, fmt.Errorf("Error creating request: %w\n", err)
	// }
	//
	// // Set headers
	// req.Header.Set("Authorization", "Bearer "+token)
	// req.Header.Set("Content-Type", "application/json")
	//
	// // Send the request
	// client := &http.Client{}
	// rawResp, err := client.Do(req)
	// if err != nil {
	//        return resp, fmt.Errorf("Error sending request: %w\n", err)
	// }
	// defer rawResp.Body.Close()
	//
	// // Read the response body
	// body, err := io.ReadAll(rawResp.Body)
	// if err != nil {
	//        return resp, fmt.Errorf("Error reading body: %w\n", err)
	// }
	//
	// // Check the response status code
	// if rawResp.StatusCode != http.StatusOK {
	//        return resp, fmt.Errorf("Unexpected status code: %d: %s", rawResp.StatusCode, string(body))
	// }

	body := []byte(strings.Replace(`
{
  "id": "c871af1e-7f86-46f4-915d-5ee54f1b19af",
  "model": "llama-3.1-sonar-huge-128k-online",
  "created": 1736598912,
  "usage": {
    "prompt_tokens": 90,
    "completion_tokens": 338,
    "total_tokens": 428
  },
  "citations": [
    "https://www.tiktok.com/discover/%D9%81%D8%B1%D8%B6-%D8%AA%D8%A3%D9%84%D9%8A%D9%81%D9%8A-%D8%B9%D8%AF%D8%AF3-%D9%81%D8%B1%D9%86%D8%B3%D9%8A%D8%A9-%D8%B3%D9%86%D8%A9-%D8%B3%D8%A7%D8%A8%D8%B9%D8%A9-%D8%A3%D8%B3%D8%A7%D8%B3%D9%8A"
  ],
  "object": "chat.completion",
  "choices": [
    {
      "index": 0,
      "finish_reason": "stop",
      "message": {
        "role": "assistant",
        "content": "Here are 5 examples of useful and relevant sentences with the proper harakats for the word \"تفاوض\" along with their translations in French:\n\n$$$json\n{\n  \"translation\": \"négocier\",\n  \"examples\": [\n    {\n      \"sentence\": \"كان علي أن يتفاوض معهم لتحقيق اتفاق.\",\n      \"translation\": \"Il devait négocier avec eux pour parvenir à un accord.\"\n    },\n    {\n      \"sentence\": \"لا بد من التفاوض مع الأطراف المعنية.\",\n      \"translation\": \"Il faut négocier avec les parties concernées.\"\n    },\n    {\n      \"sentence\": \"كانت الخطة أن يتفاوض معهم لتحقيق حل.\",\n      \"translation\": \"Le plan était de négocier avec eux pour trouver une solution.\"\n    },\n    {\n      \"sentence\": \"يجب أن نتفاوض معهم لتحقيق مصلحتنا.\",\n      \"translation\": \"Nous devons négocier avec eux pour servir nos intérêts.\"\n    },\n    {\n      \"sentence\": \"كان من الضروري التفاوض معهم لتحقيق النجاح.\",\n      \"translation\": \"Il était essentiel de négocier avec eux pour réussir.\"\n    }\n  ]\n}\n$$$\n\nThese examples are crafted to provide a variety of contexts in which the word \"تفاوض\" (to negotiate) is used, along with their translations in French. Note that the source provided in the search results does not directly offer these examples, so they have been created to meet the request."
      },
      "delta": {
        "role": "assistant",
        "content": ""
      }
    }
  ]
}
    `, "$", "`", -1))

	apiResp := PerplexityAPIResp{}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return resp, fmt.Errorf("Error deserializing api response: %w\n", err)
	}

	content := apiResp.Choices[0].Message.Content
	content = strings.Split(strings.Split(content, "```json")[1], "``")[0]

	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		return resp, fmt.Errorf("Error deserializing api response content: %w\n", err)
	}

	return resp, nil
}
