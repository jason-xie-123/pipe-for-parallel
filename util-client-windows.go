//go:build windows
// +build windows

package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

func dialPipeContext(ipcDialCtx context.Context, pipeName string) (conn net.Conn, err error) {
	pipePath := fmt.Sprintf(`\\.\pipe\%s`, pipeName)
	conn, err = winio.DialPipeContext(ipcDialCtx, pipePath)
	return
}

func isServerNotExistsErr(err error) bool {
	return os.IsNotExist(err)
}
