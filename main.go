package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/cdecker/lightningctl/coil"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

var (
	lightning  *coil.LightningRpc
	bitcoinRpc *coil.Bitcoin
)

var (
	lightningSock = flag.String("lightning-socket", "/tmp/lightning2/lightning-rpc",
		"Location of the lightning unix domain socket.")
	bitcoinRpcLoc = flag.String("bitcoin-rpc", "localhost:18332",
		"Location where bitcoind is listening for RPC calls.")
	networkParams = flag.String("network", "testnet",
		"Network to use (mainnet, testnet or regtest).")
	port = flag.Int("port", 8000, "Port to listen on for HTTP clients.")
)

func handler(w http.ResponseWriter, r *http.Request) {
	index, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
	}
	w.Write(index)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	c, err := ioutil.ReadFile("static/" + r.RequestURI[8:])
	if err != nil {
		log.Errorf("Error reading static resource: %v", err)
		http.Error(w, fmt.Sprintf("%v", err), 404)
	}
	if strings.HasSuffix(r.RequestURI, ".css") {
		w.Header().Add("Content-Type", "text/css")
	}

	w.Write(c)
}

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)

	lightning = coil.NewLightningRpc(*lightningSock)
	bitcoinRpc = coil.NewBitcoinRpc("http://rpcuser:rpcpass@localhost:18332")
	nodeRpc := coil.NewNode(lightning, bitcoinRpc)

	rpc.Register(bitcoinRpc)
	rpc.Register(lightning)
	rpc.Register(nodeRpc)

	http.HandleFunc("/", handler)
	http.Handle("/rpc/", jsonrpc2.HTTPHandler(nil))
	http.HandleFunc("/static/", staticHandler)
	http.ListenAndServe(":8000", nil)
}
