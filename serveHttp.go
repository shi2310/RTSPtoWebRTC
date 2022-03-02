package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func serveHTTP() {
	router := gin.Default()
	if _, err := os.Stat("./web"); !os.IsNotExist(err) {
		router.LoadHTMLGlob("web/templates/*")

		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{})
		})
	}
	router.POST("/receiver", reciver)
	router.StaticFS("/static", http.Dir("web/static"))
	err := router.Run(":8083")
	if err != nil {
		log.Fatalln(err)
	}
}

func reciver(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	data := c.PostForm("data")
	rtspUrl := c.PostForm("rtspUrl")

	// 声明空对象
	stream := StreamST{URL: rtspUrl}
	// 建立rtsp会话
	session, err := rtspSession(&stream)
	if err != nil {
		log.Println(err)
		return
	}
	// 建立webrtc的响应
	answerSdp, err := webrtcAnswer(data, &stream, session)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(200, answerSdp)
}
