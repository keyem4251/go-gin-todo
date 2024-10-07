package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ToDoのモデルを定義
type ToDo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Completed bool               `json:"completed" bson:"completed"`
}

var client *mongo.Client

func main() {
	// Ginのデフォルトのルーターを作成
	r := gin.Default()

	// MongoDBに接続
	var err error
	client, err = connectMongoDB()
	if err != nil {
		log.Fatal("MongoDB接続エラー:", err)
	}

	// // エンドポイントの定義
	// r.POST("/todos", createToDo)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "OK")
	})

	// サーバーを起動
	r.Run()
}

func connectMongoDB() (*mongo.Client, error) {
	// MongoDB接続のコンテキストと接続オプションの設定
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDBに接続
	// コンテナ内のMongoDBに接続
	url := os.Getenv("MONGO_URL")
	log.Println(url)
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	log.Println(clientOptions.GetURI())
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 接続確認
	pingErr := c.Ping(ctx, nil)
	if pingErr != nil {
		return nil, pingErr
	}

	fmt.Println("MongoDBに接続")
	return c, nil
}

func createToDo(c *gin.Context) {
	var todo ToDo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("todoapp").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データベースに保存できませんでした"})
		return
	}

	c.JSON(http.StatusOK, result)
}
