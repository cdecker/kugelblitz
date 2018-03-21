package lightningrpc

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/powerman/rpc-codec/jsonrpc2"
	log "github.com/sirupsen/logrus"
)

type LightningRpc struct {
	socketPath string
	conn       net.Conn
}

type Empty struct{}

type NewAddressResponse struct {
	Address string `json:"address"`
}

type GetInfoResponse struct {
	Id          string `json:"id"`
	Port        uint   `json:"port"`
	Testnet     bool   `json:"testnet"`
	Version     string `json:"version"`
	BlockHeight uint   `json:"blockheight"`
}

type ListPeersResult struct {
	Peers []Peer
}

type PeerChannel struct {
	State         string `json:"state"`
	FundingTxId   string `json:"funding_txid"`
	MsatoshiToUs  uint64 `json:"msatoshi_to_us"`
	MsatoshiTotal uint64 `json:"msatoshi_total"`
}

type Peer struct {
	State       string        `json:"state"`
	PeerId      string        `json:"id"`
	Connected   bool          `json:"connected"`
	OurAmount   int           `json:"our_amount"`
	TheirAmount int           `json:"their_amount"`
	OurFee      int           `json:"our_fee"`
	TheirFee    int           `json:"their_fee"`
	Channels    []PeerChannel `json:"channels"`
}

type Channel struct {
	From            string `json:"source"`
	To              string `json:"destination"`
	BaseFee         uint64 `json:"base_fee_millisatoshi"`
	ProportionalFee uint64 `json:"fee_per_millionth"`
	ShortChannelId  string `json:"short_channel_id"`
	Flags           uint
	LastUpdate      uint64 `json:"last_update"`
	Delay           uint
	Satoshis        uint64
	Active          bool `json:"active"`
	Public          bool `json:"public"`
}

type ListChannelsResponse struct {
	Channels []Channel `json:"channels"`
}

type AddFundsRequest struct {
	RawTransaction string `json:"rawtx"`
}

type RouteHop struct {
	NodeId  string `json:"id"`
	Amount  uint64 `json:"msatoshi"`
	Delay   uint32 `json:"delay"`
	Channel string `json:"channel"`
}

type Route struct {
	Hops []RouteHop `json:"route"`
}

type GetRouteRequest struct {
	Destination string  `json:"destination"`
	Amount      uint64  `json:"amount"`
	RiskFactor  float32 `json:"risk"`
}

type FundChannelRequest struct {
}

type DecodePayRequest struct {
	PayRequest string `json:"payrequest"`
}

type DecodePayResponse struct {
	Currency    string `json:"currency"`
	Timestamp   uint64 `json:"timestamp"`
	Expiry      uint32 `json:"expiry"`
	Payee       string `json:"payee"`
	Amount      uint64 `json:"msatoshi"`
	Description string `json:"description"`
	PaymentKey  string `json:"payment_hash"`
	Signature   string `json:"signature"`
}

func (lr *LightningRpc) call(method string, req interface{}, res interface{}) error {
	log.Debugf("Calling lightning.%s with args %v", method, req)

	clientTCP, err := jsonrpc2.Dial("unix", lr.socketPath)
	if err != nil {
		return err
	}
	defer clientTCP.Close()
	err = clientTCP.Call(method, req, res)
	if err != nil {
		log.Debugf("error calling %s: %v", method, err)
		return errors.Wrap(err, fmt.Sprintf("error calling %s", method))
	} else {
		log.Debugf("method %s returned %v", method, err)
		return nil
	}
}

func (lr *LightningRpc) NewAddress() (NewAddressResponse, error) {
	res := NewAddressResponse{}
	err := lr.call("newaddr", &Empty{}, &res)
	return res, err
}

func (lr *LightningRpc) GetInfo() (GetInfoResponse, error) {
	res := GetInfoResponse{}
	err := lr.call("getinfo", &Empty{}, &res)
	return res, err
}

func (lr *LightningRpc) ListChannels() (ListChannelsResponse, error) {
	res := ListChannelsResponse{}
	err := lr.call("listchannels", &Empty{}, &res)
	return res, err
}

type ListPeersResponse struct {
	Peers []Peer `json:"peers"`
}

func (lr *LightningRpc) ListPeers() (ListPeersResponse, error) {
	res := ListPeersResponse{}
	err := lr.call("listpeers", &Empty{}, &res)
	return res, err
}

func (lr *LightningRpc) Connect(nodeid string, address string, port uint) error {
	var params []interface{}
	params = append(params, nodeid)
	params = append(params, address)
	params = append(params, port)
	return lr.call("connect", params, &Empty{})
}

