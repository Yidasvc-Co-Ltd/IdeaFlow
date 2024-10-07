package gorm

import (
	"backend/models"
	"encoding/json"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect_To_Database_Second() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("/root/server.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func ParagraphID_self_increasing(db *gorm.DB, userID string, documentID int) int {
	db.AutoMigrate(&models.Paragraphs{})
	// 执行查询，按照创建paragarphsID降序排序，获取最新一条符合条件的记录的文档ID
	var paragraphs models.Paragraphs
	if err := db.Where("user_id = ? AND document_id=?", userID, documentID).Order("paragraph_id DESC").First(&paragraphs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//首次创建
			return 1
		} else {
			panic(err)
		}
	}
	return paragraphs.ParagraphID + 1
}

func Version_self_increasing(db *gorm.DB, userID string, documentID int, paragraphID int) int {
	db.AutoMigrate(&models.Paragraphs{})
	// 执行查询，按照创建version降序排序，获取最新一条符合条件的记录的文档ID
	var paragraphs models.Paragraphs
	if err := db.Where("user_id = ? AND document_id=? AND paragraph_id=?", userID, documentID, paragraphID).Order("paragraph_id DESC").First(&paragraphs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//首次创建
			return 1
		} else {
			panic(err)
		}
	}
	return paragraphs.Version + 1
}

func Docuemnt_paragraph_insert_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Second()
	// 迁移模型
	db.AutoMigrate(&models.Paragraphs{})
	paragarphs := ParagraphID_self_increasing(db, jsonData["userID"].(string), jsonData["documentID"].(int))
	paraGraphs := models.Paragraphs{
		ParagraphID:      paragarphs,
		Version:          1,
		DocumentID:       jsonData["documentID"].(int),
		UserID:           jsonData["userID"].(string),
		IsFirstlnChapter: "0",
		Text:             "0",
		Before:           "",
		Next:             "",
	}
	// 在数据库中创建记录
	result := db.Model(&models.Paragraphs{}).Create(&paraGraphs)
	if result.Error != nil {
		// 处理创建记录时的错误
		panic(result.Error)
	}
}
func Document_paragraph_modify_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Second()
	// 迁移模型
	db.AutoMigrate(&models.Paragraphs{})
	version := Version_self_increasing(db, jsonData["userID"].(string), jsonData["documentID"].(int), jsonData["paragraphID"].(int))
	before := []map[string]interface{}{{"paragraphID": jsonData["paragarphID"].(int), "version": version}}
	beforejson, err := json.Marshal(before)
	if err != nil {
		panic(err)
	}
	beforeString := string(beforejson)
	paragraphs := models.Paragraphs{
		ParagraphID:      jsonData["paragarphID"].(int),
		Version:          version,
		DocumentID:       jsonData["documentID"].(int),
		UserID:           jsonData["userID"].(string),
		IsFirstlnChapter: "0",
		Text:             jsonData["text"].(string),
		Before:           beforeString,
		Next:             "",
	}
	// 在数据库中创建记录
	result := db.Model(&models.Paragraphs{}).Create(&paragraphs)
	if result.Error != nil {
		// 处理创建记录时的错误
		panic(result.Error)
	}
}
func Document_paragraph_query_text_sql(jsonData map[string]interface{}) string {
	db := Connect_To_Database_Second()
	// 迁移模型
	db.AutoMigrate(&models.Paragraphs{})
	text := ""
	type QueryCondition struct {
		ParagraphID int    `json:"paragraphID"`
		Version     int    `json:"version"`
		DocumentID  int    `json"documentID"`
		UserID      string `json"userID"`
	}
	for _, jsonDataText := range jsonData["Data"].([]map[string]interface{}) {

		queryCondition := QueryCondition{
			ParagraphID: jsonDataText["paragraphID"].(int),
			Version:     jsonDataText["version"].(int),
			DocumentID:  jsonData["documentID"].(int),
			UserID:      jsonData["userID"].(string),
		}

		// 查询记录
		var result models.Paragraphs
		if err := db.Where(&queryCondition).First(&result).Error; err != nil {
			// 处理错误
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

		ArrayText := request_jsonData["text"].(string)
		text = ArrayText + result.Text
	}
	return text
}

func Document_paragraph_merge_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Second()
	//迁移模型
	db.AutoMigrate(&models.Paragraphs{})
	paragraphID := ParagraphID_self_increasing(db, jsonData["userID"].(string), jsonData["documentID"].(int))
	version := Version_self_increasing(db, jsonData["userID"].(string), jsonData["documentID"].(int), paragraphID)
	before := jsonData["data"].([]map[string]interface{})
	beforejson, err := json.Marshal(before)
	if err != nil {
		panic(err)
	}
	beforeString := string(beforejson)
	text := Document_paragraph_query_text_sql(jsonData)
	paragraphs := models.Paragraphs{
		ParagraphID:      paragraphID,
		Version:          version,
		DocumentID:       jsonData["documentID"].(int),
		UserID:           jsonData["userID"].(string),
		IsFirstlnChapter: "0",
		Text:             text,
		Before:           beforeString,
		Next:             "",
	}
	// 在数据库中创建记录
	result := db.Model(&models.Paragraphs{}).Create(&paragraphs)
	if result.Error != nil {
		// 处理创建记录时的错误
		panic(result.Error)
	}
}

