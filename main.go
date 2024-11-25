package main

import (
	"bufio"
	"bytes"
	"encoding/json"
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
const VERSION = "v0.0.2"
const BASE_URL = "https://wifi.telenor.no"

func whereIsPassword(model string) string {
	if model == "dg2200" {
		return "you can find the default password on the label under the router after the text admin"
	}

	if model == "p8702" {
		return "you can find the default password on the label under the router after the text WPA"
	}

	return "you can find the default password on the label under the router"
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
}

// Returns the username, model and friendly model name for your router
func getCPEInfo() (CPEInfo, error) {
	resp, err := http.Get(BASE_URL)
	if err != nil {
		return CPEInfo{}, err
	}

	// Important that these are case-sensitive, since "Reference:", "Model:" can be found earlier in the "const data = ..." line
	usernameKey := "reference:\""
	modelKey := "model:\""
	modelFriendlyNameKey := "friendly_name:\""

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

		return CPEInfo{username: username, model: model, modelFriendlyName: modelFriendlyName}, nil
	}

	return CPEInfo{}, errors.New("no username found in the page")
}

func getCPECookieToken(username, password string, cpe CPEInfo) ([]*http.Cookie, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("save_method", "volatile")
	writer.WriteField("redirect_to", "")
	writer.WriteField("username", username)
	writer.WriteField("password", password)
	writer.Close()

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
		return nil, errors.New("wrong password, " + whereIsPassword(cpe.model))
	}

	return resp.Cookies(), nil
}

func getRequest(endPointKey string, cookies []*http.Cookie, cpe CPEInfo) ([]byte, error) {
	endPointStr, ok := apiEndPoints[endPointKey]
	if !ok {
		return nil, errors.New("invalid endpoint")
	}

	endPoint := BASE_URL + CPE_INSIGHT_API_BASE_URL + strings.ReplaceAll(endPointStr, "${t}", cpe.username)

	request, err := http.NewRequest(http.MethodGet, endPoint, strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	request.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errors.New("Server responded: " + resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var theJSON any
	err = json.Unmarshal(body, &theJSON)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Rewritten from:
// https://github.com/kivattt/kivadiscordbridge/blob/e7cd6448056eab8ae7e1d4d2b6f331813f6b45c9/src/server/java/com/kiva/kivadiscordbridge/DiscordAPI.java#L37C34-L37C35
func sanitizeForJSON(str string) string {
	disallowedCharactersForJSON := "\"\\"

	var builder strings.Builder
	for _, c := range str {
		if strings.ContainsAny(string(c), disallowedCharactersForJSON) {
			builder.WriteRune('\\')
		}

		builder.WriteRune(c)
	}

	return builder.String()
}

func main() {
	h := flag.Bool("help", false, "display this help and exit")
	v := flag.Bool("version", false, "output version information and exit")
	list := flag.Bool("list", false, "list all API endpoints")
	endPoint := flag.String("endpoint", "", "request endpoint")
	all := flag.Bool("all", false, "request all endpoints")
	password := flag.String("password", "", "router admin password")
	output := flag.String("output", "", "output json to file")

	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.CommandLine.Init(PROGRAM_NAME, flag.ExitOnError)
	getopt.Aliases(
		"h", "help",
		"e", "endpoint",
		"l", "list",
		"a", "all",
		"p", "password",
		"o", "output",
		"v", "version",
	)

	err := getopt.CommandLine.Parse(os.Args[1:])
	if err != nil {
		os.Exit(2)
	}

	if *v {
		fmt.Println(PROGRAM_NAME, VERSION)
		fmt.Println()
		fmt.Println("CPE Insight API version supported:", CPE_INSIGHT_API_VERSION)
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

	if *endPoint == "" && !*all {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [OPTIONS]")
		fmt.Println("CPE Insight API explorer")
		fmt.Println()
		getopt.PrintDefaults()
		os.Exit(0)
	}

	if *endPoint != "" && *all {
		fmt.Println("Both --all and --endpoint were used.")
		fmt.Println("Pick one option!")
		os.Exit(1)
	}

	if *password == "" {
		fmt.Println("Password required, try:")
		fmt.Println()
		fmt.Println("    " + filepath.Base(os.Args[0]) + " --password=\"...\" " + strings.Join(os.Args[1:], " "))
		os.Exit(1)
	}

	cpe, err := getCPEInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cookies, err := getCPECookieToken(cpe.username, *password, cpe)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *endPoint != "" {
		resp, err := getRequest(*endPoint, cookies, cpe)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bytes, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if *output == "" {
			fmt.Print(string(bytes))
		} else {
			err := os.WriteFile(*output, bytes, 0o664)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	} else if *all {
		var builder strings.Builder
		builder.WriteString("{")
		for endPointKey := range apiEndPoints {
			os.Stderr.WriteString("\x1b[0;34mRequesting: " + endPointKey + "\x1b[0m\n")
			builder.WriteString("\"" + endPointKey + "\":")
			resp, err := getRequest(endPointKey, cookies, cpe)
			if err != nil {
				builder.WriteString("\"Error requesting: " + sanitizeForJSON(err.Error()) + "\"")
			} else {
				builder.Write(resp)
			}

			builder.WriteString(",")
		}
		builder.WriteString("}")

		if *output == "" {
			fmt.Println(builder.String())
		} else {
			err := os.WriteFile(*output, []byte(builder.String()), 0o664)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	} else {
		panic("No --endpoint or --all option passed")
	}
}
