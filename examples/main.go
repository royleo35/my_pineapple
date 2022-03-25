package main

import (
	"io/ioutil"
	"royleo/pineapple/src"
)

func main() {
	fileName := "code.pineapple"
	code, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	src.Execute(string(code))
}
