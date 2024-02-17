package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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
	analyzedmsg,err := msgAnalyzer.analyzeMsg(msg.Content)
	if err!=nil{
		analyzedmsg="hard coded message"
	}

	c.JSON(http.StatusOK, analyzedmsg)
}



func (msgA *MessageAnalyzer) analyzeMsg(msg string) (string,error) {
	log.Println("analyzemsg request came")
	serviceurl:=os.Getenv("MESSAGE_PROCESSING_SERVICE_URL")
	if serviceurl==""{
		log.Println("service url is empty")
		return "",fmt.Errorf("error in reading environment variable %v","MESSAGE_PROCESSING_SERVICE_URL")
	}
	requestBody,err:=json.Marshal(Message{Content: msg})
	if err!=nil{
		log.Println("marshalling failed",err)
		return "",err
	}
	resp,err:=http.Post(serviceurl,"application/json",bytes.NewBuffer(requestBody))
	if err!=nil{
		log.Printf("post to service %s failed with error %v\n ",serviceurl,err)
		return "",err
	}
	defer resp.Body.Close()
	body,err:=io.ReadAll(resp.Body)
	if err!=nil{
		log.Printf("can read from service %s with error %v\n ",serviceurl,err)
		return "",err
	}
	return string(body),nil
	
}
