package models

type Documents struct {
	DocumentID  int    `json:"documentID"`
	UserID      string `json:"userID"`
	UpdateTime  string `json:"updateTime"`
	Path        string `json:"path"`
	Paragraphs  string `json:"paragraphs"`
	IsCollected int    `json:"isCollected"`
}

type Paragraphs struct {
	ParagraphID      int    `json:"paragraphID"`
	Version          int    `json:"version"`
	DocumentID       int    `json"documentID"`
	UserID           string `json"userID"`
	IsFirstlnChapter string `json:"isFirstlnChapter"`
	Text             string `json:"text"`
	Before           string `"json:"before"`
	Next             string `json:"next"`
}

type DynFigures struct {
	DynFigureID int    `json:"dynFigureID"`
	DocumentID  int    `json"documentID"`
	UserID      string `json"userID"`
	Name        string `"json:"name"`
	CurrentTag  string `json:"currentTag"`
	CodeGenTask string `json:"codeGenTask"`
	CodeFig     string `json:"codeFig"`
	TagQueue    string `json:"tagQueue"`
}

type DynTasks struct {
	Name        string `"json:"name"`
	DynFigureID int    `json:"dynFigureID"`
	DocumentID  int    `json"documentID"`
	UserID      string `json"userID"`
	CodeShell   string `json:"codeShell"`
	TimeStart   string `json:"timeStart"`
	TimeEnd     string `json:"timeEnd"`
	Tag         string `json:"tag"`
}

type Servers struct {
	Name             string `"json:"name"`
	UserID           string `json"userID"`
	Ssh_user         string `json:"ssh_user"`
	Ip               string `json:"ip"`
	Port             string `json:"port"`
	Auth_method      string `json:"auth_method"`
	Password         string `json:"password"`
	Key              string `json:"key"`
	JumpServerName   string `json:"jumpServerName"`
	JumpServerUserID string `json:"jumpServerUserID"`
	Login_command    string `json:"login_command"`
}

type Paragraph_templates struct {
	Name   string `"json:"name"`
	UserID string `json"userID"`
	Tag    string `json:"tag"`
	Text   string `json:"Text"`
}

type File struct {
	FileID     int    `json:"fileID"`
	DocumentID int    `json:"documentID"`
	UserID     string `json:"userID"`
	Path       string `json:"path"`
	Value      string `json:"path"`
}
