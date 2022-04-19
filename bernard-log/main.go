package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const configFile = "bernard-conf.json"

type config struct {
	Host string
	Key  string
}

type newLogReq struct {
	Level   int64  `json:"level"`
	Service string `json:"service"`
	Content string `json:"content"`
}

const StdOut = 0
const StdErr = 1

func main() {
	b, err := ioutil.ReadFile(configFile)
	if errors.Is(os.ErrNotExist, err) {
		fmt.Printf("Couldn't find %s\n", configFile)
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("Error reading %s - %s", configFile, err)
		os.Exit(1)
	}

	var conf config
	err = json.Unmarshal(b, &conf)
	if err != nil {
		fmt.Printf("Error parsing %s - %s", configFile, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := newLogReq{
			Content: scanner.Text(),
			Level:   StdOut,
			Service: "TODO",
		}

		b, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Error encoding message - %s", err)
			os.Exit(1)
		}

		body := bytes.NewReader(b)
		req, err := http.NewRequest("POST", conf.Host, body)
		if err != nil {
			fmt.Printf("Error creating request - %s", err)
			os.Exit(1)
		}

		req.Header.Add("Authorization", conf.Key)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error sending request - %s", err)
			os.Exit(1)
		}

		if res.StatusCode != http.StatusOK {
			fmt.Printf("Unexpected response code - %d", res.StatusCode)
			os.Exit(1)
		}
	}
}
