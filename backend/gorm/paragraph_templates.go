package gorm

import (
	"backend/models"
	"encoding/json"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect_To_Database_Fifth() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("/root/server.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func Paragraph_templates_create_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Fifth()
	// 迁移模型
	db.AutoMigrate(&models.Paragraph_templates{})
	Paragraph_templates := models.Paragraph_templates{
		Name:   jsonData["name"].(string),
		UserID: jsonData["userID"].(string),
	}
	// 在数据库中创建记录
	result := db.Model(&models.Paragraph_templates{}).Create(&Paragraph_templates)
	if result.Error != nil {
		panic(result.Error)
	}
}

func Paragraph_templates_update_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Fifth()

	// 迁移模型
	db.AutoMigrate(&models.Paragraph_templates{})

	type UpdateCondition struct {
		Name   string `"json:"name"`
		UserID string `json"userID"`
	}
	updateCondition := UpdateCondition{
		UserID: jsonData["userID"].(string),
		Name:   jsonData["name"].(string),
	}

	// 更新记录
	result_text := db.Model(&models.Paragraph_templates{}).Where(&updateCondition).Update("text", jsonData["text"].(string))
	result_tag := db.Model(&models.Paragraph_templates{}).Where(&updateCondition).Update("tag", jsonData["tag"].(string))
	result_name := db.Model(&models.Paragraph_templates{}).Where(&updateCondition).Update("name", jsonData["new_name"].(string))

	if result_name.Error != nil {
		panic(result_name.Error)
	}
	if result_text.Error != nil {
		panic(result_text.Error)
	}
	if result_tag.Error != nil {
		panic(result_tag.Error)
	}
}

func Paragraph_templates_query_all_sql(jsonData map[string]interface{}) []map[string]interface{} {
	db := Connect_To_Database_Fifth()
	var paragraph_templates []models.Paragraph_templates
	db.Where("user_id = ?", jsonData["userID"].(string)).Find(&paragraph_templates)
	var jsonArray []string
	for _, doc := range paragraph_templates {
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
			"name": jsonMap["Name"],
		}
		jsonMaps = append(jsonMaps, jsonFilter)
	}
	return jsonMaps
}

func Paragraph_templates_query_sql(jsonData map[string]interface{}) map[string]interface{} {
	db := Connect_To_Database_Fifth()
	// 迁移模型
	db.AutoMigrate(&models.Paragraph_templates{})

	type QueryCondition struct {
		Name   string `"json:"name"`
		UserID string `json"userID"`
	}
	queryCondition := QueryCondition{
		Name:   jsonData["name"].(string),
		UserID: jsonData["userID"].(string),
	}

	// 查询记录
	var result models.Paragraph_templates
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

func Paragraph_templates_delete_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Fifth()
	// 迁移模型
	db.AutoMigrate(&models.Paragraph_templates{})
	deleteCondition := models.Paragraph_templates{
		Name:   jsonData["name"].(string),
		UserID: jsonData["userID"].(string),
	}
	// 在数据库中创建记录
	result := db.Where(&deleteCondition).Delete(&models.Paragraph_templates{})
	if result.Error != nil {
		panic(result.Error)
	}
}
