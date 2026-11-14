package network

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestRequestCurrencyRates_ParsesProviderResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("date") != "2020-03-15" || r.URL.Query().Get("apikey") != "key" {
			http.Error(w, "bad query", http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`{"meta":{"last_updated_at":"2020-03-15T23:59:59Z"},"data":{"RUB":{"code":"RUB","value":74.123456},"USD":{"code":"USD","value":1}}}`))
	}))
	defer server.Close()

	rates, validAt, err := requestCurrencyRates(context.Background(), server.URL+"/v3/historical", url.Values{"apikey": {"key"}, "date": {"2020-03-15"}})
	if err != nil {
		t.Fatal(err)
	}
	if rates["RUB"].String() != "74.123456" || rates["USD"].String() != "1" {
		t.Fatalf("курсы: %v", rates)
	}
	if !validAt.Equal(time.Date(2020, 3, 15, 23, 59, 59, 0, time.UTC)) {
		t.Fatalf("момент: %v", validAt)
	}

	// Ошибка провайдера (например, кончилась квота) возвращается ошибкой
	if _, _, err = requestCurrencyRates(context.Background(), server.URL, url.Values{}); err == nil {
		t.Fatal("ошибка провайдера не вернулась")
	}
}
