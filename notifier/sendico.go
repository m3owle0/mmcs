package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"golang.org/x/sync/errgroup"
)

const (
	SendicoBaseURL = "https://sendico.com"
)

type SendicoClient struct {
	httpClient *http.Client
	mu         sync.RWMutex
	hmacSecret string
	baseURL    string
}

func NewSendicoClient() (*SendicoClient, error) {
	client := &SendicoClient{
		httpClient: http.DefaultClient,
		baseURL:    SendicoBaseURL,
	}

	ctx := context.Background()
	if err := client.FindHMAC(ctx); err != nil {
		return nil, fmt.Errorf("failed to find HMAC: %w", err)
	}

	return client, nil
}

func (c *SendicoClient) FindHMAC(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	selection := doc.Find("script#__NUXT_DATA__")
	if selection.Length() == 0 {
		return errors.New("script tag not found")
	}

	var unstruct []any
	if err := json.Unmarshal([]byte(selection.Nodes[0].FirstChild.Data), &unstruct); err != nil {
		return err
	}

	ptr := int64(-1)
	for _, obj := range unstruct {
		switch v := obj.(type) {
		case map[string]any:
			if val, ok := v["$sapi_tokens"]; ok {
				ptr = int64(val.(float64))
				break
			}
		}
	}

	if ptr == -1 {
		return errors.New("unable to find reference to secret key")
	}

	keyPtrs := unstruct[ptr].([]any)
	secretKeys := make([]string, len(keyPtrs))
	for i, keyPtr := range keyPtrs {
		secretKeys[i] = unstruct[int64(keyPtr.(float64))].(string)
	}

	if len(secretKeys) == 0 {
		return errors.New("no secret keys found")
	}

	newSecret := decodeHMACKey(secretKeys[len(secretKeys)-1])
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hmacSecret = newSecret
	return nil
}

func (c *SendicoClient) HMACSecret() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hmacSecret
}

func (c *SendicoClient) Translate(ctx context.Context, text string) (string, error) {
	path := "/api/translate"

	request := orderedmap.New[string, any]()
	request.Set("from", "en")
	request.Set("string", text)
	request.Set("to", "ja")

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	hmac, err := buildHMAC(HMACInput{
		Secret:  c.HMACSecret(),
		Path:    path,
		Payload: request,
		Nonce:   "",
	})
	if err != nil {
		return "", err
	}

	resp, err := c.req(ctx, http.MethodPost, path, bytes.NewReader(requestJSON), hmac, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json")
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	response := struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.Data, nil
}

type SendicoShop string

const (
	SendicoYahooAuctions SendicoShop = "ayahoo"
	SendicoMercari       SendicoShop = "mercari"
	SendicoRakuma        SendicoShop = "rakuma"
	SendicoRakuten        SendicoShop = "rakuten"
	SendicoYahoo          SendicoShop = "yahoo"
)

type SendicoSearchOptions struct {
	TermJP   string
	MinPrice *int
	MaxPrice *int
}

type SendicoItem struct {
	Shop     SendicoShop `json:"shop"`
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	URL      string      `json:"url"`
	Image    string      `json:"img"`
	PriceYen int         `json:"price"`
	PriceUSD int         `json:"converted_price"`
	Labels   []string    `json:"labels"`
}

func (c *SendicoClient) Search(ctx context.Context, shop SendicoShop, opts SendicoSearchOptions) ([]SendicoItem, error) {
	path := url.URL{
		Path: fmt.Sprintf("/api/%s/items", shop),
	}

	params := orderedmap.New[string, any]()
	params.Set("global", "1")
	if opts.MaxPrice != nil {
		params.Set("max_price", fmt.Sprintf("%d", *opts.MaxPrice))
	}
	if opts.MinPrice != nil {
		params.Set("min_price", fmt.Sprintf("%d", *opts.MinPrice))
	}
	params.Set("page", "1")
	params.Set("search", opts.TermJP)

	q := path.Query()
	for pair := params.Oldest(); pair != nil; pair = pair.Next() {
		q.Add(pair.Key, fmt.Sprintf("%v", pair.Value))
	}
	path.RawQuery = q.Encode()

	hmac, err := buildHMAC(HMACInput{
		Secret:  c.HMACSecret(),
		Path:    path.Path,
		Payload: params,
		Nonce:   "",
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.req(ctx, http.MethodGet, path.String(), nil, hmac, func(req *http.Request) {
		req.Header.Set("Content-Type", "application/json")
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := struct {
		Code int `json:"code"`
		Data struct {
			Items      []SendicoItem `json:"items"`
			TotalItems int           `json:"total_items"`
		} `json:"data"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Data.Items, nil
}

func (c *SendicoClient) BulkSearch(ctx context.Context, shops []SendicoShop, opts SendicoSearchOptions) ([]SendicoItem, error) {
	items := make([]SendicoItem, 0)
	itemsMu := sync.Mutex{}

	// Use semaphore to limit concurrent requests (max 5 at a time for rate limiting)
	maxConcurrent := 5
	sem := make(chan struct{}, maxConcurrent)
	g := new(errgroup.Group)
	
	for i, shop := range shops {
		shop := shop
		i := i
		g.Go(func() error {
			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()
			
			// Reduced delay - only 200ms between requests (faster but still respectful)
			if i > 0 {
				time.Sleep(200 * time.Millisecond)
			}
			
			results, err := c.Search(ctx, shop, opts)
			if err != nil {
				log.Printf("   ⚠️  Error searching %s: %v", shop, err)
				return nil // Continue with other shops
			}

			itemsMu.Lock()
			defer itemsMu.Unlock()
			items = append(items, results...)
			return nil
		})
	}

	_ = g.Wait() // Ignore errors, we log them above
	return items, nil
}

func (c *SendicoClient) req(ctx context.Context, method, path string, body io.Reader, hmac *HMACAttributes, opts ...func(*http.Request)) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(req)
	}

	if hmac != nil {
		req.Header.Set("X-Sendico-Signature", hmac.Signature)
		req.Header.Set("X-Sendico-Nonce", hmac.Nonce)
		req.Header.Set("X-Sendico-Timestamp", fmt.Sprintf("%d", hmac.Timestamp))
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		log.Printf("   ⚠️  Sendico API error: status %d, body: %s", res.StatusCode, string(body))
		_ = res.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res, err
}
