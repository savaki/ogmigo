package shared

import (
	"fmt"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
)

type Value map[string]map[string]num.Int

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
		for assetName, amt := range assets {
			if haveAssets, ok := have[policyId]; ok {
				if haveAssets[assetName].Int64() < amt.Int64() {
					return false, fmt.Errorf("not enough %v.%v to meet demand", policyId, assetName)
				}
			} else {
				return false, fmt.Errorf("not enough %v.%v to meet demand", policyId, assetName)
			}
		}
	}
	return true, nil
}

func (v Value) AddAsset(coins ...Coin) {
	for _, coin := range coins {
		if _, ok := v[coin.assetId.PolicyID()]; !ok {
			v[coin.assetId.PolicyID()] = map[string]num.Int{}
		}
		v[coin.assetId.PolicyID()][coin.assetId.AssetName()] = v[coin.assetId.PolicyID()][coin.assetId.AssetName()].Add(coin.amount)
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
	assetId AssetID
	amount  num.Int
}

func ValueFromCoins(coins ...Coin) Value {
	var value Value
	value.AddAsset(coins...)
	return value
}
