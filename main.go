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

func main() {
	mainHandler := func(w http.ResponseWriter, r *http.Request) {
		component := components.Index()
		component.Render(r.Context(), w)
	}

	http.HandleFunc("GET /", mainHandler)

	http.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		search := r.FormValue("search")
		log.Printf("Searching for: %s", search)

		apiKey := r.FormValue("apiKey")

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

		for _, src := range sources {
			res, err := src.fn(search)
			if err != nil {
				errStr := fmt.Sprintf("An error occurred querying %s: %w", src.name, err)
				log.Printf(errStr)
				http.Error(w, errStr, 500)
				return
			}
			all = append(all, model.TranslationsAndSource{Translations: res, Source: src.name})

		}
		component := components.Result(all)
		component.Render(r.Context(), w)
	})

	log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
