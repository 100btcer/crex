package bybit_broker

import (
	. "github.com/coinrust/crex"
	"testing"
	"time"
)

func newForTest() Broker {
	b := NewBroker("https://api-testnet.bybit.com/",
		"6IASD6KDBdunn5qLpT", "nXjZMUiB3aMiPaQ9EUKYFloYNd0zM39RjRWF")
	return b
}

func TestBybitBroker_GetOpenOrders(t *testing.T) {
	b := newForTest()
	orders, err := b.GetOpenOrders("BTCUSD")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}

func TestBybitBroker_GetRecords(t *testing.T) {
	b := newForTest()
	start := time.Now().Add(-time.Hour)
	records, err := b.GetRecords("BTCUSD",
		"1", start.Unix(), 0, 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}
