package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

// Model
type Task struct {
	gorm.Model
	ClassName string `json:"className"`
	Title     string `json:"title"`
	Questions string `json:"questions"`
	Task      string `json:"task"`
	Image     string `json:"image"`
}

// Global DB
var DB *gorm.DB

func initDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("Bazaga ulanishda xatolik!")
	}
	DB.AutoMigrate(&Task{})
}

func main() {
	initDB()
	r := gin.Default()

	// CORS ruxsat
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/tasks", getTasks)
		api.POST("/tasks", createTask)
		api.PUT("/tasks/:id", updateTask)
		api.DELETE("/tasks/:id", deleteTask)
	}

	r.Run(":8080")
}

func getTasks(c *gin.Context) {
	var tasks []Task
	DB.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var newTask Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := DB.Create(&newTask)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, newTask)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if err := DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task topilmadi"})
		return
	}

	var updated Task
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ClassName = updated.ClassName
	task.Title = updated.Title
	task.Questions = updated.Questions
	task.Task = updated.Task
	task.Image = updated.Image

	DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if err := DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task topilmadi"})
		return
	}

	DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "Task o‘chirildi"})
}