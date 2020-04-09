package hbdm_swap_broker

import (
	. "github.com/coinrust/crex"
	"github.com/spf13/viper"
	"log"
	"testing"
	"time"
)

func newTestBroker() Broker {
	viper.SetConfigName("test_config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")
	baseURL := "https://api.btcgateway.pro"
	return NewBroker(baseURL, accessKey, secretKey)
}

func TestHBDMSwapBroker_GetRecords(t *testing.T) {
	b := newTestBroker()
	symbol := "BTC-USD"
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	records, err := b.GetRecords(symbol,
		"1m", start.Unix(), end.Unix(), 10)
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range records {
		t.Logf("%#v", v)
	}
}
