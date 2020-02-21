package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	json_server "github.com/pytorchtw/go-json-server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func NewPost() *json_server.Post {
	post := json_server.Post{}
	return &post
}

func init() {
	for count := 0; count < 10; count++ {
		post := NewPost()
		post.Id = count
		json_server.PostStore[count] = post
	}
}

func main() {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)

	srv := StartServer(httpServerExitDone)
	GracefulShutdown(srv, 1000000)
	httpServerExitDone.Wait()
}

func StartServer(wg *sync.WaitGroup) *http.Server {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/posts", json_server.GetPosts)

	srv := &http.Server{Addr: ":8081", Handler: router}
	go func() {
		defer wg.Done() // let main know we are done cleaning up
		log.Printf("serving...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	return srv
}

func GracefulShutdown(srv *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nshutdown with timeout: %s\n", timeout)

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error: %v\n", err)
	} else {
		log.Println("server gracefully stopped")
	}
}
