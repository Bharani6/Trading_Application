package main

import (
	"fmt"
	"stock-trading/internal/config"
	"stock-trading/internal/database"
	"stock-trading/internal/user"
)

func main() {
	config.LoadConfig("configs")
	database.ConnectDB()
	var usr user.User
	if err := database.DB.Where("email = ?", "bharanidharan@gmail.com").First(&usr).Error; err != nil {
		fmt.Println("User not found")
		return
	}
	usr.Status = "active"
	usr.Role = "admin"
	if err := database.DB.Save(&usr).Error; err != nil {
		fmt.Println("Failed to update user:", err)
		return
	}
	fmt.Println("Admin successfully activated!")
}
