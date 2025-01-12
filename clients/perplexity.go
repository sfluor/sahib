package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sahib/model"
	"strings"
)

type PerplexityClient struct {
	apiKey string
}

func (c *PerplexityClient) Query(word string) (model.PerplexityResp, error) {
	return queryPerplexity(c.apiKey, word)
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
	url := "https://api.perplexity.ai/chat/completions"

	// Create the request body using map[string]interface{}
	requestBody := map[string]interface{}{
		"model": "llama-3.1-sonar-small-128k-online",
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
		return resp, fmt.Errorf("Error serializing request body: %w\n", err)
	}

	rawResp, err := queryURL("POST", url, bytes.NewBuffer(jsonBody), map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}, false)
	if err != nil {
		return resp, fmt.Errorf("Error sending perplexity request: %w\n", err)
	}
	defer rawResp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return resp, fmt.Errorf("Error reading body: %w\n", err)
	}

	apiResp := PerplexityAPIResp{}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return resp, fmt.Errorf("Error deserializing api response: %w\n", err)
	}

	content := apiResp.Choices[0].Message.Content
	if !strings.HasPrefix(content, "{\n") {
		content = strings.Split(strings.Split(content, "```json")[1], "``")[0]
	}

	if err := json.Unmarshal([]byte(content), &resp); err != nil {
		return resp, fmt.Errorf("Error deserializing api response content: %w\n", err)
	}

	return resp, nil
}
