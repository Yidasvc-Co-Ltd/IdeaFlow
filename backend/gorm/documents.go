package gorm

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//连接数据库
func Connect_To_Database() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("/root/server.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

//同一个userID的documentID自增
func DocumentID_self_increasing(db *gorm.DB, userID string) int {
	db.AutoMigrate(&models.Documents{})
	// 执行查询，按照创建documentID降序排序，获取最新一条符合条件的记录的文档ID
	var document models.Documents
	if err := db.Where("user_id = ?", userID).Order("document_id DESC").First(&document).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//首次创建
			return 1
		} else {
			panic(err)
		}
	}
	return document.DocumentID + 1
}

func Document_create_sql(jsonData map[string]interface{}) int {
	db, err := Connect_To_Database()
	if err != nil {
		panic(err)
	}
	userID := jsonData["userID"].(string)
	path := jsonData["path"].(string)
	newDocumentID := DocumentID_self_increasing(db, userID)
	doc := models.Documents{
		DocumentID:  newDocumentID,
		UserID:      userID,
		UpdateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Path:        path,
		Paragraphs:  "",
		IsCollected: 0,
	}
	db.AutoMigrate(&models.Documents{})
	result := db.Create(&doc)
	if result.Error != nil {
		panic(result.Error)
	}
	return newDocumentID
}

func Document_delete_sql(jsonData map[string]interface{}) {
	db, err := Connect_To_Database()
	if err != nil {

		panic(err)
	}

	db.AutoMigrate(&models.Documents{})

	deleteCondition := models.Documents{DocumentID: jsonData["documentID"].(int), UserID: jsonData["userID"].(string)}

	result := db.Where(&deleteCondition).Delete(&models.Documents{})
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func Document_query_sql(jsonData map[string]interface{}) map[string]interface{} {
	db, err := Connect_To_Database()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Documents{})

	type QueryCondition struct {
		DocumentID int    `json:"documentID"`
		UserID     string `json:"userID"`
	}
	queryCondition := QueryCondition{
		DocumentID: jsonData["documentID"].(int),
		UserID:     jsonData["userID"].(string),
	}

	var result models.Documents
	if err := db.Where(&queryCondition).First(&result).Error; err != nil {

		panic(err)
	}

	// 将结果转换为 JSON 格式
	request, err := json.Marshal(result)
	if err != nil {
		// 处理 JSON 转换错误
		panic(err)
	}

	// 将 []byte 转换为 map[string]interface{}
	var request_jsonData map[string]interface{}
	err = json.Unmarshal(request, &request_jsonData)
	if err != nil {
		// 处理 JSON 解码错误
		panic(err)
	}

	return request_jsonData
}

//更新documents表的paragraphs
func Documents_paragraphs_update_sql(db *gorm.DB, jsonData map[string]interface{}) string {
	db.AutoMigrate(&models.Paragraphs{})
	var paragraph []models.Paragraphs
	errs := db.Where("user_id = ? AND document_id=?", jsonData["userID"].(string), jsonData["documentID"].(int)).Find(&paragraph).Error
	if errs == gorm.ErrRecordNotFound {
		return "1"
	}
	var jsonArray []string
	for _, doc := range paragraph {
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
		fmt.Println(jsonArray)
		jsonFilter := map[string]interface{}{
			"paragraphID": jsonMap["paragraphID"].(string),
			"version":     jsonMap["version"].(string),
		}
		jsonMaps = append(jsonMaps, jsonFilter)
	}
	jsonMapss, err := json.Marshal(jsonMaps)
	if err != nil {
		panic(err)
	}

	jsonMaps_paragraphs := string(jsonMapss)
	return jsonMaps_paragraphs
}

func Document_update_sql(jsonData map[string]interface{}) {
	db, err := Connect_To_Database()
	if err != nil {

		panic(err)
	}

	db.AutoMigrate(&models.Documents{})

	type UpdateCondition struct {
		DocumentID int    `json:"documentID"`
		UserID     string `json:"userID"`
	}
	updateCondition := UpdateCondition{
		DocumentID: jsonData["documentID"].(int),
		UserID:     jsonData["userID"].(string),
	}
	jsonMaps_paragraphs := Documents_paragraphs_update_sql(db, jsonData)
	// 更新记录
	result := db.Model(&models.Documents{}).Where(&updateCondition).Update("path", jsonData["path"].(string))
	if result.Error != nil {
		panic(err)
	}
	result_paragraphs := db.Model(&models.Documents{}).Where(&updateCondition).Update("paragraphs", jsonMaps_paragraphs)
	if result_paragraphs.Error != nil {
		panic(result_paragraphs.Error)
	}
}

func Document_update_time_sql(jsonData map[string]interface{}) {
	db, err := Connect_To_Database()
	if err != nil {
		panic(err)
	}

	// 迁移模型
	db.AutoMigrate(&models.Documents{})

	type UpdateCondition struct {
		DocumentID int    `json:"documentID"`
		UserID     string `json:"userID"`
	}
	updateCondition := UpdateCondition{
		DocumentID: int(jsonData["documentID"].(float64)),
		UserID:     jsonData["userID"].(string),
	}

	// 更新记录
	result := db.Model(&models.Documents{}).Where(&updateCondition).Update("update_time", time.Now().Format("2006-01-02 15:04:05"))
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func Document_update_is_collected_sql(jsonData map[string]interface{}) {
	db, err := Connect_To_Database()
	if err != nil {
		panic(err)
	}

	// 迁移模型
	db.AutoMigrate(&models.Documents{})

	type UpdateCondition struct {
		DocumentID int    `json:"documentID"`
		UserID     string `json:"userID"`
	}
	updateCondition := UpdateCondition{
		DocumentID: jsonData["documentID"].(int),
		UserID:     jsonData["userID"].(string),
	}

	// 更新记录
	result := db.Model(&models.Documents{}).Where(&updateCondition).Update("is_collected", jsonData["isCollected"].(float64))
	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func Document_query_all_sql(jsonData map[string]interface{}) []map[string]interface{} {
	db, err := Connect_To_Database()
	if err != nil {
		panic(err)
	}
	var documents []models.Documents
	db.Where("user_id = ?", jsonData["userID"].(string)).Find(&documents)
	var jsonArray []string
	for _, doc := range documents {
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
