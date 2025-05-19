package helper

type FilterCondition struct {
	Clause     string           `json:"clause,omitempty" bson:"clause,omitempty"`
	Conditions []ConditionGroup `json:"conditions,omitempty" bson:"condition,omitempty"`
}

type ConditionGroup struct {
	Type                 string           `json:"type,omitempty" bson:"type,omitempty"`
	Column               string           `json:"column,omitempty" bson:"column,omitempty"`
	Operator             string           `json:"operator,omitempty" bson:"operator,omitempty"`
	Value                interface{}      `json:"value,omitempty" bson:"value,omitempty"`
	Clause               string           `json:"clause,omitempty" bson:"clause,omitempty"`
	ValueType            interface{}      `json:"value_type,omitempty" bson:"value_type,omitempty"`
	ParentCollectionName string           `json:"parentCollectionName,omitempty" bson:"parentCollectionName,omitempty"`
	Conditions           []ConditionGroup `json:"conditions,omitempty" bson:"conditions,omitempty"`
}

type AggregationField struct {
	Name                 string `json:"name,omitempty" bson:"name,omitempty"`
	Field                string `json:"field_name,omitempty" bson:"field_name,omitempty"`
	ParentCollectionName string `json:"parentCollectionName,omitempty" bson:"parentCollectionName,omitempty"`
	Type                 string `json:"type,omitempty" bson:"type,omitempty"`
	Hide                 bool   `json:"hide,omitempty" bson:"hide,omitempty"`
}

type Aggregation struct {
	AggFieldName    AggregationField `json:"Agg_Field_Name,omitempty" bson:"Agg_Field_Name,omitempty"`
	AggFnName       string           `json:"Agg_Fn_Name,omitempty" bson:"Agg_Fn_Name,omitempty"`
	AggGroupByField AggregationField `json:"Agg_group_byField,omitempty" bson:"Agg_group_byField,omitempty"`
	ConvertToString bool             `json:"convert_To_String,omitempty" bson:"convert_To_String,omitempty"`
	AggColumnName   string           `json:"Agg_Column_Name,omitempty" bson:"Agg_Column_Name,omitempty"`
}

//	type CustomField struct {
//		Field           string `json:"field_name,omitempty" bson:"field_name,omitempty"`
//		Type            string `json:"type,omitempty" bson:"type,omitempty"`
//		Hide            bool   `json:"hide,omitempty" bson:"hide,omitempty"`
//		Name            string `json:"name,omitempty" bson:"name,omitempty"`
//		Reference       bool   `json:"reference,omitempty" bson:"reference,omitempty"`
//		ConvertToString bool   `json:"convert_To_String,omitempty" bson:"convert_To_String,omitempty"`
//	}
type CustomField struct {
	Name                 string `json:"name,omitempty" bson:"name,omitempty"`
	FieldName            string `json:"field_name,omitempty" bson:"field_name,omitempty"`
	ParentCollectionName string `json:"parentCollectionName,omitempty" bson:"parentCollectionName,omitempty"`
	Reference            bool   `json:"reference,omitempty" bson:"reference,omitempty"`
	Type                 string `json:"type,omitempty" bson:"type,omitempty"`
}
type CustomColumn struct {
	DataSetCustomLabelName       string        `json:"dataSetCustomLabelName,omitempty" bson:"dataSetCustomLabelName,omitempty"`
	DataSetCustomAggregateFnName string        `json:"dataSetCustomAggregateFnName,omitempty" bson:"dataSetCustomAggregateFnName,omitempty"`
	DataSetCustomField           []CustomField `json:"dataSetCustomField,omitempty" bson:"dataSetCustomField,omitempty"`
	ConvertToString              bool          `json:"convert_To_String,omitempty" bson:"convert_To_String,omitempty"`
	DataSetCustomColumnName      string        `json:"dataSetCustomColumnName,omitempty" bson:"dataSetCustomColumnName,omitempty"`
}

type SelectedListItem struct {
	HeaderName string `json:"headerName,omitempty" bson:"headerName,omitempty"`
	Field      string `json:"field,omitempty" bson:"field,omitempty"`
}

type FilterParam struct {
	ConvertToString bool        `json:"convert_To_String,omitempty" bson:"convert_To_String,omitempty"`
	ParamsName      string      `json:"parmasName,omitempty" bson:"parmasName,omitempty"`
	ParamsDataType  string      `json:"parmsDataType,omitempty" bson:"parmsDataType,omitempty"`
	DefaultValue    interface{} `json:"defaultValue,omitempty" bson:"defaultValue,omitempty"`
	Paramsvalue     interface{} `json:"paramsvalue,omitempty" bson:"paramsvalue,omitempty"`
}

type DataSetConfiguration struct {
	Id                          string                  `json:"_id,omitempty" bson:"_id,omitempty"`
	DataSetName                 string                  `json:"dataSetName,omitempty" bson:"dataSetName,omitempty"`
	DataSetDescription          string                  `json:"dataSetDescription,omitempty" bson:"dataSetDescription,omitempty"`
	DataSetJoinCollection       []DataSetJoinCollection `json:"dataSetJoinCollection,omitempty" bson:"dataSetJoinCollection,omitempty"`
	CustomColumn                []CustomColumn          `json:"CustomColumn,omitempty" bson:"CustomColumn,omitempty"`
	SelectedList                []SelectedListItem      `json:"SelectedList,omitempty" bson:"SelectedList,omitempty"`
	FilterParams                []FilterParam           `json:"FilterParams,omitempty" bson:"FilterParams,omitempty"`
	DataSetBaseCollection       string                  `json:"dataSetBaseCollection,omitempty" bson:"dataSetBaseCollection,omitempty"`
	Aggregation                 []Aggregation           `json:"Aggregation,omitempty" bson:"Aggregation,omitempty"`
	Filter                      []FilterCondition       `json:"Filter,omitempty" bson:"Filter,omitempty"`
	DataSetBaseCollectionFilter []FilterCondition       `json:"dataSetBaseCollectionFilter,omitempty, bson:"dataSetBaseCollectionFilter,omitempty"`
	Pipeline                    string                  `json:"pipeline,omitempty" bson:"pipeline,omitempty"`
	Reference_pipeline          string                  `json:"reference_pipeline,omitempty" bson:"reference_pipeline,omitempty"`
	Start                       int                     `json:"start,omitempty" bson:"start,omitempty"`
	End                         int                     `json:"end,omitempty" bson:"end,omitempty"`
}

type DataSetJoinCollection struct {
	ToCollection        string            `json:"toCollection,omitempty" bson:"toCollection,omitempty"`
	ToCollectionField   string            `json:"toCollectionField,omitempty" bson:"toCollectionField,omitempty"`
	Filter              []FilterCondition `json:"Filter,omitempty" bson:"Filter,omitempty"`
	ConvertToString     bool              `json:"convert_To_String,omitempty" bson:"convert_To_String,omitempty"`
	FromCollection      string            `json:"fromCollection,omitempty" bson:"fromCollection,omitempty"`
	FromCollectionField string            `json:"fromCollectionField,omitempty" bson:"fromCollectionField,omitempty"`
}
