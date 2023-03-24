package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// URL to fetch
const (
	fsize                      = 1 << 62
	check_url                  = "https://syscall.network/ip"
	default_environment string = "/etc/openai"
	api_key_filename    string = ".api"
	org_key_filename    string = ".org"
	models_url          string = "https://api.openai.com/v1/models"
	completions_url     string = "https://api.openai.com/v1/completions"
)

type (
	ai_client struct {
		api_key, org_id string
		use_tor         bool
		client          *http.Client
		output_file     *os.File
		args            *cmd_args
	}
	proxy func(*http.Request) (*url.URL, error)
)

func (c *ai_client) available_models() (int64, error) {
	req := fmt_req("GET", models_url)
	resp, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}
	return io.Copy(c.output_file, resp.Body)
}

func default_client(args *cmd_args, use_tor bool, api_key string, output *os.File) ai_client {
	c := ai_client{
		api_key:     api_key,
		client:      http.DefaultClient,
		output_file: output,
		args:        args,
	}
	if use_tor {
		c.client = tor_proxy()
	}
	return c
}

func fmt_req_with(method, endpoint string, body io.Reader) *http.Request {
	req, err := http.NewRequestWithContext(context.Background(), method, endpoint, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "ok-john")
	req.Header.Set("Authorization", "Bearer "+args.api_key)
	return req
}

func fmt_req(method, endpoint string) *http.Request {
	return fmt_req_with(method, endpoint, nil)
}

func tor_proxy() *http.Client {
	var torProxy string = "socks5://127.0.0.1:9050" // 9150 w/ Tor Browser
	torProxyUrl, err := url.Parse(torProxy)
	if err != nil {
		panic(err)
	}
	torTransport := &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify:     false,
		SessionTicketsDisabled: true,
		MinVersion:             0,
		Renegotiation:          tls.RenegotiateOnceAsClient,
	},
		GetProxyConnectHeader: intercept_header, Proxy: http.ProxyURL(torProxyUrl),
		MaxResponseHeaderBytes: 1 << 62}
	return &http.Client{Transport: torTransport, Timeout: time.Second * 480}
}

func intercept_header(ctx context.Context, proxyURL *url.URL, target string) (http.Header, error) {
	headers := http.Header{}
	headers.Add("iat", strconv.Itoa(int(time.Now().Unix())))
	return headers, nil
}

func (gc *ai_client) check_ip() {
	resp, err := gc.client.Do(fmt_req("GET", check_url))
	if err != nil {
		log.Fatal("Error making GET request.", err)
	}
	resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)
}
