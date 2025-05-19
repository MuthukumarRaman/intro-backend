package helper

import (
	"time"
)

var OrgList = make(map[string]Organization)

type Organization struct {
	Id     string      `json:"_id" bson:"_id"`
	Name   string      `json:"name" bson:"name"`
	Type   string      `json:"type" bson:"type"`
	Domain string      `json:"domain" bson:"domain"`
	Group  string      `json:"group" bson:"group"`
	Style  interface{} `json:"style" bson:"style"`
}

type UserToken struct {
	UserId   string `json:"user_id" bson:"user_id"`
	UserRole string `json:"user_role" bson:"user_role"`
	OrgId    string `json:"org_id" bson:"org_id"`
	OrgGroup string `json:"uo_group" bson:"uo_group"`
	Org_name string `json:"org_name" bson:"org_name"`
}

type CreatedOnData struct {
	CreatedOn time.Time `json:"created_on" bson:"created_on"`
	CreatedBy string    `json:"created_by" bson:"created_by"`
}

type Leave struct {
	Type      string    `json:"type" bson:"type"`
	EmpId     string    `json:"emp_id" bson:"emp_id"`
	StartDate time.Time `json:"start_date" bson:"start_date"`
	EndDate   time.Time `json:"end_date" bson:"end_date"`
	Status    string    `json:"status" bson:"status"`
}

type GroupSumRequest struct {
	CollectionName string   `json:"collection_name" bson:"collection_name"`
	DateColumn     string   `json:"date_column" bson:"date_column"`
	DateFormat     string   `json:"date_format" bson:"date_format"`
	GroupBy        string   `json:"group_by" bson:"group_by"`
	OutputColumns  []string `json:"output_columns" bson:"output_columns"`
	// Filter         []Filter `json:"filter,omitempty" bson:"filter,omitempty"`
}

type GroupReportRequest struct {
	CollectionName string   `json:"collection_name" bson:"collection_name"`
	GroupBy        []string `json:"group_by" bson:"group_by"`
	OutputColumns  []string `json:"output_columns" bson:"output_columns"`
	// Filter         []Filter `json:"filter,omitempty" bson:"filter,omitempty"`
}

type ReportRequest struct {
	OrgId      string    `json:"org_id" bson:"org_id"`
	EmpId      string    `json:"emp_id" bson:"emp_id"`
	EmpIds     []string  `json:"emp_ids" bson:"emp_ids"`
	Type       string    `json:"type" bson:"type"`
	DateColumn string    `json:"date_column" bson:"date_column"`
	StartDate  time.Time `json:"start_date" bson:"start_date"`
	EndDate    time.Time `json:"end_date" bson:"end_date"`
	Status     string    `json:"status" bson:"status"`
}

// type Condition struct {
// 	Column   string `json:"column" bson:"column"`
// 	Operator string `json:"operator" bson:"operator"`
// 	Type     string `json:"type" bson:"type"`
// 	Value    string `json:"value" bson:"value"`
// }

type PreSignedUploadUrlRequest struct {
	FolderPath string             `json:"folder_path" bson:"folder_path"`
	FileKey    string             `json:"file_key" bson:"file_key"`
	MetaData   map[string]*string `json:"metadata" bson:"metadata"`
}

// type Filter struct {
// 	Clause     string      `json:"clause" bson:"clause"`
// 	Conditions []Condition `json:"conditions" bson:"conditions"`
// }

type LookupQuery struct {
	Operation string        `json:"operation" bson:"operation"`
	ParentRef CollectionRef `json:"parent_collection" bson:"parent_collection"`
	ChildRef  CollectionRef `json:"child_collection" bson:"child_collection"`
}

type CollectionRef struct {
	Name    string   `json:"name" bson:"name"`
	Key     string   `json:"key" bson:"key"`
	Columns []string `json:"columns,omitempty" bson:"columns,omitempty"`
	// Filter  []Filter `json:"filter,omitempty" bson:"filter,omitempty"`
}

type EmailServerConfig struct {
	OrgId    string `json:"org_id" bson:"org_id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"user_name" bson:"user_name"`
	Password string `json:"password" bson:"password"`
}

type DeleteBookingConsignment struct {
	StartId string `json:"start_id"`
	EndId   string `json:"end_id"`
}

