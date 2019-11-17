package si

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuantityScan(t *testing.T) {
	data := []struct {
		input     interface{}
		magnitude string
		prefix    string
	}{
		{1.0, "1", ""},
		{12.0, "12", ""},
		{123.0, "123", ""},
		{1234.0, "1.234", "k"},
		{12345.0, "12.345", "k"},
		{0.000000000001, "1", "p"},
		{0.00000000001, "10", "p"},
		{0.0000000001, "100", "p"},
		{0.000000001, "1", "n"},
		{0.00000001, "10", "n"},
		{0.0000001, "100", "n"},
		{0.000001, "1", "µ"},
		{0.00001, "10", "µ"},
		{0.0001, "100", "µ"},
		{0.001, "1", "m"},
		{0.01, "10", "m"},
		{0.1, "100", "m"},
		{1.0, "1", ""},
		{10.0, "10", ""},
		{100.0, "100", ""},
		{1000.0, "1", "k"},
		{10000.0, "10", "k"},
		{100000.0, "100", "k"},
	}

	for _, d := range data {
		uut := new(Quantity)
		err := uut.Scan(d.input)
		assert.NoError(t, err)
		assert.Equal(t, d.magnitude, uut.Coeff)
		assert.Equal(t, d.prefix, uut.Prefix)
	}
}
