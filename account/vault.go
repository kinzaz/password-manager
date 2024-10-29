package account

import (
	"demo/app-demo-3/encrypter"
	"demo/app-demo-3/output"
	"encoding/json"
	"strings"
	"time"

	"github.com/fatih/color"
)

type ByteReader interface {
	Read() ([]byte, error)
}

type ByteWriter interface {
	Write([]byte)
}

type Db interface {
	ByteWriter
	ByteReader
}

type Vault struct {
	Accounts  []AccountWithTimestamp `json:"accounts"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type VaultWithDb struct {
	Vault
	db  Db
	enc encrypter.Encrypter
}

func (vault *VaultWithDb) AddAccount(acc AccountWithTimestamp) {
	vault.Accounts = append(vault.Accounts, acc)
	vault.UpdatedAt = time.Now()
	data, err := vault.ToBytes()
	encData := vault.enc.Encrypt(data)
	if err != nil {
		output.PrintError("Не удалось преобразовать")
	}
	vault.db.Write(encData)
}

func (vault *VaultWithDb) FindAccountsFromVault(str string, checker func(AccountWithTimestamp, string) bool) []AccountWithTimestamp {
	var accounts []AccountWithTimestamp
	for _, account := range vault.Accounts {
		isMatched := checker(account, str)
		if isMatched {
			accounts = append(accounts, account)
		}
	}
	return accounts
}

func (vault *VaultWithDb) DeleteAccountsByUrl(url string) bool {
	var accounts []AccountWithTimestamp
	isDeleted := false

	for _, account := range vault.Accounts {
		isMatched := strings.Contains(account.Url, url)
		if !isMatched {
			accounts = append(accounts, account)
			continue
		}
		isDeleted = true
	}
	vault.Accounts = accounts
	vault.UpdatedAt = time.Now()
	data, err := vault.ToBytes()
	encData := vault.enc.Encrypt(data)
	if err != nil {
		output.PrintError("Не удалось преобразовать")
	}
	vault.db.Write(encData)

	return isDeleted
}

func (vault *Vault) ToBytes() ([]byte, error) {
	file, err := json.Marshal(vault)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func NewVault(db Db, enc encrypter.Encrypter) *VaultWithDb {
	file, err := db.Read()

	if err != nil {
		return &VaultWithDb{
			Vault: Vault{
				Accounts:  []AccountWithTimestamp{},
				UpdatedAt: time.Now(),
			},
			db:  db,
			enc: enc,
		}
	}
	data := enc.Decrypt(file)
	var vault Vault
	err = json.Unmarshal(data, &vault)
	color.Cyan("Найдено %d аккаунтов", len(vault.Accounts))

	if err != nil {
		output.PrintError("Не удалось разобрать файл data.json")
		return &VaultWithDb{
			Vault: Vault{
				Accounts:  []AccountWithTimestamp{},
				UpdatedAt: time.Now(),
			},
			db:  db,
			enc: enc,
		}
	}

	return &VaultWithDb{
		Vault: vault,
		db:    db,
		enc:   enc,
	}
}
