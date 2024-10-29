package main

import (
	"demo/app-demo-3/account"
	"demo/app-demo-3/encrypter"
	"demo/app-demo-3/files"
	"demo/app-demo-3/output"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var menu = map[string]func(*account.VaultWithDb){
	"1": createAccount,
	"2": findAccountByUrl,
	"3": findAccountByLogin,
	"4": deleteAccount,
}

var menuVariants = []string{
	"1. Создать аккаунт", "2. Найти аккаунт по URL", "3. Найти аккаунт по логину", "4. Удалить аккаунт", "5. Выход", "Выберите вариант",
}

func main() {
	fmt.Println("__ Менеджер паролей __")
	err := godotenv.Load()
	if err != nil {
		output.PrintError("Не удалось найти env файл")
	}

	vault := account.NewVault(files.NewJsonDb("data.vault"), *encrypter.NewEncrypter())
	// vault := account.NewVault(cloud.NewCloudDb(("https://cloud.db")))

Menu:
	for {
		variant := promptData(menuVariants...)
		menuFunc := menu[variant]
		if menuFunc == nil {
			break Menu
		}
		menuFunc(vault)
	}

}

func findAccountByUrl(vault *account.VaultWithDb) {
	url := promptData("Введите url для поиска")

	accounts := vault.FindAccountsFromVault(url, func(acc account.AccountWithTimestamp, str string) bool {
		return strings.Contains(acc.Url, str)
	})
	outputResult(&accounts)
}

func findAccountByLogin(vault *account.VaultWithDb) {
	login := promptData("Введите логин для поиска")

	accounts := vault.FindAccountsFromVault(login, func(acc account.AccountWithTimestamp, str string) bool {
		return strings.Contains(acc.Login, str)
	})
	outputResult(&accounts)
}

func outputResult(accounts *[]account.AccountWithTimestamp) {
	if len(*accounts) == 0 {
		output.PrintError("Аккаунт не найден")
	}
	for _, account := range *accounts {
		account.OutputAccount()
	}
}

func deleteAccount(vault *account.VaultWithDb) {
	url := promptData("Введите url для удаления")

	isDeleted := vault.DeleteAccountsByUrl(url)

	if !isDeleted {
		output.PrintError("Аккаунты не найдены")
	}
	color.Green("Аккаунты удалены")

}

func createAccount(vault *account.VaultWithDb) {
	login := promptData("Введите логин")
	password := promptData("Введите пароль")
	url := promptData("Введите URL")

	myAccount, err := account.NewAccountWithTimestamp(login, password, url)

	if err != nil {
		output.PrintError("Неверный формат URL или логина")
		return
	}

	vault.AddAccount(*myAccount)
}

func promptData(prompt ...string) string {
	for i, line := range prompt {
		if i == len(prompt)-1 {
			fmt.Printf("%v: ", line)
		} else {
			fmt.Println(line)
		}
	}
	var res string
	fmt.Scanln(&res)
	return res
}
