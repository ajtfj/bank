package bank

import (
	"fmt"
	"sync"
)

type User struct {
	cpf      string
	name     string
	password string
}

type Account struct {
	user    *User
	balance float64
	lock    *sync.Mutex
}

type Bank struct {
	users    map[string]User
	accounts []Account
}

func (a *Account) Withdraw(value float64) (float64, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.balance < value {
		return -1, fmt.Errorf("not enough money")
	}
	a.balance -= value
	return a.balance, nil
}

func (a *Account) Deposit(value float64) float64 {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.balance += value
	return a.balance
}

func NewAccount(user *User) *Account {
	account := Account{
		user: user,
		lock: &sync.Mutex{},
	}

	return &account
}

func NewUser(cpf string, name string, password string) *User {
	user := User{
		cpf:      cpf,
		name:     name,
		password: password,
	}

	return &user
}

func NewBank() *Bank {
	bank := Bank{}

	return &bank
}

func (b *Bank) SingUp(cpf string, name string, password string) error {
	if _, ok := b.users[cpf]; ok {
		return fmt.Errorf("CPF already taken")
	}

	user := NewUser(cpf, name, password)
	acc := NewAccount(user)

	b.users[cpf] = *user
	b.accounts = append(b.accounts, *acc)

	return nil
}

func (b *Bank) SingIn(cpf string, password string) error {
	if _, ok := b.users[cpf]; !ok {
		return fmt.Errorf("CPF not found")
	}

	if b.users[cpf].password != password {
		return fmt.Errorf("wrong password")
	}

	return nil
}
