package telegram

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"golang.org/x/term"
)

type noSignUp struct{}

func (n noSignUp) SignUp(context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (n noSignUp) AcceptTermsOfService(_ context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

type SimpleAuthFlow struct {
	noSignUp    // Prevent signup
	PhoneNumber string
}

func (s SimpleAuthFlow) Phone(context.Context) (string, error) {
	return s.PhoneNumber, nil
}

func (s SimpleAuthFlow) Password(context.Context) (string, error) {
	password, err := scanLnWithoutEcho("Enter your password: ")
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func (s SimpleAuthFlow) Code(context.Context, *tg.AuthSentCode) (string, error) {
	code, err := scanLnWithoutEcho("Enter the code you received to your telegram account: ")
	if err != nil {
		return "", err
	}

	return string(code), nil
}

// Prompts the user for input without echoing the characters typed.
func scanLnWithoutEcho(s string) (string, error) {
	fmt.Print(s)
	input, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(input), nil
}
