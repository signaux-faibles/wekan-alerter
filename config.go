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
var BLACKLIST []string
var WHITELIST []string
var TEMPLATE *template.Template
var DRYRUN bool

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
	BLACKLIST = v.GetStringSlice("BLACKLIST")
	WHITELIST = v.GetStringSlice("WHITELIST")
	DRYRUN = v.GetBool("DRYRUN")

	checkConfig()

	var err error
	TEMPLATE, err = template.New("mail").ParseFiles(v.GetString("TEMPLATE"))
	if err != nil {
		panic(err)
	}
}

func checkConfig() {
	if len(WHITELIST)*len(BLACKLIST) > 0 {
		panic("WHITELIST et BLACKLIST sont mutuellement exclusives, vérifiez la configuration")
	}
}
