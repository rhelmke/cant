// Package static implements an HTTP request handler for static content. cant is a self-contained
// application. This means that static assets are embedded into the binary using go-bindata.
package static

import (
	"cant/webserver/webinterface"
	"mime"
	"net/http"
	"path"
)

// ContentHandler handles static routes by extracting assets from memory
func ContentHandler(resp http.ResponseWriter, req *http.Request) {
	// set default
	p := req.URL.Path
	if p == "" {
		p = "index.html"
	}
	// check existance and extract content from memory
	content, err := webinterface.Asset(p)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
	}
	// determine MIME type
	resp.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(p)))
	// write content into response
	resp.Write(content)
}
