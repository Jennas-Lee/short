package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "short/mainpb"
)

type Url struct {
	Url string
}

func main() {
	port := flag.Int("port", 8080, "listen port")
	uid := flag.String("uid", "localhost:8080", "uid address")
	flag.Parse()

	conn, err := grpc.NewClient(*uid, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUidGeneratorClient(conn)

	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	r.POST("/api/short", func(c *gin.Context) {
		var url Url

		if c.ShouldBindBodyWithJSON(&url) == nil {
			log.Println(url.Url)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := client.Uid(ctx, &pb.Request{})
			if err != nil {
				log.Fatalf("could not uid: %v", err)

				c.String(http.StatusInternalServerError, "there's some error.")
			} else {
				// log.Printf("%d %s\n", r.Id, base62.FormatInt(int64(r.Id)))

				c.String(http.StatusOK, "shorturl: %d", r.Id)
			}
		} else {
			c.String(http.StatusBadRequest, "Invalid url")
		}
	})

	r.GET("/:short", func(c *gin.Context) {
		url := c.Param("short")
		c.String(http.StatusOK, "we received: %s", url)
	})

	r.Run(fmt.Sprintf(":%d", *port))
}
