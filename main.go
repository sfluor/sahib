package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sahib/clients"
	"sahib/components"
	"sahib/model"
	"sync"
)

type queryFunc func(word string) (*model.Translations, error)

type source struct {
	name string
	fn   queryFunc
}

func getQueryFunc(name string, apiKey string) (queryFunc, error) {
	var fn queryFunc
	switch name {
	case model.SourceElixir:
		fn = clients.QueryElixir
	case model.SourceMaany:
		fn = clients.QueryMaany
	case model.SourcePerplexity:
		if apiKey != "" {
			client := clients.PerplexityClient{ApiKey: apiKey}
			fn = client.Query
		} else {
			return func(word string) (*model.Translations, error) {
				return &model.Translations{}, nil
			}, nil
		}
	default:
		return nil, fmt.Errorf("unknown source: " + name)
	}

	return fn, nil
}

func handleErr(res *model.Translations, err error, w http.ResponseWriter, msg string, args ...any) (*model.Translations, bool) {
	if err != nil {
		errStr := fmt.Sprintf(msg, args...)
		log.Printf(errStr)
		// http.Error(w, errStr, 500)
		return &model.Translations{Error: err.Error()}, true
	}

	return res, false
}

func isSourceEnabled(r *http.Request, source string) bool {
			return r.FormValue(source) == "on" 
}

func main() {
	if len(os.Args) < 2 {
		panic("Please provide the path to the hans wehr sqlite database")
	}

	sqlpath := os.Args[1]
	hansWehr, err := clients.NewHansWehrClient(sqlpath)
	if err != nil {
		panic(err)
	}

	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		component := components.Index()
		component.Render(r.Context(), w)
	}

	http.HandleFunc("GET /", mainHandler)

	http.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		search := r.FormValue(model.Search)
		apiKey := r.FormValue(model.ApiKey)
		log.Printf("Searching for: %s", search)

		defaultSources := []source{
			{
				name: model.SourceElixir,
				fn:   clients.QueryElixir,
			},
			{
				name: model.SourceMaany,
				fn:   clients.QueryMaany,
			},
		}

		if apiKey != "" {
			client := clients.PerplexityClient{ApiKey: apiKey}
			defaultSources = append(defaultSources, source{name: model.SourcePerplexity, fn: client.Query})
		}

		sources := make([]source, 0, len(defaultSources))

		// Remove disabled sources
		for _, source := range defaultSources {
			if isSourceEnabled(r, source.name){
				sources = append(sources, source)
			}
		}

		all := make([]model.TranslationsAndSource, len(sources))

		var wg sync.WaitGroup
		for i, src := range sources {

			name := src.name
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				res, err := src.fn(search)
				res, _ = handleErr(res, err, w, "Failed to create client for: %s: %w", name, err)
				all[idx] = model.TranslationsAndSource{Translations: res, Source: name}
			}(i)
		}

        defs := &model.Definitions{}
		wg.Add(1)
		go func() {
			defer wg.Done()
            if !isSourceEnabled(r, model.SourceWehr) {
                return
            }

			resp, err := hansWehr.Query(search)
			if err != nil {
				log.Printf("Failed to fetch hans wehr data for %s: %s", search, err)
			}
			defs = resp
		}()

		// TODO: fix error path

		wg.Wait()
		component := components.Results(all, defs)
		component.Render(r.Context(), w)
	})

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
