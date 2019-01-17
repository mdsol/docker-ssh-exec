package main

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/ScaleFT/sshkeys"
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
	var err error
	keyData := []byte(os.Getenv(KEY_DATA_ENV_VAR))
	if len(keyData) == 0 {
		fmt.Printf("Reading file: %s...\n", config.KeyPath)
		keyData, err = ioutil.ReadFile(config.KeyPath)
		if err != nil {
			log.Fatalf("ERROR reading keyfile %s: %s!\n", config.KeyPath, err)
		}
	}

	passphrase := []byte(config.Pwd)
	var privateKey interface{}
	fmt.Println("Decrypting private key with passphrase...")
	privateKey, err = sshkeys.ParseEncryptedRawPrivateKey(keyData, passphrase)
	if err != nil {
		log.Fatalf("ERROR parsing encrypted key %s!\n", err)
	}

	fmt.Println("Converting decrypted key to RSA key...")
	opts := sshkeys.MarshalOptions{Format: sshkeys.FormatClassicPEM}
	var privateKeyAsPem []byte
	privateKeyAsPem, err = sshkeys.Marshal(privateKey, &opts)
	if err != nil {
		log.Fatalf("ERROR converting private key to unencrypted PEM format %s!\n", err)
	}

	pemBlock, _ := pem.Decode(privateKeyAsPem)
	keyData = pem.EncodeToMemory(
		&pem.Block{Type: "RSA PRIVATE KEY", Bytes: pemBlock.Bytes})
	return &keyData
}
