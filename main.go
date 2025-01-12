package main

import (
	"fmt"
	"log"
	"net/http"
	"sahib/clients"
	"sahib/components"
	"sahib/model"
	"sync"
)

type queryFunc func(word string) (model.Translations, error)

type source struct {
	name string
	fn   queryFunc
}

const (
	SourceElixir     = "Elixir"
	SourceMaany      = "Maany"
	SourcePerplexity = "Perplexity"

	ApiKey = "apiKey"
	Search = "search"
)

func getQueryFunc(name string, apiKey string) (queryFunc, error) {
	var fn queryFunc
	switch name {
	case SourceElixir:
		fn = clients.QueryElixir
	case SourceMaany:
		fn = clients.QueryMaany
	case SourcePerplexity:
		if apiKey != "" {
			client := clients.PerplexityClient{ApiKey: apiKey}
			fn = client.Query
		} else {
			return func(word string) (model.Translations, error) {
				return model.Translations{}, nil
			}, nil
		}
	default:
		return nil, fmt.Errorf("unknown source: " + name)
	}

	return fn, nil
}

func handleErr(res model.Translations, err error, w http.ResponseWriter, msg string, args ...any) (model.Translations, bool) {
	if err != nil {
		errStr := fmt.Sprintf(msg, args...)
		log.Printf(errStr)
		// http.Error(w, errStr, 500)
		return model.Translations{Error: err.Error()}, true
	}
	return res, false
}

func main() {
	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		component := components.Index()
		component.Render(r.Context(), w)
	}

	http.HandleFunc("GET /", mainHandler)

	sourceNames := []string{SourceElixir, SourceMaany, SourcePerplexity}

	for _, name := range sourceNames {
		http.HandleFunc("POST /search/"+name, func(w http.ResponseWriter, r *http.Request) {
			search := r.FormValue(Search)
			apiKey := r.FormValue(ApiKey)
			fn, err := getQueryFunc(name, apiKey)
			if res, failed := handleErr(model.Translations{}, err, w, "Failed to create client for: %s: %w", name, err); failed {
				component := components.Result(name, res.Link, res.TimeMs, res.List)
				component.Render(r.Context(), w)
				return
			}

			res, err := fn(search)
			res, _ = handleErr(res, err, w, "Failed to create client for: %s: %w", name, err)
			component := components.Result(name, res.Link, res.TimeMs, res.List)
			component.Render(r.Context(), w)
		})
	}

	http.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		search := r.FormValue(Search)
		apiKey := r.FormValue(ApiKey)
		log.Printf("Searching for: %s", search)

		sources := []source{
			{
				name: "Elixir",
				fn:   clients.QueryElixir,
			},
			{
				name: "Maany",
				fn:   clients.QueryMaany,
			},
		}

		if apiKey != "" {
			client := clients.PerplexityClient{ApiKey: apiKey}
			sources = append(sources, source{name: "Perplexity", fn: client.Query})
		}

		all := make([]model.TranslationsAndSource, len(sourceNames))

		var wg sync.WaitGroup
		for i, name := range sourceNames {
			fn, err := getQueryFunc(name, apiKey)

			if res, failed := handleErr(model.Translations{}, err, w, "Failed to create client for: %s: %w", name, err); failed {
				all[i] = model.TranslationsAndSource{Translations: res, Source: name}
				continue
			}

			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				res, err := fn(search)
				res, _ = handleErr(res, err, w, "Failed to create client for: %s: %w", name, err)
				all[idx] = model.TranslationsAndSource{Translations: res, Source: name}
			}(i)
		}

		wg.Wait()
		component := components.Results(all)
		component.Render(r.Context(), w)
	})

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
