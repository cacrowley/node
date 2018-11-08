package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type FeeDistributeType int8

const (
	FeeForProposer = FeeDistributeType(0x01)
	FeeForAll      = FeeDistributeType(0x02)
	FeeFree        = FeeDistributeType(0x03)

	canceledCharacter    = "#Cxl"
	expiredCharacter     = "#Exp"
	serializeSeparator   = ";"
	amountDenomSeparator = ":"
)

type Fee struct {
	Tokens sdk.Coins
	Type   FeeDistributeType
}

func NewFee(tokens sdk.Coins, distributeType FeeDistributeType) Fee {
	return Fee{
		Tokens: tokens,
		Type:   distributeType,
	}
}

func (fee *Fee) AddFee(other Fee) {
	if other.IsEmpty() {
		return
	}

	if fee.Tokens == nil {
		fee.Tokens = other.Tokens
		fee.Type = other.Type
	} else {
		fee.Tokens = fee.Tokens.Plus(other.Tokens)
		if other.Type == FeeForAll {
			fee.Type = FeeForAll
		}
	}
}

func (fee Fee) IsEmpty() bool {
	return fee.Tokens == nil || fee.Tokens.IsEqual(sdk.Coins{})
}

// Any change of this method should communicate with services (query, explorer) developers
// More detail can be found:
// https://github.com/BiJie/BinanceChain-Doc/wiki/Fee-Calculation,-Collection-and-Distribution#publication
func (fee Fee) SerializeForPub(canceled, expired int) string {
	if fee.IsEmpty() {
		return ""
	} else {
		res := fee.serialize()
		if canceled > 0 {
			res += fmt.Sprintf("%s%s%s%d", serializeSeparator, canceledCharacter, amountDenomSeparator, canceled)
		}
		if expired > 0 {
			res += fmt.Sprintf("%s%s%s%d", serializeSeparator, expiredCharacter, amountDenomSeparator, expired)
		}
		return res
	}
}

func (fee Fee) String() string {
	return fee.serialize()
}

func (fee Fee) serialize() string {
	if fee.IsEmpty() {
		return ""
	} else {
		var buffer bytes.Buffer
		for _, coin := range fee.Tokens {
			buffer.WriteString(fmt.Sprintf("%s%s%s%s", coin.Denom, amountDenomSeparator, coin.Amount.String(), serializeSeparator))
		}
		res := buffer.String()
		return res[:len(res)-1]
	}
}
