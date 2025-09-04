package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
)

func sendDataToServer(ctx context.Context, pipeName string, pack *Package) (err error) {
	pack.UUID = uuid.New().String()

	ipcDialCtx, ipcDialCtxCancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		select {
		case <-ipcDialCtx.Done():
		default:
			if ipcDialCtxCancelFunc != nil {
				ipcDialCtxCancelFunc()
			}
		}
	}()

	var conn net.Conn
	conn, err = dialPipeContext(ipcDialCtx, pipeName)
	if err != nil {
		if isServerNotExistsErr(err) {
			err = nil

			if pack.Action == "exit" {
				return
			} else {
				conn, err = dialPipeWithRetry(ipcDialCtx, pipeName)
				if err != nil {
					return
				}
			}
		} else {
			return
		}
	}
	defer conn.Close()

	{
		err = sendPack(conn, pack)
		if err != nil {
			return
		}
	}

	{
		var responsePack *ResponsePackage

		for {
			responsePack, err = readResponsePack(conn)
			if err != nil {
				return
			}

			if responsePack != nil {
				break
			}

			continue
		}

		if responsePack.UUID != pack.UUID {
			err = fmt.Errorf("response package UUID mismatch: got %s, want %s", responsePack.UUID, pack.UUID)
			return
		}
	}

	return
}

func dialPipeWithRetry(ipcDialCtx context.Context, pipeName string) (conn net.Conn, err error) {
	checkInterval := 100 * time.Millisecond
	timer := time.NewTimer(checkInterval)
	defer timer.Stop()
	timer.Reset(checkInterval)

	for {
		select {
		case <-ipcDialCtx.Done():
			err = fmt.Errorf("timeout waiting for read process to start: %w", ipcDialCtx.Err())
			return
		case <-timer.C:
			conn, err = dialPipeContext(ipcDialCtx, pipeName)
			if err != nil {
				if isServerNotExistsErr(err) {
					// err = errServerIPCPipeNotFound
					err = nil

					timer.Reset(checkInterval)
					continue
				} else {
					return
				}
			} else {
				return
			}
		}
	}
}

func sendPack(conn net.Conn, pack *Package) (err error) {
	var packData []byte
	packData, err = json.Marshal(pack)
	if err != nil {
		return
	}

	packDataWithLength := make([]byte, 4+len(packData))
	binary.BigEndian.PutUint32(packDataWithLength[:4], uint32(len(packData)))
	copy(packDataWithLength[4:], packData)

	var writeLen int
	writeLen, err = conn.Write(packDataWithLength)
	if err != nil {
		return
	}

	if writeLen != len(packDataWithLength) {
		return fmt.Errorf("short write: wrote %d bytes, expected %d bytes", writeLen, len(packDataWithLength))
	}

	return
}

func readResponsePack(conn net.Conn) (responsePack *ResponsePackage, err error) {
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	var responsePackLength uint32
	err = binary.Read(conn, binary.BigEndian, &responsePackLength)
	if err != nil {
		if !os.IsTimeout(err) {
			if err == io.EOF {
				err = nil
				return
			}
			return
		}
		err = nil
		return
	}

	if responsePackLength > 0 {
		responsePackData := make([]byte, responsePackLength)

		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		var n int
		n, err = io.ReadFull(conn, responsePackData)
		if err != nil {
			if !os.IsTimeout(err) {
				if err == io.EOF {
					err = nil
					return
				}
				return
			}

			err = nil
			return
		}
		if n != int(responsePackLength) {
			err = fmt.Errorf("short read from pipe: got %d, want %d", n, responsePackLength)
			return
		}

		var responsePackObj ResponsePackage
		err = json.Unmarshal(responsePackData, &responsePackObj)
		if err != nil {
			return
		}

		responsePack = &responsePackObj
		return
	} else {
		err = fmt.Errorf("invalid response package length: %d", responsePackLength)
		return
	}
}
