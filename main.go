package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kivattt/getopt"
)

const PROGRAM_NAME = "cpe-insight"
const VERSION = "v0.0.1"
const BASE_URL = "https://wifi.telenor.no"

//const BASE_URL = "http://localhost:8080"

func whereIsPassword(model string) string {
	if model == "dg2200" {
		return "You can find the default password on the label under the router after the text admin."
	}

	if model == "p8702" {
		return "You can find the default password on the label under the router after the text WPA."
	}

	return "You can find the default password on the label under the router."
}

func isOnlyDigits(str string) bool {
	for _, c := range str {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

type CPEInfo struct {
	username          string
	model             string
	modelFriendlyName string
	txnId             string
}

// Returns the username, model, friendly model name and txnId for your router
func getCPEInfo() (CPEInfo, error) {
	resp, err := http.Get(BASE_URL)
	if err != nil {
		return CPEInfo{}, err
	}

	// Important that these are case-sensitive, since "Reference:", "Model:" can be found earlier in the "const data = ..." line
	usernameKey := "reference:\""
	modelKey := "model:\""
	modelFriendlyNameKey := "friendly_name:\""
	txnIDKey := "txnId:\""

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if !strings.HasPrefix(strings.TrimLeft(line, " \t"), "const data =") {
			continue
		}

		getValue := func(key string) (string, error) {
			startIndex := strings.Index(line, key)
			if startIndex == -1 {
				return "", errors.New("key not found")
			}

			start := startIndex + len(key)
			if start >= len(line) {
				return "", errors.New("key not found, out of range error")
			}

			end := len(line) - 1
			if end < 0 {
				panic("should not happen! We already ensured len(line) > 0")
			}

			data := line[start:end]

			nextQuote := strings.Index(data, "\"")
			if nextQuote == -1 {
				return "", errors.New("closing quote not found for username")
			}

			return data[:nextQuote], nil
		}

		username, err := getValue(usernameKey)
		if err != nil {
			return CPEInfo{}, err
		}

		if !isOnlyDigits(username) {
			return CPEInfo{}, errors.New("invalid username found, contained non-digit characters")
		}

		model, err := getValue(modelKey)
		if err != nil {
			return CPEInfo{}, err
		}

		modelFriendlyName, err := getValue(modelFriendlyNameKey)
		if err != nil {
			return CPEInfo{}, err
		}

		txnId, err := getValue(txnIDKey)
		if err != nil {
			return CPEInfo{}, err
		}

		return CPEInfo{username: username, model: model, modelFriendlyName: modelFriendlyName, txnId: txnId}, nil
	}

	return CPEInfo{}, errors.New("no username found in the page")
}

func getCPECookieToken(username, password, txnId string) ([]*http.Cookie, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("save_method", "volatile")
	writer.WriteField("redirect_to", "")
	writer.WriteField("username", username)
	writer.WriteField("password", password)
	writer.Close()

	//request, err := http.NewRequest("POST", BASE_URL + "/login?/login=&txn=" + txnId, body)
	request, err := http.NewRequest("POST", BASE_URL+"/login?/login=", body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("Origin", BASE_URL) // Required because of sveltekit: https://stackoverflow.com/q/73790956

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("wrong password")
	}

	return resp.Cookies(), nil
}

func customerGet(cookies []*http.Cookie, cpe CPEInfo) (any, error) {
	request, err := http.NewRequest(http.MethodGet, BASE_URL+"/v1/"+cpe.username+"/customer", strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	return resp.Body, nil
}

func main() {
	h := flag.Bool("help", false, "display this help and exit")
	v := flag.Bool("version", false, "output version information and exit")
	list := flag.Bool("list", false, "list all API endpoints")
	//endpoint := flag.String("endpoint", "", "request endpoint")
	//all := flag.Bool("all", false, "request all endpoints")
	password := flag.String("password", "", "router admin password")
	//output := flag.String("output", "", "output json to file")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init(PROGRAM_NAME, flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		//"e", "endpoint",
		"l", "list",
		//"a", "all",
		"p", "password",
		//"o", "output",
		"v", "version",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(2)
	}

	if *v {
		fmt.Println(PROGRAM_NAME, VERSION)
		os.Exit(0)
	}

	if *h {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS]")
		fmt.Println("CPE Insight API explorer")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	if *list {
		PrintEndPoints()
		os.Exit(0)
	}

	fmt.Println("Getting router information... This may take a while")

	cpe, err := getCPEInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Router model:", cpe.modelFriendlyName)
	fmt.Println("txnId:", cpe.txnId)
	fmt.Println()
	fmt.Println(whereIsPassword(cpe.model))
	fmt.Println("Username:", cpe.username)

	fmt.Println("Logging in...")

	cookies, err := getCPECookieToken(cpe.username, *password, cpe.txnId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(cookies)

	customerGet(cookies, cpe)
}
