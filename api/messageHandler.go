package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"simpleMsgSvc/pkg"

	"github.com/gin-gonic/gin"
)

type MessageAnalyzer struct {
}

var msgAnalyzer MessageAnalyzer

type Message struct {
	Content string `json:"content"`
}

func init() {
	msgAnalyzer = MessageAnalyzer{}
}
func SimpleMessagesAPI() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		pingHandler(c)
	})
	router.POST("/message", func(c *gin.Context) {
		messageHandler(c)
	})
	router.Run(":8080")
}
func pingHandler(c *gin.Context) {
	log.Println("request on ping endpoint")
	c.String(http.StatusOK, "ping message to simple messages service")
}
func messageHandler(c *gin.Context) {
	log.Println("request on message end point")

	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	val, err := pkg.GetValue(msg.Content)
	if err != nil {
		fmt.Println("error in reading the value from redis cache ", err)
	}
	var analyzedmsg string
	if val != "" {
		log.Printf("value for %s is %s ", msg.Content, val)
		analyzedmsg = val
	} else {
		analyzedmsg, err = msgAnalyzer.analyzeMsg(c, msg.Content)
		if err != nil {
			analyzedmsg = fmt.Sprintf("error in /message end point %s", err)
		}
	}
	//write to redis cache
	err = pkg.SetValue(msg.Content, analyzedmsg)
	if err != nil {
		fmt.Println("error in setting value in redis cache ", err)
	}
	c.JSON(http.StatusOK, analyzedmsg)
}

func (msgA *MessageAnalyzer) analyzeMsg(c *gin.Context, msg string) (string, error) {
	serviceURL := os.Getenv("MESSAGE_PROCESSING_SERVICE_URL")
	if serviceURL == "" {
		log.Println("service url is empty")
		return "", fmt.Errorf("error in reading environment variable %v", "MESSAGE_PROCESSING_SERVICE_URL")
	}
	requestBody, err := json.Marshal(Message{Content: msg})
	if err != nil {
		log.Println("marshalling failed", err)
		return "", err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", serviceURL, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("failed to create request to service %s with error %v\n", serviceURL, err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Copy tracing headers
	for _, h := range []string{"x-request-id", "x-b3-traceid", "x-b3-spanid", "x-b3-parentspanid", "x-b3-sampled", "x-b3-flags", "x-ot-span-context"} {
		if value := c.GetHeader(h); value != "" {
			req.Header.Set(h, value)
		}
	}

	// Execute the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("post to service %s failed with error %v\n", serviceURL, err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("can't read from service %s with error %v\n", serviceURL, err)
		return "", err
	}
	return string(body), nil
}
