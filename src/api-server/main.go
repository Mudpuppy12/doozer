package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
	"github.com/RichardKnop/machinery/v1/signatures"
)

type user struct {
	Name     string `json:"username" form:"username" query:"usernamename"`
	Password string `json:"password" form:"password" query:"password"`
}

type uuid struct {
	UUID string `json:"uuid" form:"uuid" query:"uuid"`
}

type numbers struct {
	Numbers string `json:"numbers" form:"numbers" query:"numbers"`
}

var (
	broker        string
	resultBackend string
	exchange      string
	exchangeType  string
	defaultQueue  string
	bindingKey    string
	server        *machinery.Server
	task0         signatures.TaskSignature
	cnf           config.Config
)

func sAtoi(stslice string) []int {
	strs := strings.Split(stslice, ",")
	ary := make([]int, len(strs))
	for i := range ary {
		ary[i], _ = strconv.Atoi(strs[i])
	}
	return ary
}

func init() {

	viper.SetConfigName("config") // no need to include file extension
	viper.AddConfigPath("/Users/denn8098/GoProjects/doozer/src/api-server")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}

	broker = viper.GetString("dozer.broker")
	resultBackend = viper.GetString("dozer.result_backend")
	exchange = viper.GetString("dozer.exchange")
	exchangeType = viper.GetString("dozer.exchange_type")
	defaultQueue = viper.GetString("dozer.default_queue")
	bindingKey = viper.GetString("dozer.binding_key")

	cnf = config.Config{
		Broker:        broker,
		ResultBackend: resultBackend,
		Exchange:      exchange,
		ExchangeType:  exchangeType,
		DefaultQueue:  defaultQueue,
		BindingKey:    bindingKey,
	}

	server, err = machinery.NewServer(&cnf)
	errors.Fail(err, "Could not initialize server")
}

func login(c echo.Context) (err error) {

	u := new(user)

	if err = c.Bind(u); err != nil {
		return
	}

	username := u.Name
	password := u.Password

	if username == "mudpuppy" && password == "dirtypaws" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "Mudpuppy"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func apiTask(c echo.Context) (err error) {

	u := new(uuid)

	if err = c.Bind(u); err != nil {
		return
	}

	task := u.UUID

	tasknames, err := server.GetBackend().GetState(task)

	if err != nil {
		return c.String(http.StatusBadRequest, "Error: Task not found")
	}

	if tasknames.State == "PENDING" {
		return c.String(http.StatusOK, "Status : PENDING")
	}

	if tasknames.State == "RECEIVED" {
		return c.String(http.StatusOK, "Status : RECEIVED")
	}

	if tasknames.State == "STARTED" {
		return c.String(http.StatusOK, "Status : STARTED")
	}

	if tasknames.State == "FAILURE" {
		return c.String(http.StatusOK, "Status : FAILURE")
	}

	if tasknames.State == "SUCCESS" {
		result := fmt.Sprintf("%v", tasknames.Result.Value)
		return c.String(http.StatusOK, "Status : SUCCESS\nResult : "+result)
	}

	return c.String(http.StatusBadRequest, "ERROR: Something broken")

}

func apiAdd(c echo.Context) (err error) {

	u := new(numbers)
	//var args []signatures.TaskArg

	if err = c.Bind(u); err != nil {
		return
	}

	nbrs := sAtoi(u.Numbers)

	args := []signatures.TaskArg{}

	for _, v := range nbrs {
		args = append(args, signatures.TaskArg{Type: "int64", Value: v})
	}

	task0 = signatures.TaskSignature{
		Name: "add",
		Args: args,
	}

	asyncResult, err := server.SendTask(&task0)
	errors.Fail(err, "Could not send task")

	result, err := asyncResult.GetWithTimeout(5000000000, 1)

	if err != nil { // Handle errors reading the config file
		taskState := asyncResult.GetState()
		return c.String(http.StatusOK, "Defered! "+taskState.TaskUUID+"")
	}

	r := fmt.Sprintf("%v", result.Interface())
	return c.String(http.StatusOK, "Result: "+r+"")

}

func apiMul(c echo.Context) (err error) {

	u := new(numbers)
	//var args []signatures.TaskArg

	if err = c.Bind(u); err != nil {
		return
	}

	nbrs := sAtoi(u.Numbers)

	args := []signatures.TaskArg{}

	for _, v := range nbrs {
		args = append(args, signatures.TaskArg{Type: "int64", Value: v})
	}

	task0 = signatures.TaskSignature{
		Name: "multiply",
		Args: args,
	}

	asyncResult, err := server.SendTask(&task0)
	errors.Fail(err, "Could not send task")

	result, err := asyncResult.GetWithTimeout(5000000000, 1)

	if err != nil { // Handle errors reading the config file
		taskState := asyncResult.GetState()
		return c.String(http.StatusOK, "Defered! "+taskState.TaskUUID+"")
	}

	r := fmt.Sprintf("%v", result.Interface())
	return c.String(http.StatusOK, "Result: "+r+"")

}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)
	r.POST("/add", apiAdd)
	r.POST("/mul", apiMul)
	r.POST("/tasks", apiTask)

	e.Logger.Fatal(e.Start(":1323"))
}
