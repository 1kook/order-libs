package cex_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/goonma/order-libs/cex"
	"github.com/goonma/order-libs/defines"
	"github.com/goonma/order-libs/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestClient_GetLastClosedCandles(t *testing.T) {
	err := godotenv.Load(os.ExpandEnv("/config/.env"))
	if err != nil {
		err := godotenv.Load(os.ExpandEnv(".env"))
		if err != nil {
			panic(err)
		}
	}

	c := cex.NewBinanceClient(os.Getenv("BIN_API_KEY"), os.Getenv("BIN_API_SECRET"), false)
	data, err := c.GetLastClosedCandle(context.Background(), "BTCUSDT")
	if err != nil {
		t.Error(err)
		return
	}

	b, _ := json.MarshalIndent(data, "", "")
	fmt.Println(string(b))
}

func TestClient_GetLastClosedCandlesByCEX(t *testing.T) {
	err := godotenv.Load(os.ExpandEnv("/config/.env"))
	if err != nil {
		err := godotenv.Load(os.ExpandEnv(".env"))
		if err != nil {
			panic(err)
		}
	}

	opts := []cex.ClientOption{
		cex.WithApiKeyAndSecretKey(os.Getenv("BIN_API_KEY"), os.Getenv("BIN_API_SECRET")),
	}

	c, err := cex.NewCEXClient(defines.CEXCodeBinance, opts...)
	require.NoError(t, err)

	data, err := c.GetLastClosedCandle(context.Background(), "BTCUSDT")
	if err != nil {
		t.Error(err)
		return
	}

	b, _ := json.MarshalIndent(data, "", "")
	fmt.Println(string(b))
}

func TestClient_CreateOrderBuyCEX(t *testing.T) {
	var b []byte
	var err error
	var data *models.MarketOrder

	err = godotenv.Load(os.ExpandEnv("/config/.env"))
	if err != nil {
		err := godotenv.Load(os.ExpandEnv(".env"))
		if err != nil {
			panic(err)
		}
	}

	opts := []cex.ClientOption{
		cex.WithApiKeyAndSecretKey(os.Getenv("BIN_API_KEY"), os.Getenv("BIN_API_SECRET")),
	}

	c, err := cex.NewCEXClient(defines.CEXCodeBinance, opts...)
	require.NoError(t, err)

	pairToken := "BNBUSDT"
	clientOrderId := fmt.Sprintf("%s-%v", pairToken, time.Now().UTC().Unix())

	data, err = c.CreateOrderBuy(context.Background(), pairToken, "b-"+clientOrderId, 10)
	if err != nil {
		t.Error(err)
		return
	}

	b, _ = json.MarshalIndent(data, " ", "")
	fmt.Println(string(b))

	execQty, _ := strconv.ParseFloat(data.ExecutedQuantity, 64)
	data, err = c.CreateOrderSell(context.Background(), pairToken, "s-"+clientOrderId, execQty)
	if err != nil {
		t.Error(err)
		return
	}

	b, _ = json.MarshalIndent(data, " ", "")
	fmt.Println(string(b))
}
