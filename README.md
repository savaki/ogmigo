ogmigo
-------------------------

`ogmigo` is a go client for [ogmios](https://ogmios.dev).  

This library is under heavy development, use at your own risk.

### Example

```go
package example

import (
	"context"
	"github.com/savaki/ogmigo"
)

func example(ctx context.Context) error {
	client, err := ogmigo.NewChainSyncClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	response, err := client.ReadNext(ctx) // data is an instance of json.RawMessage
	if err != nil {
		return err
	}
	
	// do something with response
	
	return nil
}
```

### Submodules

`ogmigo` imports `ogmios` as a submodule for testing purposes.  To fetch the submodules,

```bash
git submodule update --init --recursive
```