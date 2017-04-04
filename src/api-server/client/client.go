package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// Post like a Form data Values

type auth struct {
	Username string
	Password string
}

type token struct {
	Token string `json:"token"`
}

func loginJSON(host string, port string, username string, password string) string {

	url := fmt.Sprintf("http://%s:%s/login", host, port)

	cred := auth{username, password}
	jsonStr, _ := json.Marshal(cred)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var t = new(token)
	err = json.Unmarshal(body, &t)

	if err != nil {
		fmt.Println("whoops:", err)
	}
	return t.Token
}
func main() {

	viper.SetConfigName("config") // no need to include file extension
	viper.AddConfigPath("/Users/denn8098/GoProjects/doozer/src/api-server/client/")

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	host := viper.GetString("config.host")
	port := viper.GetString("config.port")
	username := viper.GetString("config.username")
	password := viper.GetString("config.password")

	token := loginJSON(host, port, username, password)

	if token == "" {
		panic(fmt.Errorf("Can't get Auth token. Check username and password in config"))
	}

	fmt.Println(token)

}
