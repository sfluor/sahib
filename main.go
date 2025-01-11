package main

import (
	"fmt"
	"log"
	"net/http"
	"sahib/components"
)


func main() {
    mainHandler := func(w http.ResponseWriter, r *http.Request) {
		component := components.Index()
		component.Render(r.Context(), w)
	}

    http.HandleFunc("GET /", mainHandler)

    http.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
        search := r.FormValue("search")
        log.Printf("Searching for: %s", search)
        res := "Fake result: " + search

        component := components.Result(res)
        component.Render(r.Context(), w)
    })

    log.Print("Listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


    // <footer>
    //   @InputGroup("registerKey", "", "Your Perplexity API Key", "Register API Key", templ.Attributes{})
    //   </footer>

func prompt(word string) string {
    count := 5
    target := "french"
return fmt.Sprintf(`Give me %d examples of useful and relevant sentences from medias or stories with the proper harakats.

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
