//// API client - A Work in progress learning to write a client using authentication
/// JWT and command line config file reading.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "1.0"

// Structure for our authentication for JSON
type auth struct {
	Username string
	Password string
}

//Structure for the token return
type token struct {
	Token string `json:"token"`
}

type uuid struct {
	UUID string `json:"uuid"`
}

type numbers struct {
	Numbers string `json:"Numbers" form:"Numbers" query:"Numbers"`
}

// Here we start defining our client commands w/ Cobra
// mainCmd is what is issued when someone just types client with no arguments

var mainCmd = &cobra.Command{
	Use:   "client",
	Short: "Dozer api client",
	Long:  "Simple client to interact with Dozer API service.",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version.",
	Long:  "The version of the dispatch service.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add api call.",
	Long:  "Call the api to task a worker to add numbers.",
	Run: func(cmd *cobra.Command, args []string) {

		host := viper.GetString("config.host")
		port := viper.GetString("config.port")
		username := viper.GetString("config.username")
		password := viper.GetString("config.password")

		token := loginJSON(host, port, username, password)
		numberStr := (viper.GetString("i"))

		// If we didn't get a token back, then error out
		if token == "" {
			log.Fatal(fmt.Errorf("Can't get Auth token. Check username and password in config file"))
		}

		goAdd(host, port, token, numberStr)

	},
}

var mulCmd = &cobra.Command{
	Use:   "mul",
	Short: "Multiply api call.",
	Long:  "Call the api to task a worker to muliply numbers.",
	Run: func(cmd *cobra.Command, args []string) {

		host := viper.GetString("config.host")
		port := viper.GetString("config.port")
		username := viper.GetString("config.username")
		password := viper.GetString("config.password")

		token := loginJSON(host, port, username, password)
		numberStr := (viper.GetString("i"))

		// If we didn't get a token back, then error out
		if token == "" {
			log.Fatal(fmt.Errorf("Can't get Auth token. Check username and password in config file"))
		}

		goMul(host, port, token, numberStr)

	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print a JWT token.",
	Long:  "Print out a JWT token from a successful login",
	Run: func(cmd *cobra.Command, args []string) {

		host := viper.GetString("config.host")
		port := viper.GetString("config.port")
		username := viper.GetString("config.username")
		password := viper.GetString("config.password")

		token := loginJSON(host, port, username, password)

		if token == "" {
			log.Fatal(fmt.Errorf("Can't get Auth token. Check username and password in config file"))
		}

		fmt.Printf("Your JWT Token is :[%s]\n", token)
	},
}

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Lookup a task uuid.",
	Long:  "Looku a task uuid and return the result",
	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("config.host")
		port := viper.GetString("config.port")
		username := viper.GetString("config.username")
		password := viper.GetString("config.password")

		token := loginJSON(host, port, username, password)
		goLookup(host, port, token)

	},
}

/////// * API function calls below /////

func goLookup(host string, port string, tk string) {
	url := fmt.Sprintf("http://%s:%s/restricted/tasks", host, port)

	id := uuid{viper.GetString("uuid")}
	jsonStr, _ := json.Marshal(id)

	auth := fmt.Sprintf("Bearer %s", tk)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
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

func goAdd(host string, port string, tk string, nbrs string) {
	url := fmt.Sprintf("http://%s:%s/restricted/add", host, port)

	auth := fmt.Sprintf("Bearer %s", tk)

	data := numbers{nbrs}

	jsonStr, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
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

func goMul(host string, port string, tk string, nbrs string) {
	url := fmt.Sprintf("http://%s:%s/restricted/mul", host, port)

	auth := fmt.Sprintf("Bearer %s", tk)

	data := numbers{nbrs}

	jsonStr, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
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

func runCmd() {
	host := viper.GetString("config.host")
	port := viper.GetString("config.port")
	username := viper.GetString("config.username")
	password := viper.GetString("config.password")

	token := loginJSON(host, port, username, password)
	goRestricted(host, port, token)
}

func init() {

	viper.SetConfigName("config") // no need to include file extension
	viper.AddConfigPath("/Users/denn8098/GoProjects/doozer/src/api-server/client/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}

	// Adding commands into the client

	mainCmd.AddCommand(versionCmd)
	mainCmd.AddCommand(addCmd)
	mainCmd.AddCommand(mulCmd)
	mainCmd.AddCommand(tokenCmd)
	mainCmd.AddCommand(lookupCmd)

	lookupFlags := lookupCmd.Flags()
	addFlags := addCmd.Flags()
	mulFlags := mulCmd.Flags()

	lookupFlags.String("uuid", "", "uuid of task you want to lookup.")
	viper.BindPFlag("uuid", lookupFlags.Lookup("uuid"))

	addFlags.String("i", "", "integers to add")
	viper.BindPFlag("i", addFlags.Lookup("i"))

	mulFlags.String("i", "", "integers to multiply")
	viper.BindPFlag("i", mulFlags.Lookup("i"))

}

func main() {
	mainCmd.Execute()

}
