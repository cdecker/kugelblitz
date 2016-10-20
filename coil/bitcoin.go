package coil

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Bitcoin struct {
	network string
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

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

func (b *Bitcoin) call(method string, params []interface{}, res interface{}) error {
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

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := http.Post("http://rpcuser:rpcpass@localhost:18332",
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

func (b *Bitcoin) GetInfo(_ *Empty, response *GetBInfoResponse) error {
	return b.call("getinfo", nil, response)
}

func NewBitcoinRpc() *Bitcoin {
	return &Bitcoin{
		network: "-testnet",
	}
}

func (bc *Bitcoin) exec(method string, args []string) (string, error) {
	a := append([]string{bc.network, method}, args...)
	c := exec.Command("/usr/local/bin/bitcoin-cli", a...)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error calling bitcoin rpc: %s", err)
	}
	return strings.TrimSpace(string(out[:])), nil
}

/*
func (bc *Bitcoin) ConnectPeer(rpc *LightningRpc, host string, port uint) error {
	log.Debug("Connecting to %s:%d", host, port)
	addr, err := rpc.NewAddress()
	fmt.Println(addr)
	if err != nil {
		return err
	}

	c := exec.Command("/usr/local/bin/bitcoin-cli", "-testnet", "sendtoaddress", addr, "0.01")
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error sending funds to P2SH: %s", err)
	}
	outs := strings.TrimSpace(string(out[:]))
	c = exec.Command("/usr/local/bin/bitcoin-cli", "-testnet", "getrawtransaction", outs)
	out, err = c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error getting raw transaction: %s", err)
	}
	outs = strings.TrimSpace(string(out[:]))
	fmt.Printf("%s, %v\n", out, err)

	err = rpc.Connect(host, port, outs)
	if err != nil {
		return fmt.Errorf("Error connecting: %s", err)
	}
	return nil
}
*/
