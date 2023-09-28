package ogmigo

import (
	"context"
	"testing"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
)

func TestNewLoggingStore(t *testing.T) {
	p := chainsync.PointStruct{
		BlockNo: 123,
		Hash:    "hash",
		Slot:    456,
	}

	ctx := context.Background()
	store := NewLoggingStore(DefaultLogger)
	_ = store.Save(ctx, p.Point())

	pp, err := store.Load(ctx)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if got, want := len(pp), 0; got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}
