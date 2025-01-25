package clients

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

const ContentType = "Content-Type"

func elapsed(start time.Time) string {
	return time.Now().Sub(start).Truncate(10 * time.Millisecond).String()
}

func queryURL(
	typ string,
	url string,
	body io.Reader,
	headers map[string]string,
	withHTTP2 bool) (*http.Response, error) {
	req, err := http.NewRequest(typ, url, body)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:131.0) Gecko/20100101 Firefox/131.0")
	req.Header.Set("Accept", "*/*")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	client := &http.Client{}
	if withHTTP2 {
		// http1 doesn't work with the maany website
		client.Transport = &http2.Transport{}
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}

	if res.StatusCode != 200 {
		text, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status code error: %d %s\n%s", res.StatusCode, res.Status, truncate(string(text), 256))
	}

	return res, nil
}

func truncate(text string, max int) string {
    if len(text) > max {
        return text[:max]
    }

    return text
}
