package shared

import (
	"fmt"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
)

type Value map[string]map[string]num.Int

var ErrInsufficientFunds = fmt.Errorf("insufficient funds")

func Add(a Value, b Value) Value {
	result := Value{}
	for policyId, assets := range a {
		for assetName, amt := range assets {
			if _, ok := result[policyId]; !ok {
				result[policyId] = map[string]num.Int{}
			}
			result[policyId][assetName] = amt
		}
	}
	for policyId, assets := range b {
		for assetName, amt := range assets {
			if _, ok := result[policyId]; !ok {
				result[policyId] = map[string]num.Int{}
			}
			result[policyId][assetName] = result[policyId][assetName].Add(amt)
		}
	}

	return result
}

func Subtract(a Value, b Value) Value {
	result := Value{}
	for policyId, assets := range a {
		for assetName, amt := range assets {
			if _, ok := result[policyId]; !ok {
				result[policyId] = map[string]num.Int{}
			}
			result[policyId][assetName] = amt
		}
	}
	for policyId, assets := range b {
		for assetName, amt := range assets {
			if _, ok := result[policyId]; !ok {
				result[policyId] = map[string]num.Int{}
			}
			result[policyId][assetName] = result[policyId][assetName].Sub(amt)
		}
	}

	return result
}

func Enough(have Value, want Value) (bool, error) {
	for policyId, assets := range want {
		if haveAssets, ok := have[policyId]; ok {
			for assetName, amt := range assets {
				if haveAssets[assetName].LessThan(amt) {
					return false, fmt.Errorf("not enough %v (%v) to meet demand (%v): %w", assetName, have[policyId][assetName].String(), amt, ErrInsufficientFunds)
				}
			}
		}
	}
	return true, nil
}

func LessThan(a, b Value) bool {
	for policy, policyMap := range b {
		for asset, amt := range policyMap {
			if a[policy] != nil && !a[policy][asset].LessThan(amt) {
				return false
			}
		}
	}

	return true
}

func GreaterThan(a, b Value) bool {
	for policy, policyMap := range b {
		for asset, amt := range policyMap {
			if a[policy] != nil && !a[policy][asset].GreaterThan(amt) {
				return false
			}
		}
	}

	return true
}

func Equal(a, b Value) bool {
	for policy, policyMap := range b {
		for asset, amt := range policyMap {
			if a[policy] != nil && !a[policy][asset].Equal(amt) {
				return false
			}
		}
	}

	return true
}

func (v Value) AddAsset(coins ...Coin) {
	for _, coin := range coins {
		if _, ok := v[coin.AssetId.PolicyID()]; !ok {
			v[coin.AssetId.PolicyID()] = map[string]num.Int{}
		}
		v[coin.AssetId.PolicyID()][coin.AssetId.AssetName()] = v[coin.AssetId.PolicyID()][coin.AssetId.AssetName()].Add(coin.Amount)
	}
}

func (v Value) AdaLovelace() num.Int {
	return v.AssetAmount(AdaAssetID)
}

func (v Value) AssetAmount(asset AssetID) num.Int {
	if nested, ok := v[asset.PolicyID()]; ok {
		return nested[asset.AssetName()]
	}
	return num.Int64(0)
}

func (v Value) AssetsExceptAda() (Value, uint32) {
	policies := Value{}
	var cnt uint32 = 0
	for policy, tokenMap := range v {
		if policy == AdaPolicy {
			continue
		}
		policies[policy] = map[string]num.Int{}
		for token, quantity := range tokenMap {
			policies[policy][token] = quantity
			cnt++
		}
	}
	return policies, cnt
}

func (v Value) AssetsExceptAdaCount() uint32 {
	var cnt uint32 = 0
	for policy, tokenMap := range v {
		if policy == AdaPolicy {
			continue
		}
		cnt += uint32(len(tokenMap))
	}
	return cnt
}

func (v Value) IsAdaPresent() bool {
	if v[AdaPolicy] != nil {
		if v[AdaPolicy][AdaAsset].GreaterThan(num.Uint64(0)) {
			return true
		}
	}

	return false
}

type Coin struct {
	AssetId AssetID
	Amount  num.Int
}

func ValueFromCoins(coins ...Coin) Value {
	value := Value{}
	value.AddAsset(coins...)
	return value
}
