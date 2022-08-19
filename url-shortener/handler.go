package main

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathUrl
	if err := yaml.Unmarshal(yml, &pathUrls); err != nil {
		return nil, err
	}

	pathsToUrls := make(map[string]string)
	for _, pathUrl := range pathUrls {
		pathsToUrls[pathUrl.Path] = pathUrl.Url
	}

	return MapHandler(pathsToUrls, fallback), nil
}

type pathUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
