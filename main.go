package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// 关闭console日志
	gin.DisableConsoleColor()
	// 记录为文件日志
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// 设置模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/sign", func(c *gin.Context) {
		// 接收四个参数 接口参数 + nonce + timeStamp + apiSecret
		body := c.DefaultQuery("body", "")
		// 随机数
		Nonce := c.DefaultQuery("nonce", "")
		TimeStamp := c.DefaultQuery("timeStamp", "")
		apiSecret := c.DefaultQuery("apiSecret", "")

		if body == "" || Nonce == "" || TimeStamp == "" || apiSecret == "" {
			c.JSON(200, gin.H{
				"code":   400,
				"data":   nil,
				"reason": "error",
			})
			return
		}
		timeStamp, _ := strconv.Atoi(TimeStamp)
		signStr := body + Nonce + strconv.FormatInt(int64(timeStamp), 10)
		// fmt.Printf("签名字符串： %s \n", signStr)
		sign := sign(signStr, apiSecret)
		// fmt.Printf("签名结果：%s \n", sign)
		// gin.LoggerWithWriter(f, "签名字符串： %s \n", signStr)
		c.JSON(200, gin.H{
			"code":   0,
			"data":   map[string]string{"sign": sign},
			"reason": "success",
		})
	})
	r.GET("/", func(c *gin.Context) {
		time.Sleep(time.Second * 5)
		c.String(200, "hello")
	})
	// _ = r.Run(":8080") // listen and serve on 0.0.0.0:8080
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

// sign 签名
func sign(str string, privateKey string) string {
	str = toHash256(str)
	getPrivateKey := stringToPrivateKey(privateKey)
	bys := []byte(str)
	r, s, err := ecdsa.Sign(rand.Reader, getPrivateKey, bys)
	if err != nil {
		return ""
	}
	rb := hex.EncodeToString(r.Bytes())
	sb := hex.EncodeToString(s.Bytes())
	return rb + sb
}

// toHash256 对请求内容hash
func toHash256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return string(hash.Sum([]byte(nil)))
}

// stringTOPrivateKey string->私钥
func stringToPrivateKey(str string) *ecdsa.PrivateKey {
	var newCurve = elliptic.P256()
	prv, _ := hex.DecodeString(str)
	if len(prv) == 0 {
		return nil
	}
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = newCurve
	privateKey.D = new(big.Int).SetBytes(prv)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(prv)
	return privateKey
}
