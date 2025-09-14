package models

type Stock struct {
	Symbol   string `yaml:"symbol"`
	Exchange string `yaml:"exchange"`
	Segment  string `yaml:"segment"`
	Name     string `yaml:"name"`
}

type StocksConfig struct {
	Stocks []Stock `yaml:"stocks"`
}
