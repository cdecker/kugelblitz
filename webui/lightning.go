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
	RPC  *lr.LightningRpc
}

func NewLightning(lrpc *lr.LightningRpc) Lightning {
	return Lightning{
		lrpc: lrpc,
		RPC:  lrpc,
	}
}

func (l *Lightning) IsAlive(_ *lr.Empty) bool {
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
	*res = response
	return err
}

func (l *Lightning) AddFunds(rawtx string) error {
	return l.lrpc.AddFunds(rawtx)
}
func (l *Lightning) FundChannel(nodeid string, capacity uint64) error {
	return l.lrpc.FundChannel(nodeid, capacity)
}

func (l *Lightning) Connect(req *lr.ConnectRequest, _ *lr.Empty) error {
	return l.lrpc.Connect(req.NodeId, req.Host, req.Port)
}

type PaymentRequestInfoRequest struct {
	Destination string `json:"destination"`
}

type PaymentRequestInfoResponse struct {
	Hops        []lr.RouteHop `json:"route"`
	PaymentHash string        `json:"paymenthash"`
	Amount      uint64        `json:"amount"`
}

func (l *Lightning) GetPaymentRequestInfo(req *PaymentRequestInfoRequest, res *PaymentRequestInfoResponse) error {
	var route lr.Route
	fmt.Println(req.Destination)

	req2 := lr.DecodePayRequest{
		PayRequest: req.Destination,
	}

	res2 := lr.DecodePayResponse{}
	err := l.lrpc.DecodePay(&req2, &res2)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", res2)

	routeReq := &lr.GetRouteRequest{
		Destination: res2.Payee,
		Amount:      res2.Amount,
		RiskFactor:  1,
	}
	fmt.Printf("%#v\n", routeReq)
	res.Amount = res2.Amount
	res.PaymentHash = res2.PaymentKey
	err = l.GetRoute(routeReq, &route)
	if err != nil {
		return err
	}
	res.Hops = route.Hops
	return nil
}
