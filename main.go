package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var slackEP string

type Endpoints struct {
	Endpoints []Endpoint `json:"endpoints"`
}
type Endpoint struct {
	Endpoint string `json:"endpoint"`
	Id       string `json:"id"`
	Secret   string `json:"secret"`
}

func main() {
	port := ":8000"
	fmt.Println("MicroPinger 1.0 Starting")
	if strings.ToLower(os.Getenv("mode")) == "server" {
		serv(port)
	} else if strings.ToLower(os.Getenv("mode")) == "client" {
		slackEP = os.Getenv("slack")
		if slackEP == "" {
			fmt.Println("Missing Slack endpoint")
			os.Exit(1)
		}
		ping()
	} else {
		fmt.Println("mode environment variable not found")
		os.Exit(1)
	}
}

func serv(p string) {
	fmt.Printf("Waiting for requests on port %v\n", p)
	http.HandleFunc("/", reply)
	err := http.ListenAndServe(p, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func reply(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong")
	fmt.Printf("pong: %v\n", req.RemoteAddr)
	defer req.Body.Close()
}

func readJson() Endpoints {
	ep, err := os.Open("/endpoints/endpoints.json")
	if err != nil {
		fmt.Println(err)
	}
	defer ep.Close()
	var endpoints Endpoints
	b, _ := ioutil.ReadAll(ep)
	json.Unmarshal(b, &endpoints)
	return endpoints
}

func ping() {
	ep := readJson()
	for _, endPoint := range ep.Endpoints {
		fmt.Printf("Pinging: %s\n", endPoint)
		r, err := fetch(endPoint)
		if err != nil {
			signal(&endPoint.Endpoint, &r)
		}
	}
}

func fetch(e Endpoint) (int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	req, err := http.NewRequest("GET", e.Endpoint, nil)
	req.Header.Set("x-ibm-client-id", e.Id)
	req.Header.Set("x-ibm-client-secret", e.Secret)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("Failed: " + e.Endpoint)
	} else {
		defer resp.Body.Close()
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Failed: %v, %v\n", e.Endpoint, resp.StatusCode)
		return resp.StatusCode, errors.New("Failed: " + e.Endpoint)
	}
	return resp.StatusCode, nil
}

func signal(endPoint *string, errorCode *int) {
	body, _ := json.Marshal(map[string]string{
		"text": *endPoint + " Return code: " + strconv.Itoa(*errorCode),
	})
	fmt.Printf("Signalling %s, return code: %d \n", *endPoint, *errorCode)
	req, err := http.Post(slackEP, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	defer req.Body.Close()
}
