package main

import (
	"testing"
)

func TestGenKey(t *testing.T) {
	key := genKey()
	t.Log(key)
}

func TestStoreUserInfo(t *testing.T) {
	storeUserInfo("yourUsername", "yourPassword")
}

func TestLoadUserInfo(t *testing.T) {
	loadUserInfo()
}

// func TestToSave(t *testing.T) {
// 	account := Account{
// 		Username: "yourUsername",
// 		Password: "yourPassword",
// 	}
// 	toSave(&account)
// }
