package entity

type Currency struct {
	ID       int64  `db:"id"`
	Code     string `db:"code"`
	FullName string `db:"name"`
	Sign     string `db:"sign"`
}

const (
	CurrencyCodeMaxLen     = 3
	CurrencySignMaxLen     = 3
	CurrencyFullNameMinLen = 3
	CurrencyFullNameMaxLen = 40
)
