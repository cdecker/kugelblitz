package coil

import (
	"net"
	"net/rpc/jsonrpc"
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
	Id      string `json:"id"`
	Port    uint   `json:"port"`
	Testnet bool   `json:"testnet"`
	Version string `json:"version"`
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
	client, err := jsonrpc.Dial("unix", lr.socketPath)
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Call(method, req, res)
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

func (lr *LightningRpc) Connect(req *ConnectRequest, res *GetPeersResponse) error {
	client, err := jsonrpc.Dial("unix", lr.socketPath)
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.Call("getpeers", Empty{}, res)
	return err
}

type Route struct {
}

type ConnectRequest struct {
	Host         string `json:"host"`
	Port         uint   `json:"port"`
	FundingTxHex string `json:"tx"`
	Async        bool   `json:"async"`
}

type ConnectResponse struct {
}

func NewLightningRpc(socketPath string) *LightningRpc {
	return &LightningRpc{
		socketPath: socketPath,
	}
}
