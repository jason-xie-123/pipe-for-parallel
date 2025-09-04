//go:build darwin
// +build darwin

package main

import (
	"context"
	"fmt"
	"net"
	"strings"
)

func dialPipeContext(ipcDialCtx context.Context, pipeName string) (conn net.Conn, err error) {
	socketPath := fmt.Sprintf("/tmp/%s.sock", pipeName)
	d := net.Dialer{}
	conn, err = d.DialContext(ipcDialCtx, "unix", socketPath)
	return
}

func isServerNotExistsErr(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if strings.Contains(opErr.Err.Error(), "connection refused") {
			return true
		} else if strings.Contains(opErr.Err.Error(), "no such file or directory") {
			return true
		}
	}

	return false
}
