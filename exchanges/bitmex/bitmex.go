package bitmex

import (
	. "github.com/coinrust/crex"
	"github.com/frankrap/bitmex-api"
	"github.com/frankrap/bitmex-api/swagger"
	"strings"
	"time"
)

// BitMEX the BitMEX exchange
type BitMEX struct {
	client *bitmex.BitMEX
	symbol string
}

func (b *BitMEX) GetName() (name string) {
	return "bitmex"
}

func (b *BitMEX) GetBalance(currency string) (result Balance, err error) {
	var margin swagger.Margin
	margin, err = b.client.GetMargin()
	if err != nil {
		return
	}
	result.Equity = float64(margin.MarginBalance)
	result.Available = float64(margin.AvailableMargin)
	result.RealizedPnl = float64(margin.RealisedPnl)
	result.UnrealisedPnl = float64(margin.UnrealisedPnl)
	return
}

func (b *BitMEX) GetOrderBook(symbol string, depth int) (result OrderBook, err error) {
	var ret bitmex.OrderBook
	ret, err = b.client.GetOrderBook(depth, symbol)
	if err != nil {
		return
	}
	for _, v := range ret.Asks {
		result.Asks = append(result.Asks, Item{
			Price:  v.Price,
			Amount: v.Amount,
		})
	}
	for _, v := range ret.Bids {
		result.Bids = append(result.Bids, Item{
			Price:  v.Price,
			Amount: v.Amount,
		})
	}
	result.Time = ret.Timestamp
	return
}

func (b *BitMEX) GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error) {
	//@param "binSize" (string) Time interval to bucket by. Available options: [1m,5m,1h,1d].
	var binSize string
	if strings.HasSuffix(period, "m") {
		binSize = period
	} else if strings.HasSuffix(period, "h") {
		binSize = period
	} else if strings.HasSuffix(period, "d") {
		binSize = period
	} else {
		binSize = period + "m"
	}
	var o []swagger.TradeBin
	o, err = b.client.GetBucketed(symbol,
		binSize,
		false,
		"",
		"",
		float32(limit),
		-1,
		false,
		time.Unix(from, 0),
		time.Unix(end, 0))
	if err != nil {
		return
	}
	for _, v := range o {
		records = append(records, Record{
			Symbol:    v.Symbol,
			Timestamp: v.Timestamp,
			Open:      v.Open,
			High:      v.High,
			Low:       v.Low,
			Close:     v.Close,
			Volume:    float64(v.Volume),
		})
	}
	return
}

func (b *BitMEX) SetContractType(currencyPair string, contractType string) (err error) {
	b.symbol = currencyPair
	return
}

func (b *BitMEX) GetContractID() (symbol string, err error) {
	return b.symbol, nil
}

func (b *BitMEX) SetLeverRate(value float64) (err error) {
	return
}

func (b *BitMEX) OpenLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size)
}

func (b *BitMEX) OpenShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size)
}

func (b *BitMEX) CloseLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Sell, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *BitMEX) CloseShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error) {
	return b.PlaceOrder(symbol, Buy, orderType, price, size, OrderReduceOnlyOption(true))
}

func (b *BitMEX) PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64,
	size float64, opts ...OrderOption) (result Order, err error) {
	params := ParseOrderParameter(opts...)
	var side string
	var _orderType string
	if direction == Buy {
		side = bitmex.SIDE_BUY
	} else if direction == Sell {
		side = bitmex.SIDE_SELL
	}
	if orderType == OrderTypeLimit {
		_orderType = bitmex.ORD_TYPE_LIMIT
	} else if orderType == OrderTypeMarket {
		_orderType = bitmex.ORD_TYPE_MARKET
	} else if orderType == OrderTypeStopLimit {
		_orderType = bitmex.ORD_TYPE_STOP_LIMIT
	} else if orderType == OrderTypeStopMarket {
		_orderType = bitmex.ORD_TYPE_STOP
	}
	var execInst string
	if params.PostOnly {
		execInst = "ParticipateDoNotInitiate"
	}
	if params.ReduceOnly {
		if execInst != "" {
			execInst += ","
		}
		execInst += "ReduceOnly"
	}
	var order swagger.Order
	order, err = b.client.PlaceOrder(side,
		_orderType, params.StopPx, price, int32(size), "", execInst, symbol)
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *BitMEX) GetOpenOrders(symbol string) (result []Order, err error) {
	var ret []swagger.Order
	ret, err = b.client.GetOrders(symbol)
	if err != nil {
		return
	}
	for _, v := range ret {
		result = append(result, b.convertOrder(&v))
	}
	return
}

