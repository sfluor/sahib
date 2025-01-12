package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)


func tashkil(sentences []string) ([]string, error) {
    url := "https://www.tashkil.net/api/openai/tashkil"

	requestBody := map[string]interface{}{
        "hasPremium": false,
        "targetLanguage": "French",
        "messages": []map[string]interface{}{
            {
                "role": "user",
                "content": strings.Join(sentences, "\n"),
            },
        },
	}

	// Serialize the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return sentences, fmt.Errorf("Error serializing request body: %w\n", err)
	}

	rawResp, err := queryURL("POST", url, bytes.NewBuffer(jsonBody), nil, false)
	if err != nil {
		return sentences, fmt.Errorf("Error sending perplexity request: %w\n", err)
	}
	defer rawResp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return sentences, fmt.Errorf("Error reading body: %w\n", err)
	}

    out := strings.Split(string(body), "\n")
    if len(out) != len(sentences) {
        return sentences, fmt.Errorf("wrong number of sentences, expected %d got %d", len(sentences), len(out))
    }

    return out, nil
}
