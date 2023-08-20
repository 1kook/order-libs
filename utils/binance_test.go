package utils_test

import (
	"fmt"
	"testing"

	"github.com/goonma/order-libs/utils"
	"github.com/stretchr/testify/assert"
)

func TestInitLotSize(t *testing.T) {
	lsz := utils.InitBinanceLotSizes()
	fmt.Printf("%v", lsz)
}

func TestGetBinanceAppliedLotSizeAmount(t *testing.T) {
	tokenAmount := 13.45672343

	rs := utils.GetBinanceAppliedLotSizeAmount("OPUSDT", tokenAmount)
	assert.Equal(t, 13.45, rs)

	rs = utils.GetBinanceAppliedLotSizeAmount("BNBUSDT", tokenAmount)
	assert.Equal(t, 13.456, rs)

	rs = utils.GetBinanceAppliedLotSizeAmount("BTCUSDT", tokenAmount)
	assert.Equal(t, 13.45672, rs)
}
