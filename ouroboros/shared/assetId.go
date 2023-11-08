package shared

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

type AssetID string

// AdaAssetID maps to Coins in the v5 Value struct.
const (
	AdaPolicy        = "ada"
	AdaAsset         = "lovelace"
	AdaAssetIDString = AdaPolicy + "." + AdaAsset
	AdaAssetID       = AssetID(AdaAssetIDString)
)

func FromSeparate(policy string, assetName string) AssetID {
	return AssetID(fmt.Sprintf("%v.%v", policy, assetName))
}

func (a AssetID) HasPolicyID(s string) bool {
	return len(s) == 56 && strings.HasPrefix(string(a), s)
}

func (a AssetID) HasAssetID(re *regexp.Regexp) bool {
	return re.MatchString(string(a))
}

func (a AssetID) IsZero() bool {
	return a == ""
}

func (a AssetID) MatchAssetName(re *regexp.Regexp) ([]string, bool) {
	if assetName := a.AssetName(); len(assetName) > 0 {
		ss := re.FindStringSubmatch(assetName)
		return ss, len(ss) > 0
	}
	return nil, false
}

func (a AssetID) String() string {
	return string(a)
}

func (a AssetID) AssetName() string {
	s := string(a)
	if index := strings.Index(s, "."); index > 0 {
		return s[index+1:]
	}
	return ""
}

func (a AssetID) AssetNameUTF8() (string, bool) {
	if data, err := hex.DecodeString(a.AssetName()); err == nil {
		if utf8.Valid(data) {
			return string(data), true
		}
	}
	return "", false
}

func (a AssetID) PolicyID() string {
	s := string(a)
	if index := strings.Index(s, "."); index > 0 {
		return s[:index]
	}
	return s // Assets with empty-string name come back as just the policy ID
}
