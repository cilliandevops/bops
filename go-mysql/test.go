// @Time    : 7/28/2023 11:22 AM
// @Author  : Cillian
// @Email   : cilliandevops@gmail.com
// Website  : www.cillian.website
// Have a good day!

package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ServerInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	IP      string `json:"ip"`
	OS      string `json:"os"`
	Owner   string `json:"owner"`
	Comment string `json:"comment"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "YOUR_MYSQL_USERNAME:YOUR_MYSQL_PASSWORD@tcp(YOUR_MYSQL_HOST:YOUR_MYSQL_PORT)/YOUR_MYSQL_DATABASE")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	router := gin.Default()
	router.GET("/serverinfo", getServerInfoList)
	router.POST("/serverinfo", createServerInfo)
	router.PUT("/serverinfo/:id", updateServerInfo)
	router.DELETE("/serverinfo/:id", deleteServerInfo)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Error starting the server:", err)
	}
}

func getServerInfoList(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM server_info")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch server info"})
		return
	}
	defer rows.Close()

	serverInfoList := []ServerInfo{}
	for rows.Next() {
		var serverInfo ServerInfo
		if err := rows.Scan(&serverInfo.ID, &serverInfo.Name, &serverInfo.IP, &serverInfo.OS, &serverInfo.Owner, &serverInfo.Comment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch server info"})
			return
		}
		serverInfoList = append(serverInfoList, serverInfo)
	}

	c.JSON(http.StatusOK, serverInfoList)
}

func createServerInfo(c *gin.Context) {
	var newServerInfo ServerInfo
	if err := c.ShouldBindJSON(&newServerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	result, err := db.Exec("INSERT INTO server_info (name, ip, os, owner, comment) VALUES (?, ?, ?, ?, ?)", newServerInfo.Name, newServerInfo.IP, newServerInfo.OS, newServerInfo.Owner, newServerInfo.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create server info"})
		return
	}

	newID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get new server info ID"})
		return
	}

	newServerInfo.ID = int(newID)
	c.JSON(http.StatusOK, newServerInfo)
}

func updateServerInfo(c *gin.Context) {
	serverInfoID := c.Param("id")

	var updatedServerInfo ServerInfo
	if err := c.ShouldBindJSON(&updatedServerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	_, err := db.Exec("UPDATE server_info SET name = ?, ip = ?, os = ?, owner = ?, comment = ? WHERE id = ?", updatedServerInfo.Name, updatedServerInfo.IP, updatedServerInfo.OS, updatedServerInfo.Owner, updatedServerInfo.Comment, serverInfoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update server info"})
		return
	}

	c.Status(http.StatusOK)
}

func deleteServerInfo(c *gin.Context) {
	serverInfoID := c.Param("id")

	_, err := db.Exec("DELETE FROM server_info WHERE id = ?", serverInfoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete server info"})
		return
	}

	c.Status(http.StatusOK)
}
