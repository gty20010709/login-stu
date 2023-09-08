package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	uuid "github.com/google/uuid"
	term "golang.org/x/term"
)

const HelpInfo = `login-stu v0.0.2

Description:
	本程序用于在无头环境下登录校园网。
	或者在桌面环境中，通过双击登录，简化流程。

Usage:
	login-stu [username] [password] [help]

Examples:
	login-stu yourUsername yourPassword
`

const LoginUrl = "https://a.stu.edu.cn:444/ac_portal/login.php"

type Account struct {
	Username string
	Password string
}

func main() {
	args := os.Args[1:]

	account, err := getAccountInfo(args)
	if err != nil {
		fmt.Println("读取配置文件错误。 Error: ", err)
		fmt.Print("请参照下方说明使用本程序： \n\n")
		fmt.Print(HelpInfo)
		os.Exit(1)
	}

	resp, err := loginStu(account)
	if err != nil {
		fmt.Println(err)
		fmt.Println("请确定电脑已经连上校园网，然后再尝试登录。")
	}
	fmt.Println(resp)

	var input string

	fmt.Println("Please press Enter to exit.")
	fmt.Scanln(&input)
}

func getAccountInfo(args []string) (*Account, error) {
	if len(args) == 1 && args[0] == "help" {
		fmt.Print(HelpInfo)
		os.Exit(1)
	} else if len(args) == 2 {

		account := Account{
			Username: args[0],
			Password: args[1],
		}

		toSave(&account)

		return &account, nil
	} else {
		fmt.Print("未指定用户名和密码，尝试查找配置文件。try to find config file....\n\n")
		account, err := loadUserInfo()
		if err != nil {
			fmt.Println("读取配置文件错误。 Error: ", err)
			fmt.Println("请手动输入帐号密码。 ")
			var username, password string
			fmt.Print("请输入校园网帐号: ")
			fmt.Scanln(&username)
			fmt.Print("请输入校园网密码(不会显示出来): ")
			// fmt.Scanln(&password)
			_pw, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			password = string(_pw)
			// fmt.Println(password)

			fmt.Println()

			account := Account{
				Username: username,
				Password: password,
			}
			toSave(&account)
			return &account, nil
		}

		return account, nil

	}

	return nil, nil
}

func toSave(accout *Account) {
	username := accout.Username
	password := accout.Password
	for {
		var choice string
		fmt.Print("是否保存登录信息？(yes/no): ")
		fmt.Scanln(&choice)
		if strings.ToLower(choice) == "yes" {
			storeUserInfo(username, password)
			break
		} else if strings.ToLower(choice) == "no" {
			break
		} else {
			fmt.Println("输入有误,请重新输入 (yes/no)")
		}
	}
}

func loginStu(account *Account) (string, error) {
	username := account.Username
	password := account.Password
	// 创建一个 HTTP 客户端，并忽略SSL 错误
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	formData := url.Values{
		"opr":         {"pwdLogin"},
		"userName":    {username},
		"pwd":         {password},
		"ipv4or6":     {""},
		"rememberPwd": {"0"},
	}

	resp, err := client.Post(LoginUrl, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		// fmt.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

func storeUserInfo(username, password string) error {
	key := genKey()
	encryptedUsername, err := EncryptMessage(username, key)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	encryptedPassword, err := EncryptMessage(password, key)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	encryptedInfo := encryptedUsername + "." + encryptedPassword + "." + key
	fmt.Println("加密后的登录信息：" + encryptedInfo)

	// write to file
	configFilePath, err := getConfigFilePath()
	if err != nil {
		//
		return err
	}
	configFile, err := os.Create(configFilePath)
	if _, err := io.WriteString(configFile, encryptedInfo); err != nil {
		//
		return err
	}
	defer configFile.Close()

	return nil
}

func loadUserInfo() (*Account, error) {

	configFilePath, err := getConfigFilePath()
	if err != nil {
		//
		return nil, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		//
		return nil, err
	}

	defer configFile.Close()

	encryptedInfo, err := io.ReadAll(configFile)
	if err != nil {
		//
		return nil, err
	}

	parts := strings.Split(string(encryptedInfo), ".")
	encryptedUsername, encryptedPassword, key := parts[0], parts[1], parts[2]
	// fmt.Println(encryptedUsername + "." + encryptedPassword + "." + key)

	username, err := DecryptMessage(encryptedUsername, key)
	if err != nil {
		//
		return nil, err
	}

	password, err := DecryptMessage(encryptedPassword, key)
	if err != nil {
		//
		return nil, err
	}

	return &Account{
		Username: username,
		Password: password,
	}, nil
}

func getConfigFilePath() (string, error) {
	executor, err := os.Executable()
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	runPath := path.Dir(executor)
	configFilePath := filepath.Join(runPath, "config")
	return configFilePath, nil
}

func genKey() string {

	uuidStr := uuid.New().String()
	key := strings.ReplaceAll(uuidStr, "-", "")
	return key
}
