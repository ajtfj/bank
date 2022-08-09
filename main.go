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

type Deposit struct {
	Value float64 `json:"value"`
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
	json.Unmarshal(body, &signUp)
	jrBank.SignUp(signUp.CPF, signUp.Name, signUp.Password)
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

	deposit := Deposit{}
	json.Unmarshal(body, &deposit)
	user.Account.Deposit(deposit.Value)
}

func main() {
	http.HandleFunc("/signup", SignUpHandler)

	http.HandleFunc("/deposit", DepositHandler)

	http.ListenAndServe(":8081", nil)
}
