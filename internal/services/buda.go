package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mtavano/buda-go"
	"github.com/pkg/errors"
)

type Balance struct {
	ID                    string   `json:"id"`
	Amount                []string `json:"amount"`
	AvailableAmount       []string `json:"available_amount"`
	FrozenAmount          []string `json:"frozen_amount"`
	PendingWithdrawAmount []string `json:"pending_withdraw_amount"`
	AccountID             int      `json:"account_id"`
	TotalFiatAmount       float64  `json:"total_fiat_amount"`
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Buda struct {
	secret  string
	key     string
	baseURL string
	buda    *buda.Buda

	client HTTPClient
}

func NewBuda(baseURL, apiKey, apiSecret string) *Buda {
	return &Buda{
		secret:  apiSecret,
		key:     apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (b *Buda) GetBalance() ([]*Balance, error) {
	res, err := b.makeRequest(http.MethodGet, "/balances", nil, true)
	if err != nil {
		return nil, err
	}

	payload := &struct {
		Balances []*Balance `json:"balances"`
	}{}
	err = b.scanBody(res, payload)
	if err != nil {
		return nil, err
	}

	for _, balance := range payload.Balances {
		if strings.ToLower(balance.ID) == "clp" {
			continue
		}
		ticker, err := b.GetTicker(strings.ToLower(balance.ID) + "-clp")
		if err != nil {
			return nil, err
		}
		balance.TotalFiatAmount, err = strconv.ParseFloat(ticker.LastPrice[0], 64)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(balance.Amount[0], 64)
		if err != nil {
			return nil, err
		}
		balance.TotalFiatAmount = balance.TotalFiatAmount * amount
	}

	return payload.Balances, nil
}

// Ticker ...
type Ticker struct {
	LastPrice         []string `json:"last_price"`
	MaxBid            []string `json:"max_bid"`
	MinAsk            []string `json:"min_ask"`
	PriceVariation24H string   `json:"price_variation_24h"`
	PriceVariation7D  string   `json:"price_variation_7d"`
	Volume            []string `json:"volume"`
}

func (b *Buda) GetTicker(pair string) (*Ticker, error) {
	if pair == "" {
		return nil, errors.New("market pair cannot be empty")
	}

	url := fmt.Sprintf("/markets/%s/ticker", pair)
	res, err := b.makeRequest(http.MethodGet, url, nil, false)
	if err != nil {
		return nil, err
	}

	payload := &struct {
		Ticker *Ticker `json:"ticker"`
	}{}

	err = b.scanBody(res, payload)
	if err != nil {
		return nil, err
	}

	return payload.Ticker, nil
}

func (b *Buda) makeRequest(method, path string, body io.Reader, private bool) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", b.baseURL, path)
	fmt.Println(">>> url", url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "buda: Buda.makeRequest http.NewRequest error")
	}
	req.Header.Set("Content-Type", "application/json")

	if private {
		err = b.authenticate(req)
		if err != nil {
			return nil, errors.Wrap(err, "buda: authenticateRequest error")
		}
	}

	response, err := b.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "buda: httpClient.Do error")
	}

	// Leemos el body original
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "buda: io.ReadAll error")
	}

	// Creamos un nuevo buffer con el contenido
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	// Reemplazamos el body original con un NopCloser
	response.Body = io.NopCloser(bodyBuffer)

	// Ahora podemos leer el contenido para logging
	fmt.Println(string(bodyBytes))

	return response, nil
}

func (b *Buda) scanBody(res *http.Response, scanner interface{}) error {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, scanner)
}

func (b *Buda) authenticate(req *http.Request) error {
	valArray, err := createValArray(req)
	if err != nil {
		return nil
	}

	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	valArray = append(valArray, nonce)
	sign := createSign(valArray, b.secret)

	req.Header.Add("X-SBTC-APIKEY", b.key)
	req.Header.Add("X-SBTC-NONCE", nonce)
	req.Header.Add("X-SBTC-SIGNATURE", sign)

	return nil
}

func createValArray(req *http.Request) ([]string, error) {
	var params []string

	params = append(params, req.Method)
	params = append(params, req.URL.RequestURI())

	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		b := req.Body
		body, err := io.ReadAll(b)
		if err != nil {
			return nil, err
		}
		params = append(params, base64.StdEncoding.EncodeToString(body))
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	return params, nil
}

func createSign(valArray []string, secret string) string {
	h := hmac.New(sha512.New384, []byte(secret))
	rawStr := strings.Join(valArray, " ")
	h.Write([]byte(rawStr))
	return hex.EncodeToString(h.Sum(nil))
}
