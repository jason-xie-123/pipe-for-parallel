//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"net"
	"os"
)

func createPipeServer(pipeName string) (listener net.Listener, err error) {
	pipePath := fmt.Sprintf("/tmp/%s.sock", pipeName)
	if _, err := os.Stat(pipePath); err == nil {
		os.Remove(pipePath)
	}

	listener, err = net.Listen("unix", pipePath)

	return
}
