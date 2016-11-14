package webui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cdecker/kugelblitz/static"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/prometheus/common/log"
)

func init() {
	http.HandleFunc("/", handler)
	http.Handle("/rpc/", jsonrpc2.HTTPHandler(nil))
	http.HandleFunc("/static/", staticHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	index, err := static.Asset("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
	}
	w.Write(index)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	c, err := static.Asset(r.RequestURI[8:])
	if err != nil {
		log.Errorf("Error reading static resource: %v", err)
		http.Error(w, fmt.Sprintf("%v", err), 404)
	}
	if strings.HasSuffix(r.RequestURI, ".css") {
		w.Header().Add("Content-Type", "text/css")
	}

	w.Write(c)
}
