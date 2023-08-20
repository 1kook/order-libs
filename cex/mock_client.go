package cex

import (
	"context"
	"fmt"
	"time"

	"github.com/goonma/order-libs/defines"
	"github.com/goonma/order-libs/models"
)

type MockClient struct {
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) GetLastClosedCandle(ctx context.Context, tokenPair string) (*models.ClosedCandle, error) {
	closePrice := ctx.Value("close_price")
	return &models.ClosedCandle{
		Open:      fmt.Sprint(closePrice),
		Close:     fmt.Sprint(closePrice),
		CloseTime: time.Now().UTC(),
		OpenTime:  time.Now().UTC(),
	}, nil
}

func (c *MockClient) CreateOrderBuy(ctx context.Context, tokenPair string, clientOrderId string, amount float64) (*models.MarketOrder, error) {

	respCode := fmt.Sprintf("%v", ctx.Value("resp_code"))
	err := getResponseError(respCode)
	if err != nil {
		return nil, err
	}

	tokenAmount := float64(0)
	if amt, ok := ctx.Value("token_amount").(float64); ok {
		tokenAmount = amt
	}

	tnow := time.Now()
	tunix := tnow.UnixMilli()

	return &models.MarketOrder{
		Symbol:        tokenPair,
		OrderID:       fmt.Sprint(tunix),
		ClientOrderID: clientOrderId,
		// Price:            fmt.Sprintf("%f", amount/tokenAmount),
		QuoteQuantity:    fmt.Sprintf("%f", amount),
		OrigQuantity:     fmt.Sprintf("%.10f", tokenAmount),
		ExecutedQuantity: fmt.Sprintf("%.10f", tokenAmount),
		Status:           defines.MarketOrderStatusFilled,
		Type:             defines.OrderTypeMarket,
		Side:             defines.OrderSideBuy,
		TransactTime:     tnow,
		Timestamp:        tnow,
		Raw:              nil,
	}, nil
}

func (c *MockClient) CreateOrderSell(ctx context.Context, tokenPair string, clientOrderId string, tokenAmount float64) (*models.MarketOrder, error) {
	respCode := fmt.Sprintf("%v", ctx.Value("resp_code"))
	err := getResponseError(respCode)
	if err != nil {
		return nil, err
	}

	amount := float64(0)
	if amt, ok := ctx.Value("price_amount").(float64); ok {
		amount = amt
	}

	tnow := time.Now()
	tunix := tnow.UnixMilli()

	return &models.MarketOrder{
		Symbol:        tokenPair,
		OrderID:       fmt.Sprint(tunix),
		ClientOrderID: clientOrderId,
		// Price:            fmt.Sprintf("%f", amount/tokenAmount),
		QuoteQuantity:    fmt.Sprintf("%f", amount),
		OrigQuantity:     fmt.Sprintf("%.10f", tokenAmount),
		ExecutedQuantity: fmt.Sprintf("%.10f", tokenAmount),
		Status:           defines.MarketOrderStatusFilled,
		Type:             defines.OrderTypeMarket,
		Side:             defines.OrderSideSell,
		TransactTime:     tnow,
		Timestamp:        tnow,
		Raw:              nil,
	}, nil
}

func (c *MockClient) GetOrderById(ctx context.Context, marketOrderId string) (*models.MarketOrder, error) {
	return nil, nil
}

func (c *MockClient) GetOrderByClientOrderId(ctx context.Context, clientOrderId string) (*models.MarketOrder, error) {
	return nil, nil
}

func getResponseError(respCode string) error {
	switch respCode {
	case "2":
		return defines.MarketOrderErrorNotEnoughBalance
	case "3":
		return defines.MarketOrderErrorFilterFailureLotSize
	case "4":
		return defines.CexErrorInvalidAPIKey
	case "5":
		return defines.CexErrorInvalidAPIKey
	case "1", "":
		return nil
	}
	return defines.UnknownError
}
