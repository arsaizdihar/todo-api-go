package main

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
	"todo-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	ID int `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Done bool `json:"done"`
}

type TodoInput struct {
	Title string `json:"title" binding:"required"`
}

func getTodos(c *gin.Context) {
	var todos []Todo

	db.Find(&todos)
	c.JSON(http.StatusOK, todos)
}

func getTodo(c *gin.Context) {

	var todo Todo

	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "id not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func toggleTodo(c *gin.Context) {

	var todo Todo

	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "id not found"})
		return
	}

	todo.Done = !todo.Done
	db.Save(&todo)

	c.JSON(http.StatusOK, todo)
}

func updateTodo(c *gin.Context) {

	var todo Todo

	input := c.MustGet(gin.BindKey).(*TodoInput)

	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "id not found"})
		return
	}

	db.Model(&todo).Updates(Todo{Title: input.Title})

	c.JSON(http.StatusOK, todo)

}

func deleteTodo(c *gin.Context) {

	var todo Todo

	if err := db.First(&todo, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "id not found"})
		return
	}
	db.Delete(&todo)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}


func createTodo(c *gin.Context) {
	input := c.MustGet(gin.BindKey).(*TodoInput)
	todo := Todo{Title: input.Title}
	db.Create(&todo)

	c.JSON(http.StatusCreated, todo)
}

var db *gorm.DB
var err error

func main() {
	godotenv.Load()

	dsn := os.Getenv("DSN")
	

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Todo{})
	

	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	r.Use(middleware.Errors())

	r.GET("/api/todos", getTodos)
	r.POST("/api/todos", gin.Bind(TodoInput{}), createTodo)
	r.GET("/api/todos/:id", getTodo)
	r.POST("/api/todos/toggle/:id", toggleTodo)
	r.PUT("/api/todos/:id", updateTodo)
	r.DELETE("/api/todos/:id", deleteTodo)

	r.Run("localhost:" + os.Getenv("PORT"))
}