package cex

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/common"
	"github.com/goonma/order-libs/defines"
	"github.com/goonma/order-libs/models"
	"github.com/goonma/order-libs/utils"
)

type BinanceClient struct {
	client *binance.Client
}

func NewBinanceClient(apiKey, apiSecret string, useTestNet bool) *BinanceClient {
	binance.UseTestnet = useTestNet
	return &BinanceClient{
		client: binance.NewClient(apiKey, apiSecret),
	}
}

func (c *BinanceClient) GetLastClosedCandle(ctx context.Context, tokenPair string) (*models.ClosedCandle, error) {
	klines, err := c.client.NewKlinesService().
		Symbol(tokenPair).
		Interval("1s").
		EndTime(time.Now().UnixMilli()).
		Limit(1).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	if len(klines) == 0 {
		return nil, nil
	}

	return &models.ClosedCandle{
		Open:      klines[0].Open,
		High:      klines[0].High,
		Low:       klines[0].Low,
		Close:     klines[0].Close,
		OpenTime:  time.UnixMilli(klines[0].OpenTime),
		CloseTime: time.UnixMilli(klines[0].CloseTime),

		Raw: klines[0],
	}, nil
}

func (c *BinanceClient) CreateOrderBuy(ctx context.Context, tokenPair string, clientOrderId string, amount float64) (*models.MarketOrder, error) {

	// Get last close candle
	lcc, err := c.GetLastClosedCandle(ctx, tokenPair)
	if err != nil {
		return nil, fmt.Errorf("Cannot get last closed price [%s]: %s", tokenPair, err.Error())
	}

	currentTokenPrice, err := strconv.ParseFloat(lcc.Close, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid last closed price [%s]: %s", tokenPair, err.Error())
	}

	buyTokenAmount := amount / currentTokenPrice
	buyTokenAmount = utils.GetBinanceAppliedLotSizeAmount(tokenPair, buyTokenAmount)

	// amountStr := fmt.Sprintf("%f", amount)
	amountStr := fmt.Sprintf("%f", buyTokenAmount)
	resp, err := c.client.NewCreateOrderService().
		NewClientOrderID(clientOrderId).
		Symbol(tokenPair).
		Type(binance.OrderTypeMarket).
		Side(binance.SideTypeBuy).
		Quantity(amountStr).
		// Fixes for error: "Quote order qty market orders are not supported for this symbol."
		// QuoteOrderQty(amountStr).
		Do(ctx)

	if err != nil {
		return nil, convertCexErrorFromBinance(err)
	}

	executedQuantity, err := strconv.ParseFloat(resp.ExecutedQuantity, 64)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse ExecutedQuantity: %s", err.Error())
	}

	tokenPrice := resp.Price
	if len(resp.Fills) > 0 {
		// just support quote asset usdt
		token := strings.ReplaceAll(tokenPair, "USDT", "")
		for _, item := range resp.Fills {
			if token != item.CommissionAsset {
				continue
			}
			tokenPrice = item.Price
			commission, _ := strconv.ParseFloat(item.Commission, 64)
			executedQuantity = executedQuantity - commission
		}
	}

	strQuantity := resp.ExecutedQuantity
	if executedQuantity > 0 {
		strQuantity = fmt.Sprint(executedQuantity)
	}

	return &models.MarketOrder{
		Symbol:           tokenPair,
		OrderID:          fmt.Sprint(resp.OrderID),
		ClientOrderID:    resp.ClientOrderID,
		Price:            tokenPrice,
		OrigQuantity:     resp.OrigQuantity,
		QuoteQuantity:    resp.CummulativeQuoteQuantity,
		ExecutedQuantity: strQuantity,

		Side:   defines.OrderSideBuy,
		Type:   defines.OrderType(resp.Type),
		Status: defines.MarketOrderStatus(resp.Status),

		Timestamp:    time.UnixMilli(resp.TransactTime),
		TransactTime: time.UnixMilli(resp.TransactTime),

		Raw: resp,
	}, nil
}

func (c *BinanceClient) CreateOrderSell(ctx context.Context, tokenPair string, clientOrderId string, tokenAmount float64) (*models.MarketOrder, error) {

	// update lot size amount
	tokenAmount = utils.GetBinanceAppliedLotSizeAmount(tokenPair, tokenAmount)

	amountStr := fmt.Sprintf("%f", tokenAmount)
	resp, err := c.client.NewCreateOrderService().
		NewClientOrderID(clientOrderId).
		Symbol(tokenPair).
		Type(binance.OrderTypeMarket).
		Side(binance.SideTypeSell).
		Quantity(amountStr).
		Do(ctx)

	if err != nil {
		return nil, convertCexErrorFromBinance(err)
	}

	executedQuantity, err := strconv.ParseFloat(resp.ExecutedQuantity, 64)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse ExecutedQuantity: %s", err.Error())
	}

	quoteQuantity, err := strconv.ParseFloat(resp.CummulativeQuoteQuantity, 64)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse CummulativeQuoteQuantity: %s", err.Error())
	}

	tokenPrice := fmt.Sprintf("%f", quoteQuantity/executedQuantity)
	return &models.MarketOrder{
		Symbol:           tokenPair,
		OrderID:          fmt.Sprint(resp.OrderID),
		ClientOrderID:    resp.ClientOrderID,
		Price:            tokenPrice,
		OrigQuantity:     resp.OrigQuantity,
		QuoteQuantity:    resp.CummulativeQuoteQuantity,
		ExecutedQuantity: resp.ExecutedQuantity,

		Side:   defines.OrderSideSell,
		Type:   defines.OrderType(resp.Type),
		Status: defines.MarketOrderStatus(resp.Status),

		Timestamp:    time.UnixMilli(resp.TransactTime),
		TransactTime: time.UnixMilli(resp.TransactTime),

		Raw: resp,
	}, nil
}

func (c *BinanceClient) GetOrderById(ctx context.Context, marketOrderId string) (*models.MarketOrder, error) {
	return nil, nil
}

func (c *BinanceClient) GetOrderByClientOrderId(ctx context.Context, clientOrderId string) (*models.MarketOrder, error) {
	return nil, nil
}

func convertCexErrorFromBinance(err error) error {
	var rErr *defines.CexError

	if err == nil {
		return nil
	}

	if !common.IsAPIError(err) {
		rErr = defines.UnknownError
		rErr.SetRawError(err)
		return rErr
	}

	exErr := err.(*common.APIError)
	switch exErr.Code {
	case -2014:
		rErr = defines.CexErrorInvalidAPIKey
	case -1013:
		rErr = defines.MarketOrderErrorFilterFailureLotSize
	case -2010:
		rErr = defines.MarketOrderErrorNotEnoughBalance
	default:
		rErr = defines.UnknownError
	}

	rErr.SetRawError(err)
	return rErr
}
