package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/cdecker/kugelblitz/bitcoin"
	"github.com/cdecker/kugelblitz/lightningrpc"
	"github.com/cdecker/kugelblitz/webui"
	log "github.com/sirupsen/logrus"
)

var (
	lightningRpc *lightningrpc.LightningRpc
	bitcoinRpc   *bitcoin.BitcoinRpc
)

var (
	lightningSock    = flag.String("lightning-socket", filepath.Join(os.Getenv("HOME"), ".lightning/lightning-rpc"), "Location of the lightning unix domain socket.")
	bitcoinRpcLoc    = flag.String("bitcoin-rpc", "localhost:18332", "Location where bitcoind is listening for RPC calls.")
	bitcoinRpcUser   = flag.String("bitcoin-user", "rpcuser", "")
	bitcoinRpcPass   = flag.String("bitcoin-pass", "rpcpass", "")
	manageBitcoind   = flag.Bool("start-bitcoind", true, "Do you want kugelblitz to manage bitcoind if it is not running?")
	manageLightningd = flag.Bool("start-lightningd", true, "Do you want kugelblitz to manage lightning if it is not running?")
	bitcoinDatadir   = flag.String("bitcoin-datadir", filepath.Join(os.Getenv("HOME"), ".bitcoin"), "Where is the bitcoind datadir?")
	networkParams    = flag.String("network", "testnet", "Network to use (mainnet, testnet or regtest).")
	port             = flag.Int("port", 19735, "Port to listen on for HTTP clients.")
	debug            = flag.Bool("debug", false, "Be very verbose")
)

func maybeStartBitcoind(bitcoinRpc *bitcoin.BitcoinRpc) bitcoin.BitcoinD {
	if !*manageBitcoind {
		log.Debug("Not checking for running bitcoind since we don't manage it.")
		return nil
	} else if bitcoinRpc.IsAlive() {
		log.Debug("No need to run bitcoind, it's already alive")
		return nil
	} else {
		log.Debug("Starting bitcoind")

		b := bitcoin.NewBitcoinD(bitcoin.BitcoinDOpts{
			Datadir: *bitcoinDatadir,
		})
		b.Start()
		return b
	}

}

func main() {
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	lrpc := lightningrpc.NewLightningRpc(*lightningSock)
	lightningRpc := webui.NewLightning(lrpc)

	bitcoinRpc = bitcoin.NewBitcoinRpc(
		fmt.Sprintf("http://%s:%s@localhost:18332", *bitcoinRpcUser, *bitcoinRpcPass))
	nodeRpc := bitcoin.NewNode(&lightningRpc, bitcoinRpc)

	b := maybeStartBitcoind(bitcoinRpc)
	if b != nil {
		defer b.Stop()
	}

	rpc.Register(bitcoinRpc)
	rpc.Register(&lightningRpc)
	rpc.Register(nodeRpc)

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
