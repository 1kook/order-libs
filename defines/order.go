package defines

type (
	OrderSide string
	OrderType string

	MarketOrderStatus string
)

const (
	// For trading or market order
	OrderTypeLimit  OrderType = "LIMIT"
	OrderTypeMarket OrderType = "MARKET"

	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"

	// For market order
	MarketOrderStatusNew             MarketOrderStatus = "NEW"
	MarketOrderStatusPartiallyFilled MarketOrderStatus = "PARTIALLY_FILLED"
	MarketOrderStatusFilled          MarketOrderStatus = "FILLED"
	MarketOrderStatusCanceled        MarketOrderStatus = "CANCELED"
	MarketOrderStatusPendingCancel   MarketOrderStatus = "PENDING_CANCEL"
	MarketOrderStatusRejected        MarketOrderStatus = "REJECTED"
	MarketOrderStatusExpired         MarketOrderStatus = "EXPIRED"

	OrderStatusNew           = "NEW"
	OrderStatusComplete      = "COMPLETE"
	OrderStatusFailPushEx    = "FAIL_PUSH_EX"
	OrderStatusPendingPushEx = "PENDING_PUSH_EX"

	// ?????
	OrderStatusStopLoss = "STOP_LOSS"
)
