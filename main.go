package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ajtfj/bank/bank"
)

type SignUpBody struct {
	CPF      string `json:"cpf"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type DepositRequestBody struct {
	Value float64 `json:"value"`
}

type WithdrawRequestBody struct {
	Value float64 `json:"value"`
}

type BalanceResponseBody struct {
	Balance float64 `json:"balance"`
}

type Credentials struct {
	CPF      string `json:"cpf"`
	Password string `json:"password"`
}

var (
	jrBank = bank.NewBank()
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusInternalServerError)
		return
	}

	signUp := SignUpBody{}
	if err := json.Unmarshal(body, &signUp); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if err := jrBank.SignUp(signUp.CPF, signUp.Name, signUp.Password); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetAuthenticatedUser(r *http.Request) (*bank.User, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	credentials := Credentials{}
	json.Unmarshal(body, &credentials)
	user, err := jrBank.GetUser(credentials.CPF)
	if err != nil {
		return nil, err
	}

	if err := user.Authenticate(credentials.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusInternalServerError)
		return
	}

	deposit := DepositRequestBody{}
	if err := json.Unmarshal(body, &deposit); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	balance := user.Account.Deposit(deposit.Value)
	balannceResponseBody := BalanceResponseBody{
		Balance: balance,
	}
	jsonResponseBody, err := json.Marshal(balannceResponseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponseBody)
}

func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusInternalServerError)
		return
	}

	withdraw := WithdrawRequestBody{}
	if err := json.Unmarshal(body, &withdraw); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := user.Account.Withdraw(withdraw.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	balanceResponseBody := BalanceResponseBody{
		Balance: balance,
	}
	jsonResponseBody, err := json.Marshal(balanceResponseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponseBody)
}

func main() {
	http.HandleFunc("/signup", SignUpHandler)

	http.HandleFunc("/deposit", DepositHandler)

	http.HandleFunc("/withdraw", WithdrawHandler)

	http.ListenAndServe(":8081", nil)
}
