
# Ogmigo v6 migration guide

This document discusses how to upgrade code that uses Ogmigo that is connected to Ogmios v5 instances and must connect to Ogmios v6. While not completely pain-free, it is possible to safely migrate with a minimum of pain.

Note that this document is a supplement to [the original Ogmios migration guide](https://github.com/CardanoSolutions/ogmios/blob/master/architectural-decisions/accepted/017-api-version-6-major-rewrite.md). It’s highly recommended that you read that particular guide before reading this guide, especially if your code handles a lot of individual values. There are also Ogmios API schema pages for [v5](https://ogmios.dev/api/v5.6/) and [v6](https://ogmios.dev/api/) that are worth consulting if all else fails.

# Caveats

There are a handful of caveats that should be considered before reading the main document.

* The Ogmigo v6 upgrade adds no new major functionality beyond v6 struct support and a compatibility layer that allows both v5 and v6 Ogmigo JSON/DB/at-rest data to be unmarshalled into v6 structs. If anything wasn’t supported in Ogmigo v5, it won’t be supported in Ogmigo v6.
* Support for the Byron era is limited in Ogmigo v6, and may not work as expected. In fact, as of this writing (Nov. 2023), critical functionality, such as the _compatibility_ module, assumes Byron _isn’t_ supported. If you must have Byron support, further Ogmigo work will be required.
* Any attempt to use a v6-enabled Ogmigo library will automatically break the code using Ogmigo. Many fundamental structs (e.g., _Block_) have been altered to assume v6 structs. If somebody absolutely must use v5 structs, they’ll have to import the _v5_ module and use those structs. Even then, it is highly recommended to use the _compatibility_ module whenever possible, thereby assisting in a smooth transition to the default (v6) code.

# Support for v6

As mentioned in the caveats, the v6 structs are now the default in the Ogmigo library. Code will have to be written with this in mind if working directly with v6. The Ogmios migration guide discusses the gist of what changed between v5 and v6. In general, other than changing specific references within the structs, there aren’t major changes required to understand the new structs.

In addition, Ogmigo automatically handles the switch from JSON-WSP to JSON-RPC 2.0. You should just be able to let Ogmigo handle everything automatically, simply pointing it to the appropriate Ogmios endpoint.

## Major Struct Changes

While accurate overall, the Ogmios architectural decision doc linked above has minor differences compared to what's actually in the code and structs. If you’re using v5 structs and want to see where to look in v6 structs for specific data, you should look at the _compatibility_ module and see how v5 is translated to v6.

That said, there are a few general caveats to note.

* A couple of more esoteric v6 struct fields can’t be populated by v5, at least not without adding a lot of code for little gain (e.g., consensus proofs). If there are gaps, get in touch but you may have to work out a solution if you must have any missing data.
* Some data in v5 is encoded in Base64 and switched in v6 to Base16. The _compatibility_ module automatically handles conversion. Still, if your code directly uses these entries, you’ll need to adjust as necessary.
* `Value` JSON types don’t always carry over between v5 and v6, and aren’t always clearly defined by the API docs. We’ve looked at both real-world data and Ogmios unit tests, and have done the best we can to determine what’s allowed. It is still recommended that you double check any values critical to your functionality. If there are discrepancies, please get in touch.

### Creating Value Instances

When creating a new `Value`, there are some recommended methods to do so. If creating a `Value` from scratch, using the `ValueFromCoins()` call is recommended. (`Coin` attaches asset IDs to asset amounts.) You can pass in as many `Coin` instances as you wish. If you wish to add to a currently existing `Value` or change a current asset amount, the `AddAsset()` call is recommended.

Example code can be found below.

```go
// Create some Ada, matching was considered `Coins` in v5 Value instances.
oneAda := shared.ValueFromCoins(shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1_000_000)})

// Create some Ada and another asset.
assetId := shared.FromSeparate("policyX", "assetY")
assetZ := shared.ValueFromCoins(
    shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1_000_000)},
    shared.Coin{AssetId: assetId, Amount: num.Int64(2_000_000)},
)

// Add to a previously existing Value.
var v shared.Value
for _, a := range out.Amount {
    if a.Unit == "lovelace" {
        c, ok := num.New(a.Quantity)
        if !ok {
            return chainsync.TxOut{}, fmt.Errorf("invalid quantity: %w", err)
        }
    v.AddAsset(shared.Coin{AssetId: shared.AdaAssetID, Amount: c})
    } else {
        policyId := a.Unit[:56]
        assetName := a.Unit[56:]
        assetId := shared.AssetID(fmt.Sprintf("%v.%v", policyId, assetName))
        c, ok := num.New(a.Quantity)
        if !ok {
            return chainsync.TxOut{}, fmt.Errorf("invalid quantity: %w", err)
        }
        v.AddAsset(shared.Coin{AssetId: assetId, Amount: c})
    }
}
```

# Compatibility Module

A _compatibility_ module has been added to the code. The main purpose of the module is to act as a drop-in replacement when interacting with JSON, a DB, or some other form of at-rest data. (Ogmigo requests and responses are one example.) Let’s say you want to get the next block. You can use the compatibility module to receive the response from a v5 or v6 Ogmios instance, and place it in a v6 struct. It is highly recommended that you drop in any compatibility-related structs while still on v5, confirming all functionality still works as expected, and then migrating to v6.

There are a couple of caveats that are discussed below.

* The responses have changed between v5 and v6. Specific details are discussed in subsections but your code will need to adjust slightly in order to handle these changes.
* `Value` compatibility is pretty straightforward. It's just important to note that v6 `Value` structs compress Ada `Coins` from v5 into a double-nested map, alongside all other assets. (The asset ID for ADA is "ada.lovelace".) Take care when attempting to read any assets other than ADA. An example of how to read all non-ADA assets can be seen below.

```go
assetId1 := shared.FromSeparate("policyX", "assetY")
assetId2 := shared.FromSeparate("party", "anchor")
totalValue := shared.ValueFromCoins(
    shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1_000_000)},
    shared.Coin{AssetId: assetId1, Amount: num.Int64(2_000_000)},
    shared.Coin{AssetId: assetId2, Amount: num.Int64(3_000_000)},
)
for policy, policyMap := range totalValue.AssetsExceptAda() {
    for asset, amt := range policyMap {
        fmt.Printf(" - %v %v\n", amt.String(), shared.FromSeparate(policy, asset))
    }
}
```

## FindIntersect (v5) / findIntersection (v6)

When catching up to the tip of the blockchain, the code has changed slightly. Instead of checking the result type to see if it’s _IntersectionFound_ or _IntersectionNotFound_, you’ll just get the result and check to see if there’s an error. If there’s an error, you shouldn’t proceed and should process the error as is appropriate.

### FindIntersect example (v5)

An example of how `FindIntersect` works is seen below.

```go
func (h Handler) Handle(ctx context.Context, data []byte) error {
    ctx = h.logger.WithContext(ctx)

    var response chainsync.Response
    if err := json.Unmarshal(data, &response); err != nil {
        return err
    }

    if response.Result == nil {
        return nil
    }

    if err := h.publishToStream(ctx, response.Result); err != nil {
        return err
    }

    switch {
    case response.Result.IntersectionFound != nil:
        intersection := response.Result.IntersectionFound.Point
        // Process the intersection.

    case response.Result.IntersectionNotFound != nil:
        tip := response.Result.IntersectionNotFound.Tip
        err := chainsync.ResultError{Code: 1000, Message: "Intersection not found", Data: &tip}
        return err
    }

    return nil
}
```

### findIntersection example (v6)

The following example takes the previous code and updates it to use v6 (`findIntersection`).

```go
func (h Handler) Handle(ctx context.Context, data []byte) error {
    ctx = h.logger.WithContext(ctx)

    var response compatibility.CompatibleResponsePraos
    if err := json.Unmarshal(data, &response); err != nil {
        return err
    }

    if response.Result == nil {
        return nil
    }

    if err := h.publishToStream(ctx, response.Result); err != nil {
        return err
    }

    if response.Method == chainsync.FindIntersectionMethod {
        result := response.MustFindIntersectResult()
        if result.Error != nil {
            intersection := result.Intersection
            // Process the intersection.
        } else {
            return result.Error
        }
    }

    return nil
}
```

## RequestNext (v5) / nextBlock (v6)

If you’re requesting the next block, the response format has changed. This means that Ogmigo isn’t a drop-in replacement in this case. However, it’s easy enough to shift gears. You can check the response’s _Method_ field to see how to cast the response's _Result_ field. From there, you can check the _Direction_ field of the result in order to determine how to interpret the rest of the result.

### RequestNext example (v5)

An example of how `RequestNext` works is seen below.

```go
func (h Handler) Handle(ctx context.Context, data []byte) error {
    ctx = h.logger.WithContext(ctx)

    var response chainsync.Response
    if err := json.Unmarshal(data, &response); err != nil {
        return err
    }

    if response.Result == nil {
        return nil
    }

    if err := h.publishToStream(ctx, response.Result); err != nil {
        return err
    }

    switch {
    case response.Result.RollForward != nil:
        ps := response.Result.RollForward.Block.PointStruct()
        if err := h.dao.RollForward(ctx, ps); err != nil {
            return err
        }

    case response.Result.RollBackward != nil:
        if ps, ok := response.Result.RollBackward.Point.PointStruct(); ok {
            if err := h.dao.RollBackward(ctx, ps.Slot); err != nil {
                return err
            }
        }
    }

    return nil
}
```

### nextBlock example (v6)

The following example takes the previous code and updates it to use v6 (`nextBlock`).

```go
func (h Handler) Handle(ctx context.Context, data []byte) error {
    ctx = h.logger.WithContext(ctx)

    var response compatibility.CompatibleResponsePraos
    if err := json.Unmarshal(data, &response); err != nil {
        return err
    }

    if response.Result == nil {
        return nil
    }

    if err := h.publishToStream(ctx, response.Result); err != nil {
        return err
    }

    if response.Method == chainsync.NextBlockMethod {
        result := response.MustNextBlockResult()
        direction := result.Direction
        switch direction {
        case chainsync.RollForwardString:
            ps := result.Block.PointStruct()
            if err := h.dao.RollForward(ctx, ps); err != nil {
                return err
            }

        case chainsync.RollBackwardString:
            if ps, ok := result.Point.PointStruct(); ok {
                if err := h.dao.RollBackward(ctx, ps.Slot); err != nil {
                    return err
                }
            }

        default:
            return fmt.Errorf("invalid direction for next block - %v", direction)
        }
    }

    return nil
}
```