func (b *BitMEX) GetOrder(symbol string, id string) (result Order, err error) {
	var ret swagger.Order
	ret, err = b.client.GetOrder(id, symbol)
	if err != nil {
		return
	}
	result = b.convertOrder(&ret)
	return
}

func (b *BitMEX) CancelOrder(symbol string, id string) (result Order, err error) {
	var order swagger.Order
	order, err = b.client.CancelOrder(id)
	if err != nil {
		return
	}
	result = b.convertOrder(&order)
	return
}

func (b *BitMEX) CancelAllOrders(symbol string) (err error) {
	_, err = b.client.CancelAllOrders(symbol)
	return
}

func (b *BitMEX) AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error) {
	var resp swagger.Order
	resp, err = b.client.AmendOrder2(id, "", "", 0, float32(size), 0, 0, price, 0, 0, "")
	if err != nil {
		return
	}
	result = b.convertOrder(&resp)
	return
}

func (b *BitMEX) GetPositions(symbol string) (result []Position, err error) {
	var ret swagger.Position
	ret, err = b.client.GetPosition(symbol)
	if err != nil {
		return
	}
	result = []Position{
		{
			Symbol:    symbol,
			OpenTime:  time.Time{},
			OpenPrice: ret.AvgEntryPrice,
			Size:      float64(ret.CurrentQty),
			AvgPrice:  ret.AvgCostPrice,
		},
	}
	return
}

func (b *BitMEX) convertOrder(order *swagger.Order) (result Order) {
	result.ID = order.OrderID
	result.Symbol = order.Symbol
	result.Price = order.Price
	result.StopPx = order.StopPx
	result.Size = float64(order.OrderQty)
	result.Direction = b.convertDirection(order.Side)
	result.Type = b.convertOrderType(order.OrdType)
	result.AvgPrice = order.AvgPx
	result.FilledAmount = float64(order.CumQty)
	if strings.Contains(order.ExecInst, "ParticipateDoNotInitiate") {
		result.PostOnly = true
	}
	if strings.Contains(order.ExecInst, "ReduceOnly") {
		result.ReduceOnly = true
	}
	result.Status = b.orderStatus(order)
	return
}

func (b *BitMEX) convertDirection(side string) Direction {
	switch side {
	case bitmex.SIDE_BUY:
		return Buy
	case bitmex.SIDE_SELL:
		return Sell
	default:
		return Buy
	}
}

func (b *BitMEX) convertOrderType(orderType string) OrderType {
	switch orderType {
	case bitmex.ORD_TYPE_LIMIT:
		return OrderTypeLimit
	case bitmex.ORD_TYPE_MARKET:
		return OrderTypeMarket
	case bitmex.ORD_TYPE_STOP_LIMIT:
		return OrderTypeStopLimit
	case bitmex.ORD_TYPE_STOP:
		return OrderTypeStopMarket
	default:
		return OrderTypeLimit
	}
}

func (b *BitMEX) orderStatus(order *swagger.Order) OrderStatus {
	orderState := order.OrdStatus
	switch orderState {
	case bitmex.OS_NEW:
		return OrderStatusNew
	case bitmex.OS_PARTIALLY_FILLED:
		return OrderStatusPartiallyFilled
	case bitmex.OS_FILLED:
		return OrderStatusFilled
	case bitmex.OS_CANCELED:
		return OrderStatusCancelled
	case bitmex.OS_REJECTED:
		return OrderStatusRejected
	default:
		return OrderStatusCreated
	}
}

func (b *BitMEX) SubscribeTrades(market Market, callback func(trades []Trade)) error {
	return ErrNotImplemented
}

func (b *BitMEX) SubscribeLevel2Snapshots(market Market, callback func(ob *OrderBook)) error {
	return ErrNotImplemented
}

func (b *BitMEX) SubscribeOrders(market Market, callback func(orders []Order)) error {
	return ErrNotImplemented
}

func (b *BitMEX) SubscribePositions(market Market, callback func(positions []Position)) error {
	return ErrNotImplemented
}

func (b *BitMEX) RunEventLoopOnce() (err error) {
	return
}

func NewBitMEX(params *Parameters) *BitMEX {
	baseUri := "www.bitmex.com"
	if params.Testnet {
		baseUri = "testnet.bitmex.com"
	}
	client := bitmex.New(baseUri, params.AccessKey, params.SecretKey)
	if strings.HasPrefix(params.ProxyURL, "socks5:") {
		socks5Proxy := strings.ReplaceAll(params.ProxyURL, "socks5:", "")
		client.SetProxy(socks5Proxy)
	} else if strings.HasPrefix(params.ProxyURL, "http://") {
		client.SetHttpProxy(params.ProxyURL)
	}
	return &BitMEX{
		client: client,
	}
}
