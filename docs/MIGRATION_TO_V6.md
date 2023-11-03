
# Ogmigo v6 migration guide

This document discusses how to upgrade code that uses Ogmigo that is connected to Ogmios v5 instances and must connect to Ogmios v6. While not completely pain-free, it is possible to safely migrate with a minimum of pain.

Note that this document is a supplement to [the original Ogmios migration guide](https://github.com/CardanoSolutions/ogmios/blob/master/architectural-decisions/accepted/017-api-version-6-major-rewrite.md). It’s highly recommended that you read that particular guide before reading this guide, especially if your code handles a lot of individual values. There are also Ogmios API schema pages for [v5](https://ogmios.dev/api/v5.6/) and [v6](https://ogmios.dev/api/) that are worth consulting if all else fails.

# Caveats

There are a handful of caveats that should be considered before reading the main document.

* The Ogmigo v6 upgrade adds no new functionality beyond v6 struct support and a compatibility layer that allows both v5 and v6 Ogmigo JSON output to be unmarshalled into v6 structs. If anything wasn’t supported in Ogmigo v5, it won’t be supported in Ogmigo v6.
* Support for the Byron era is limited in Ogmigo v6, and may not work as expected. In fact, as of this writing (Nov. 2023), critical functionality, such as the _compatibility_ module, assumes Byron _isn’t_ supported. If you must have Byron support, further Ogmigo work will be required.
* Any attempt to use a v6-enabled Ogmigo library will automatically break the code using Ogmigo. Many fundamental structs (e.g., _Block_) have been altered to assume v6 structs. If somebody absolutely must use v5 structs, they’ll have to import the _v5_ module and use those structs. Even then, it is highly recommended to use the _compatibility_ module whenever possible, thereby assisting in a smooth transition to the default (v6) code.

# Support for v6

As mentioned in the caveats, the v6 structs are now the default in the Ogmigo library. Code will have to be written with this in mind if working directly with v6. The Ogmios migration guide discusses the gist of what changed between v5 and v6. In general, other than changing specific references within the structs, there aren’t major changes required to understand the new structs.

Note that in particular, Ogmigo automatically handles the switch from JSON-WSP to JSON-RPC 2.0. You should just be able to let Ogmigo handle everything automatically, simply pointing it to the appropriate Ogmios endpoint.

## Major Struct Changes

The Ogmios architectural decision doc linked above isn’t fully correct when discussing how various structs have changed. If you’re using v5 structs and want to see where to look in v6 structs for specific data, you should look at the _compatibility_ module and see how v5 is translated to v6.

That said, there are a few general caveats to note.

* A couple of more esoteric v6 struct fields can’t be populated by v5, at least not without adding a lot of code for little gain (e.g., consensus proofs). If there are gaps, get in touch but you may have to fill them yourself.
* Some data in v5 is encoded in Base64 and switched in v6 to Base16. The _compatibility_ module automatically handles conversion. Still, if your code directly uses these entries, you’ll need to adjust as necessary.
* Value types don’t always carry over between v5 and v6, and aren’t always properly defined by the API docs. We’ve looked at both real-world data and Ogmios unit tests, and have done the best we can to determine what’s allowed. It is still recommended that you double check any values critical to your functionality. If there are discrepancies, please get in touch.

# Compatibility Module

A _compatibility_ module has been added to the code. The main purpose of the module is to act as a drop-in replacement for Ogmigo requests and responses. Let’s say you want to get the next block. You can use the compatibility module to receive the response from a v5 or v6 Ogmios instance, and place it in a v6 struct. It is highly recommended that you drop in any compatibility-related structs while still on v5, confirming all functionality still works as expected, and then migrating to v6.

There are a couple of caveats that are discussed below.

* The responses have changed between v5 and v6. Specific details are discussed in subsections but your code will need to adjust slightly in order to handle these changes.
* The compatibility module doesn’t cover all possible use cases. For example, there may be data stored in a database using v5-oriented records. In such cases, you’ll need to manually convert the data.

## FindIntersect (v5) / findIntersection (v6)

When catching up to the tip of the blockchain, the code has changed slightly. Instead of checking the result type to see if it’s _IntersectionFound_ or _IntersectionNotFound_, you’ll just get the result and check to see if there’s an error. If there’s an error, you shouldn’t proceed and should process the error as is appropriate.

### FindIntersect example (v5)

Note that interesting portions have been bolded, and should be compared against v6 code.

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
        //

    case response.Result.IntersectionNotFound != nil:
        tip := response.Result.IntersectionNotFound.Tip
        err := chainsync.ResultError{Code: 1000, Message: "Intersection not found", Data: &tip}
        return err
    }

    return nil
}
```

### findIntersection example (v6)

The following example takes the previous code and updates it to use v6. Interesting portions have been bolded.

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
        //
        } else {
            return result.Error
        }
    }

    return nil
}
```

## RequestNext (v5) / nextBlock (v6)

If you’re requesting the next block, the response format has changed. This means that Ogmigo isn’t a drop-in replacement in this case. However, it’s easy enough to shift gears. You can check the response’s Method field. From there, you can check the _Direction_ field of the response’s _Result_ field which has been cast as an appropriate struct. An example can be seen below.

### RequestNext example (v5)

Note that interesting portions have been bolded, and should be compared against v6 code.

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

The following example takes the previous code and updates it to use v6. Interesting portions have been bolded.

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
