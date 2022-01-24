package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/spf13/viper"
)

var MONGO string
var DB string
var SMTPHOST string
var SMTPPORT string
var SMTPFROM string
var WHITELIST []string
var TEMPLATE *template.Template

func loadConfig() {
	v := viper.New()
	v.SetConfigName("wekan-alerter")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}
	MONGO = v.GetString("MONGO")
	SMTPHOST = v.GetString("SMTPHOST")
	SMTPPORT = v.GetString("SMTPPORT")
	SMTPFROM = v.GetString("SMTPFROM")
	DB = v.GetString("DB")
	WHITELIST = v.GetStringSlice("WHITELIST")

	var err error
	TEMPLATE, err = template.New("mail").ParseFiles(v.GetString("TEMPLATE"))
	if err != nil {
		panic(err)
	}
}
