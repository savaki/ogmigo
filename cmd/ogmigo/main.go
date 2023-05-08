// Copyright 2021 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"github.com/thuannguyen2010/ogmigo"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
)

var opts struct {
	DB     string
	Ogmios string
	Points cli.StringSlice
	Tick   int64
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "ogmios",
			Usage:       "ogmios websocket endpoint",
			Value:       "ws://3.21.245.207:1337", // mainnet2
			EnvVars:     []string{"OGMIOS"},
			Destination: &opts.Ogmios,
		},
		&cli.StringSliceFlag{
			Name:        "point",
			Aliases:     []string{"p"},
			Usage:       "initial starting point in the form {slot}/{hash} e.g.",
			EnvVars:     []string{"POINT"},
			Destination: &opts.Points,
		},
		&cli.Int64Flag{
			Name:        "tick",
			Usage:       "display progress every tick slots",
			EnvVars:     []string{"TICK"},
			Value:       5e3,
			Destination: &opts.Tick,
		},
	}
	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func action(_ *cli.Context) error {
	client := ogmigo.New(
		ogmigo.WithEndpoint("ws://localhost:1337"),
		ogmigo.WithLogger(ogmigo.DefaultLogger),
	)

	ctx := context.Background()
	redeemer, err := client.EvaluateTx(ctx, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(redeemer)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Kill, os.Interrupt)

	<-stop

	return nil
}
