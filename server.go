package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
)

func server(config Config) {

	// open receive port
	readSocket := openUDPSocket(`r`, net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: config.UDPPort,
	})
	keyData := readKeyData(&config)
	go serveHTTP(keyData, config.HTTPPort)
	fmt.Printf("Listening on UDP port %d...\n", config.UDPPort)
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
			_, err = writeSocket.Write(*keyData)
			if err != nil {
				fmt.Printf("ERROR writing data to socket:%s!\n", err)
			}
			writeSocket.Close()
		}
	}
}

func readKeyData(config *Config) *[]byte {
	// var keyData []byte
	var err error
	keyData := []byte(os.Getenv(KEY_DATA_ENV_VAR))
	if len(keyData) == 0 {
		fmt.Printf("Reading file: %s...\n", config.KeyPath)
		keyData, err = ioutil.ReadFile(config.KeyPath)
		if err != nil {
			log.Fatalf("ERROR reading keyfile %s: %s!\n", config.KeyPath, err)
		}
	}
	pemBlock, _ := pem.Decode(keyData)
	if pemBlock != nil {
		if x509.IsEncryptedPEMBlock(pemBlock) {
			fmt.Println("Decrypting private key with passphrase...")
			decoded, err := x509.DecryptPEMBlock(pemBlock, []byte(config.Pwd))
			if err == nil {
				header := `PRIVATE KEY` // default key type in header
				matcher := regexp.MustCompile("-----BEGIN (.*)-----")
				if matches := matcher.FindSubmatch(keyData); len(matches) > 1 {
					header = string(matches[1])
				}
				keyData = pem.EncodeToMemory(
					&pem.Block{Type: header, Bytes: decoded})
			} else {
				fmt.Printf("Error decrypting PEM-encoded secret: %s\n", err)
			}
		}
	}
	return &keyData
}
