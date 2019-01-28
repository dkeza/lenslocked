package models

import "strings"

const (
	// ErrorNotFound is default record not found error
	ErrorNotFound modelError = "models: resource not found"
	// ErrIDInvalid is invalid id
	ErrIDInvalid modelError = "models: invalid ID"
	// ErrorPasswordIncorrect error
	ErrorPasswordIncorrect modelError = "models: incorrect password provided"
	// ErrEmailRequired returns error when E-Mail address is not provided
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailInvalid returns error when rmail is invalid
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailTaken returns error
	ErrEmailTaken modelError = "models: email address is already taken"
	//ErrPasswordTooShort returns error
	ErrPasswordTooShort modelError = "models: password must be minimum 8 characters long"
	//ErrPasswordRequired returns error
	ErrPasswordRequired modelError = "models: password is required"
	// ErrTitleRequired returns error
	ErrTitleRequired modelError = "models: title is required"
	// ErrTokenInvalid returns error
	ErrTokenInvalid modelError = "models: token provided is not valid"
	// ErrRememberRequired returns error
	ErrRememberRequired privateError = "models: remember token is required"
	// ErrRememberTooShort returns error
	ErrRememberTooShort privateError = "models: remember token must be at least 32 bytes"
	// ErrUserIDRequired returns error
	ErrUserIDRequired privateError = "models: user ID is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
