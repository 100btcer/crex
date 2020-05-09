package backtest

import (
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/dataloader"
	"github.com/coinrust/crex/log"
	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/datatypes"
	"os"
	"path/filepath"
	"time"
)

const (
	OriginEChartsJs = "https://go-echarts.github.io/go-echarts-assets/assets/echarts.min.js"
	MyEChartsJs     = "https://cdnjs.cloudflare.com/ajax/libs/echarts/4.7.0/echarts.min.js"

	OriginEChartsBulmaCss = "https://go-echarts.github.io/go-echarts-assets/assets/bulma.min.css"
	MyEChartsBulmaCss     = "https://cdnjs.cloudflare.com/ajax/libs/bulma/0.8.2/css/bulma.min.css"
)

type PlotData struct {
	NameItems []string
	Prices    []float64
	Equities  []float64
}

type Backtest struct {
	data      *dataloader.Data
	symbol    string
	strategy  Strategy
	exchanges []ExchangeSim
	outputDir string
	logs      LogItems
}

const SimpleDateTimeFormat = "2006-01-02 15:04:05.000"

// NewBacktest Create backtest
// data: The data
// outputDir: 日志输出目录
func NewBacktest(data *dataloader.Data, symbol string, strategy Strategy, exchanges []ExchangeSim, outputDir string) *Backtest {
	b := &Backtest{
		data:      data,
		symbol:    symbol,
		strategy:  strategy,
		outputDir: outputDir,
	}
	b.exchanges = exchanges
	var exs []Exchange
	for _, v := range exchanges {
		exs = append(exs, v)
	}
	strategy.Setup(TradeModeBacktest, exs...)
	b.logs = LogItems{}

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	logger := NewBtLogger(b,
		filepath.Join(outputDir, "result.log"),
		log.DebugLevel,
		false)
	log.SetLogger(logger)

	return b
}

// SetData Set data for backtest
func (b *Backtest) SetData(data *dataloader.Data) {
	b.data = data
}

// GetTime get current time
func (b *Backtest) GetTime() time.Time {
	if b.data == nil {
		return time.Now()
	}
	return b.data.GetOrderBook().Time
}

// Run Run backtest
func (b *Backtest) Run() {
	b.data.Reset()

	// Init
	b.strategy.OnInit()

	for {
		b.strategy.OnTick()

		b.runEventLoopOnce()

		b.addItemStats()

		if !b.data.Next() {
			break
		}
	}

	// Exit
	b.strategy.OnExit()

	log.Sync()
}

func (b *Backtest) runEventLoopOnce() {
	for _, exchange := range b.exchanges {
		exchange.RunEventLoopOnce()
	}
}

func (b *Backtest) addItemStats() {
	ob := b.data.GetOrderBook()
	tm := ob.Time
	update := false
	timestamp := time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute()+1, 0, 0, time.UTC)
	var lastItem *LogItem

	if len(b.logs) > 0 {
		lastItem = b.logs[len(b.logs)-1]
		if timestamp.Unix() == lastItem.Time.Unix() {
			update = true
			return
		}
	}
	var item *LogItem
	if update {
		item = lastItem
		item.RawTime = ob.Time
		item.Ask = ob.AskPrice()
		item.Bid = ob.BidPrice()
		item.Stats = nil
		b.fetchItemStats(item)
	} else {
		item = &LogItem{
			Time:    timestamp,
			RawTime: ob.Time,
			Ask:     ob.AskPrice(),
			Bid:     ob.BidPrice(),
			Stats:   nil,
		}
		b.fetchItemStats(item)
		b.logs = append(b.logs, item)
		//log.Printf("%v / %v", tick.Timestamp, timestamp)
	}
}

func (b *Backtest) fetchItemStats(item *LogItem) {
	n := len(b.exchanges)
	for i := 0; i < n; i++ {
		balance, err := b.exchanges[i].GetBalance("BTC")
		if err != nil {
			panic(err)
		}
		item.Stats = append(item.Stats, LogStats{
			Balance: balance.Available,
			Equity:  balance.Equity,
		})
	}
}

func (b *Backtest) GetLogs() LogItems {
	return b.logs
}

// ComputeStats Calculating Backtest Statistics
func (b *Backtest) ComputeStats() (result *Stats) {
	result = &Stats{}

	if len(b.logs) == 0 {
		return
	}

	logs := b.logs

	n := len(logs)

	result.Start = logs[0].Time
	result.End = logs[n-1].Time
	result.Duration = result.End.Sub(result.Start)
	result.EntryPrice = logs[0].Price()
	result.ExitPrice = logs[n-1].Price()
	result.EntryEquity = logs[0].TotalEquity()
	result.ExitEquity = logs[n-1].TotalEquity()
	result.BaHReturn = (result.ExitPrice - result.EntryPrice) / result.EntryPrice * result.EntryEquity
	result.BaHReturnPnt = (result.ExitPrice - result.EntryPrice) / result.EntryPrice
	result.EquityReturn = result.ExitEquity - result.EntryEquity
	result.EquityReturnPnt = result.EquityReturn / result.EntryEquity

	return
}

