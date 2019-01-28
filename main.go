package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func handle(handler Handler, c *routing.Context, okCode int, badRequestBody string) {
	c.SetContentType("application/json")
	c.Response.Header.Set("Connection", "Keep-Alive")

	err := handler.ServeHTTP(c)

	if err != nil {
		code := err.(*Error).Code
		c.SetStatusCode(code)

		if code == 400 || code == 404 {
			c.WriteString(badRequestBody)
		}

		if code == 500 {
			log.Println(err.Error())
		}
		// log.Println(err.Error())
	} else {
		c.SetStatusCode(okCode)
	}
}

func top() {
	// cmd := exec.Command("top", "-b", "-n", "1")
	cmd := exec.Command("cat", "/sys/fs/cgroup/memory/memory.usage_in_bytes")
	stdout, err := cmd.Output()
	if err == nil {
		log.Println("usage: ", string(stdout))
	}
}

func topLoop() {
	top()
	time.Sleep(60 * time.Second)
	go topLoop()
}

func main() {
	start := time.Now()
	log.Println("started!")

	go topLoop()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Println("caught sig: ", sig)
		os.Exit(0)
	}()

	storage := &MongoStorage{
		Uri:      MongoUri,
		Database: "hl",
	}
	storage.Init()

	go func() {
		err := Parse(DataFile, OptionsFile, storage, false)
		if err != nil {
			panic(err)
		}
		log.Println("Parsed", time.Now().Sub(start).Seconds())
	}()

	router := routing.New()
	router.NotFound(func(c *routing.Context) error {
		c.SetStatusCode(404)
		return nil
	})

	afh := &AccountsFilterHandler{storage}
	router.Get("/accounts/filter/", func(c *routing.Context) error {
		handle(afh, c, 200, "{}")
		return nil
	})

	agh := &AccountsGroupHandler{storage}
	router.Get("/accounts/group/", func(c *routing.Context) error {
		handle(agh, c, 200, "{}")
		return nil
	})

	arh := &AccountsRecommendHandler{storage}
	router.Get("/accounts/<id>/recommend/", func(c *routing.Context) error {
		handle(arh, c, 200, "{}")
		return nil
	})

	ash := &AccountsSuggestHandler{storage}
	router.Get("/accounts/<id>/suggest/", func(c *routing.Context) error {
		handle(ash, c, 200, "{}")
		return nil
	})

	newAccountHandler := &AccountsNewHandler{storage}
	router.Post("/accounts/new/", func(c *routing.Context) error {
		handle(newAccountHandler, c, 201, "")
		return nil
	})

	updateAccountHandler := &AccountsUpdateHandler{storage}
	router.Post("/accounts/<id>/", func(c *routing.Context) error {
		handle(updateAccountHandler, c, 202, "")
		return nil
	})

	addLikesHandler := &AccountsLikesHandler{storage}
	router.Post("/accounts/likes/", func(c *routing.Context) error {
		handle(addLikesHandler, c, 202, "")
		return nil
	})

	fasthttp.ListenAndServe(":80", router.HandleRequest)
}
