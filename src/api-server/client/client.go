package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// Structure for our authentication for JSON
type auth struct {
	Username string
	Password string
}

//Structure for the token return
type token struct {
	Token string `json:"token"`
}

// This function will hopefully display a welcome message
// based on the authentication token provided in login

func goRestricted(host string, port string, tk string) {
	url := fmt.Sprintf("http://%s:%s/restricted", host, port)

	auth := fmt.Sprintf("Bearer %s", tk)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

func goAdd(host string, port string, tk string) {
	url := fmt.Sprintf("http://%s:%s/restricted/add", host, port)

	auth := fmt.Sprintf("Bearer %s", tk)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

func listTasks(host string, port string, tk string) {
	url := fmt.Sprintf("http://%s:%s/restricted/tasks", host, port)

	auth := fmt.Sprintf("Bearer %s", tk)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

// This function will log you in via Json payload and return an auth token
// if successfull

func loginJSON(host string, port string, username string, password string) string {

	url := fmt.Sprintf("http://%s:%s/login", host, port)

	cred := auth{username, password}
	jsonStr, _ := json.Marshal(cred)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var t = new(token)
	err = json.Unmarshal(body, &t)

	if err != nil {
		log.Fatal(err)
	}
	return t.Token
}
func main() {

	viper.SetConfigName("config") // no need to include file extension
	viper.AddConfigPath("/Users/denn8098/GoProjects/doozer/src/api-server/client/")

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}

	host := viper.GetString("config.host")
	port := viper.GetString("config.port")
	username := viper.GetString("config.username")
	password := viper.GetString("config.password")

	token := loginJSON(host, port, username, password)

	// If we didn't get a token back, then error out
	if token == "" {
		log.Fatal(fmt.Errorf("Can't get Auth token. Check username and password in config file"))
	}

	//goRestricted(host, port, token)
	goAdd(host, port, token)
	listTasks(host, port, token)
}
