package main

import (
	"demo/app-demo-3/account"
	"fmt"

	"github.com/fatih/color"
)

func main() {
	fmt.Println("__ Менеджер паролей __")
	vault := account.NewVault()

Menu:
	for {
		variant := getMenu()
		switch variant {
		case 1:
			createAccount(vault)
		case 2:
			findAccount(vault)
		case 3:
			deleteAccount(vault)
		default:
			break Menu
		}
	}

}

func findAccount(vault *account.Vault) {
	url := promptData("Введите url для поиска")

	accounts := vault.FindAccountsByUrl(url)
	if len(accounts) == 0 {
		color.Red("Аккаунт не найден")
	}
	for _, account := range accounts {
		account.OutputAccount()
	}
}

func deleteAccount(vault *account.Vault) {
	url := promptData("Введите url для удаления")

	isDeleted := vault.DeleteAccountsByUrl(url)

	if !isDeleted {
		color.Red("Аккаунт не найден")
	}
	color.Green("Аккаунты удалены")

}

func createAccount(vault *account.Vault) {
	login := promptData("Введите логин")
	password := promptData("Введите пароль")
	url := promptData("Введите URL")

	myAccount, err := account.NewAccountWithTimestamp(login, password, url)

	if err != nil {
		fmt.Println("Неверный формат URL или логина")
		return
	}

	vault.AddAccount(*myAccount)
}

func getMenu() int {
	var variant int

	fmt.Println("Выберите вариант: ")
	fmt.Println("1. Создать аккаунт")
	fmt.Println("2. Найти аккаунт")
	fmt.Println("3. Удалить аккаунт")
	fmt.Println("4. Выход")

	fmt.Scan(&variant)
	return variant
}

func promptData(prompt string) string {
	fmt.Println(prompt + ": ")
	var res string
	fmt.Scanln(&res)
	return res
}
