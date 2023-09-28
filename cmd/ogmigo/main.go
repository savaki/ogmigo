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
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync/atomic"

	"github.com/SundaeSwap-finance/ogmigo"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/urfave/cli/v2"
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
		ogmigo.WithEndpoint(opts.Ogmios),
		ogmigo.WithLogger(ogmigo.DefaultLogger),
	)

	var (
		ctx    = context.Background()
		re     = regexp.MustCompile(`^(\d+)/([a-zA-Z0-9]+)$`)
		points chainsync.Points
	)

	for _, s := range opts.Points.Value() {
		match := re.FindStringSubmatch(s)
		if len(match) != 3 {
			return fmt.Errorf("ogmigo: failed to parse point, %v", s)
		}
		slot, _ := strconv.ParseUint(match[1], 10, 64)
		points = append(points, chainsync.PointStruct{
			Hash: match[2],
			Slot: slot,
		}.Point())
	}

	var counter int64
	var callback ogmigo.ChainSyncFunc = func(ctx context.Context, data []byte) error {
		if v := atomic.AddInt64(&counter, 1); v%opts.Tick != 0 {
			return nil
		}

		var response chainsync.Response
		if err := json.Unmarshal(data, &response); err != nil {
			return err
		}
		if response.Result == nil {
			return nil
		}
		if response.Result.RollForward == nil {
			return nil
		}

		ps := response.Result.RollForward.Block.PointStruct()
		fmt.Printf("slot=%v hash=%v block=%v\n", ps.Slot, ps.Hash, ps.BlockNo)

		return nil
	}
	closer, err := client.ChainSync(ctx, callback,
		ogmigo.WithPoints(points...),
		ogmigo.WithReconnect(true),
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
