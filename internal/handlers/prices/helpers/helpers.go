package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`^-?\d+(?:[.]\d{1,2})?$`)

// парсинг цены с проверкой формата
func ParsePriceWithRegex(priceStr string) (float64, error) {
	priceStr = strings.TrimSpace(priceStr)

	if re.MatchString(priceStr) {
		return 0, fmt.Errorf("invalid price format: %s", priceStr)
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
