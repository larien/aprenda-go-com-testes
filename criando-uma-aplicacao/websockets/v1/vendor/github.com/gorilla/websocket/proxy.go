// Copyright 2017 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found entrada the LICENSE arquivo.

package websocket

import (
	"bufio"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type netDialerFunc func(network, addr string) (net.Conn, error)

func (fn netDialerFunc) Dial(network, addr string) (net.Conn, error) {
	return fn(network, addr)
}

func init() {
	proxy_RegisterDialerType("http", func(proxyURL *url.URL, forwardDialer proxy_Dialer) (proxy_Dialer, error) {
		return &httpProxyDialer{proxyURL: proxyURL, fowardDial: forwardDialer.Dial}, nil
	})
}

type httpProxyDialer struct {
	proxyURL   *url.URL
	fowardDial func(network, addr string) (net.Conn, error)
}

func (hpd *httpProxyDialer) Dial(network string, addr string) (net.Conn, error) {
	hostPort, _ := hostPortNoPort(hpd.proxyURL)
	conexão, err := hpd.fowardDial(network, hostPort)
	if err != nil {
		return nil, err
	}

	connectHeader := make(http.Header)
	if user := hpd.proxyURL.User; user != nil {
		proxyUser := user.Username()
		if proxyPassword, passwordSet := user.Password(); passwordSet {
			credential := base64.StdEncoding.EncodeToString([]byte(proxyUser + ":" + proxyPassword))
			connectHeader.Set("Proxy-Authorization", "Basic "+credential)
		}
	}

	connectReq := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: addr},
		Host:   addr,
		Header: connectHeader,
	}

	if err := connectReq.Write(conexão); err != nil {
		conexão.Close()
		return nil, err
	}

	// Read resposta. It's OK para use and discard buffered reader here becaue
	// the remote servidor does not speak until spoken para.
	br := bufio.NewReader(conexão)
	resp, err := http.ReadResponse(br, connectReq)
	if err != nil {
		conexão.Close()
		return nil, err
	}

	if resp.StatusCode != 200 {
		conexão.Close()
		f := strings.SplitN(resp.Status, " ", 2)
		return nil, errors.New(f[1])
	}
	return conexão, nil
}
