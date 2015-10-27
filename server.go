package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func server(config Config) {

	// open receive port
	readSocket := openUDPSocket(`r`, net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: config.Port,
	})
	fmt.Printf("Listening on UDP port %d...\n", config.Port)
	defer readSocket.Close()

	// main loop
	for {
		data := make([]byte, UDP_MSG_SIZE)
		size, clientAddr, err := readSocket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("Error reading from receive port: ", err)
		}
		clientMsg := data[0:size]
		if string(clientMsg) == KEY_REQUEST_TEXT {
			fmt.Printf("Received key request from %s, sending key.\n",
				clientAddr.IP)
			// reply to the client on the same port
			writeSocket := openUDPSocket(`w`, net.UDPAddr{
				IP:   clientAddr.IP,
				Port: clientAddr.Port + 1,
			})

			response := os.Getenv(KEY_DATA_ENV_VAR)
			if response == `` {
				keyData, err := ioutil.ReadFile(config.KeyPath)
				if err == nil {
					response = string(keyData)
				} else {
					response = fmt.Sprintf("ERROR reading keyfile %s: %s!",
						config.KeyPath, err)
					fmt.Println(response)
				}
			}
			_, err = writeSocket.Write([]byte(response))
			if err != nil {
				fmt.Printf("ERROR writing data to socket:%s!\n", err)
			}
			writeSocket.Close()
		}
	}
}
