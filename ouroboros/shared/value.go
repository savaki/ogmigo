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

// Maps to Coins in v5 Value struct.
func NewValueAda(coins num.Int) Value {
	value := Value{}
	value.AddAsset(AdaAssetID, coins)
	return value
}

func NewValueAsset(asset AssetID, coins num.Int) Value {
	value := Value{}
	value.AddAsset(asset, coins)
	return value
}

func (v Value) AddAda(coins num.Int) {
	v.AddAsset(AdaAssetID, coins)
}

func (v Value) AddAsset(asset AssetID, coins num.Int) {
	if _, ok := v[asset.PolicyID()]; !ok {
		v[asset.PolicyID()] = map[string]num.Int{}
	}
	v[asset.PolicyID()][asset.AssetName()] = coins
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
