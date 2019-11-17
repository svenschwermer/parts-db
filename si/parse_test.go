package si

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteSpaces(t *testing.T) {
	data := []struct {
		in       string
		expected string
	}{
		{"", ""},
		{"\t1.234\u202FkΩ\r", "1.234kΩ"},
		{"1\u202F\u200BF", "1F"},
	}

	for _, d := range data {
		out := deleteSpaces(d.in)
		assert.Equal(t, d.expected, out)
	}
}

func TestParse(t *testing.T) {
	q, err := Parse("20R")
	require.NoError(t, err)
	assert.Equal(t, decimal.NewFromFloat(20), q.Mag)
}
