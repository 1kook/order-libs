package cex

import (
	"context"
	"fmt"

	"github.com/goonma/order-libs/defines"
	"github.com/goonma/order-libs/models"
)

type clientOption struct {
	ApiKey     string
	SecretKey  string
	UseTestNet bool
}

type ClientOption struct {
	apply func(*clientOption)
}

type CEXClient interface {
	GetLastClosedCandle(ctx context.Context, tokenPair string) (*models.ClosedCandle, error)
	CreateOrderBuy(ctx context.Context, tokenPair string, clientOrderId string, amount float64) (*models.MarketOrder, error)
	CreateOrderSell(ctx context.Context, tokenPair string, clientOrderId string, tokenAmount float64) (*models.MarketOrder, error)
	GetOrderById(ctx context.Context, marketOrderId string) (*models.MarketOrder, error)
	GetOrderByClientOrderId(ctx context.Context, clientOrderId string) (*models.MarketOrder, error)
}

func NewCEXClient(code defines.CEXCode, opts ...ClientOption) (CEXClient, error) {
	opt := &clientOption{}
	if len(opts) > 0 {
		opt = applyClientOption(opts)
	}

	switch code {
	case defines.CEXCodeBinance:
		if len(opt.ApiKey) == 0 || len(opt.SecretKey) == 0 {
			return nil, fmt.Errorf("Invalid API key and Secret key")
		}
		return NewBinanceClient(opt.ApiKey, opt.SecretKey, opt.UseTestNet), nil
	case defines.CEXCodeMock:
		return NewMockClient(), nil
	}

	return nil, fmt.Errorf("Not found CEX client")
}

func WithUseTestNet() ClientOption {
	return ClientOption{
		apply: func(opt *clientOption) {
			opt.UseTestNet = true
		},
	}
}

func WithApiKeyAndSecretKey(apiKey string, secretKey string) ClientOption {
	return ClientOption{
		apply: func(opt *clientOption) {
			opt.ApiKey = apiKey
			opt.SecretKey = secretKey
		},
	}
}

func applyClientOption(callOptions []ClientOption) *clientOption {
	if len(callOptions) == 0 {
		return &clientOption{}
	}

	optCopy := &clientOption{}
	for _, f := range callOptions {
		f.apply(optCopy)
	}
	return optCopy
}
