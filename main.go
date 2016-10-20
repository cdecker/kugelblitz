package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcrpcclient"
	"github.com/cdecker/lightningctl/coil"
)

var (
	lightning *coil.LightningRpc
	connCfg   = &btcrpcclient.ConnConfig{
		Host:         "localhost:18332",
		User:         "rpcuser",
		Pass:         "rpcpass",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
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
	c, err := ioutil.ReadFile(r.RequestURI[8:])
	if err != nil {
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
	bitcoinRpc = coil.NewBitcoinRpc()

	rpc.Register(bitcoinRpc)
	rpc.Register(lightning)

	http.HandleFunc("/", handler)
	http.HandleFunc("/rpc/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err := rpc.ServeRequest(serverCodec)
		if err != nil {
			log.Printf("Error while serving JSON request: %v", err)
			http.Error(w, "Error while serving JSON request, details have been logged.", 500)
			return
		}

	}))
	http.HandleFunc("/static/", staticHandler)
	http.ListenAndServe(":8000", nil)
}
