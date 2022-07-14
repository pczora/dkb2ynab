package config

type Header int

type Config struct {
	SkipLines       int
	DateColumn      int
	PayeeColumn     int
	MemoColumn      int
	AmountColumn    int
	NormalizeAmount bool
	DateFormat      string
}
