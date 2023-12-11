// @Time    : 8/1/2023 6:05 PM
// @Author  : Cillian
// @Email   : cilliandevops@gmail.com
// Website  : www.cillian.website
// Have a good day!

package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func main() {
	// 连接到数据库
	dsn := "user:password@tcp(hostname:port)/dbname"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}
	defer db.Close()

	// 自动迁移（创建表结构）
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("数据迁移失败：", err)
	}

	// 创建Gin实例
	r := gin.Default()

	// 添加路由
	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db.Create(&user)
		c.JSON(200, user)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		err := db.First(&user, id).Error
		if err != nil {
			c.JSON(404, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(200, user)
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		err := db.First(&user, id).Error
		if err != nil {
			c.JSON(404, gin.H{"error": "用户不存在"})
			return
		}

		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		db.Model(&user).Updates(newUser)
		c.JSON(200, user)
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		err := db.First(&user, id).Error
		if err != nil {
			c.JSON(404, gin.H{"error": "用户不存在"})
			return
		}

		db.Delete(&user)
		c.JSON(200, gin.H{"message": "用户已删除"})
	})

	// 启动Gin服务器
	r.Run(":8080")
}
