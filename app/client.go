package app

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/rpc"
)

type (
	Client struct {
		endpoint  string
		tlsConfig *tls.Config
	}
	APIRequest struct {
		File   string
		Params map[string]string
	}
	APIResponse struct {
		Data string
	}
)

func New(c *Config) (*Client, error) {
	client := new(Client)
	client.endpoint = c.Endpoint
	client.loadCertificates(c)
	return client, nil
}

func (c *Client) Run() error {
	conn, err := tls.Dial("tcp", c.endpoint, c.tlsConfig)
	if err != nil {
		return fmt.Errorf("Failed to connect to remote endpoint: %s -> %s", c.endpoint, err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", c.endpoint)
	rpcClient := rpc.NewClient(conn)

	params := make(map[string]string)

	params["Name"] = "b*"
	req := &APIRequest{
		File:   "Get-Services.ps1",
		Params: params,
	}
	res := APIResponse{}
	fmt.Printf("Sending RPC request %q to %s\n", req, c.endpoint)
	if err := rpcClient.Call("API.Execute", req, &res); err != nil {
		return fmt.Errorf("Failed to call API.Execute for %s -> %s", req.File, err)
	}
	fmt.Printf("API Response: %s\n", res.Data)
	return nil
}

func (c *Client) Close() {
}

func (c *Client) loadCertificates(config *Config) error {
	cert, err := tls.LoadX509KeyPair(config.CrtPath, config.KeyPath)
	if err != nil {
		return fmt.Errorf("Unable to load certificates crt: %s key: %s -> %s", config.CrtPath, config.KeyPath, err)
	}
	if len(cert.Certificate) != 2 {
		return fmt.Errorf("CRT file should contain 2 certificates, Client and CA certificate")
	}
	ca, err := x509.ParseCertificate(cert.Certificate[1])
	if err != nil {
		return fmt.Errorf("Unable to parse CA certificate -> %s", err)
	}
	for i, certificate := range cert.Certificate {
		pCert, err := x509.ParseCertificate(certificate)
		if err != nil {
			fmt.Printf("Failed to parse certificate[%d] -> %s\n", i, err)
			continue
		}
		fmt.Printf("[%d] CN: %s Org: %s Serial: %s\n", i, pCert.Subject.CommonName, pCert.Subject.Organization[0], pCert.Subject.SerialNumber)
	}
	pool := x509.NewCertPool()
	pool.AddCert(ca)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ClientCAs:    pool,
	}
	c.tlsConfig = tlsConfig
	return nil
}
