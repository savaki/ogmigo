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
	"encoding/json"
	"fmt"
	"github.com/thuannguyen2010/ogmigo"
	"github.com/thuannguyen2010/ogmigo/ouroboros/chainsync"
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
		ogmigo.WithEndpoint("ws://172.0.0.1:1337"),
		ogmigo.WithLogger(ogmigo.DefaultLogger),
	)

	var (
		ctx    = context.Background()
		points chainsync.Points
	)
	points = []chainsync.Point{
		chainsync.PointStruct{
			BlockNo: 1009189,
			Hash:    "b95423b8778536cf9a53e8855cef2bd8702f79ffa3c918b8857437cfb238474d",
			Slot:    23003752,
		}.Point(),
	}
	//useV6 := false
	useV6 := true
	var callback ogmigo.ChainSyncFunc = func(ctx context.Context, data []byte) error {
		var response chainsync.Response
		if useV6 {
			var resV6 chainsync.ResponseV6
			if err := json.Unmarshal(data, &resV6); err != nil {
				return err
			}
			response = resV6.ConvertToV5()
		} else {
			if err := json.Unmarshal(data, &response); err != nil {
				return err
			}
		}

		if response.Result == nil {
			return nil
		}
		if response.Result.RollForward != nil {
			ps := response.Result.RollForward.Block.PointStruct()
			fmt.Println("blockNo", ps.BlockNo, "hash", ps.Hash, "slot", ps.Slot)
		}
		if response.Result.RollBackward != nil {
			return nil
		}

		//ps := response.Result.RollForward.Block.PointStruct()
		//fmt.Printf("slot=%v hash=%v block=%v\n", ps.Slot, ps.Hash, ps.BlockNo)

		return nil
	}
	closer, err := client.ChainSync(ctx, callback,
		ogmigo.WithPoints(points...),
		ogmigo.WithReconnect(true),
		ogmigo.WithUseV6(useV6),
	)

	if err != nil {
		return err
	}
	defer closer.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Kill, os.Interrupt)

	<-stop

	return nil
}
