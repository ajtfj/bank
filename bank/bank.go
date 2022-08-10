package bank

import (
	"fmt"
	"sync"
)

type Account struct {
	balance float64
	mu      *sync.Mutex
}

func NewAccount() *Account {
	account := Account{
		mu: &sync.Mutex{},
	}

	return &account
}

func (a *Account) Withdraw(value float64) (float64, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.balance < value {
		return -1, fmt.Errorf("not enough money")
	}
	a.balance -= value
	return a.balance, nil
}

func (a *Account) Deposit(value float64) float64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.balance += value
	return a.balance
}

type User struct {
	cpf      string
	name     string
	password string
	Account  *Account
}

func NewUser(cpf string, name string, password string, acc *Account) *User {
	user := User{
		cpf:      cpf,
		name:     name,
		password: password,
		Account:  acc,
	}

	return &user
}

func (u *User) Authenticate(pass string) error {
	if u.password != pass {
		return fmt.Errorf("wrong password")
	}

	return nil
}

type Bank struct {
	users map[string]User
}

func NewBank() *Bank {
	bank := Bank{
		users: make(map[string]User),
	}

	return &bank
}

func (b *Bank) SignUp(cpf string, name string, password string) error {
	if _, ok := b.users[cpf]; ok {
		return fmt.Errorf("CPF already taken")
	}

	acc := NewAccount()
	user := NewUser(cpf, name, password, acc)

	b.users[cpf] = *user

	return nil
}

func (b *Bank) SignIn(cpf string, password string) error {
	if _, ok := b.users[cpf]; !ok {
		return fmt.Errorf("CPF not found")
	}

	if b.users[cpf].password != password {
		return fmt.Errorf("wrong password")
	}

	return nil
}

func (b *Bank) GetUser(cpf string) (*User, error) {
	user, ok := b.users[cpf]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}
