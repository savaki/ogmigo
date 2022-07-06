ogmigo
-------------------------

`ogmigo` is a go client for [ogmios](https://ogmios.dev).

This library is under heavy development, use at your own risk.

### Example

```go
package example

import (
	"context"
	"github.com/SundaeSwap-finance/ogmigo"
)

func example(ctx context.Context) error {
	var callback ogmigo.ChainSyncFunc = func(ctx context.Context, data []byte) error {
		// do work
		return nil
	}

	client := ogmigo.New(ogmigo.WithEndpoint("ws://example.com:1337"))
	closer, err := client.ChainSync(ctx, callback)
	if err != nil {
		return err
	}
	if err := closer.Close(); err != nil {
		return err
	}

	return nil
}
```

### Submodules

`ogmigo` imports `ogmios` as a submodule for testing purposes. To fetch the submodules,

```bash
git submodule update --init --recursive
```


