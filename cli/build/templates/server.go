package templates

var ServerTemplate = `
	package main

	import (
    "net/http"
		"github.com/go-chi/chi/v5"
		"github.com/go-chi/chi/v5/middleware"

		{{ range .Imports }}
    {{ .HandlerPackage }}  "{{ .PackagePath }}"
		{{ end }}
	)

	func main() {
		router := chi.NewRouter()
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
    router.Use(middleware.StripSlashes)

		{{ range .Routes }}
    router.Method("{{ .Method }}", "{{ .Pattern }}", http.HandlerFunc({{ .Handler }}.{{ .Method }}))
		{{ end}}
		
		http.ListenAndServe(":"+ "{{ .Port }}", router)
	}
`
