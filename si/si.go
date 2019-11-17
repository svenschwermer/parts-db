package si

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	bigNegOne = big.NewInt(-1)
	bigZero   = big.NewInt(0)
	bigTen    = big.NewInt(10)
)

// Quantity is a physical quantity
type Quantity struct {
	Mag    decimal.Decimal
	Coeff  string
	Prefix string
	Unit   string
	Valid  bool
}

// Scan implements sql.Scanner
func (q *Quantity) Scan(src interface{}) error {
	if src == nil {
		q.Valid = false
		return nil
	}

	err := q.Mag.Scan(src)
	if err != nil {
		return err
	}

	decPlaces := decPlaces(q.Mag.Coefficient())
	normExp := q.Mag.Exponent() + decPlaces - 1
	exp := normExp
	if exp < 0 {
		exp = ((exp - 2) / 3) * 3
	} else {
		exp = (exp / 3) * 3
	}

	q.Prefix, err = prefix(exp)
	if err != nil {
		return err
	}

	i := normExp - exp + 1
	q.Coeff = q.Mag.Coefficient().String()
	if i < decPlaces {
		q.Coeff = q.Coeff[0:i] + "." + q.Coeff[i:]
	} else {
		q.Coeff += strings.Repeat("0", int(i-decPlaces))
	}

	q.Valid = true
	return nil
}

func (q *Quantity) String() string {
	if !q.Valid {
		return ""
	}
	return fmt.Sprintf("%v\u202F%s%s", q.Coeff, q.Prefix, q.Unit)
}

func decPlaces(x *big.Int) int32 {
	var result int32
	x = &*x
	if x.Sign() == -1 {
		x.Mul(x, bigNegOne)
	}
	for x.Cmp(bigZero) != 0 {
		result++
		x.Div(x, bigTen)
	}
	return result
}

func prefix(exp int32) (string, error) {
	switch exp {
	case -15:
		return "f", nil
	case -12:
		return "p", nil
	case -9:
		return "n", nil
	case -6:
		return "µ", nil
	case -3:
		return "m", nil
	case 0:
		return "", nil // zero width space
	case 3:
		return "k", nil
	case 6:
		return "M", nil
	case 9:
		return "G", nil
	case 12:
		return "T", nil
	case 15:
		return "P", nil
	default:
		return "", fmt.Errorf("prefix unknown for exponent %d", exp)
	}
}
