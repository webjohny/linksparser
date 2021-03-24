package mysql

import (
	"database/sql"
)

type Site struct {
	Id sql.NullInt64 `db:"id" json:"id"`
	Language sql.NullString `db:"language" json:"language"`
	Status sql.NullInt32 `db:"status" json:"status"`
	Theme sql.NullString `db:"theme" json:"theme"`
	Domain sql.NullString `db:"domain" json:"domain"`
	Login sql.NullString `db:"login" json:"login"`
	Password sql.NullString `db:"password" json:"password"`
	Info sql.NullString `db:"info" json:"info"`
	MoreTags sql.NullString `db:"more_tags" json:"more_tags"`
	Extra sql.NullString `db:"extra" json:"extra"`
	SymbMicroMarking sql.NullString `db:"symb_micro_marking" json:"symb_micro_marking"`
	CountRows sql.NullInt64 `db:"count_rows" json:"count_rows"`
	CreatedAt sql.NullString `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at" json:"deleted_at"`
}

type Config struct {
	Id sql.NullInt64 `db:"id" json:"id"`
	FlickrKey sql.NullString `db:"flickr_key" json:"flickr_key"`
	FlickrSecret sql.NullString `db:"flickr_secret" json:"flickr_secret"`
	Antigate sql.NullString `db:"antigate" json:"antigate"`
	Language sql.NullString `db:"language" json:"language"`
	Variants sql.NullString `db:"variants" json:"variants"`
	Extra sql.NullString `db:"extra" json:"extra"`
	CreatedAt sql.NullString `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at" json:"deleted_at"`
}

type Cat struct {
	Id sql.NullInt64 `db:"id" json:"id"`
	SiteId sql.NullInt64 `db:"site_id" json:"site_id"`
	Title sql.NullString `db:"title" json:"title"`
	CreatedAt sql.NullString `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at" json:"deleted_at"`
}

type Uagent struct {
	Id sql.NullInt64 `db:"id" json:"id"`
	Sign sql.NullString `db:"sign" json:"sign"`
	Status sql.NullInt32 `db:"status" json:"status"`
	Timeout sql.NullString `db:"timeout" json:"timeout"`
	CreatedAt sql.NullString `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at" json:"deleted_at"`
}

type Result struct {
	Id sql.NullInt64 `db:"id" json:"id"`
	TaskId sql.NullInt64 `db:"task_id" json:"task_id"`
	CatId sql.NullInt64 `db:"cat_id" json:"cat_id"`
	SiteId sql.NullInt64 `db:"site_id" json:"site_id"`
	Keyword sql.NullString `db:"keyword" json:"keyword"`
	Domain sql.NullString `db:"domain" json:"domain"`
	Links sql.NullString `db:"links" json:"links"`
	Author sql.NullString `db:"author" json:"author"`
	Content sql.NullString `db:"content" json:"content"`
	Text sql.NullString `db:"text" json:"text"`
	CreatedAt sql.NullString `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at" json:"deleted_at"`
}

type Task struct {
	Id sql.NullInt64 `db:"id" json:"id"`
	Log sql.NullString `db:"log" json:"log"`
	LogLast sql.NullString `db:"log_last" json:"log_last"`
	ParentId sql.NullInt64 `db:"parent_id" json:"parent_id"`
	SiteId sql.NullInt64 `db:"site_id" json:"site_id"`
	CatId sql.NullInt64 `db:"cat_id" json:"cat_id"`
	Cat sql.NullString `db:"cat" json:"cat"`
	Keyword sql.NullString `db:"keyword" json:"keyword"`
	TryCount sql.NullInt32 `db:"try_count" json:"try_count"`
	ErrorsCount sql.NullInt32 `db:"errors_count" json:"errors_count"`
	Status sql.NullInt32 `db:"status" json:"status"`
	Error sql.NullString `db:"error" json:"error"`
	Stream sql.NullInt64 `db:"stream" json:"stream"`
	Timeout sql.NullString `db:"timeout" json:"timeout"`
	CreatedAt sql.NullString `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	ParsedAt sql.NullString `db:"parsed_at" json:"parsed_at"`
}

type Proxy struct {
	Id sql.NullInt64 `json:"id"`
	Type sql.NullString `json:"type"`
	Host sql.NullString `json:"host"`
	Port sql.NullString `json:"port"`
	Login sql.NullString `json:"login"`
	Password sql.NullString `json:"password"`
	Agent sql.NullString `json:"agent"`
	Status sql.NullInt64 `json:"status"`
	Stream sql.NullString `json:"stream"`
	Timeout sql.NullString `json:"timeout"`
	CreatedAt sql.NullString `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at" json:"deleted_at"`
}

type Image struct {
	Id sql.NullInt64 `json:"id" db:"id"`
	SiteId sql.NullInt64 `json:"site_id" db:"site_id"`
	SourceId sql.NullString `json:"source_id" db:"source_id"`
	Url sql.NullString `json:"url" db:"url"`
	Author sql.NullString `json:"author" db:"author"`
	ShortUrl sql.NullString `json:"short_url" db:"short_url"`
	Keyword sql.NullString `json:"keyword" db:"keyword"`
	Source sql.NullBool `json:"source" db:"source"`
	Status sql.NullBool `json:"status" db:"status"`
	CreatedAt sql.NullString `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at" json:"deleted_at"`
}

