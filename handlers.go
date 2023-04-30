package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var indexTemplate = `<!DOCTYPE html>
<html>
<head>
{{- range $search := . }}
    <link rel="search" type="application/opensearchdescription+xml" title="{{ $search.Name }}" href="/osxml/{{ $search.Name }}">
{{- end }}
</head>
<body>
{{- range $search := . }}
	<img width="{{ $search.Image.Width }}" height="{{ $search.Image.Height }}" src="{{ $search.Image.URL | url }}" />
	<b>{{ $search.Name }}</b>
	<br />
{{- end }}
</body>
</html>`

var funcMap = template.FuncMap{
	"url": func(s string) template.URL {
		return template.URL(s)
	},
}

// Index serves the HTML we need to detect the search plugins to be added.
func (s *Server) Index(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	t := template.Must(template.New("index").Funcs(funcMap).Parse(indexTemplate))
	if err := t.Execute(w, s.searchEngines); err != nil {
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
		fmt.Printf("ERROR: %v\n", err)

		return
	}
}

// OSXML produces the OpenSearchDescription XML for the search engine specified by "name" in the URL.
func (s *Server) OSXML(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	osdName := ps.ByName("name")

	engine := s.searchEngines.GetByShortName(osdName)

	xmlBytes, err := xml.MarshalIndent(engine, "", "  ")
	if err != nil {
		http.Error(w, "failed to generate XML", http.StatusInternalServerError)
		fmt.Printf("ERROR: %v\n", err)

		return
	}

	w.Header().Set("Content-Type", "application/opensearchdescription+xml")
	w.WriteHeader(http.StatusOK)
	w.Write(xmlBytes)
}
