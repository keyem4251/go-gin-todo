package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Completed bool               `json:"completed" bson:"completed"`
}

func main() {
	// MongoDBに接続
	client, err := connectDB()
	if err != nil {
		log.Fatal("MongoDB接続エラー")
	}
	collection = createCollection(client)

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
	r.GET("/todos", getTodos)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/totos/:id", deleteTodo)

	// サーバーをポート8080で起動
	r.Run(":8080")
}

func connectDB() (*mongo.Client, error) {
	// MongoDBに接続
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB接続オプションの設定
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")

	// MongoDBに接続
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB接続エラー: %v", err)
		return nil, err
	}

	// MongoDBにPingして接続確認
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB接続に失敗しました: %v", err)
		return nil, err
	}

	return client, nil
}

func createCollection(c *mongo.Client) *mongo.Collection {
	collection := c.Database("todoapp").Collection("todos")
	return collection
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

	result, err := collection.InsertOne(ctx, todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo作成中にエラー"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func getTodos(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo取得中にエラー"})
		return
	}
	defer cursor.Close(ctx)

	var todos []Todo
	if err := cursor.All(ctx, &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo解析中にエラー"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
		return
	}

	var todo Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効な入力"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": todo})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo更新中にエラー"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo削除中にエラー"})
		return
	}
	c.JSON(http.StatusOK, result)
}
