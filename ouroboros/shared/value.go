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
				if haveAssets[assetName].BigInt().Cmp(amt.BigInt()) == -1 {
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
			if a[policy] != nil && a[policy][asset].BigInt().Cmp(amt.BigInt()) != -1 {
				return false
			}
		}
	}

	return true
}

func GreaterThan(a, b Value) bool {
	for policy, policyMap := range b {
		for asset, amt := range policyMap {
			if a[policy] != nil && a[policy][asset].BigInt().Cmp(amt.BigInt()) != 1 {
				return false
			}
		}
	}

	return true
}

func Equal(a, b Value) bool {
	for policy, policyMap := range b {
		for asset, amt := range policyMap {
			if a[policy] != nil && a[policy][asset].BigInt().Cmp(amt.BigInt()) != 0 {
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

func (v Value) AssetsExceptAda() Value {
	policies := Value{}
	for policy, tokenMap := range v {
		if policy == AdaPolicy {
			continue
		}
		policies[policy] = map[string]num.Int{}
		for token, quantity := range tokenMap {
			policies[policy][token] = quantity
		}
	}
	return policies
}

type Coin struct {
	AssetId AssetID
	Amount  num.Int
}

func ValueFromCoins(coins ...Coin) Value {
	var value Value
	value.AddAsset(coins...)
	return value
}
