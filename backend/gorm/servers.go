package gorm

import (
	"backend/models"
	"encoding/json"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect_To_Database_Sixth() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("/root/server.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func Servers_create_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Sixth()
	// 迁移模型
	db.AutoMigrate(&models.Servers{})
	servers := models.Servers{
		Name:             jsonData["name"].(string),
		UserID:           jsonData["userID"].(string),
		Ssh_user:         "",
		Ip:               "",
		Port:             "",
		Auth_method:      "0",
		Password:         "",
		Key:              "",
		JumpServerName:   "",
		JumpServerUserID: "",
		Login_command:    "",
	}
	// 在数据库中创建记录
	result := db.Model(&models.Servers{}).Create(&servers)
	if result.Error != nil {
		panic(result.Error)
	}
}

func Servers_update_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Sixth()

	// 迁移模型
	db.AutoMigrate(&models.Servers{})

	type UpdateCondition struct {
		Name   string `"json:"name"`
		UserID string `json"userID"`
	}
	updateCondition := UpdateCondition{
		UserID: jsonData["userID"].(string),
		Name:   jsonData["name"].(string),
	}

	// 更新记录
	result_name := db.Model(&models.Servers{}).Where(&updateCondition).Update("name", jsonData["new_name"].(string))
	result_ssh_user := db.Model(&models.Servers{}).Where(&updateCondition).Update("ssh_user", jsonData["ssh_user"].(string))
	result_ip := db.Model(&models.Servers{}).Where(&updateCondition).Update("ip", jsonData["ip"].(string))
	result_port := db.Model(&models.Servers{}).Where(&updateCondition).Update("port", jsonData["port"].(string))
	result_auth_method := db.Model(&models.Servers{}).Where(&updateCondition).Update("auth_method", jsonData["auth_method"].(string))
	result_password := db.Model(&models.Servers{}).Where(&updateCondition).Update("password", jsonData["password"].(string))
	result_key := db.Model(&models.Servers{}).Where(&updateCondition).Update("key", jsonData["key"].(string))
	result_jumpServerName := db.Model(&models.Servers{}).Where(&updateCondition).Update("jump_server_name", jsonData["jumpServerName"].(string))
	result_jumpServerUserID := db.Model(&models.Servers{}).Where(&updateCondition).Update("jump_server_user_id", jsonData["jumpServerUserID"].(string))
	result_login_command := db.Model(&models.Servers{}).Where(&updateCondition).Update("login_command", jsonData["login_command"].(string))

	if result_name.Error != nil {
		// 处理错误
		panic(result_name.Error)
	} else if result_ssh_user.Error != nil {
		panic(result_ssh_user.Error)
	} else if result_ip.Error != nil {
		panic(result_ip.Error)
	} else if result_port.Error != nil {
		panic(result_port.Error)
	} else if result_auth_method.Error != nil {
		panic(result_auth_method.Error)
	} else if result_password.Error != nil {
		panic(result_password.Error)
	} else if result_key.Error != nil {
		panic(result_key.Error)
	} else if result_jumpServerName.Error != nil {
		panic(result_jumpServerName.Error)
	} else if result_jumpServerUserID.Error != nil {
		panic(result_jumpServerUserID.Error)
	} else if result_login_command.Error != nil {
		panic(result_login_command.Error)
	}
}

func Servers_query_all(jsonData map[string]interface{}) []map[string]interface{} {

	db := Connect_To_Database_Sixth()
	var servers []models.Servers
	db.Where("user_id = ? ", jsonData["userID"].(string)).Find(&servers)

	var jsonArray []string
	for _, doc := range servers {
		jsonDateConversion, err := json.Marshal(doc)
		if err != nil {
			panic(err)
		}
		jsonArray = append(jsonArray, string(jsonDateConversion))
	}
	// 创建一个空的 []map[string]interface{} 类型的切片，用于存储转换后的数据
	var jsonMaps []map[string]interface{}
	// 遍历每个 JSON 字符串
	for _, jsonString := range jsonArray {
		// 创建一个 map 用于存储解析后的 JSON 数据
		var jsonMap map[string]interface{}

		// 解析 JSON 字符串到 map[string]interface{} 类型
		errs := json.Unmarshal([]byte(jsonString), &jsonMap)
		if errs != nil {
			panic(errs)
		}

		// 将解析后的 map 添加到 jsonMaps 切片中
		jsonFilter := map[string]interface{}{
			"name":             jsonMap["Name"],
			"userID":           jsonMap["UserID"].(string),
			"ssh_user":         jsonMap["ssh_user"].(string),
			"ip":               jsonMap["ip"].(string),
			"port":             jsonMap["port"].(string),
			"auth_method":      jsonMap["auth_method"].(string),
			"password":         jsonMap["password"].(string),
			"key":              jsonMap["key"].(string),
			"jumpServerName":   jsonMap["jumpServerName"].(string),
			"jumpServerUserID": jsonMap["jumpServerUserID"].(string),
			"login_command":    jsonMap["login_command"].(string),
		}
		jsonMaps = append(jsonMaps, jsonFilter)
	}
	return jsonMaps
}

func Servers_query_sql(jsonData map[string]interface{}) map[string]interface{} {
	db := Connect_To_Database_Third()
	// 迁移模型
	db.AutoMigrate(&models.Servers{})

	type QueryCondition struct {
		Name   string `json:"name"`
		UserID string `json:"userID"`
	}
	queryCondition := QueryCondition{
		Name:   jsonData["name"].(string),
		UserID: jsonData["userID"].(string),
	}

	// 查询记录
	var result models.Servers
	if err := db.Where(&queryCondition).First(&result).Error; err != nil {
		panic(err)
	}

	// 将结果转换为 JSON 格式
	request, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	// 将 []byte 转换为 map[string]interface{}
	var request_jsonData map[string]interface{}
	err = json.Unmarshal(request, &request_jsonData)
	if err != nil {
		panic(err)
	}

	return request_jsonData
}

func Servers_delete_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Sixth()
	// 迁移模型
	db.AutoMigrate(&models.Servers{})
	deleteCondition := models.Servers{
		Name:   jsonData["name"].(string),
		UserID: jsonData["userID"].(string),
	}
	// 在数据库中创建记录
	result := db.Where(&deleteCondition).Delete(&models.Servers{})
	if result.Error != nil {

		panic(result.Error)
	}
}
