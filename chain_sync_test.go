package ogmigo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/savaki/ogmigo/ouroboros/chainsync"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_ChainSync(t *testing.T) {
	endpoint := os.Getenv("OGMIOS")
	//endpoint = "ws://100.99.230.19:11337"
	if endpoint == "" {
		t.SkipNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var (
		p       = message.NewPrinter(language.English)
		counter int64
		read    int64
	)

	client := New(WithEndpoint(endpoint))
	var callback ChainSyncFunc = func(ctx context.Context, data []byte) error {
		var response chainsync.Response
		decoder := json.NewDecoder(bytes.NewReader(data)) // use decoder to check for unknown fields
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&response)
		if err != nil {
			fmt.Println(string(data))
			t.Fatalf("got %v; want nil", err)
		}

		read += int64(len(data))
		if v := atomic.AddInt64(&counter, 1); v%1e3 == 0 {
			var blockNo uint64
			if response.Result != nil && response.Result.RollForward != nil {
				if ps, ok := response.Result.RollForward.Tip.PointStruct(); ok {
					blockNo = ps.BlockNo
				}
			}
			log.Printf("read: block=%v, n=%v, read=%v", blockNo, p.Sprintf("%d", v), p.Sprintf("%d", read))
		}

		return nil
	}

	wait, err := client.ChainSync(ctx, echoStore{}, callback)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if err := wait(); err != nil {
		t.Fatalf("got %v; want nil", err)
	}
}

type echoStore struct {
}

func (e echoStore) Save(_ context.Context, p chainsync.Point) error {
	fmt.Print("Save => ")
	return json.NewEncoder(os.Stdout).Encode(p)
}

func (e echoStore) Load(context.Context) (chainsync.Points, error) {
	return nil, nil
}
