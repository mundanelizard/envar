package command

import (
	"fmt"

	"github.com/mundanelizard/envi/pkg/cli"
)

func Signup() *cli.Command {
	user := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "user",
			Usage:    "Share user is required - ie 'envi signup -user='username'",
			Required: true,
		},
	}

	password := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "password",
			Usage:    "Share user is required - ie 'envi signup -password='password'",
			Required: true,
		},
	}

	return &cli.Command{
		Name:   "signup",
		Action: handleSignup,
		Flags:  []cli.Flagger{user, password},
	}
}

func handleSignup(values *cli.ActionArgs, args []string) {
	username, _ := values.GetString("user")
	password, _ := values.GetString("password")

	err := srv.CreateAccount(username, password)
	if err != nil {
		logger.Fatal(err)
		return
	}

	fmt.Println("Account creation was successful!")
}
