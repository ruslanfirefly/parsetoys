package main

import (
	"log"
	goCharset "code.google.com/p/go-charset/charset"
	"strings"
	"io/ioutil"
	"fmt"
	"github.com/ungerik/go-mail/email"
)

func error_log(err error) {
	if err != nil {
		email.InitGmail(gmailAdress, gmailPassword)
		letter := email.NewBriefMessage("Parse error", fmt.Sprint(err), tomail)
		letter.Send()
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
