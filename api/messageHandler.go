package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)
type MessageAnalyzer struct {
}

var msgAnalyzer MessageAnalyzer

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
	type Message struct {
		Content string `json:"content"`
	}
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	analyzedmsg := msgAnalyzer.analyzeMsg(msg.Content)

	c.JSON(http.StatusOK, analyzedmsg)
}



func (msgA *MessageAnalyzer) analyzeMsg(msg string) string {
	return " Hallo " + "" + msg
}
