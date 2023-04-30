package main

import (
	"fmt"

	"github.com/pkg/browser"
)

func main() {
	opts := &ServerOpts{
		Host: "localhost",
		Port: 5030,
	}

	server, err := NewServer(opts)
	if err != nil {
		panic(err)
	}

	err = server.LoadSearchEngines(SearchEngines{
		NewOpenSearchDescription(&OpenSearchDescriptionOpts{
			ShortName:   "Search Go packages",
			Description: "Search Go packages on pkg.go.dev",
			Tags:        []string{"golang @golang whaaazaaa"},
			URL: &OpenSearchURL{
				Template: "https://pkg.go.dev/search?q={searchTerms}",
				Type:     "text/html",
				Method:   "get",
			},
			Image: &OpenSearchImage{
				FaviconDomain: "go.dev",
				Width:         16,
				Height:        16,
			},
		}).FillFavicon(),
	})
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("http://%s:%d", opts.Host, opts.Port)

	if opts.TLSEnabled {
		url = fmt.Sprintf("https://%s:%d", opts.Host, opts.Port)
	}

	fmt.Printf("Starting server on %s...\n", url)

	browser.OpenURL(url)

	if err := server.Start(); err != nil {
		panic(err)
	}
}
