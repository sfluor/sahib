package main

import (
	"fmt"
	"log"
	"net/http"
	"sahib/clients"
	"sahib/components"
	"sahib/model"
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

func handleErr(err error, w http.ResponseWriter, msg string, args ...any) bool {
	if err != nil {
		errStr := fmt.Sprintf(msg, args...)
		log.Printf(errStr)
		http.Error(w, errStr, 500)
		return true
	}
	return false
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
			if handleErr(err, w, "Failed to create client for: %s: %w", name, err) {
				return
			}

			res, err := fn(search)
			if handleErr(err, w, "Failed to query: %s: %w", name, err) {
				return
			}

			component := components.Result(name, res.List)
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

		all := []model.TranslationsAndSource{}

		for _, name := range sourceNames {
			fn, err := getQueryFunc(name, apiKey)
			if handleErr(err, w, "Failed to create client for: %s: %w", name, err) {
				return
			}

			res, err := fn(search)
			if handleErr(err, w, "Failed to query: %s: %w", name, err) {
				return
			}

			all = append(all, model.TranslationsAndSource{Translations: res, Source: name})
		}

		component := components.Results(all)
		component.Render(r.Context(), w)
	})

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
