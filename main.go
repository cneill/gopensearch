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
			Tags:        []string{"golang"},
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
		NewOpenSearchDescription(&OpenSearchDescriptionOpts{
			ShortName:   "Search deps.dev",
			Description: "Search for package dependencies/versions on deps.dev",
			Tags:        []string{"vulnerabilities"},
			URL: &OpenSearchURL{
				Template: "https://deps.dev/search?q={searchTerms}",
				Type:     "text/html",
				Method:   "get",
			},
			Image: &OpenSearchImage{
				FaviconDomain: "deps.dev",
				Width:         16,
				Height:        16,
			},
		}).FillFavicon(),
		NewOpenSearchDescription(&OpenSearchDescriptionOpts{
			ShortName:   "Search kubernetes.io",
			Description: "Search for documentation on the Kubernetes site",
			Tags:        []string{"k8s"},
			URL: &OpenSearchURL{
				Template: "https://kubernetes.io/search/?q={searchTerms}",
				Type:     "text/html",
				Method:   "get",
			},
			Image: &OpenSearchImage{
				FaviconDomain: "kubernetes.io",
				Width:         16,
				Height:        16,
			},
		}).FillFavicon(),
		NewOpenSearchDescription(&OpenSearchDescriptionOpts{
			ShortName:   "Terraform module search",
			Description: "Search for Terraform modules on registry.terraform.io",
			Tags:        []string{"terraform"},
			URL: &OpenSearchURL{
				Template: "https://registry.terraform.io/search/modules?q={searchTerms}",
				Type:     "text/html",
				Method:   "get",
			},
			Image: &OpenSearchImage{
				FaviconDomain: "terraform.io",
				Width:         16,
				Height:        16,
			},
		}).FillFavicon(),
		NewOpenSearchDescription(&OpenSearchDescriptionOpts{
			ShortName:   "Terraform provider search",
			Description: "Search for Terraform providers on registry.terraform.io",
			Tags:        []string{"terraform"},
			URL: &OpenSearchURL{
				Template: "https://registry.terraform.io/search/providers?q={searchTerms}",
				Type:     "text/html",
				Method:   "get",
			},
			Image: &OpenSearchImage{
				FaviconDomain: "terraform.io",
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
