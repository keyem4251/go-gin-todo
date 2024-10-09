package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Completed bool               `json:"completed" bson:"completed"`
}

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

	// Ginのデフォルトルーターを作成
	r := gin.Default()
	log.Println("起動テスト")
	// ルートエンドポイントの設定
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "GinとMongoDBが正常に動作しています",
		})
	})

	r.POST("/todos", createTodo)

	// サーバーをポート8080で起動
	r.Run(":8080")
}

// todoを作成
func createTodo(c *gin.Context) {
	var todo Todo

	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な入力"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database("todoapp").Collection("todos")
	result, err := collection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo作成中にエラー"})
		return
	}
	c.JSON(http.StatusOK, result)
}
