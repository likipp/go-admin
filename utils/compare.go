package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CompareHashAndPassword(e string, p string) (bool, error) {
	fmt.Println(e, p)
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		log.Print(err.Error(), "这里有一个错误")
		return false, err
	}
	return true, nil
}
