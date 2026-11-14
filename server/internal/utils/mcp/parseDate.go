package mcp

import (
	"fmt"
	"time"
)

// DateLayout — формат дней в аргументах и ответах инструментов
const DateLayout = "2006-01-02"

// ParseOptionalDate разбирает необязательный день yyyy-MM-dd
func ParseOptionalDate(field string, value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	date, err := time.Parse(DateLayout, *value)
	if err != nil {
		return nil, fmt.Errorf("%s: ожидается дата yyyy-MM-dd, получено %q", field, *value)
	}
	return &date, nil
}
