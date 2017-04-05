package main

import (
	"fmt"
	"log"
	"net/http"
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

var (
	broker                                          string
	resultBackend                                   string
	exchange                                        string
	exchangeType                                    string
	defaultQueue                                    string
	bindingKey                                      string
	server                                          *machinery.Server
	task0, task1, task2, task3, task4, task5, task6 signatures.TaskSignature
	cnf                                             config.Config
)

func init() {

	viper.SetConfigName("config") // no need to include file extension
	viper.AddConfigPath("/Users/denn8098/GoProjects/doozer/src/api-server")

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

func apiTask(c echo.Context) error {

	tasknames, _ := server.GetBackend().GetState("task_d7e33c78-6888-4e87-8a26-18a83b98fa95")
	result := fmt.Sprintf("%v", tasknames.Result.Value)
	return c.String(http.StatusOK, "Results:"+result)

}

func apiAdd(c echo.Context) error {
	//user := c.Get("user").(*jwt.Token)
	//claims := user.Claims.(jwt.MapClaims)
	//name := claims["name"].(string)

	task0 = signatures.TaskSignature{
		Name: "add",
		Args: []signatures.TaskArg{
			{
				Type:  "int64",
				Value: 1,
			},
			{
				Type:  "int64",
				Value: 1,
			},
		},
	}

	asyncResult, err := server.SendTask(&task0)
	fmt.Printf("%v", asyncResult)
	errors.Fail(err, "Could not send task")

	result, err := asyncResult.GetWithTimeout(5000000000, 1)

	if err != nil { // Handle errors reading the config file
		taskState := asyncResult.GetState()
		fmt.Printf("Current state of %v task is:\n", taskState.TaskUUID)
		fmt.Println(taskState.State)
		return c.String(http.StatusOK, "Defered! "+taskState.TaskUUID+"")
	}

	zippy := fmt.Sprintf("1 + 1 = %v", result.Interface())
	return c.String(http.StatusOK, "Add! "+zippy+"")

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
	r.GET("/add", apiAdd)
	r.GET("/tasks", apiTask)

	e.Logger.Fatal(e.Start(":1323"))
}
