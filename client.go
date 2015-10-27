package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func client(config Config) {

	// open send port
	writeSocket := openUDPSocket(`w`, net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255), // (broadcast IPv4)
		Port: config.Port,
	})
	defer writeSocket.Close()

	// open receive port on send port + 1
	_, porttxt, _ := net.SplitHostPort(writeSocket.LocalAddr().String())
	port, _ := strconv.Atoi(porttxt)
	readSocket := openUDPSocket(`r`, net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port + 1,
	})
	defer readSocket.Close()

	// listen for reply: first start 2 channels: dataCh, and errCh
	data, errors := make(chan []byte), make(chan error)
	go func(dataCh chan []byte, errCh chan error) {
		keyData := make([]byte, UDP_MSG_SIZE)
		n, _, err := readSocket.ReadFromUDP(keyData)
		if err != nil {
			errCh <- err
		}
		dataCh <- keyData[0:n]
	}(data, errors)

	// send key request
	fmt.Println("Broadcasting UDP key request...")
	_, err := writeSocket.Write([]byte(KEY_REQUEST_TEXT))
	if err != nil {
		fmt.Println("ERROR sending key request: ", err)
		os.Exit(101)
	}

	// now start the timeout channel
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(config.Wait) * time.Second)
		timeout <- true
	}()

	// now wait for a reply, an error, or a timeout
	reply := ``
	select {
	case bytes := <-data:
		reply = string(bytes)
		if strings.HasPrefix(reply, `ERROR`) == true {
			fmt.Println("Received error from server:", reply)
			os.Exit(102)
		}
		fmt.Println("Got key from server.")
	case err := <-errors:
		fmt.Println("Error reading from receive port:", err)
		os.Exit(103)
	case <-timeout:
		fmt.Println("WARNING: timed out waiting for response from key server.")
	}

	// create key dir and file
	keyWritten := false // keep track of whether the key was written
	if reply != `` {
		fmt.Printf("Writing key to %s\n", config.KeyPath)
		err = os.MkdirAll(filepath.Dir(config.KeyPath), 0700)
		if err != nil {
			fmt.Printf("ERROR creating directory %s: %s\n", config.KeyPath, err)
			os.Exit(104)
		}
		err = ioutil.WriteFile(config.KeyPath, []byte(reply), 0600)
		if err != nil {
			fmt.Printf("ERROR writing keyfile %s: %s\n", config.KeyPath, err)
			os.Exit(105)
		}
		keyWritten = true
	}
	// defer close and deletion of keyfile
	// from here on, set exitCode and call return instead of os.Exit()
	exitCode := 0
	defer func() {
		if keyWritten == true {
			fmt.Printf("Deleting key file %s...\n", config.KeyPath)
			if err := os.Remove(config.KeyPath); err != nil {
				fmt.Printf("ERROR deleting keyfile '%s': %v\n",
					config.KeyPath, err)
				exitCode = 106
				return
			}
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	// run command
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	cmdText := strings.Join(flag.Args(), " ")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running command:", cmdText)
	if err := cmd.Start(); err != nil {
		fmt.Printf("ERROR starting command '%s': %v\n", cmdText, err)
		exitCode = 107
		return
	}

	if err = cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
				fmt.Printf("ERROR: command '%s' exited with status %d\n",
					cmdText, exitCode)
			} else {
				fmt.Printf("ERROR: command '%s' exited with unknown status",
					cmdText)
				exitCode = 108 // problem getting command's exit status?
			}
			return
		} else {
			fmt.Printf("ERROR waiting on command '%s': %v\n", cmdText, err)
			exitCode = 109
			return
		}
	}

	fmt.Println("Command completed successfully.")
}
