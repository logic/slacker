package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTickers(t *testing.T) {
	tickerHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"quoteResponse\":{\"result\":[{\"symbol\":\"TEST\",\"longName\":\"Test\"},{\"symbol\":\"MOAR\",\"longName\":\"Moar\"}]}}")
	}

	ts := httptest.NewServer(http.HandlerFunc(tickerHandler))
	defer ts.Close()

	apiYahooFinance = ts.URL
	tickers, err := GetTickers([]string{"TEST", "TEST1"})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if tickers[0].Symbol != "TEST" || tickers[0].LongName != "Test" ||
		tickers[1].Symbol != "MOAR" || tickers[1].LongName != "Moar" {
		t.Errorf("result mismatch")
	}
}

func TestYahooDown(t *testing.T) {
	apiYahooFinance = "http://127.0.0.1:8888/"
	_, err := GetTickers([]string{"TEST"})
	if err == nil {
		t.Errorf("GetTickers didn't catch unreachable API endpoint")
	}
}

func TestYahooErrorStatus(t *testing.T) {
	tickerHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error", http.StatusBadRequest)
	}

	ts := httptest.NewServer(http.HandlerFunc(tickerHandler))
	defer ts.Close()

	apiYahooFinance = ts.URL
	_, err := GetTickers([]string{"TEST"})
	if err == nil {
		t.Errorf("GetTickers didn't catch error from API endpoint")
	}
}

func TestYahooErrorContentType(t *testing.T) {
	tickerHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "error")
	}

	ts := httptest.NewServer(http.HandlerFunc(tickerHandler))
	defer ts.Close()

	apiYahooFinance = ts.URL
	_, err := GetTickers([]string{"TEST"})
	if err == nil {
		t.Errorf("GetTickers didn't catch bad content type")
	}
}

func TestYahooBadRecordCount(t *testing.T) {
	tickerHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		fmt.Fprint(w, "this is a failure\nso bad\n")
	}

	ts := httptest.NewServer(http.HandlerFunc(tickerHandler))
	defer ts.Close()

	apiYahooFinance = ts.URL
	_, err := GetTickers([]string{"TEST"})
	if err == nil {
		t.Errorf("GetTickers didn't catch wrong number of records")
	}
}

func TestYahooBadFieldCount(t *testing.T) {
	tickerHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		fmt.Fprint(w, "this is a failure\nso bad\n")
	}

	ts := httptest.NewServer(http.HandlerFunc(tickerHandler))
	defer ts.Close()

	apiYahooFinance = ts.URL
	_, err := GetTickers([]string{"TEST", "TEST1"})
	if err == nil {
		t.Errorf("GetTickers didn't catch wrong number of fields")
	}
}

func TestYahooBadURL(t *testing.T) {
	apiYahooFinance = ":"
	_, err := GetTickers(nil)
	if err == nil {
		t.Errorf("GetTickers didn't catch invalid URL")
	}
}
