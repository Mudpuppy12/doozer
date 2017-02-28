// curl is a simple cURL replacement.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Post like a Form data Values

func postform() {
	apiUrl := "http://localhost:1323"
	resource := "/login"

	data := url.Values{}

	data.Set("username", "mudpuppy")
	data.Set("password", "dirtypaws")

	u, _ := url.ParseRequestURI(apiUrl)

	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Println(resp.Status)
}

// Example using Json input

func postjson() {
	url := "http://localhost:1323/login"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"username":"mudpuppy", "password":"dirtypaws"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main() {
	postform()
}
