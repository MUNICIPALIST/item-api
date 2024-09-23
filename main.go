package main

import (
	_ "github.com/MUNICIPALIST/item-api/docs" // Замените на ваш модуль
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
}

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=4321 dbname=itemAPI port=5432 sslmode=disable TimeZone=Asia/Almaty"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Faild to connect to database!")
	}

	if err = DB.AutoMigrate(&Item{}); err != nil {
		log.Fatalf("Не удалось выполнить миграцию базы данных: %v", err)
	}

	DB.AutoMigrate(&Item{})

}

// @title           Item API
// @version         1.0
// @description     Это API для управления предметами.
// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

func main() {
	InitDB()
	r := gin.Default()

	r.POST("/items", CreateItem)
	r.GET("/items", GetItems)
	r.GET("/items/:id", GetItem)
	r.PUT("/items/:id", UpdateItem)
	r.DELETE("/items/:id", DeleteItem)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}

// CreateItem godoc
// @Summary      Создать новый Item
// @Description  Создает новый Item
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        item  body      Item  true  "Item для создания"
// @Success      200   {object}  Item
// @Failure      400   {object}  ErrorResponse
// @Router       /items [post]
func CreateItem(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, ErrorResponse{Error: "Item not found"})
		return
	}
	DB.Create(&item)
	c.JSON(200, item)
}

// GetItems godoc
// @Summary      Получить список всех Items
// @Description  Возвращает список всех Items
// @Tags         items
// @Produce      json
// @Success      200  {array}   Item
// @Router       /items [get]
func GetItems(c *gin.Context) {
	var items []Item
	DB.Find(&items)
	c.JSON(200, items)
}

// GetItem godoc
// @Summary      Получить Item по ID
// @Description  Возвращает Item по заданному ID
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "ID Item"
// @Success      200  {object}  Item
// @Failure      404  {object}  ErrorResponse
// @Router       /items/{id} [get]
func GetItem(c *gin.Context) {
	id := c.Param("id")
	var item Item
	if err := DB.First(&item, id).Error; err != nil {
		c.JSON(404, ErrorResponse{Error: "Item not found"})
		return
	}
	c.JSON(200, item)
}

// UpdateItem godoc
// @Summary      Обновить существующий Item
// @Description  Обновляет данные существующего Item по ID
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        id    path      int   true  "ID Item"
// @Param        item  body      Item  true  "Обновленные данные Item"
// @Success      200   {object}  SuccessResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Router       /items/{id} [put]
func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item Item
	if err := DB.First(&item, id).Error; err != nil {
		c.JSON(404, ErrorResponse{Error: "Item not found"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, ErrorResponse{Error: "Item not found"})
		return
	}
	DB.Save(&item)
	c.JSON(200, item)
}

// DeleteItem godoc
// @Summary      Удалить Item
// @Description  Удаляет существующий Item по ID
// @Tags         items
// @Param        id   path      int  true  "ID Item"
// @Success      200  {object}  SuccessResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /items/{id} [delete]

func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	var item Item
	if err := DB.First(&item, id).Error; err != nil {
		c.JSON(404, ErrorResponse{Error: "Item not found"})
		return
	}
	DB.Delete(&item)
	c.JSON(200, SuccessResponse{Message: "Item deleted"})
}
