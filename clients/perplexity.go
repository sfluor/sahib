package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sahib/model"
	"strings"
	"time"
)

type PerplexityResp struct {
	Translation string                  `json:"translation"`
	Examples    []PerplexityRespExample `json:"examples"`
}

type PerplexityRespExample struct {
	Sentence    string `json:"sentence"`
	Translation string `json:"translation"`
}

type PerplexityClient struct {
	ApiKey string
}

func (c *PerplexityClient) Query(word string) (*model.Translations, error) {
	return queryPerplexity(c.ApiKey, word)
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

func queryPerplexity(token string, word string) (*model.Translations, error) {
	result := &model.Translations{}
	resp := PerplexityResp{}
	url := "https://api.perplexity.ai/chat/completions"

	start := time.Now()
	defer func() {
		result.Elapsed = elapsed(start)
	}()

	// Create the request body using map[string]interface{}
	requestBody := map[string]interface{}{
		"model": "sonar-pro",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Don't repeat yourself, be precise and concise.",
			},
			{
				"role":    "user",
				"content": prompt(word),
			},
		},
		"max_tokens":               "1000",
		"temperature":              0.2,
		"top_p":                    0.9,
		"search_domain_filter":     []string{"perplexity.ai"},
		"return_images":            false,
		"return_related_questions": false,
		"search_recency_filter":    "month",
		"top_k":                    0,
		"stream":                   false,
		"presence_penalty":         0,
		"frequency_penalty":        1,
	}

	// Serialize the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return result, fmt.Errorf("Error serializing request body: %w\n", err)
	}

	rawResp, err := queryURL("POST", url, bytes.NewBuffer(jsonBody), map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}, false)
	if err != nil {
		return result, fmt.Errorf("Error sending perplexity request: %w\n", err)
	}
	defer rawResp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return result, fmt.Errorf("Error reading body: %w\n", err)
	}

	apiResp := PerplexityAPIResp{}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return result, fmt.Errorf("Error deserializing api response: %w\n", err)
	}

	if len(apiResp.Choices) == 0 {
		return result, fmt.Errorf("Empty response received from perplexity: %s\n", body)
	}

    log.Printf("Perplexity response: %s",apiResp.Choices[0].Message.Content )
	content := extractJSON(apiResp.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		return result, fmt.Errorf("Error deserializing api response content: %w\n", err)
	}

	result.List = append(result.List, model.Translation{
		Arabic:      word,
		Translation: resp.Translation,
	})

	for _, row := range resp.Examples {
		result.List = append(result.List, model.Translation{
			Arabic:      row.Sentence,
			Translation: row.Translation,
		})
	}

	return result, nil
}

func extractJSON(input string) string {
    start := strings.Index(input, "{")
    if start == -1 {
        return ""
    }

    depth := 0
    for i := start; i < len(input); i++ {
        switch input[i] {
        case '{':
            depth++
        case '}':
            depth--
            if depth == 0 {
                return input[start : i+1]
            }
        }
    }
    return ""
}
