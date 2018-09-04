package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

const fileName = "account.json"

type Account struct {
	Amount   int64 `json:"amount"`
	IsClosed bool
	IsOpened bool
	mu       sync.Mutex
}

func Open(initialDeposit int64) *Account {

	if initialDeposit < 0 {
		return nil
	}
	var a Account
	a.Amount = initialDeposit
	a.IsOpened = true
	return &a
}

func (a *Account) Close() (payout int64, ok bool) {

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.IsClosed {
		return 0, false
	}
	a.IsClosed = true
	return a.Amount, true
}

func (a *Account) Balance() (balance int64, ok bool) {

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.IsClosed {
		return 0, false
	}
	return a.Amount, true
}

func (a *Account) Deposit(amount int64) (newBalance int64, ok bool) {

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.IsClosed {
		return 0, false
	}
	newBalance = a.Amount + amount

	if newBalance < 0 {
		return 0, false
	}
	a.Amount = newBalance
	return a.Amount, true
}

func Load() (*Account, error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, 0755)
	defer f.Close()

	if err != nil {
		return nil, fmt.Errorf("Account is not created")
	}
	byteValue, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}

	var account Account

	json.Unmarshal(byteValue, &account)

	if !account.IsOpened {
		return nil, fmt.Errorf("Account is not created")
	}

	return &account, nil
}

func (a *Account) Store() error {
	byteValue, err := json.Marshal(a)

	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, byteValue, 0755)

	if err != nil {
		return err
	}
	return err
}
