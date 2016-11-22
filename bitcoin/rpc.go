package bitcoin

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/cdecker/kugelblitz/lightningrpc"
)

type BitcoinRpc struct {
	network string
	url     string
}

type Node struct {
	lightningRpc *lightningrpc.LightningRpc
	bitcoinRpc   *BitcoinRpc
	fundingAddr  string
}

type GetBInfoResponse struct {
	Version         uint32  `json:"version"`
	Protocolversion uint    `json:"protocolversion"`
	Walletversion   uint    `json:"walletversion"`
	Balance         float32 `json:"balance"`
	Blocks          uint    `json:"blocks"`
	Timeoffset      int     `json:"Timeoffset"`
	Connections     uint    `json:"connections"`
	Difficulty      float32 `json:"difficulty"`
	Testnet         bool    `json:"testnet"`
	Keypoololdest   uint    `json:"keypoololdest"`
	Keypoolsize     uint    `json:"keypoolsize"`
	Paytxfee        float32 `json:"paytxfee"`
	Relayfee        float32 `json:"relayfee"`
	Errors          string  `json:"errors"`
}

type ConnectPeerRequest struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Capacity uint64 `json:"capacity"`
	Async    bool   `json:"async"`
}

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

type SendToAddressRequest struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type TxReference struct {
	TransactionId string `json:"txid"`
}

type GetRawTransactionResponse struct {
	RawTransaction string `json:"rawtx"`
}

type Address struct {
	Addr string `json:"addr"`
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

func (b *BitcoinRpc) SendToAddress(req *SendToAddressRequest, res *TxReference) error {
	var params []interface{}
	params = append(params, req.Address)
	params = append(params, req.Amount)

	return b.call("sendtoaddress", params, &res.TransactionId)
}

func (b *BitcoinRpc) GetRawTransaction(req *TxReference, res *GetRawTransactionResponse) error {
	var params []interface{}
	params = append(params, req.TransactionId)

	return b.call("getrawtransaction", params, &res.RawTransaction)
}

func (b *BitcoinRpc) GetNewAddress(req *lightningrpc.Empty, res *string) error {
	var params []interface{}
	return b.call("getnewaddress", params, res)
}

func (b *BitcoinRpc) call(method string, params []interface{}, res interface{}) error {
	request := map[string]interface{}{
		"method": method,
		"id":     0,
		"params": params,
	}

	response := struct {
		Result interface{} `json:"result"`
		Error  interface{} `json:"error"`
		Id     uint        `json:"id"`
	}{
		Result: res,
	}

	log.Debugf("Calling bitcoin.%s with args %v", method, params)

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := http.Post(b.url,
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}
	if response.Error == nil {
		return nil
	} else {
		return fmt.Errorf("Error reported by server: %v", response.Error)
	}
}

func (b *BitcoinRpc) GetInfo(_ *lightningrpc.Empty, response *GetBInfoResponse) error {
	return b.call("getinfo", nil, response)
}

func NewBitcoinRpc(url string) *BitcoinRpc {
	return &BitcoinRpc{
		network: "-testnet",
		url:     url,
	}
}

func (br *BitcoinRpc) IsAlive() bool {
	return br.GetInfo(&lightningrpc.Empty{}, &GetBInfoResponse{}) == nil
}

func (bc *BitcoinRpc) exec(method string, args []string) (string, error) {
	a := append([]string{bc.network, method}, args...)
	c := exec.Command("/usr/local/bin/bitcoin-cli", a...)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error calling bitcoin rpc: %s", err)
	}
	return strings.TrimSpace(string(out[:])), nil
}

func NewNode(lrpc *lightningrpc.LightningRpc, brpc *BitcoinRpc) *Node {
	return &Node{
		lightningRpc: lrpc,
		bitcoinRpc:   brpc,
	}
}

func (n *Node) ConnectPeer(req *ConnectPeerRequest, res *lightningrpc.Empty) error {
	log.Debugf("Connecting to %s:%d", req.Host, req.Port)
	var addrResp lightningrpc.NewAddressResponse
	err := n.lightningRpc.NewAddress(&lightningrpc.Empty{}, &addrResp)
	if err != nil {
		log.Error(err)
		return err
	}

	var sendResp TxReference
	err = n.bitcoinRpc.SendToAddress(&SendToAddressRequest{
		Address: addrResp.Address,
		Amount:  fmt.Sprintf("%f", float64(req.Capacity)*10e-8),
	}, &sendResp)
	if err != nil {
		log.Error(err)
		return err
	}

	var rawResp GetRawTransactionResponse
	err = n.bitcoinRpc.GetRawTransaction(&sendResp, &rawResp)
	if err != nil {
		log.Error(err)
		return err
	}

	// Finally we need to tell lightningd to connect to that node
	// with the funds provided
	connReq := &lightningrpc.ConnectRequest{
		Host:         req.Host,
		Port:         req.Port,
		FundingTxHex: rawResp.RawTransaction,
	}
	if req.Async {
		go n.lightningRpc.Connect(connReq, &lightningrpc.Empty{})
	} else {
		return n.lightningRpc.Connect(connReq, &lightningrpc.Empty{})
	}

	return nil
}

type FundingAddr struct {
	Addr string `json:"addr"`
}

func (n *Node) GetFundingAddr(req *lightningrpc.Empty, res *Address) error {
	var err error
	if n.fundingAddr == "" {
		err = n.bitcoinRpc.GetNewAddress(req, &n.fundingAddr)
	}
	res.Addr = n.fundingAddr
	return err
}

type KugelblitzInfo struct {
}

func (n *Node) GetInfo(req *lightningrpc.Empty, res *KugelblitzInfo) error {
	return nil
}
