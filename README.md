ogmigo
-------------------------

`ogmigo` is a go client for [ogmios](https://ogmios.dev).  

This library is under heavy development, use at your own risk.

### Example

```go
func example() {
  client, err := New(ctx)
  if err != nil { return err }
  defer client.Close()
  
  data, err := client.ReadNext(ctx) // data is an instance of json.RawMessage
  if err != nil { return err }

  var response chainsync.Response
  err = json.Unmarshal(data, &response)
}
```
