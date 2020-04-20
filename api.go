package crex

import "net/http"

type Parameters struct {
	DebugMode  bool
	HttpClient *http.Client
	ProxyURL   string // socks5://127.0.0.1:1080 | http://127.0.0.1:1080
	AccessKey  string
	SecretKey  string
	Passphrase string
	Testnet    bool
	WebSocket  bool // Enable websocket option
}

type ApiOption func(p *Parameters)

func ApiDebugModeOption(debugMode bool) ApiOption {
	return func(p *Parameters) {
		p.DebugMode = debugMode
	}
}

func ApiHttpClientOption(httpClient *http.Client) ApiOption {
	return func(p *Parameters) {
		p.HttpClient = httpClient
	}
}

func ApiProxyURLOption(proxyURL string) ApiOption {
	return func(p *Parameters) {
		p.ProxyURL = proxyURL
	}
}

func ApiAccessKeyOption(accessKey string) ApiOption {
	return func(p *Parameters) {
		p.AccessKey = accessKey
	}
}

func ApiSecretKeyOption(secretKey string) ApiOption {
	return func(p *Parameters) {
		p.SecretKey = secretKey
	}
}

func ApiPassPhraseOption(passPhrase string) ApiOption {
	return func(p *Parameters) {
		p.Passphrase = passPhrase
	}
}

func ApiTestnetOption(testnet bool) ApiOption {
	return func(p *Parameters) {
		p.Testnet = testnet
	}
}

func ApiWebSocketOption(enabled bool) ApiOption {
	return func(p *Parameters) {
		p.WebSocket = enabled
	}
}

type OrderParameter struct {
	StopPx     float64
	PostOnly   bool
	ReduceOnly bool
	PriceType  string
}

// 订单选项
type OrderOption func(p *OrderParameter)

// 触发价格选项
func OrderStopPxOption(stopPx float64) OrderOption {
	return func(p *OrderParameter) {
		p.StopPx = stopPx
	}
}

// 被动委托选项
func OrderPostOnlyOption(postOnly bool) OrderOption {
	return func(p *OrderParameter) {
		p.PostOnly = postOnly
	}
}

// 只减仓选项
func OrderReduceOnlyOption(reduceOnly bool) OrderOption {
	return func(p *OrderParameter) {
		p.ReduceOnly = reduceOnly
	}
}

// OrderPriceType 选项
func OrderPriceTypeOption(priceType string) OrderOption {
	return func(p *OrderParameter) {
		p.PriceType = priceType
	}
}

// Exchange 交易所接口
type Exchange interface {
	WebSocket

	// 获取 Exchange 名称
	GetName() (name string)

	// 获取账号余额
	GetBalance(currency string) (result Balance, err error)

	// 获取订单薄(OrderBook)
	GetOrderBook(symbol string, depth int) (result OrderBook, err error)

	// 获取K线数据
	// period: 数据周期. 分钟或者关键字1m(minute) 1h 1d 1w 1M(month) 1y 枚举值：1 3 5 15 30 60 120 240 360 720 "5m" "4h" "1d" ...
	GetRecords(symbol string, period string, from int64, end int64, limit int) (records []Record, err error)

	// 设置合约类型
	// currencyPair: 交易对，如: BTC-USD(OKEX) BTC(HBDM)
	// contractType: W1,W2,Q1,Q2
	SetContractType(currencyPair string, contractType string) (err error)

	// 获取当前设置的合约ID
	GetContractID() (symbol string, err error)

	// 设置杠杆大小
	SetLeverRate(value float64) (err error)

	// 开多
	OpenLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error)

	// 开空
	OpenShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error)

	// 平多
	CloseLong(symbol string, orderType OrderType, price float64, size float64) (result Order, err error)

	// 平空
	CloseShort(symbol string, orderType OrderType, price float64, size float64) (result Order, err error)

	// 下单
	PlaceOrder(symbol string, direction Direction, orderType OrderType, price float64, size float64,
		opts ...OrderOption) (result Order, err error)

	// 获取活跃委托单列表
	GetOpenOrders(symbol string) (result []Order, err error)

	// 获取委托信息
	GetOrder(symbol string, id string) (result Order, err error)

	// 撤销全部委托单
	CancelAllOrders(symbol string) (err error)

	// 撤销单个委托单
	CancelOrder(symbol string, id string) (result Order, err error)

	// 修改委托
	AmendOrder(symbol string, id string, price float64, size float64) (result Order, err error)

	// 获取持仓
	GetPositions(symbol string) (result []Position, err error)

	// 运行一次(回测系统调用)
	RunEventLoopOnce() (err error) // Run sim match for backtest only
}

func ParseOrderParameter(opts ...OrderOption) *OrderParameter {
	p := &OrderParameter{
		StopPx:     0,
		PostOnly:   false,
		ReduceOnly: false,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}
