package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	// MongoDBに接続
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDBの接続URIを取得
	mongoURI := "mongodb://mongo:27017"

	// MongoDB接続オプションの設定
	clientOptions := options.Client().ApplyURI(mongoURI)

	// MongoDBに接続
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB接続エラー: %v", err)
	}

	// MongoDBにPingして接続確認
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB接続に失敗しました: %v", err)
	}
	fmt.Println("MongoDBに正常に接続されました！")

	// Ginのデフォルトルーターを作成
	r := gin.Default()
	log.Println("起動テスト")
	// ルートエンドポイントの設定
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "GinとMongoDBが正常に動作しています",
		})
	})

	// サーバーをポート8080で起動
	r.Run(":8080")
}
