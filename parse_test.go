package xml2json

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Product struct {
	ID      int     `json:"id"`
	Price   float64 `json:"price"`
	Deleted bool    `json:"deleted"`
}

const (
	productString = `
	<?xml version="1.0" encoding="UTF-8"?>	
		<id>42</id>
		<price>13.32</price>
		<deleted>true</deleted>
		`
)

func TestJSTypeParsing(t *testing.T) {
	xml := strings.NewReader(productString)
	jsBuf, err := Convert(xml)
	assert.NoError(t, err, "could not parse test xml")
	product := Product{}
	err = json.Unmarshal(jsBuf.Bytes(), &product)
	assert.NoError(t, err, "could not unmarshal test json")
	assert.Equal(t, 42, product.ID, "price should match")
	assert.Equal(t, 13.32, product.Price, "price should match")
	assert.Equal(t, true, product.Deleted, "price should match")
}
