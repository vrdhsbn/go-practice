package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// album の構造体
// json:"id" みたいなのはJSONにシリアライズされるときにプロパティ名として用いられる
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// 初期データ
// サイズを指定していないのでこれは配列ではなくスライス
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	// Gin router を初期化
	router := gin.Default()

	// エンドポイントをGET/POSTそれぞれに定義する（ハンドラ関数と紐付ける）
	// :を付けるとパスパラメータとして解釈される
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.GET("/", welcome)

	// HTTPサーバを起動する
	// Dockerで動かす場合はlocalhostではなく0.0.0.0と記載する
	router.Run("0.0.0.0:8080")
}

// GETでアクセスした場合の処理。albums の内容をJSONとして返す
func getAlbums(c *gin.Context) {
	// 第1引数はクライアントに送るステータスコード（ここでは200）
	// ちなみに IndentedJSON の代わりに JSON でもほぼ違いはないみたいだ
	// インデントされていてデバッグしやすいとかその程度の違い
	c.IndentedJSON(http.StatusOK, albums)
}

// gin.Context についての公式チュートリアルからのメモ：
// gin.Context is the most important part of Gin. It carries request details, validates and serializes JSON, and more. (Despite the similar name, this is different from Go’s built-in context package.)

// POSTでアクセスした場合の処理。postされたデータを albums に追加する
func postAlbums(c *gin.Context) {
	var newAlbum album

	// リクエストボディのJSONと newAlbum の構造をバインドする（マッピングさせる）
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// 既存の albums スライスに newAlbum を追加して、成功した旨（201）を返す
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// たとえばこんなデータをpostしてみる
// {"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}

// GETでアクセスした場合の処理。idがマッチする album を返す
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// albums スライスから、idがマッチするレコードを探してJSONとして返す
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	// マッチするレコードがなかった場合のレスポンス
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// ルートにアクセスしたときのメッセージ表示
func welcome(c *gin.Context) {
	fmt.Fprint(c.Writer, "welcome!")
}
