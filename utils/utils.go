package utils

import (
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// FormatCurrency formats a float64 value to BRL currency
func FormatCurrency(value float64) string {
	p := message.NewPrinter(language.BrazilianPortuguese)
	return p.Sprintf("%s %.2f", currency.BRL.String(), value)
}
