package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/SundaeSwap-finance/ogmigo/v6"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/compatibility"
)

func main() {
	var callback ogmigo.ChainSyncFunc = func(ctx context.Context, data []byte) error {

		// Quick-and-dirty way to distinguish b/w 2 different responses.
		var response compatibility.CompatibleResponsePraos
		if err := json.Unmarshal(data, &response); err != nil {
			fmt.Printf("Failed Unmarshal: %v\n", err)
			return nil
		}

		switch response.Method {
		case chainsync.FindIntersectionMethod, chainsync.FindIntersectMethod:
			fmt.Printf("CompatibleResultFindIntersection - %v\n", response.Result)
		case chainsync.NextBlockMethod, chainsync.RequestNextMethod:
			fmt.Printf("CompatibleResultNextBlock - %v\n", response.Result)
		default:
			fmt.Printf("Unsupported method: %v\n", response.Method)
		}

		return nil
	}

	debugPtr := flag.Bool("debug", false, "Print debug statements")
	flag.Parse()
	ctx := context.Background()

	ogmios_addr := "ws://100.72.33.61:1339"
	my_client := ogmigo.New(ogmigo.WithEndpoint(ogmios_addr))
	closer, err := my_client.ChainSync(ctx, callback)
	if err != nil {
		fmt.Printf("Failed ChainSync Open: %v\n", err)
		return
	}

	point, err := my_client.ChainTip(ctx)
	if err != nil {
		fmt.Printf("Failed ChainTip: %v\n", err)
		return
	}
	if *debugPtr {
		fmt.Printf("GOT THE TIP - %v\n", point.String())
	}

	epoch, err := my_client.CurrentEpoch(ctx)
	if err != nil {
		fmt.Printf("Failed CurrentEpoch: %v\n", err)
		return
	}
	if *debugPtr {
		fmt.Printf("GOT THE EPOCH - %v\n", epoch)
	}

	params, err := my_client.CurrentProtocolParameters(ctx)
	if err != nil {
		fmt.Printf("Failed CurrentEpoch: %v\n", err)
		return
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, params, "", "  "); err != nil {
		fmt.Printf("Failed Marshal: %v\n", err)
		return
	}
	if *debugPtr {
		fmt.Printf("PARAMS - %v\n", dst.String())
	}

	summaries, err := my_client.EraSummaries(ctx)
	if err != nil {
		fmt.Printf("Failed EraSummaries: %v\n", err)
		return
	}
	if *debugPtr {
		fmt.Printf("GOT THE ERA SUMMARIES - %v - %v\n", summaries.Summaries[len(summaries.Summaries)-1].Start.Slot, summaries.Summaries[len(summaries.Summaries)-1].End.Slot)
	}

	start, err := my_client.EraStart(ctx)
	if err != nil {
		fmt.Printf("Failed EraStart: %v\n", err)
		return
	}
	if *debugPtr {
		fmt.Printf("GOT THE ERA START - %v\n", start.Slot.Uint64())
	}

	// Caveat emptor. Due to internal changes, UTXO queries are supported on a
	// best-effort basis, and are very slow even when they do work.
	// utxos_addr, err := my_client.UtxosByAddress(ctx, "addr1v8ua3ne8pp050eyfyhavazzlkdh2e38urw38jlnl55nkwccuzg2m5", "addr1qyaj05kuw0c4amuqgyxr6arpgjnns65fcc4af6kqwt3wznfmylfdcul3tmhcqsgv846xz3988p4gn33t6n4vquhzu9xswpk9uz")
	// if err != nil {
	// 	fmt.Printf("Failed UtxosByAddress: %v", err)
	// 	return
	// }
	// if *debugPtr {
	// 	fmt.Printf("GOT THE ADDRESSES - ", utxos_addr[1].Address, " - ", utxos_addr[1].Index)
	// }

	// id1 := chainsync.UtxoTxID{ID: "16fe7982c416714c22af503165eb9a49eaa55575b9c2e9deb2c400ed4592da03"}
	// query1 := chainsync.TxInQuery{Transaction: id1, Index: 0}

	// utxos_txin, err := my_client.UtxosByTxIn(ctx, query1)
	// if err != nil {
	// 	fmt.Printf("Failed UtxosByTxIn: %v", err)
	// 	return
	// }
	// if *debugPtr {
	// 	fmt.Printf("GOT THE TXINS - ", utxos_txin[0].Address, " - ", utxos_txin[0].Index)
	// }

	// signed_tx_cbor := "84a300818258205acac2562442ac231184fa4f0a73b443be17cca9808ad7429859ea8caff8a92f010182a200583900a5a8a7844561173c280175cb94f683bf413f523afec46ff9999b39a246c0f86efe59ed5cefb045586333b4e6c96a1a776d42d7aa28142258011a02faf080a200583900bb0d2cc0d7f7b80d3c0d7a7ac441f3865ffd297613c67f06951eb7fa619e09c3ce68a7e2bd51443bf36f83ba547199f8d6851a92a5b7ee1f011b00000001cb2bb6b"
	// err = my_client.SubmitTx(ctx, signed_tx_cbor)
	// if err != nil {
	// 	fmt.Printf("Failed SubmitTx: ", err)
	// } else if *debugPtr {
	// 	fmt.Printf("Successful TX submission - ", signed_tx_cbor)
	// }
	time.Sleep(2 * time.Second)

	if err := closer.Close(); err != nil {
		fmt.Printf("Failed ChainSync Close: %v\n", err)
	}
}
