package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type authResponse struct {
	AccesToken string      `json:"access_token"`
	TokenType  string      `json:"token_type"`
	ExpiresIn  json.Number `json:"expires_in"`
	Scope      string      `json:"scope"`
	UUID       string      `json:"uid"`
	AccountId  string      `json:"account_id"`
}

type readFilesResponse struct {
	entries map[string]string `json:"entries"`
	cursor  string
	hasMore bool
}

func main() {
	var client_id string
	client_id = "6f0oswzr0go1qek"

	fmt.Println("\033[33mOpen this link to start the initialization process:")
	fmt.Println("\033[37mhttps://www.dropbox.com/oauth2/authorize?client_id=" + client_id + "&response_type=code")

	var accessCode string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		accessCode = scanner.Text()

		if len(accessCode) > 0 {
			break
		}
	}

	client := &http.Client{}
	data := url.Values{
		"code":       {accessCode},
		"grant_type": {"authorization_code"},
	}

	req, _ := http.NewRequest(http.MethodPost, "https://api.dropboxapi.com/oauth2/token", strings.NewReader(data.Encode()))
	req.SetBasicAuth("6f0oswzr0go1qek", "81hyir0csho9iei")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer req.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var result authResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
	}

	fmt.Println("\033[33mAccess Token:")
	fmt.Println("\033[37m" + result.AccesToken)

	requestBody, err := json.Marshal(map[string]interface{}{
		"path":                                "",
		"recursive":                           false,
		"include_media_info":                  false,
		"include_deleted":                     false,
		"include_has_explicit_shared_members": false,
		"include_mounted_folders":             true,
		"include_non_downloadable_files":      true,
	})
	if err != nil {
		fmt.Println(err)
	}

	req, err = http.NewRequest(http.MethodPost, "https://api.dropboxapi.com/2/files/list_folder", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "Bearer "+result.AccesToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &res); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
	}

	requestBody, err = json.Marshal("this is a test")
	if err != nil {
		fmt.Println(err)
	}

	req, err = http.NewRequest(http.MethodPost, "https://content.dropboxapi.com/2/files/upload", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "Bearer "+result.AccesToken)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Dropbox-API-Arg", `{"path": "/test.txt", "mode":"overwrite"}`)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\033[33mFile created\033[37m")
}
