package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	URL   = "https://server.duinocoin.com/v3/users/"
	QUERY = "?limit=5"
)

type Result struct {
	Result struct {
		Balance struct {
			Balance  float64 `json:"balance"`
			Created  string  `json:"created"`
			Username string  `json:"username"`
			Verified string  `json:"verified"`
		} `json:"balance"`
		Miners []struct {
			Accepted   int         `json:"accepted"`
			Algorithm  string      `json:"algorithm"`
			Diff       int         `json:"diff"`
			Hashrate   float64     `json:"hashrate"`
			Identifier string      `json:"identifier"`
			Ki         int         `json:"ki"`
			Pool       string      `json:"pool"`
			Rejected   int         `json:"rejected"`
			Sharetime  float64     `json:"sharetime"`
			Software   string      `json:"software"`
			Threadid   string      `json:"threadid"`
			Username   string      `json:"username"`
			Wd         interface{} `json:"wd"`
		} `json:"miners"`
		Prices struct {
			Bch      float64 `json:"bch"`
			Justswap float64 `json:"justswap"`
			Max      float64 `json:"max"`
			Nano     float64 `json:"nano"`
			Nodes    float64 `json:"nodes"`
			Pancake  float64 `json:"pancake"`
			Sushi    float64 `json:"sushi"`
			Trx      float64 `json:"trx"`
			Xmg      float64 `json:"xmg"`
		} `json:"prices"`
		Server       string        `json:"server"`
		Transactions []interface{} `json:"transactions"`
	} `json:"result"`
	Success bool `json:"success"`
}

func main() {
	checkBalance()
}

func checkBalance() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered :", r)
			checkBalance()
		}
	}()
	log.Println("Enter your DUCO username: ")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	name, _ := reader.ReadString('\n')
	// convert CRLF to LF
	name = strings.Replace(name, "\n", "", -1)
	log.Println("Ok, connection...")

	var balance float64
	for {
		time.Sleep(4 * time.Second)
		resp, err := http.Get(URL + name + QUERY)
		if err != nil {
			log.Println(err.Error())
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
		}
		var result Result
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Println("ERROR: %s", err.Error())
		}

		srcOk := "\xE2\x9C\x85"
		r, _ := utf8.DecodeRuneInString(srcOk)

		srcProc := "\xF0\x9F\x9A\x80"
		r1, _ := utf8.DecodeRuneInString(srcProc)

		if balance < result.Result.Balance.Balance {
			diff := result.Result.Balance.Balance - balance
			balance = result.Result.Balance.Balance
			fmt.Printf("%v Balance + %f, Current -> %f\n", string(r), diff, balance)
		} else {
			fmt.Printf("%v no changes \n", string(r1))
		}
	}
}
