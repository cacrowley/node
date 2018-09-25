package store

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/BiJie/BinanceChain/wire"
)

// queryOrderBook queries the store for the serialized order book for a given pair.
func queryOrderBook(cdc *wire.Codec, ctx context.CoreContext, pair string) (*[]byte, error) {
	bz, err := ctx.Query(fmt.Sprintf("dex/orderbook/%s", pair))
	if err != nil {
		return nil, err
	}
	return &bz, nil
}

// decodeOrderBook decodes the order book to a set of OrderBookLevel structs
func decodeOrderBook(cdc *wire.Codec, bz *[]byte) (*OrderBook, error) {
	var ob OrderBook
	err := cdc.UnmarshalBinary(*bz, &ob)
	if err != nil {
		return nil, err
	}
	return &ob, nil
}

// GetOrderBook decodes the order book from the serialized store
func GetOrderBook(cdc *wire.Codec, ctx context.CoreContext, pair string) (*OrderBook, error) {
	bz, err := queryOrderBook(cdc, ctx, pair)
	if err != nil {
		return nil, err
	}
	if bz == nil {
		return nil, nil
	}
	book, err := decodeOrderBook(cdc, bz)
	return book, err
}

func queryOpenOrders(cdc *wire.Codec, ctx context.CoreContext, pair string, addr string) (*[]byte, error) {
	if bz, err := ctx.Query(fmt.Sprintf("dex/openorders/%s/%s", pair, addr)); err != nil {
		return nil, err
	} else {
		return &bz, nil
	}
}

func DecodeOpenOrders(cdc *wire.Codec, bz *[]byte) ([]OpenOrder, error) {
	openOrders := make([]OpenOrder, 0)
	if err := cdc.UnmarshalBinary(*bz, &openOrders); err != nil {
		return nil, err
	} else {
		return openOrders, nil
	}
}

func GetOpenOrders(cdc *wire.Codec, ctx context.CoreContext, pair string, addr string) ([]OpenOrder, error) {
	if bz, err := queryOpenOrders(cdc, ctx, pair, addr); err != nil {
		return nil, err
	} else if bz == nil {
		return []OpenOrder{}, nil
	} else {
		openOrders, err := DecodeOpenOrders(cdc, bz)
		return openOrders, err
	}
}
