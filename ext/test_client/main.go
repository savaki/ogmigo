package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/SundaeSwap-finance/ogmigo/v6"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
)

func main() {
	var callback ogmigo.ChainSyncFunc = func(ctx context.Context, data []byte) error {
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, data, "", "  "); err != nil {
			fmt.Println("Failed Marshal: %v", err)
		}

		// Do whatever here.
		return nil
	}

	ctx := context.Background()

	ogmios_addr := "ws://localhost:1339"
	my_client := ogmigo.New(ogmigo.WithEndpoint(ogmios_addr))
	closer, err := my_client.ChainSync(ctx, callback)
	if err != nil {
		fmt.Println("Failed ChainSync: %v", err)
		return
	}

	point, err := my_client.ChainTip(ctx)
	if err != nil {
		fmt.Println("Failed ChainTip: %v", err)
		return
	}
	println("GOT THE TIP - ", point.String())

	epoch, err := my_client.CurrentEpoch(ctx)
	if err != nil {
		fmt.Println("Failed CurrentEpoch: %v", err)
		return
	}
	println("GOT THE EPOCH - ", epoch)

	params, err := my_client.CurrentProtocolParameters(ctx)
	if err != nil {
		fmt.Println("Failed CurrentEpoch: %v", err)
		return
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, params, "", "  "); err != nil {
		fmt.Println("Failed Marshal: %v", err)
		return
	}
	println("PARAMS - ", dst.String())

	summaries, err := my_client.EraSummaries(ctx)
	if err != nil {
		fmt.Println("Failed EraSummaries: %v", err)
		return
	}
	println("GOT THE ERA SUMMARIES - ", summaries.Summaries[len(summaries.Summaries)-1].Start.Slot, " - ", summaries.Summaries[len(summaries.Summaries)-1].End.Slot)

	start, err := my_client.EraStart(ctx)
	if err != nil {
		fmt.Println("Failed EraStart: %v", err)
		return
	}
	println("GOT THE ERA START - ", start.Slot.Uint64())

	utxos_addr, err := my_client.UtxosByAddress(ctx, "addr1q92hfcy7rwmegfwmul2zvlnru3wpnwleu30ujvcpx2hh84mvk4uyun72qraqj8e7hpfluuzxa7zhmcmxcyhqu8yuu6aqr4cpl3", "addr1vyjdya9lmyfmysdgej3qc6thwat33lwm3umk8uwnvhy9g3gzvd2jp")
	if err != nil {
		fmt.Println("Failed UtxosByAddress: %v", err)
		return
	}
	println("GOT THE ADDRESSES - ", utxos_addr[1].Address, " - ", utxos_addr[1].Index)

	id1 := chainsync.UtxoTxID{ID: "16fe7982c416714c22af503165eb9a49eaa55575b9c2e9deb2c400ed4592da03"}
	query1 := chainsync.TxInQuery{Transaction: id1, Index: 0}

	utxos_txin, err := my_client.UtxosByTxIn(ctx, query1)
	if err != nil {
		fmt.Println("Failed UtxosByTxIn: %v", err)
		return
	}
	println("GOT THE TXINS - ", utxos_txin[0].Address, " - ", utxos_txin[0].Index)

	signed_tx_cbor := "84a300818258205acac2562442ac231184fa4f0a73b443be17cca9808ad7429859ea8caff8a92f010182a200583900a5a8a7844561173c280175cb94f683bf413f523afec46ff9999b39a246c0f86efe59ed5cefb045586333b4e6c96a1a776d42d7aa28142258011a02faf080a200583900bb0d2cc0d7f7b80d3c0d7a7ac441f3865ffd297613c67f06951eb7fa619e09c3ce68a7e2bd51443bf36f83ba547199f8d6851a92a5b7ee1f011b00000001cb2bb6b"
	err = my_client.SubmitTx(ctx, signed_tx_cbor)
	if err != nil {
		fmt.Println("Failed SubmitTx: ", err)
	} else {
		fmt.Println("Successful TX submission")
	}

	if err := closer.Close(); err != nil {
		fmt.Println("Failed ChainSync Close: %v", err)
	}
}
