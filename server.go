package main

import (
	"encoding/json"
	"fmt"
	"go_test/bank-account"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ParamsInitialAmount struct {
	InitialAmount int64
}

type ParamsAmount struct {
	Amount int64 `json:"amount"`
}

type ErrorResult struct {
	Error string `json:"error"`
}

func parseInitialAmount(text []byte) (int64, error) {

	p := &ParamsInitialAmount{}
	err := json.Unmarshal(text, p)

	if err != nil {
		return 0, fmt.Errorf("parseInitialAmount error")
	}
	return p.InitialAmount, nil
}

func parseAmount(text []byte) (int64, error) {

	p := &ParamsAmount{}
	err := json.Unmarshal(text, p)

	if err != nil {
		return 0, fmt.Errorf("parseAmount error")
	}
	return p.Amount, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(400)
		sendResult(w, ErrorResult{Error: err.Error()})
		return
	}
	switch r.Method {
	case http.MethodPut:
		initialAmount, err := parseInitialAmount(body)

		if err != nil {
			w.WriteHeader(400)
			sendResult(w, ErrorResult{Error: err.Error()})
			return
		}
		account := account.Open(initialAmount)

		if account == nil {
			w.WriteHeader(400)
			sendResult(w, ErrorResult{Error: "Not enough money"})
			return
		}

		if !accountStore(w, account) {
			return
		}
		return
	case http.MethodPost:
		account, err := account.Load()

		if err != nil {
			w.WriteHeader(500)
			sendResult(w, ErrorResult{Error: err.Error()})
			return
		}

		if account.IsClosed {
			w.WriteHeader(400)
			sendResult(w, ErrorResult{Error: "Account is closed"})
			return
		}
		amount, err := parseAmount(body)

		if err != nil {
			w.WriteHeader(400)
			sendResult(w, ErrorResult{Error: err.Error()})
			return
		}
		_, ok := account.Deposit(amount)

		if !ok {
			w.WriteHeader(400)
			sendResult(w, ErrorResult{Error: "Not enough money"})
			return
		}

		if !accountStore(w, account) {
			return
		}
		return
	case http.MethodDelete:
		account, err := account.Load()

		if err != nil {
			w.WriteHeader(500)
			sendResult(w, ErrorResult{Error: err.Error()})
			return
		}
		_, ok := account.Close()

		if !ok {
			w.WriteHeader(400)
			sendResult(w, ErrorResult{Error: "Account is closed"})
			return
		}

		if !accountStore(w, account) {
			return
		}
		return
	case http.MethodGet:
		account, err := account.Load()

		if err != nil {
			w.WriteHeader(400)
			sendResult(w, ErrorResult{Error: err.Error()})
			return
		}

		if account.IsClosed {
			w.WriteHeader(400)
			sendResult(w, ErrorResult{Error: "Account is closed"})
			return
		}
		amount := ParamsAmount{
			Amount: account.Amount,
		}
		sendResult(w, amount)
		return
	}
}

func accountStore(w http.ResponseWriter, account *account.Account) bool {
	err := account.Store()

	if err != nil {
		w.WriteHeader(500)
		sendResult(w, ErrorResult{Error: err.Error()})
	}
	return err == nil
}

func sendResult(w http.ResponseWriter, data interface{}) {
	result, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	contentLength := len(string(result)) + 1
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))

	fmt.Fprintln(w, string(result))
}

func main() {
	http.HandleFunc("/account", handler)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