type Screen struct {
	ID         string      `json:"_id" bson:"_id"`
	Type       string      `json:"type" bson:"type" validate:"required"`
	Name       string      `json:"name" bson:"name" validate:"required"`
	Config     interface{} `json:"config" bson:"config" validate:"required"`
	CreatedOn  time.Time   `json:"created_on" bson:"created_on" validate:omitempty"`
	CreatedBy  string      `json:"created_by" bson:"created_by" validate:omitempty"`
	Updated_by string      `json:"updated_by" bson:"updated_by" validate:omitempty"`
	Updated_on time.Time   `json:"updated_on" bson:"updated_on" validate:omitempty"`
	Status     string      `json:"status" bson:"status" validate:omitempty""`
}
type model_config struct {
	ModelName      string    `json:"model_name" bson:"model_name"`
	CollectionName string    `json:"collection_name" bson:"collection_name"`
	IsCollection   string    `json:"is_collection" bson:"is_collection"`
	CreatedOn      time.Time `json:"created_on" bson:"created_on" validate:"omitempty"`
	CreatedBy      string    `json:"created_by" bson:"created_by" validate:"omitempty"`
	Status         string    `json:"status" bson:"status" validate:"omitempty"`
}

type Link struct {
	Req         string
	URL         string
	Expiration  time.Time
	CreatedDate time.Time
	Appcode     string
}

type InsertDataResponse struct {
	ValidationErrors map[string]string
	InsertionError   error
}

type FieldValuePair struct {
	FieldName  string      `json:"fieldname" bson:"fieldname"`
	FieldValue interface{} `json:"fieldvalue" bson:"fieldvalue"`
}

//	type PaginationRequest struct {
//		Start            int               `json:"start,omitempty" bson:"start,omitempty" validate:"omitempty"`
//		End              int               `json:"end,omitempty" bson:"end,omitempty" validate:"omitempty"`
//		CreatedBy        string            `json:"createdby,omitempty"`
//		CreatedOn        time.Time         `json:"createdon,omitempty"`
//		FilterColumns    []FieldValuePair  `json:"filterColumns,omitempty" bson:"filterColumns,omitempty" validate:"omitempty"`
//		Filter           []FilterCondition `json:"filter,omitempty" bson:"filter,omitempty" validate:"omitempty"`
//		Sort             []SortCriteria    `json:"sort,omitempty" bson:"sort,omitempty" validate:"omitempty"`
//		Status           string            `json:"status,omitempty" bson:"status,omitempty" validate:"omitempty"`
//		Groupname        string            `json:"group_name,omitempty" bson:"group_nam,omitempty" validate:"omitempty"`
//		GroupDescription string            `json:"groupDescription,omitempty" bson:"groupDescription,omitempty"`
//		GroupType        string            `json:"grouptype,omitempty" bson:"grouptype,omitempty"`
//		FilterParams     []FilterParam     `json:"FilterParams,omitempty" bson:"FilterParams,omitempty"`
//	}
type PaginationRequest struct {
	Start            int               `json:"start,omitempty" bson:"start,omitempty" validate:"omitempty"`
	End              int               `json:"end,omitempty" bson:"end,omitempty" validate:"omitempty"`
	CreatedBy        string            `json:"createdby,omitempty"`
	CreatedOn        time.Time         `json:"createdon,omitempty"`
	FilterColumns    []FieldValuePair  `json:"filterColumns,omitempty" bson:"filterColumns,omitempty" validate:"omitempty"`
	Filter           []FilterCondition `json:"filter,omitempty" bson:"filter,omitempty" validate:"omitempty"`
	Sort             []SortCriteria    `json:"sort,omitempty" bson:"sort,omitempty" validate:"omitempty"`
	Status           string            `json:"status,omitempty" bson:"status,omitempty" validate:"omitempty"`
	FilterParam      []FilterParam     `json:"filterParams,omitempty" bson:"filterParams,omitempty"`
	Groupname        string            `json:"group_name,omitempty" bson:"group_nam,omitempty" validate:"omitempty"`
	GroupDescription string            `json:"groupDescription,omitempty" bson:"groupDescription,omitempty"`
	GroupType        string            `json:"grouptype,omitempty" bson:"grouptype,omitempty"`
}
type SortCriteria struct {
	Sort  string `json:"sort"`
	ColID string `json:"colId"`
}
