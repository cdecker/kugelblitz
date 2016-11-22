package webui

import (
	"fmt"

	lr "github.com/cdecker/kugelblitz/lightningrpc"
)

// This file takes care of handing incoming JSON-RPC requests over to
// the matching JSON-RPC call in lightning. These are just the shims
// that are exposed over the jsonrpc2

type Lightning struct {
	lrpc *lr.LightningRpc
}

func NewLightning(lrpc *lr.LightningRpc) Lightning {
	return Lightning{
		lrpc: lrpc,
	}
}

func (l *Lightning) IsAlive() bool {
	_, err := l.lrpc.GetInfo()
	return err == nil
}

func (l *Lightning) Close(req *lr.PeerReference, res *lr.Empty) error {
	return l.lrpc.Close(req.PeerId)
}

func (l *Lightning) GetInfo(_ *lr.Empty, res *lr.GetInfoResponse) error {
	info, err := l.lrpc.GetInfo()
	*res = info
	return err
}

func (l *Lightning) GetPeers(_ *lr.Empty, res *lr.GetPeersResponse) error {
	peers, err := l.lrpc.GetPeers()
	*res = peers
	return err
}

func (l *Lightning) GetRoute(req *lr.GetRouteRequest, res *lr.Route) error {
	route, err := l.lrpc.GetRoute(req.Destination, req.Amount, req.RiskFactor)
	*res = route
	return err
}

func (l *Lightning) NewAddress(_ *lr.Empty, res *lr.NewAddressResponse) error {
	addr, err := l.lrpc.NewAddress()
	*res = addr
	return err
}

func (l *Lightning) SendPayment(req *lr.SendPaymentRequest, res *lr.SendPaymentResponse) error {
	response, err := l.lrpc.SendPayment(req.Route, req.PaymentHash)
	fmt.Printf("%#v\n", response)
	*res = response
	return err
}

func (l *Lightning) Connect(req *lr.ConnectRequest, _ *lr.Empty) error {
	return l.lrpc.Connect(req.Host, req.Port, req.FundingTxHex)
}
