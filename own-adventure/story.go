package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

var htmlHandlerTemplate = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Choose Your Own Adventure</title>
    <link rel="stylesheet" href="css/style.css">
</head>

<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
    <p>{{.}}</p>
    {{end}}
    <ul>
        {{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>
</html>`

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(r *http.Request) string
}
type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func defaultPathFn(r *http.Request) string {
	path := r.URL.Path
	if path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, nil, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	chapter, ok := h.s[path]
	if !ok {
		http.Error(w, "Chapter not found.", http.StatusNotFound)
		return
	}
	h.renderTemplate(w, chapter)
}

func (h handler) renderTemplate(w http.ResponseWriter, chapter Chapter) {
	if h.t == nil {
		h.t = template.Must(template.New("").Parse(htmlHandlerTemplate))
	}
	err := h.t.Execute(w, chapter)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func JsonStory(filename string) (Story, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(f)
	var story Story
	if err = d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}
