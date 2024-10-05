package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Ginのデフォルトのルーターを作成
	r := gin.Default()

	// MongoDBに接続
	err := connectMongoDB()
	if err != nil {
		log.Fatal("MongoDB接続エラー:", err)
	}

	// /pingへGETでハンドラーを作成
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "OK",
		})
	})

	// サーバーを起動
	r.Run()
}

func connectMongoDB() error {
	// MongoDB接続のコンテキストと接続オプションの設定
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDBに接続
	// コンテナ内のMongoDBに接続
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 接続確認
	pingErr := client.Ping(ctx, nil)
	if pingErr != nil {
		return pingErr
	}

	fmt.Println("MongoDBに接続")
	return nil
}
