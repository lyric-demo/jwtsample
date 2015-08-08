package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	// PORT 监听端口
	PORT = 9999
	// SECUREKEY 密钥
	SECUREKEY = "7FAA2DAD-0C8A-4A4D-A098-52D0CCE4A6FA"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"Status": 9, "Error": err.Error()})
	})

	e.ServeFile("/login", "static/login.html")
	e.ServeFile("/", "static/index.html")

	e.Post("/login", Login)
	api := e.Group("/api", func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			tokenString := c.Request().Header.Get("x-access-token")
			if tokenString == "" {
				return errors.New("Token isvalid")
			}
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if token.Claims["uid"] != "admin" {
					return nil, errors.New("User isvalid")
				}
				return []byte(SECUREKEY), nil
			})
			if err != nil {
				return err
			}
			c.Set("Uid", token.Claims["uid"])
			return h(c)
		}
	})
	{
		api.Get("/data", GetData)
	}

	log.Printf("Server is running at %d port.\n", PORT)
	e.Run(fmt.Sprintf(":%d", PORT))
}

// Login 登录
func Login(c *echo.Context) error {
	userName := c.Form("UserName")
	pwd := c.Form("Password")
	if (userName == "admin" || userName == "000000") && pwd == "admin" {
		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims["uid"] = userName
		token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
		tokenStr, err := token.SignedString([]byte(SECUREKEY))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"Token": tokenStr, "Status": 1})
	}
	return errors.New("登录错误！")
}

// GetData 获取测试数据
func GetData(c *echo.Context) error {
	fmt.Println("===> Request user:", c.Get("Uid"))
	data := make([]string, 100)
	for i := 0; i < 100; i++ {
		data[i] = fmt.Sprintf("Display Content %d", i+1)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"Status": 1,
		"Data":   data,
	})
}
