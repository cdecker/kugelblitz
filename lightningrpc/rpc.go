package lightningrpc

import (
	"net"
	"net/rpc/jsonrpc"

	log "github.com/Sirupsen/logrus"
	"github.com/powerman/rpc-codec/jsonrpc2"
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

type GetPeersResult struct {
	Peers []Peer
}

type Peer struct {
	State       string `json:"state"`
	PeerId      string `json:"peerid"`
	Connected   bool   `json:"connected"`
	OurAmount   int    `json:"our_amount"`
	TheirAmount int    `json:"their_amount"`
	OurFee      int    `json:"our_fee"`
	TheirFee    int    `json:"their_fee"`
}

func (lr *LightningRpc) call(method string, req interface{}, res interface{}) error {
	log.Debugf("Calling lightning.%s with args %v", method, req)

	clientTCP, err := jsonrpc2.Dial("unix", lr.socketPath)
	if err != nil {
		return err
	}
	defer clientTCP.Close()
	err = clientTCP.Call(method, req, res)
	return err
}

func (lr *LightningRpc) NewAddress(_ *Empty, res *NewAddressResponse) error {
	client, err := jsonrpc.Dial("unix", lr.socketPath)
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Call("newaddr", Empty{}, res)
	return err
}

func (lr *LightningRpc) IsAlive() bool {
	return lr.GetInfo(&Empty{}, &GetInfoResponse{}) == nil
}

func (lr *LightningRpc) GetInfo(req *Empty, res *GetInfoResponse) error {
	return lr.call("getinfo", req, res)
}

type Channel struct {
	From            string `json: "from"`
	To              string `json: "to"`
	BaseFee         uint   `json:"base_fee"`
	ProportionalFee uint   `json:"proportional_fee"`
}

type GetChannelsResponse struct {
	Channels []Channel `json:"channels"`
}

func (lr *LightningRpc) GetChannels(_ *Empty, res *GetChannelsResponse) error {
	client, err := jsonrpc.Dial("unix", lr.socketPath)
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Call("getchannels", Empty{}, res)
	return err
}

type GetPeersResponse struct {
	Peers []Peer `json:"peers"`
}

func (lr *LightningRpc) GetPeers(_ *Empty, res *GetPeersResponse) error {
	client, err := jsonrpc.Dial("unix", lr.socketPath)
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Call("getpeers", Empty{}, res)
	return err
}

func (lr *LightningRpc) Connect(req *ConnectRequest, res *Empty) error {
	var params []interface{}
	params = append(params, req.Host)
	params = append(params, req.Port)
	params = append(params, req.FundingTxHex)
	return lr.call("connect", params, res)
}

type PeerReference struct {
	PeerId string `json:"peerid"`
}

func (lr *LightningRpc) Close(req *PeerReference, res *Empty) error {
	var params []interface{}
	params = append(params, req.PeerId)
	return lr.call("close", params, res)
}

type RouteHop struct {
	NodeId string `json:"id"`
	Amount uint64 `json:"msatoshi"`
	Delay  uint32 `json:"delay"`
}

type Route struct {
	Hops []RouteHop `json:"route"`
}

type GetRouteRequest struct {
	Destination string  `json:"destination"`
	Amount      uint64  `json:"amount"`
	RiskFactor  float32 `json:"risk"`
}

func (lr *LightningRpc) GetRoute(req *GetRouteRequest, res *Route) error {
	var params []interface{}
	params = append(params, req.Destination)
	params = append(params, req.Amount)
	params = append(params, req.RiskFactor)
	return lr.call("getroute", params, res)
}

type SendPaymentRequest struct {
	Route       []RouteHop `json:"route"`
	PaymentHash string     `json:"paymenthash"`
}

type SendPaymentResponse struct {
	PaymentKey string `json:"preimage"`
}

func (lr *LightningRpc) SendPayment(req *SendPaymentRequest, res *SendPaymentResponse) error {
	var params []interface{}
	params = append(params, req.Route)
	params = append(params, req.PaymentHash)
	return lr.call("sendpay", params, res)
}

type Node struct {
	Id   string `json:"nodeid"`
	Port int    `json:"port"`
	Ip   string `json:"hostname"`
}

type GetNodesResponse struct {
	Nodes []Node `json:"nodes"`
}

func (lr *LightningRpc) GetNodes(req *Empty, res *GetNodesResponse) error {
	return lr.call("getnodes", req, res)
}

type ConnectRequest struct {
	Host         string `json:"host"`
	Port         uint   `json:"port"`
	FundingTxHex string `json:"tx"`
}

func NewLightningRpc(socketPath string) *LightningRpc {
	return &LightningRpc{
		socketPath: socketPath,
	}
}
