# gopensearch

This implements something to create custom OpenSearch search engines and add them to your browser.

## Usage

There are no CLI flags at this point, that'll come in the future. For now, modify [`main.go`](./main.go) as necessary to
add the search engines you want and run `gopensearch`. It will launch your browser and let you add each of the search
engines you provided by right-clicking the URL bar and adding them. Once added, I highly recommend adding smart keywords
in the `about:preferences#search` page like e.g. `@gopkg` for `pkg.go.dev`.

```bash
go build && ./gopensearch
```

Spec: https://github.com/dewitt/opensearch/blob/master/opensearch-1-1-draft-6.md
