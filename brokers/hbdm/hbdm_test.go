package hbdm

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
	"time"
)

func newForTest() Broker {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")
	return New(accessKey, secretKey, true)
}

func TestHBDM_GetAccountSummary(t *testing.T) {
	b := newForTest()
	accountSummary, err := b.GetAccountSummary("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", accountSummary)
}

func TestHBDM_GetOrderBook(t *testing.T) {
	b := newForTest()
	b.SetContractType("BTC", "W1")
	ob, err := b.GetOrderBook("BTC200327", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ob)
}

func TestHBDM_GetRecords(t *testing.T) {
	b := newForTest()
	b.SetContractType("BTC", ContractTypeW1)
	symbol, err := b.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", symbol)
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	records, err := b.GetRecords(symbol,
		"1m", start.Unix(), end.Unix(), 10)
	if err != nil {
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}

func TestHBDM_GetContractID(t *testing.T) {
	b := newForTest()
	b.SetContractType("BTC", ContractTypeW1)
	symbol, err := b.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", symbol)
}

func TestHBDM_GetOpenOrders(t *testing.T) {
	b := newForTest()
	b.SetContractType("BTC", ContractTypeW1)
	symbol, err := b.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("symbol: %v", symbol)

	orders, err := b.GetOpenOrders(symbol)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestHBDM_GetOrder(t *testing.T) {
	b := newForTest()
	b.SetContractType("BTC", ContractTypeW1)
	symbol, err := b.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}

	order, err := b.GetOrder(symbol, "694901372910391296")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestHBDM_PlaceOrder(t *testing.T) {
	b := newForTest()
	b.SetLeverRate(10)
	b.SetContractType("BTC", ContractTypeW1)
	symbol, err := b.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}

	order, err := b.PlaceOrder(symbol,
		Buy,
		OrderTypeLimit,
		3000,
		0,
		1,
		false,
		false,
		nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestHBDM_PlaceOrder2(t *testing.T) {
	b := newForTest()
	b.SetLeverRate(10)
	b.SetContractType("BTC", ContractTypeW1)
	symbol, err := b.GetContractID()
	if err != nil {
		t.Error(err)
		return
	}

	order, err := b.PlaceOrder(symbol,
		Sell,
		OrderTypeMarket,
		3000,
		0,
		1,
		false,
		true,
		nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}
