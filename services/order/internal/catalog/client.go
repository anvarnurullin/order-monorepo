package catalog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Product struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	SKU          string    `json:"sku"`
	Price        float64   `json:"price"`
	QtyAvailable int       `json:"qty_available"`
	CreatedAt    time.Time `json:"created_at"`
}

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) GetProduct(id int64) (*Product, error) {
	url := fmt.Sprintf("%s/api/v1/products/%d", c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product not found, status: %d", resp.StatusCode)
	}

	var p Product
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (c *Client) DecreaseQty(id int64, qty int) error {
	url := fmt.Sprintf("%s/api/v1/products/%d/decrease", c.BaseURL, id)

	body := struct {
		Quantity int `json:"quantity"`
	}{Quantity: qty}

	reqBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to decrease the qty, status: %d", resp.StatusCode)
	}

	return nil
}