func (b *Backtest) priceLine(plotData *PlotData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "价格", Width: "1270px", Height: "500px"},
		charts.ToolboxOpts{Show: true},
		charts.TooltipOpts{Show: true, Trigger: "axis", TriggerOn: "mousemove|click"},
		charts.DataZoomOpts{Type: "slider", Start: 0, End: 100},
		charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}, Scale: true},
	)

	line.AddXAxis(plotData.NameItems)
	line.AddYAxis("price", plotData.Prices,
		charts.MPNameTypeItem{Name: "最大值", Type: "max"},
		charts.MPNameTypeItem{Name: "最小值", Type: "min"},
		charts.MPStyleOpts{Label: charts.LabelTextOpts{Show: true}},
		//charts.LineOpts{Smooth: true, YAxisIndex: 0},
	)

	return line
}

func (b *Backtest) equityLine(plotData *PlotData) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "净值", Width: "1270px", Height: "400px"},
		charts.ToolboxOpts{Show: true},
		charts.TooltipOpts{Show: true, Trigger: "axis", TriggerOn: "mousemove|click"},
		charts.DataZoomOpts{Type: "slider", Start: 0, End: 100},
		charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}, Scale: true},
	)

	line.AddXAxis(plotData.NameItems)

	line.AddYAxis("equity", plotData.Equities,
		charts.MPNameTypeItem{Name: "最大值", Type: "max"},
		charts.MPNameTypeItem{Name: "最小值", Type: "min"},
		charts.MPStyleOpts{Label: charts.LabelTextOpts{Show: true}},
		//charts.LineOpts{Smooth: true, YAxisIndex: 0},
	)

	return line
}

// Plot Output backtest results
func (b *Backtest) Plot() {
	var plotData PlotData

	for _, v := range b.logs {
		plotData.NameItems = append(plotData.NameItems, v.Time.Format(SimpleDateTimeFormat))
		plotData.Prices = append(plotData.Prices, v.Price())
		plotData.Equities = append(plotData.Equities, v.TotalEquity())
	}

	p := charts.NewPage()
	p.Add(b.priceLine(&plotData), b.equityLine(&plotData))

	filename := filepath.Join(b.outputDir, "result.html")
	f, err := os.Create(filename)
	if err != nil {
		log.Error(err)
	}

	replaceJSAssets(&p.JSAssets)
	replaceCssAssets(&p.CSSAssets)

	p.Render(f)
}

// 替换Js资源，使用cdn加速资源，查看网页更快
func replaceJSAssets(jsAssets *datatypes.OrderedSet) {
	for i := 0; i < len(jsAssets.Values); i++ {
		if jsAssets.Values[i] == OriginEChartsJs {
			jsAssets.Values[i] = MyEChartsJs
		}
	}
}

func replaceCssAssets(cssAssets *datatypes.OrderedSet) {
	for i := 0; i < len(cssAssets.Values); i++ {
		if cssAssets.Values[i] == OriginEChartsBulmaCss {
			cssAssets.Values[i] = MyEChartsBulmaCss
		}
	}
}

func (b *Backtest) PlotOld() {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "回测", Width: "1270px", Height: "600px"},
		charts.ToolboxOpts{Show: true},
		charts.ToolboxOpts{Show: true},
		charts.TitleOpts{Title: "回测"},
		charts.TooltipOpts{Show: true, Trigger: "axis", TriggerOn: "mousemove|click"},
		charts.DataZoomOpts{Type: "slider", Start: 0, End: 100},
		//charts.LegendOpts{Right: "80%"},
		//charts.SplitLineOpts{Show: true},
		//charts.SplitAreaOpts{Show: true},
	)
	nameItems := make([]string, 0)
	prices := make([]float64, 0)
	equities := make([]float64, 0)

	for _, v := range b.logs {
		nameItems = append(nameItems, v.Time.Format(SimpleDateTimeFormat))
		prices = append(prices, v.Price())
		equities = append(equities, v.TotalEquity())
	}

	line.AddXAxis(nameItems)
	line.AddYAxis("price", prices,
		charts.MPNameTypeItem{Name: "最大值", Type: "max"},
		charts.MPNameTypeItem{Name: "最小值", Type: "min"},
		charts.MPStyleOpts{Label: charts.LabelTextOpts{Show: true}},
	//charts.LineOpts{Smooth: true, YAxisIndex: 0},
	)

	line.AddYAxis("equity", equities,
		charts.MPNameTypeItem{Name: "最大值", Type: "max"},
		charts.MPNameTypeItem{Name: "最小值", Type: "min"},
		charts.MPStyleOpts{Label: charts.LabelTextOpts{Show: true}},
	//charts.LineOpts{Smooth: true, YAxisIndex: 0},
	)

	line.SetGlobalOptions(charts.YAxisOpts{SplitLine: charts.SplitLineOpts{Show: true}, Scale: true})

	filename := filepath.Join(b.outputDir, "result.html")
	f, err := os.Create(filename)
	if err != nil {
		log.Error(err)
	}
	line.Render(f)
}
