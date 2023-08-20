package utils

import (
	"context"
	_ "embed"
	"strconv"
	"strings"

	"github.com/adshao/go-binance/v2"
)

// https://api.binance.us/api/v3/exchangeInfo
var (
	mapLotSizes = make(map[string]*binance.LotSizeFilter, 0)
)

func init() {
	InitBinanceLotSizes()
}

func InitBinanceLotSizes() map[string]*binance.LotSizeFilter {
	if len(mapLotSizes) > 0 {
		return mapLotSizes
	}

	svc := binance.NewClient("", "").NewExchangeInfoService()
	res, err := svc.Do(context.Background())
	if err != nil {
		panic("cannot get lot size from binance")
	}

	for _, item := range res.Symbols {
		// Just support quote token USDT
		if !strings.HasSuffix(item.Symbol, "USDT") {
			continue
		}
		mapLotSizes[item.Symbol] = item.LotSizeFilter()
	}

	return mapLotSizes
}

func GetBinanceAppliedLotSizeAmount(symbol string, tokenAmount float64) float64 {
	digit := 0

	lotSizeFilter := mapLotSizes[symbol]
	if lotSizeFilter == nil {
		return tokenAmount
	}

	stepSize, _ := strconv.ParseFloat(lotSizeFilter.StepSize, 64)
	parts := strings.Split(strconv.FormatFloat(stepSize, 'f', -1, 64), ".")
	if len(parts) >= 2 {
		digit = len(parts[1])
	}

	return Round(tokenAmount, digit)
}
