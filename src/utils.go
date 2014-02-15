package main

import (
	"log"
	goCharset "code.google.com/p/go-charset/charset"
	"strings"
	"io/ioutil"
	"fmt"
)

func error_log(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func encode_string(str string, charset string) string {
	r, err := goCharset.NewReader(charset, strings.NewReader(str))
	error_log(err)
	result, err := ioutil.ReadAll(r)
	error_log(err)
	fmt.Printf("%s\n", result)
	return string(result)
}
