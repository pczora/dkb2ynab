package config

type DkbConfig struct {
	DateMapping   int
	PayeeMapping  int
	MemoMapping   int
	AmountMapping int
}

func NewDkbConfig() Config {
	config := Config{
		SkipLines:       7,
		DateColumn:      1,
		PayeeColumn:     3,
		MemoColumn:      4,
		AmountColumn:    7,
		NormalizeAmount: true,
		DateFormat:      "02.01.2006",
	}
	return config
}
