package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Character struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Stats string
}

var collection *mgo.Collection
var session *mgo.Session

func init() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	collection = session.DB("master-tools-db").C("characters")
}

func main() {
	defer session.Close()
	router := gin.Default()

	v1 := router.Group("/api/characters")
	{
		v1.POST("/", createCharacter)
		v1.GET("/", fetchAllCharacters)
		v1.GET("/:id", fetchSingleCharacter)
		v1.PUT("/:id", updateCharacter)
		v1.DELETE("/:id", deleteCharacter)
	}

	router.Run()
}

func createCharacter(context *gin.Context) {
	stats := context.PostForm("stats")
	fmt.Printf("stats: %v\n", stats)
	character := Character{bson.NewObjectId(), stats}
	fmt.Printf("character: %v\n", character)
	err := collection.Insert(&character)
	if err != nil {
		log.Fatal(err)
	}

	context.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "character created successfully",
	})
}

func fetchAllCharacters(context *gin.Context) {
	var characters []Character
	err := collection.Find(nil).All(&characters)
	if err != nil {
		log.Fatal(err)
	}

	if len(characters) <= 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no character found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   characters,
	})
}

func fetchSingleCharacter(context *gin.Context) {
	character := Character{}
	id := bson.ObjectIdHex(context.Param("id"))
	err := collection.FindId(id).One(&character)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		context.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "character not found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   character,
	})
}

func updateCharacter(context *gin.Context) {
	id := bson.ObjectIdHex(context.Param("id"))
	stats := context.PostForm("stats")

	err := collection.UpdateId(id, bson.M{"Stats": stats})

	if err != nil {
		fmt.Println("Error: " + err.Error())
		context.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "character not found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "character updated successfully!",
	})
}

func deleteCharacter(context *gin.Context) {
	id := bson.ObjectIdHex(context.Param("id"))
	err := collection.RemoveId(id)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		context.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "character not found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "character deleted successfully!",
	})
}
