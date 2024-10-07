package gorm

import (
	"backend/models"
	"encoding/json"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func File_id_all(jsonData map[string]interface{}) []map[string]interface{} {

	db, err := gorm.Open(sqlite.Open("/root/server.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var file []models.File
	db.Where("document_id = ? AND user_id = ?", jsonData["documentID"].(int), jsonData["userID"].(string)).Find(&file)
	var jsonArray []string
	for _, doc := range file {
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
		jsonMaps = append(jsonMaps, jsonMap)
	}
	return jsonMaps
}

func Picture_upload_sql(jsonData map[string]interface{}) {
	dsn := "root:root@tcp(yidaproserver.ltd:4501)/server"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 迁移模型
	db.AutoMigrate(&models.File{})
	file := models.File{
		UserID:     jsonData["userID"].(string),
		DocumentID: jsonData["documentID"].(int),
		Path:       jsonData["path"].(string),
		Value:      jsonData["value"].(string),
	}
	// 在数据库中创建记录
	result := db.Model(&models.File{}).Create(&file)
	if result.Error != nil {
		panic(result.Error)
	}
}

func Picture_upload_query_sql(jsonData map[string]interface{}) models.File {
	dsn := "root:root@tcp(yidaproserver.ltd:4501)/server"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 迁移模型
	db.AutoMigrate(&models.File{})
	file := models.File{
		UserID:     jsonData["userID"].(string),
		DocumentID: jsonData["documentID"].(int),
		Path:       jsonData["path"].(string),
		Value:      jsonData["value"].(string),
	}
	var picture models.File
	db.Where(file).Find(&picture)
	return picture
}
