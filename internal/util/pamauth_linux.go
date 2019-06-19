package util

import (
	"fmt"
	"github.com/msteinert/pam"
	"log"
	"strings"
)

func CheckUser(username string) bool {
	// get and check current user
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	// ToDo(frennkie) do not fatal here
	if currentUser.Uid == "0" {
		log.Fatalf("Fatal: can not be run as root")
	}

	// ToDo(frennkie) do not fatal here
	if currentUser.Username != username {
		log.Fatalf("Fatal: can only authenticate for same user")
	}

	return true
}

func CheckUserPassword(username, password string) bool {
	t, err := pam.StartFunc("", "", func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			return strings.TrimSuffix(string(password), "\r"), nil
		case pam.PromptEchoOn:
			return strings.TrimSuffix(string(username), "\r"), nil
		case pam.ErrorMsg:
			fmt.Println("## at ErrorMsg")
			log.Print(msg)
			return "", nil
		case pam.TextInfo:
			fmt.Println("## at TextInfo")
			fmt.Println(msg)
			return "", nil
		}
		return "", errors.New("unrecognized message style")
	})

	// ToDo(frennkie) do not fatal here
	if err != nil {
		log.Fatalf("Start: %s", err.Error())
	}

	// ToDo(frennkie) do not fatal here
	err = t.Authenticate(0)
	if err != nil {
		log.Fatalf("Authenticate: %s", err.Error())
	}

	return true
}
