package account

import (
	"errors"
	"math/rand/v2"
	"net/url"
	"time"

	"github.com/fatih/color"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789*!")

type Account struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Url      string `json:"url"`
}

type AccountWithTimestamp struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Account
}

func (acc Account) OutputAccount() {
	color.Cyan(acc.Login)
	color.Cyan(acc.Password)
	color.Cyan(acc.Url)
}

func (acc *Account) generatePassword(n int) {
	res := make([]rune, n)
	for i := range res {
		res[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	acc.Password = string(res)
}

// func NewAccount(login, password, urlString string) (*account, error) {
// 	if login == "" {
// 		return nil, errors.New("INVALID_LOGIN")
// 	}

// 	_, err := url.ParseRequestURI(urlString)
// 	if err != nil {
// 		return nil, errors.New("INVALID_URL")
// 	}

// 	newAcc := &account{
// 		login:    login,
// 		password: password,
// 		url:      urlString,
// 	}

// 	if newAcc.password == "" {
// 		newAcc.generatePassword(12)
// 	}
// 	return newAcc, nil
// }

func NewAccountWithTimestamp(login, password, urlString string) (*AccountWithTimestamp, error) {
	if login == "" {
		return nil, errors.New("INVALID_LOGIN")
	}

	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, errors.New("INVALID_URL")
	}

	newAcc := &AccountWithTimestamp{

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Account: Account{
			Url:      urlString,
			Login:    login,
			Password: password,
		},
	}

	if password == "" {
		newAcc.generatePassword(12)
	}
	return newAcc, nil
}
