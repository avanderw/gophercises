package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	filename := flag.String("file", "gopher.json", "JSON file with story")
	flag.Parse()
	fmt.Printf("Using JSON file: %s\n", *filename)

	story, err := JsonStory(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	tpl := template.Must(template.New("").Parse(htmlHandlerTemplate))
	h := NewHandler(story, WithTemplate(tpl))
	fmt.Printf("Listening on http://localhost:%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))

}
