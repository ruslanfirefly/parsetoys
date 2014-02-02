package main

import (
	"log"
	charset "code.google.com/p/go-charset/charset"
	"strings"
	"io/ioutil"
	"fmt"
)

func error_log(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func encode_string(str string) string {
	r, err := charset.NewReader("windows-1251", strings.NewReader(str))
	error_log(err)
	result, err := ioutil.ReadAll(r)
	error_log(err)
	fmt.Printf("%s\n", result)
	return string(result)
}