type PeerReference struct {
	PeerId string `json:"peerid"`
}

func (lr *LightningRpc) Close(peerId string) error {
	var params []interface{}
	params = append(params, peerId)
	return lr.call("close", params, &Empty{})
}

func (lr *LightningRpc) GetRoute(destination string, amount uint64, riskfactor float32) (Route, error) {
	var params []interface{}
	params = append(params, destination)
	params = append(params, amount)
	params = append(params, riskfactor)
	res := Route{}
	err := lr.call("getroute", params, &res)
	return res, err
}

type SendPaymentRequest struct {
	Route       []RouteHop `json:"route"`
	PaymentHash string     `json:"paymenthash"`
}

type SendPaymentResponse struct {
	PaymentKey string `json:"preimage"`
}

func (lr *LightningRpc) SendPayment(route []RouteHop, paymentHash string) (SendPaymentResponse, error) {
	var params []interface{}
	params = append(params, route)
	params = append(params, paymentHash)
	res := SendPaymentResponse{}
	err := lr.call("sendpay", params, &res)
	return res, err
}

type NodeAddress struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

type Node struct {
	Id         string        `json:"nodeid"`
	Addresses  []NodeAddress `json:"addresses"`
	Color      string        `json:"color"`
	Alias      string        `json:"alias"`
	LastUpdate uint32        `json:"last_timestamp"`
}

type ListNodesResponse struct {
	Nodes []Node `json:"nodes"`
}

func (lr *LightningRpc) ListNodes() (ListNodesResponse, error) {
	res := ListNodesResponse{}
	err := lr.call("listnodes", &Empty{}, &res)
	return res, err
}

type ConnectRequest struct {
	Host   string `json:"host"`
	Port   uint   `json:"port"`
	NodeId string `json:"nodeid"`
}

type Invoice struct {
	PaymentHash string `json:"rhash"`
	PaymentKey  string `json:"paymentKey"`
	Label       string `json:"label"`
}

func (lr *LightningRpc) Invoice(amount uint64, label string) (Invoice, error) {
	var params []interface{}
	params = append(params, amount)
	params = append(params, label)
	res := Invoice{}
	err := lr.call("invoice", params, &res)
	return res, err
}

func (lr *LightningRpc) AddFunds(rawtx string) error {
	var params []interface{}
	params = append(params, rawtx)
	res := Empty{}
	return lr.call("addfunds", params, &res)
}

func (lr *LightningRpc) FundChannel(nodeid string, capacity uint64) error {
	var params []interface{}
	params = append(params, nodeid)
	params = append(params, capacity)
	res := Empty{}
	return lr.call("fundchannel", params, &res)
}

func (lr *LightningRpc) DecodePay(req *DecodePayRequest, res *DecodePayResponse) error {
	var params []interface{}
	params = append(params, req.PayRequest)
	return lr.call("decodepay", params, res)
}

type ListInvoiceResp struct {
	invoices []Invoice
}

func (lr *LightningRpc) ListInvoice(_ *Empty, res *ListInvoiceResp) error {
	// TODO implement
	return nil
}

type Payment struct {
	Id          uint
	PaymentHash string `json:"payment_hash"`
	Destination string
	Msatoshi    uint64
	Timestamp   uint64
	Status      string
}
type ListPaymentsResp struct {
	Payments []Payment `json:"payments"`
}

type ListFundsOutput struct {
	TransactionID string `json:"txid"`
	OutputIndex   uint32 `json:"output"`
	Value         uint64 `json:"value"`
}

type ListFundsChannel struct {
	PeerId              string `json:"peer_id"`
	ChannelSatoshi      uint64 `json:"channel_sat"`
	ChannelTotalSatoshi uint64 `json:"channel_total_sat"`
}

type ListFundsResponse struct {
	Outputs  []ListFundsOutput
	Channels []ListFundsChannel
}

func (lr *LightningRpc) ListFunds(_ *Empty, res *ListFundsResponse) error {
	var params []interface{}
	params = append(params, nil)
	params = append(params, nil)
	return lr.call("listfunds", params, res)
}

func (lr *LightningRpc) ListPayments(_ *Empty, res *ListPaymentsResp) error {
	var params []interface{}
	params = append(params, nil)
	params = append(params, nil)
	return lr.call("listpayments", params, res)
}
func NewLightningRpc(socketPath string) *LightningRpc {
	return &LightningRpc{
		socketPath: socketPath,
	}
}

func (lr *LightningRpc) Stop() error {
	var params []interface{}
	res := Empty{}
	return lr.call("stop", params, &res)
}
