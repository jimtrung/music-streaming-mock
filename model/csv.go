package model

type Account struct {
	Username string `csv:"name"`
	Password string `csv:"password"`
}
