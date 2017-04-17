//  Created by Elliott Polk on 17/04/2017
//  Copyright © 2017. All rights reserved.
//  beam/cmd/beam/beam.go
//
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/elliottpolk/beam/client"
	"github.com/elliottpolk/beam/log"
	"github.com/elliottpolk/beam/server"

	"github.com/urfave/cli"
)

const (
	AddrFlag       string = "addr"
	DirFlag        string = "dir"
	FromFlag       string = "from"
	ToFlag         string = "to"
	BlockFlag      string = "block"
	ConcurrentFlag string = "concurrent"

	DefaultAddr  string = ":8888"
	DefaultDir   string = "."
	DefaultFrom  string = ""
	DefaultTo    string = ""
	DefaultBlock int64  = 0
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "beam server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, a", AddrFlag),
					Value: DefaultAddr,
					Usage: "",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, d", DirFlag),
					Value: DefaultDir,
					Usage: "directory to share",
				},
			},
			Action: func(context *cli.Context) {
				context.Command.VisibleFlags()
				log.Fatal(server.ListenAndServe(context.String(AddrFlag), context.String(DirFlag)))
			},
		},
		{
			Name:  "ls",
			Usage: "list out the remote files",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, a", AddrFlag),
					Value: DefaultAddr,
					Usage: "",
				},
			},
			Action: func(context *cli.Context) {
				context.Command.VisibleFlags()

				from := fmt.Sprintf("%s/ls", context.String(AddrFlag))
				from = strings.TrimPrefix(from, "http://")
				from = strings.TrimPrefix(from, "https://")
				from = fmt.Sprintf("http://%s", from)

				res, err := http.Get(from)
				if err != nil {
					log.Fatal(err)
				}
				defer res.Body.Close()

				msg, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
				}

				if code := res.StatusCode; code != http.StatusOK {
					log.Fatalf("server responded with status code %v and message %s", code, string(msg))
				}

				log.Info(string(msg))
			},
		},
		{
			Name:  "get",
			Usage: "retrieve remote file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, a", AddrFlag),
					Value: DefaultAddr,
					Usage: "",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, f", FromFlag),
					Value: DefaultFrom,
					Usage: "remote file to copy",
				},
				cli.StringFlag{
					Name:  fmt.Sprintf("%s, t", ToFlag),
					Value: DefaultTo,
					Usage: "file to write to",
				},
				cli.Int64Flag{
					Name:  fmt.Sprintf("%s, b", BlockFlag),
					Value: DefaultBlock,
					Usage: "block to begin downloading at",
				},
			},
			Action: func(context *cli.Context) {
				context.Command.VisibleFlags()

				addr := context.String(AddrFlag)
				from := context.String(FromFlag)
				to := context.String(ToFlag)
				block := context.Int64(BlockFlag)

				if len(to) < 1 {
					to = from
				}

				if err := client.DialAndDownloadAt(addr, from, to, block); err != nil {
					log.Fatal(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
