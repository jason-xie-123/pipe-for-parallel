package main

import (
	"context"
	"fmt"
	"os"
	packageVersion "pipe-for-parallel/version"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "pipe-for-parallel",
		Usage:   "CLI tool to pipe-for-parallel scripts",
		Version: packageVersion.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "action",
				Usage: "support action: read / write / exit",
			},
			&cli.StringFlag{
				Name:  "pipe",
				Usage: "pipe name",
			},
			&cli.StringFlag{
				Name:  "message",
				Usage: "the write message",
			},
		},
		Action: func(c *cli.Context) (err error) {
			action := c.String("action")
			if len(action) == 0 {
				return fmt.Errorf("invalid value for --action: %s. Valid options are 'read', 'write', 'exit'", action)
			}

			pipeName := c.String("pipe")
			if len(pipeName) == 0 {
				return fmt.Errorf("invalid value for --pipe: %s", pipeName)
			}

			message := c.String("message")

			ctx, cancel := context.WithCancel(context.Background())
			defer func() {
				select {
				case <-ctx.Done():
					return
				default:
					cancel()
				}
			}()

			switch action {
			case "read":
				err = startServer(ctx, pipeName)
				return
			case "write", "exit":
				err = sendDataToServer(ctx, pipeName, &Package{
					Action:  action,
					Message: message,
				})
				return
			default:
				err = fmt.Errorf("invalid value for --action: %s. Valid options are 'read', 'write', 'exit'", action)
				return
			}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
