package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func startServer(ctx context.Context, pipeName string) (err error) {
	var listener net.Listener

	listener, err = createPipeServer(pipeName)
	if err != nil {
		return
	}
	defer listener.Close()

	packChan := make(chan *Package, 20)

	go func() {
		for {
			select {
			case <-ctx.Done():
				listener.Close()
				return
			case pack := <-packChan:
				switch pack.Action {
				case "write":
					fmt.Println(pack.Message)
				case "exit":
					listener.Close()
				default:
					fmt.Printf("unknown action: %s\n", pack.Action)
				}
			}
		}
	}()

	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				err = nil
				return
			}
			return
		} else {
			go func(conn net.Conn) (err error) {
				defer func() {
					conn.Close()
				}()

				for {
					var pack *Package
					pack, err = tryReadPack(conn)
					if err != nil {
						return
					}

					if pack != nil {
						packChan <- pack

						err = sendResponsePack(conn, pack)
						if err != nil {
							err = fmt.Errorf("failed to send response pack: %w", err)
							return
						}
						return
					}
				}
			}(conn)
		}
	}
}

func tryReadPack(conn net.Conn) (pack *Package, err error) {
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	var packLength uint32
	err = binary.Read(conn, binary.BigEndian, &packLength)
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

	if packLength > 0 {
		packData := make([]byte, packLength)
		for {
			conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			var n int
			n, err = io.ReadFull(conn, packData)
			if err != nil {
				if !os.IsTimeout(err) {
					if err == io.EOF {
						err = nil
						continue
					}
					return
				}
				continue
			}
			if n != int(packLength) {
				err = fmt.Errorf("short read from pipe: got %d, want %d", n, packLength)
				return
			}

			break
		}

		var packObj Package
		err = json.Unmarshal(packData, &packObj)
		if err != nil {
			return
		}

		pack = &packObj
		return
	} else {
		err = fmt.Errorf("invalid package length: %d", packLength)
		return
	}
}

func sendResponsePack(conn net.Conn, pack *Package) (err error) {
	var responsePack ResponsePackage = ResponsePackage{
		UUID: pack.UUID,
	}
	var responsePackData []byte
	responsePackData, err = json.Marshal(responsePack)
	if err != nil {
		return
	}

	responsePackDataWithLength := make([]byte, 4+len(responsePackData))
	binary.BigEndian.PutUint32(responsePackDataWithLength[:4], uint32(len(responsePackData)))
	copy(responsePackDataWithLength[4:], responsePackData)

	var writeLen int
	writeLen, err = conn.Write(responsePackDataWithLength)
	if err != nil {
		return
	}

	if writeLen != len(responsePackDataWithLength) {
		err = fmt.Errorf("short write: wrote %d bytes, expected %d bytes", writeLen, len(responsePackDataWithLength))
		return
	}

	return
}
