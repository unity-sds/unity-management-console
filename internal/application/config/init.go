package config

import (
	"time"

	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/retriever/httpretriever"
)

func InitApplication() (f *ffclient.GoFeatureFlag) {
	return initFeatureFlags()
}

func initFeatureFlags() (f *ffclient.GoFeatureFlag) {
	ff, err := ffclient.New(ffclient.Config{
		PollingInterval: 3 * time.Second,
		Retriever: &httpretriever.Retriever{
			URL:     "https://gist.githubusercontent.com/buggtb/7596f2079d66590c74db51996433ab39/raw/6709c4b01fe9360c0fa0e62f47e2498f33333e25/flag-file.yaml",
			Timeout: 2 * time.Second,
		},
	})
	defer ffclient.Close()

	if err != nil {

	}
	return ff
}
