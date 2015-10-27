package main

import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_RECV_PORT = 1067
	CLIENT_TIMEOUT   = 3 // seconds
	KEY_REQUEST_TEXT = `Key, plz!`
	UDP_MSG_SIZE     = 4096 // the effective max key size
	KEY_DATA_ENV_VAR = `DOCKER-SSH-KEY`
)

// Package version & timestamp - interpolated by goxc
const VERSION = "0.5.1"
const SOURCE_DATE = "2015-10-26T06:23:34-07:00"

func main() {
	config := newConfig()
	if config.Server {
		server(config)
	} else {
		client(config)
	}
}

// Opens a UDP read socket at UDPAddr, and returns the connection object.
// mode should be "r" or "w"
// Exits and prints the error if one occurs.
func openUDPSocket(mode string, addr net.UDPAddr) (socket *net.UDPConn) {
	var err error
	if mode == `w` {
		socket, err = net.DialUDP("udp4", nil, &addr)
	} else {
		socket, err = net.ListenUDP("udp4", &addr)
	}
	if err != nil {
		fmt.Println("Error opening receive port: ", err)
		os.Exit(0)
	}
	return socket
}
