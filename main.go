package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ajtfj/bank/bank"
)

type SignUpBody struct {
	CPF      string `json:"cpf"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var (
	jrBank = bank.NewBank()
)

func SignUphHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read request body")
	}

	signUp := SignUpBody{}
	json.Unmarshal(body, &signUp)
	jrBank.SignUp(signUp.CPF, signUp.Name, signUp.Password)
}

func main() {
	http.HandleFunc("/signup", SignUphHandler)

	http.ListenAndServe(":8081", nil)
}
