package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/BiJie/BinanceChain/common/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/common"
)

const (
	Decimals       int8  = 8
	MaxTotalSupply int64 = 9000000000000000000 // 90 billions with 8 decimal digits

	DotBSuffix = ".B"
)

type Token struct {
	Name        string      `json:"name"`
	Symbol      string      `json:"symbol"`
	TotalSupply int64       `json:"total_supply"`
	Owner       sdk.Address `json:"owner"`
}

func NewToken(name, symbol string, totalSupply int64, owner sdk.Address) Token {
	return Token{
		Name:        name,
		Symbol:      symbol,
		TotalSupply: totalSupply,
		Owner:       owner,
	}
}

func (token *Token) IsOwner(addr sdk.Address) bool { return bytes.Equal(token.Owner, addr) }
func (token Token) String() string {
	return fmt.Sprintf("{Name: %v, Symbol: %v, TotalSupply: %v, Owner: %X}",
		token.Name, token.Symbol, token.TotalSupply, token.Owner)
}

func ValidateSymbol(symbol string) error {
	if len(symbol) == 0 {
		return errors.New("token symbol cannot be empty")
	}

	if len(symbol) > 8 {
		return errors.New("token symbol is too long")
	}

	if strings.HasSuffix(symbol, DotBSuffix) {
		symbol = strings.TrimSuffix(symbol, DotBSuffix)
	}

	if !utils.IsAlphaNum(symbol) {
		return errors.New("token symbol should be alphanumeric")
	}

	return nil
}

func GenerateTokenAddress(token Token, sequence int64) (sdk.Address, error) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(sequence))
	secret := append(token.Owner, b...)
	priv := makePrivKey(secret)
	return priv.PubKey().Address(), nil
}

func makePrivKey(secret common.HexBytes) crypto.PrivKey {
	privKey := crypto.GenPrivKeyEd25519FromSecret(secret)
	return privKey
}