func Document_paragraph_query_sql(jsonData map[string]interface{}) map[string]interface{} {
	db := Connect_To_Database_Second()
	// 迁移模型
	db.AutoMigrate(&models.Paragraphs{})
	type QueryCondition struct {
		ParagraphID int    `json:"paragraphID"`
		Version     int    `json:"version"`
		DocumentID  int    `json"documentID"`
		UserID      string `json"userID"`
	}
	queryCondition := QueryCondition{
		ParagraphID: jsonData["paragraphID"].(int),
		Version:     jsonData["version"].(int),
		DocumentID:  jsonData["documentID"].(int),
		UserID:      jsonData["userID"].(string),
	}

	// 查询记录
	var result models.Paragraphs
	if err := db.Where(&queryCondition).First(&result).Error; err != nil {
		// 处理错误
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

func Docuemnt_paragraph_version_query_sql(jsonData map[string]interface{}) []map[string]interface{} {
	db := Connect_To_Database_Second()

	var paragraphs []models.Paragraphs
	db.Where("user_id = ? AND document_id AND paragraph_id", jsonData["userID"].(string), jsonData["documentID"].(int), jsonData["paragraphID"].(int)).Find(&paragraphs)
	var jsonArray []string
	for _, doc := range paragraphs {
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
			"version": jsonMap["version"].(int),
			"before":  jsonMap["before"].(string),
			"next":    jsonMap["next"].(string),
		}
		jsonMaps = append(jsonMaps, jsonFilter)
	}
	return jsonMaps
}

func Document_paragraph_version_jump_sql(jsonData map[string]interface{}) map[string]interface{} {
	db := Connect_To_Database_Second()
	// 迁移模型
	db.AutoMigrate(&models.Paragraphs{})
	type QueryCondition struct {
		ParagraphID int    `json:"paragraphID"`
		Version     int    `json:"version"`
		DocumentID  int    `json"documentID"`
		UserID      string `json"userID"`
	}
	queryCondition := QueryCondition{
		ParagraphID: jsonData["paragraphID"].(int),
		Version:     jsonData["version"].(int),
		DocumentID:  jsonData["documentID"].(int),
		UserID:      jsonData["userID"].(string),
	}

	// 查询记录
	var result models.Paragraphs
	if err := db.Where(&queryCondition).First(&result).Error; err != nil {
		// 处理错误
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

func Document_paragraph_separate_sql(jsonData map[string]interface{}) {
	db := Connect_To_Database_Second()
	// 迁移模型
	db.AutoMigrate(&models.Paragraphs{})
	paragraphID_first := ParagraphID_self_increasing(db, jsonData["userID"].(string), jsonData["documentID"].(int))
	paragraphID_second := ParagraphID_self_increasing(db, jsonData["userID"].(string), jsonData["documentID"].(int))
	before := []map[string]interface{}{{"paragraphID": jsonData["paragarphID"].(int), "version": jsonData["version"].(int)}}
	beforejson, err := json.Marshal(before)
	if err != nil {
		panic(err)
	}
	beforeString := string(beforejson)
	request := Document_paragraph_query_sql(jsonData)
	request_text := request["text"].(string)
	text_first := request_text[:jsonData["first_length"].(int)]
	text_second := request_text[:len(request_text)-jsonData["first_length"].(int)]
	paraGraphs_first := models.Paragraphs{
		ParagraphID:      paragraphID_first,
		Version:          1,
		DocumentID:       jsonData["documentID"].(int),
		UserID:           jsonData["userID"].(string),
		IsFirstlnChapter: "0",
		Text:             text_first,
		Before:           "",
		Next:             beforeString,
	}
	// 在数据库中创建记录
	result_first := db.Model(&models.Paragraphs{}).Create(&paraGraphs_first)
	if result_first.Error != nil {
		// 处理创建记录时的错误
		panic(result_first.Error)
	}

	paraGraphs_second := models.Paragraphs{
		ParagraphID:      paragraphID_second,
		Version:          1,
		DocumentID:       jsonData["documentID"].(int),
		UserID:           jsonData["userID"].(string),
		IsFirstlnChapter: "0",
		Text:             text_second,
		Before:           "",
		Next:             beforeString,
	}
	// 在数据库中创建记录
	result_second := db.Model(&models.Paragraphs{}).Create(&paraGraphs_second)
	if result_second.Error != nil {
		// 处理创建记录时的错误
		panic(result_second.Error)
	}
}

func Document_paragraph_query_all_text_sql(jsonData map[string]interface{}) string {
	db := Connect_To_Database_Second()

	var paragraphs []models.Paragraphs
	db.Where("user_id = ? AND document_id AND paragraph_id", jsonData["userID"].(string), jsonData["documentID"].(int), jsonData["paragraphID"].(int)).Find(&paragraphs)
	var jsonArray []string
	for _, doc := range paragraphs {
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
			"text": jsonMap["text"].(string),
		}
		jsonMaps = append(jsonMaps, jsonFilter)
	}

	text := ""
	for _, value := range jsonMaps {
		text += value["text"].(string) + "\n"

	}
	return text
}
