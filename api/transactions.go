package api

import (
	"encoding/json"
	"time"
)

type getTransactionsResponse struct {
	Data []Transaction `json:"data"`
	Errors
}

type postTransactionsResponse struct {
	Data Transaction `json:"data"`
	Errors
}

// Transaction JSON
type Transaction struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	Amount       Money     `json:"amount"`
	NativeAmount Money     `json:"native_amount"`
	Created      time.Time `json:"created_at"`
	Updated      time.Time `json:"updated_at"`
	To           struct {
		Resource string `json:"resource"`
		Email    string `json:"email"`
		ID       string `json:"id"`
	} `json:"to"`
}

// EmptyTransaction result
var EmptyTransaction = Transaction{}

// EmptyTransactions result
var EmptyTransactions = []Transaction{}

// Request JSON
type Request struct {
	Type        string `json:"type"`
	To          string `json:"to"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Idem        string `json:"idem"`
}

// Transactions for a given account
func (c *Client) Transactions(accountID string) ([]Transaction, error) {
	acc, err := c.findAccount(accountID)
	if err != nil {
		return EmptyTransactions, err
	}
	res, err := c.Get("/accounts/" + acc.ID + "/transactions")
	if err != nil {
		return EmptyTransactions, err
	}
	defer res.Body.Close()
	var result getTransactionsResponse
	return result.Data, json.NewDecoder(res.Body).Decode(&result)
}

// Transfer money from an acc to another
func (c *Client) Transfer(
	from, to, amount, currency, description string,
) (Transaction, error) {
	acc, err := c.findAccount(from)
	if err != nil {
		return EmptyTransaction, err
	}
	request := Request{
		Type:        "transfer",
		To:          to,
		Amount:      amount,
		Currency:    currency,
		Description: description,
	}
	return c.transactionsPost(acc, request)
}

// Send money from an acc to a btc addr
func (c *Client) Send(
	from, to, amount, currency, description string,
) (Transaction, error) {
	acc, err := c.findAccount(from)
	if err != nil {
		return EmptyTransaction, err
	}
	request := Request{
		Type:        "send",
		To:          to,
		Amount:      amount,
		Currency:    currency,
		Description: description,
	}
	return c.transactionsPost(acc, request)
}

func (c *Client) transactionsPost(
	acc Account, request Request,
) (Transaction, error) {
	res, err := c.Post("/accounts/"+acc.ID+"/transactions", request)
	if err != nil {
		return EmptyTransaction, err
	}
	defer res.Body.Close()
	var result postTransactionsResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if res.StatusCode == 201 {
		return result.Data, err
	}
	return EmptyTransaction, c.newAPIError(res.Status, result.Errors)
}
