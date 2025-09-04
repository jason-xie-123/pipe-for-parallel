//go:build windows
// +build windows

package main

import (
	"errors"
	"fmt"
	"net"
	"syscall"

	"github.com/Microsoft/go-winio"
)

func createPipeServer(pipeName string) (listener net.Listener, err error) {
	cfg := &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GA;;;BU)",
		InputBufferSize:    4096,
		OutputBufferSize:   4096,
	}

	pipePath := fmt.Sprintf(`\\.\pipe\%s`, pipeName)

	listener, err = winio.ListenPipe(pipePath, cfg)
	if err != nil {
		if errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
			err = fmt.Errorf("failed to listen on pipe %s: access denied, please make sure no other process is using the same named pipe and you have permission to access it: %w", pipePath, err)
		}
		return
	}

	return
}
