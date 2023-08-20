package models

import (
	"time"

	"github.com/goonma/order-libs/defines"
)

type ClosedCandle struct {
	Open  string `json:"open"`
	High  string `json:"high"`
	Low   string `json:"low"`
	Close string `json:"close"`

	OpenTime  time.Time   `json:"openTime"`
	CloseTime time.Time   `json:"closeTime"`
	Raw       interface{} `json:"-"`
}

// define create market order response
type MarketOrder struct {
	Symbol        string `json:"symbol"`
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`

	Price string `json:"price"`

	// Order total amount or CummulativeQuoteQuantity
	QuoteQuantity string `json:"quote_quantity"`

	// Token amount
	OrigQuantity string `json:"origQty"`
	// Token amount after commission
	ExecutedQuantity string `json:"executedQty"`

	Status defines.MarketOrderStatus `json:"status"`
	Type   defines.OrderType         `json:"type"`
	Side   defines.OrderSide         `json:"side"`

	TransactTime time.Time `json:"transactTime"`
	Timestamp    time.Time `json:"timestamp"`

	// raw response from cex sdk client
	Raw interface{} `json:"-"`
}
