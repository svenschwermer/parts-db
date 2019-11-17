package si

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shopspring/decimal"
)

func Parse(s string) (*Quantity, error) {
	s = deleteSpaces(s)
	i := strings.IndexFunc(s, unicode.IsLetter)
	if i == -1 {
		return nil, errors.New("no prefix/unit found")
	}
	mag, err := decimal.NewFromString(s[:i])
	if err != nil {
		return nil, err
	}

	prefix, unit, err := splitPrefixUnit(s[i:])
	if err != nil {
		return nil, err
	}
	exp, err := getExponent(prefix)
	if err != nil {
		return nil, err
	}

	q := &Quantity{
		Mag:    mag.Mul(decimal.NewFromFloat(10).Pow(decimal.NewFromFloat(float64(exp)))),
		Coeff:  s[:i],
		Prefix: prefix,
		Unit:   unit,
		Valid:  true,
	}
	return q, nil
}

func deleteSpaces(s string) string {
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) && r != '\u200B' {
			result = append(result, r)
		}
	}
	return string(result)
}

func splitPrefixUnit(s string) (prefix, unit string, err error) {
	if len(s) == 1 {
		unit = s
	} else {
		p, i := utf8.DecodeRuneInString(s)
		prefix = string([]rune{p})
		unit = s[i:]
	}

	switch strings.ToLower(unit) {
	case "ohm", "r":
		unit = "Ω"
	}
	return
}

func getExponent(prefix string) (int32, error) {
	switch prefix {
	case "f":
		return -15, nil
	case "p":
		return -12, nil
	case "n":
		return -9, nil
	case "µ":
		return -6, nil
	case "m":
		return -3, nil
	case "", "\u200B":
		return 0, nil
	case "k":
		return 3, nil
	case "M":
		return 6, nil
	case "G":
		return 9, nil
	case "T":
		return 12, nil
	case "P":
		return 15, nil
	default:
		return 0, fmt.Errorf("exponent unknown for prefix %s", prefix)
	}
}
