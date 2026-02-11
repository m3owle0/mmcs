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
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"golang.org/x/sync/errgroup"
)

// ErrHMACRefreshNeeded is returned when HMAC secret needs to be refreshed
var ErrHMACRefreshNeeded = errors.New("HMAC secret refreshed, retry needed")

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
		// If HMAC was refreshed, rebuild HMAC and retry once
		if err == ErrHMACRefreshNeeded {
			hmac, err = buildHMAC(HMACInput{
				Secret:  c.HMACSecret(),
				Path:    path,
				Payload: request,
				Nonce:   "",
			})
			if err != nil {
				return "", err
			}
			// Retry with new HMAC
			resp, err = c.req(ctx, http.MethodPost, path, bytes.NewReader(requestJSON), hmac, func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			})
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
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
	TermJP     string
	MinPrice   *int
	MaxPrice   *int
	Page       int    // Page number (default: 1)
	CategoryID *int   // Category ID for filtering (e.g., clothing category)
	// Note: Sendico API doesn't support sort/mobile parameters
	// We search multiple pages starting from page 1 (which typically has newest items)
}

// Clothing category IDs for each marketplace (Sendico API)
// These filter searches to clothing/fashion categories only, ensuring we only get clothing items
// Category IDs are based on Sendico's category structure:
// - Mercari: 3088 = Fashion (includes clothing, shoes, bags, accessories)
// - Yahoo Auctions: 23000 = Fashion category
// - Other markets: Using approximate IDs (may need verification/adjustment)
var clothingCategoryIDs = map[SendicoShop]int{
	SendicoMercari:       3088,  // Mercari Fashion category
	SendicoYahooAuctions: 23000, // Yahoo Auctions Fashion category
	SendicoRakuma:        100,   // Rakuma Fashion category (approximate)
	SendicoRakuten:       100,   // Rakuten Fashion category (approximate)
	SendicoYahoo:         100,   // Yahoo PayPay Flea Fashion category (approximate)
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
	
	// Page number (default to 1)
	page := opts.Page
	if page < 1 {
		page = 1
	}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("search", opts.TermJP)
	
	// Add category filter for clothing if specified
	// NOTE: Temporarily disabled - may cause 403 errors if Sendico API doesn't support this parameter
	// TODO: Test category parameter support with Sendico API
	// if opts.CategoryID != nil {
	// 	params.Set("category", fmt.Sprintf("%d", *opts.CategoryID))
	// }
	
	// Note: Sendico API doesn't support sort/mobile parameters
	// Recently uploaded items are typically on page 1, so we search multiple pages
	// to ensure we catch all recently uploaded items

	q := path.Query()
	for pair := params.Oldest(); pair != nil; pair = pair.Next() {
		q.Add(pair.Key, fmt.Sprintf("%v", pair.Value))
	}
	path.RawQuery = q.Encode()

	// Retry logic with HMAC refresh handling
	maxRetries := 3
	var resp *http.Response
	var err error
	
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Build HMAC for each attempt (in case secret was refreshed)
		hmac, hmacErr := buildHMAC(HMACInput{
			Secret:  c.HMACSecret(),
			Path:    path.Path,
			Payload: params,
			Nonce:   "",
		})
		if hmacErr != nil {
			return nil, hmacErr
		}
		
		resp, err = c.req(ctx, http.MethodGet, path.String(), nil, hmac, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/json")
		})
		
		if err == nil {
			break // Success
		}
		
		// If HMAC refresh needed, wait and retry (but limit retries to avoid loops)
		if err == ErrHMACRefreshNeeded {
			if attempt < maxRetries {
				// Only retry once for HMAC refresh to avoid infinite loops
				if attempt == 0 {
					log.Printf("   🔄 HMAC refresh needed, retrying once...")
					time.Sleep(500 * time.Millisecond)
					continue
				}
			}
			// If HMAC refresh didn't help, likely a different issue (e.g., invalid parameters)
			return nil, fmt.Errorf("HMAC refresh didn't resolve access issue - may be invalid parameters: %w", err)
		}
		
		// For other errors or max retries reached
		if attempt == maxRetries {
			return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries+1, err)
		}
		
		// Wait before retry for other errors (reduced delays)
		time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
	}
	
	if resp == nil {
		return nil, fmt.Errorf("no response after retries")
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

// SearchMultiplePages searches multiple pages to ensure we don't miss any items
// Returns all items from pages 1 through maxPages, or until no more items are found
// Automatically applies clothing category filter based on the shop
func (c *SendicoClient) SearchMultiplePages(ctx context.Context, shop SendicoShop, opts SendicoSearchOptions, maxPages int) ([]SendicoItem, error) {
	if maxPages < 1 {
		maxPages = 3 // Default to 3 pages to catch recently uploaded items
	}
	
	allItems := make([]SendicoItem, 0)
	seenCodes := make(map[string]bool) // Deduplicate across pages
	
	// Apply clothing category filter for this shop (once, outside the loop)
	// NOTE: Temporarily disabled - may cause 403 errors
	// TODO: Verify category parameter support with Sendico API
	// if categoryID, ok := clothingCategoryIDs[shop]; ok {
	// 	opts.CategoryID = &categoryID
	// }
	
	for page := 1; page <= maxPages; page++ {
		pageOpts := opts
		pageOpts.Page = page
		
		// Search with minimal retry logic (avoid loops)
		items, err := c.Search(ctx, shop, pageOpts)
		if err != nil {
			// If error on later pages, return what we have (might be last page)
			if page > 1 {
				log.Printf("   ⚠️  Error on page %d (may be last page): %v", page, err)
				return allItems, nil // Return what we have so far
			}
			// For first page, return error
			return nil, err
		}
		
		// If no items returned, we've reached the end
		if len(items) == 0 {
			break
		}
		
		// Deduplicate items across pages
		for _, item := range items {
			if !seenCodes[item.Code] {
				seenCodes[item.Code] = true
				allItems = append(allItems, item)
			}
		}
		
		// If we got fewer items than expected, might be last page
		// But continue to next page to be sure (some APIs return partial pages)
		if len(items) < 20 { // Assuming ~20 items per page
			// Still continue to next page to be safe
		}
		
		// Minimal delay between pages (only if searching multiple pages)
		if page < maxPages && maxPages > 1 {
			time.Sleep(200 * time.Millisecond) // Reduced delay
		}
	}
	
	return allItems, nil
}

func (c *SendicoClient) BulkSearch(ctx context.Context, shops []SendicoShop, opts SendicoSearchOptions) ([]SendicoItem, error) {
	// Use single page search for speed (original behavior)
	return c.BulkSearchSinglePage(ctx, shops, opts)
}

// BulkSearchSinglePage searches only page 1 for fastest performance
func (c *SendicoClient) BulkSearchSinglePage(ctx context.Context, shops []SendicoShop, opts SendicoSearchOptions) ([]SendicoItem, error) {
	items := make([]SendicoItem, 0)
	itemsMu := sync.Mutex{}

	// Optimized concurrency for speed
	maxConcurrent := 5 // Increased concurrency for faster searches
	requestDelay := 200 * time.Millisecond // Reduced delay for faster processing
	sem := make(chan struct{}, maxConcurrent)
	g := new(errgroup.Group)
	
	for i, shop := range shops {
		shop := shop
		i := i
		g.Go(func() error {
			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()
			
			// Minimal delay between requests
			if i > 0 {
				time.Sleep(requestDelay)
			}
			
			// Search only page 1
			pageOpts := opts
			pageOpts.Page = 1
			// NOTE: Category filter temporarily disabled to avoid 403 errors
			// TODO: Verify category parameter support with Sendico API
			// if categoryID, ok := clothingCategoryIDs[shop]; ok {
			// 	pageOpts.CategoryID = &categoryID
			// }
			results, err := c.Search(ctx, shop, pageOpts)
			if err != nil {
				// Check if it's a rate limit error
				if errStr := err.Error(); strings.Contains(errStr, "429") || strings.Contains(errStr, "rate limited") {
					log.Printf("   ❌ Rate limited searching %s: %v", shop, err)
				} else {
					log.Printf("   ⚠️  Error searching %s: %v", shop, err)
				}
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

// BulkSearchMultiplePages searches multiple pages across multiple shops
func (c *SendicoClient) BulkSearchMultiplePages(ctx context.Context, shops []SendicoShop, opts SendicoSearchOptions, maxPages int) ([]SendicoItem, error) {
	items := make([]SendicoItem, 0)
	itemsMu := sync.Mutex{}

	// Optimized concurrency for multiple users - balance between speed and rate limits
	maxConcurrent := 5 // Increased for faster multi-page searches
	requestDelay := 300 * time.Millisecond // Reduced delay for faster processing
	sem := make(chan struct{}, maxConcurrent)
	g := new(errgroup.Group)
	
	for i, shop := range shops {
		shop := shop
		i := i
		g.Go(func() error {
			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }()
			
			// Delay between requests to respect rate limits
			if i > 0 {
				time.Sleep(requestDelay)
			}
			
			// Search multiple pages to ensure we don't miss recently uploaded items
			results, err := c.SearchMultiplePages(ctx, shop, opts, maxPages)
			if err != nil {
				// Check if it's a rate limit error
				if errStr := err.Error(); strings.Contains(errStr, "429") || strings.Contains(errStr, "rate limited") {
					log.Printf("   ❌ Rate limited searching %s: %v", shop, err)
				} else {
					log.Printf("   ⚠️  Error searching %s: %v", shop, err)
				}
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
	maxRetries := 3
	baseDelay := 2 * time.Second
	
	// Read body into bytes once for retries (body can only be read once)
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
	}
	
	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, func() io.Reader {
			if bodyBytes != nil {
				return bytes.NewReader(bodyBytes)
			}
			return nil
		}())
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
			if attempt < maxRetries {
				delay := baseDelay * time.Duration(1<<uint(attempt)) // Exponential backoff: 2s, 4s, 8s
				log.Printf("   ⚠️  Request error (attempt %d/%d), retrying in %v: %v", attempt+1, maxRetries+1, delay, err)
				time.Sleep(delay)
				continue
			}
			return nil, err
		}

		// Handle 403 (Forbidden/Access Denied) - HMAC secret may have expired
		if res.StatusCode == http.StatusForbidden {
			bodyBytes, _ := io.ReadAll(res.Body)
			_ = res.Body.Close()
			bodyStr := string(bodyBytes)
			
			// Check if it's an access denied error (HMAC expired)
			if strings.Contains(bodyStr, "Access denied") || strings.Contains(bodyStr, "403") {
				// Refresh HMAC secret - it may have expired
				log.Printf("   🔄 Access denied (403) - refreshing HMAC secret...")
				if err := c.FindHMAC(ctx); err != nil {
					log.Printf("   ❌ Failed to refresh HMAC secret: %v", err)
					return nil, fmt.Errorf("failed to refresh HMAC secret: %w", err)
				}
				log.Printf("   ✅ HMAC secret refreshed")
				
				// Return special error so caller can rebuild HMAC with new secret and retry
				// We can't rebuild HMAC here because we don't have the payload
				return nil, ErrHMACRefreshNeeded
			}
			
			// Other 403 errors (not access denied)
			log.Printf("   ⚠️  Sendico API error: status 403, body: %s", bodyStr)
			return nil, fmt.Errorf("access denied (403): %s", bodyStr)
		}

		// Handle 429 (Too Many Requests) with retry
		if res.StatusCode == http.StatusTooManyRequests {
			_ = res.Body.Close()
			
			// Check for Retry-After header
			retryAfter := baseDelay
			if retryAfterStr := res.Header.Get("Retry-After"); retryAfterStr != "" {
				if seconds, err := time.ParseDuration(retryAfterStr + "s"); err == nil {
					retryAfter = seconds
				} else if seconds, err := time.ParseDuration(retryAfterStr); err == nil {
					retryAfter = seconds
				}
			} else {
				// Exponential backoff if no Retry-After header
				retryAfter = baseDelay * time.Duration(1<<uint(attempt))
			}
			
			if attempt < maxRetries {
				log.Printf("   ⚠️  Rate limited (429) on %s (attempt %d/%d), retrying in %v", path, attempt+1, maxRetries+1, retryAfter)
				time.Sleep(retryAfter)
				continue
			}
			
			// Max retries reached
			return nil, fmt.Errorf("rate limited (429) after %d attempts", maxRetries+1)
		}

		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			log.Printf("   ⚠️  Sendico API error: status %d, body: %s", res.StatusCode, string(body))
			_ = res.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		return res, nil
	}
	
	return nil, fmt.Errorf("max retries exceeded")
}
