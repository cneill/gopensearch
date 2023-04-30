package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	NS    = "http://a9.com/-/spec/opensearch/1.1/"
	MOZNS = "http://www.mozilla.org/2006/browser/search/"
)

// OpenSearchDescriptionOpts has the user-configurable pieces of an OpenSearch description document.
type OpenSearchDescriptionOpts struct {
	ShortName     string           `xml:"ShortName"`
	LongName      string           `xml:"LongName"`
	Description   string           `xml:"Description"`
	Tags          []string         `xml:"Tags"` // space separated
	URL           *OpenSearchURL   `xml:"Url"`
	Image         *OpenSearchImage `xml:"Image,omitempty"`
	Developer     string           `xml:"Developer,omitempty"`
	InputEncoding string           `xml:"InputEncoding,omitempty"`
	// Query         *OpenSearchQuery `xml:"Query"`
}

type OpenSearchImage struct {
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Type   string `xml:"type,attr,omitempty"`
	URL    string `xml:",chardata"`

	// These are not part of the actual OpenSearchDescription, they configure other behavior

	// FaviconDomain tells us the domain to get the favicon for if it is non-empty
	FaviconDomain string `xml:"-"`
}

func (o *OpenSearchImage) Validate() error {
	if o.Width == 0 {
		return fmt.Errorf("no width supplied for the image")
	} else if o.Height == 0 {
		return fmt.Errorf("no height was supplied for the image")
	}

	return nil
}

type OpenSearchURL struct {
	Template string `xml:"template,attr"`
	Type     string `xml:"type,attr"`
	Method   string `xml:"method,attr"`
}

func (o *OpenSearchURL) Validate() error {
	if o.Template == "" {
		return fmt.Errorf("No Template was supplied")
	}

	// TODO: validate Type

	return nil
}

// TODO: OpenSearchQuery ?

func (o *OpenSearchDescription) Name() string {
	if o.LongName != "" {
		return o.LongName
	}

	return o.ShortName
}

func (o *OpenSearchDescriptionOpts) Validate() error {
	switch {
	// case len(o.ShortName) > 16:
	// 	return fmt.Errorf("supplied 'ShortName' was longer than 16 characters")
	case len(o.LongName) > 48:
		return fmt.Errorf("supplied 'LongName' was longer than 48 characters")
	case len(o.Description) > 1024:
		return fmt.Errorf("supplied 'Description' was longer than 1024 characters")
	case len(strings.Join(o.Tags, " ")) > 256:
		return fmt.Errorf("total length of supplied 'Tags' was greater than 256 characters")
	case len(o.Developer) > 64:
		return fmt.Errorf("supplied 'Developer' was longer than 64 characters")
	}

	if err := o.URL.Validate(); err != nil {
		return fmt.Errorf("failed to validate URL: %w", err)
	}

	if err := o.Image.Validate(); err != nil {
		return fmt.Errorf("failed to validate Image: %w", err)
	}

	return nil
}

// OpenSearchDescription contains the data used to populate an OpenSearch description document in XML.
type OpenSearchDescription struct {
	XMLName xml.Name `xml:"OpenSearchDescription"`
	NS      string   `xml:"xmlns,attr"`
	MOZNS   string   `xml:"xmlns:moz,attr"`

	*OpenSearchDescriptionOpts

	err error
}

// NewOpenSearchDescription returns an object that is ready to be XML encoded for use by OpenSearch clients.
func NewOpenSearchDescription(opts *OpenSearchDescriptionOpts) *OpenSearchDescription {
	return &OpenSearchDescription{
		NS:    NS,
		MOZNS: MOZNS,

		OpenSearchDescriptionOpts: opts,
	}
}

func (o *OpenSearchDescription) Err() error {
	if o.err != nil {
		return o.err
	}

	return nil
}

// Validate returns any errors in the object's fields.
func (o *OpenSearchDescription) Validate() error {
	if err := o.Err(); err != nil {
		return err
	}

	return o.OpenSearchDescriptionOpts.Validate()
}

func (o *OpenSearchDescription) FillFavicon() *OpenSearchDescription {
	var faviconDomain string

	if o.Image != nil && o.Image.FaviconDomain != "" {
		faviconDomain = o.Image.FaviconDomain
	} else {
		// parse the "template" attr to get the URL's domain
		parsedTemplate, err := url.Parse(o.URL.Template)
		if err != nil {
			o.err = fmt.Errorf("failed to parse URL template %q to grab favicon: %w", o.URL.Template, err)
			return o
		}

		faviconDomain = parsedTemplate.Hostname()
	}

	// use DuckDuckGo's service for grabbing favicons
	// ex: https://icons.duckduckgo.com/ip3/go.dev.ico
	iconURL := fmt.Sprintf("https://icons.duckduckgo.com/ip3/%s.ico", faviconDomain)
	resp, err := http.Get(iconURL)
	if err != nil {
		o.err = fmt.Errorf("failed to retrieve favicon URL %q from DDG: %w", iconURL, err)
		return o
	}

	defer resp.Body.Close()
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		o.err = fmt.Errorf("failed to read favicon contents from DDG (URL: %q): %w", iconURL, err)
		return o
	}

	// encode the returned icon as a base64 data URL and set the Image URL to that value
	encoded := base64.StdEncoding.EncodeToString(contents)

	contentType := resp.Header.Get("Content-Type")

	if contentType == "" {
		contentType = "image/x-icon"
	}

	if o.Image == nil {
		o.Image = &OpenSearchImage{}
	}

	o.Image.Type = contentType
	o.Image.URL = fmt.Sprintf("data:%s;base64,%s", contentType, encoded)

	return o
}

type SearchEngines []*OpenSearchDescription

func (s SearchEngines) Validate() error {
	for _, engine := range s {
		if err := engine.Validate(); err != nil {
			return fmt.Errorf("search engine validation failed: %w", err)
		}
	}

	return nil
}

func (s SearchEngines) GetByShortName(shortName string) *OpenSearchDescription {
	for _, engine := range s {
		if engine.ShortName == shortName {
			return engine
		}
	}

	return nil
}
