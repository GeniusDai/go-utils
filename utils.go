package utils

import (
	"log"
	"os"
)

func ErrPanic(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func ErrMsg(err error) {
	if err != nil {
		log.Println(err)
	}
}

func ErrFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true
	} else {
		return false
	}
}

