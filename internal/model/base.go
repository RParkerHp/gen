package model

import (
	"bytes"
	"strings"

	"gorm.io/gen/field"
)

const (
	// DefaultModelPkg ...
	DefaultModelPkg = "model"
)

// Status sql status
type Status int

const (
	// UNKNOWN ...
	UNKNOWN Status = iota
	// SQL ...
	SQL
	// DATA ...
	DATA
	// VARIABLE ...
	VARIABLE
	// IF ...
	IF
	// ELSE ...
	ELSE
	// WHERE ...
	WHERE
	// SET ...
	SET
	// FOR ...
	FOR
	// END ...
	END
	// TRIM ...
	TRIM
)

// SourceCode source code
type SourceCode int

const (
	// Struct ...
	Struct SourceCode = iota
	// Table ...
	Table
	// Object ...
	Object
)

// GormKeywords ...
var GormKeywords = KeyWord{
	words: []string{
		"UnderlyingDB", "UseDB", "UseModel", "UseTable", "Quote", "Debug", "TableName", "WithContext",
		"As", "Not", "Or", "Build", "Columns", "Hints",
		"Distinct", "Omit",
		"Select", "Where", "Order", "Group", "Having", "Limit", "Offset",
		"Join", "LeftJoin", "RightJoin",
		"Save", "Create", "CreateInBatches",
		"Update", "Updates", "UpdateColumn", "UpdateColumns",
		"Find", "FindInBatches", "First", "Take", "Last", "Pluck", "Count",
		"Scan", "ScanRows", "Row", "Rows",
		"Delete", "Unscoped",
		"Scopes",
	},
}

// DOKeywords ...
var DOKeywords = KeyWord{
	words: []string{
		"Alias", "TableName", "WithContext",
	},
}

// GenKeywords ...
var GenKeywords = KeyWord{
	words: []string{
		"generateSQL", "whereClause", "setClause",
	},
}

// KeyWord ...
type KeyWord struct {
	words []string
}

// FullMatch full match
func (g *KeyWord) FullMatch(word string) bool {
	for _, item := range g.words {
		if word == item {
			return true
		}
	}
	return false
}

// Contain contain
func (g *KeyWord) Contain(text string) bool {
	for _, item := range g.words {
		if strings.Contains(text, item) {
			return true
		}
	}
	return false
}

var (
	defaultDataType             = "string"
	dataType        dataTypeMap = map[string]dataTypeMapping{
		"numeric":    func(string) string { return "int32" },
		"integer":    func(string) string { return "int32" },
		"int":        func(string) string { return "int32" },
		"smallint":   func(string) string { return "int32" },
		"mediumint":  func(string) string { return "int32" },
		"bigint":     func(string) string { return "int64" },
		"float":      func(string) string { return "float32" },
		"real":       func(string) string { return "float64" },
		"double":     func(string) string { return "float64" },
		"decimal":    func(string) string { return "float64" },
		"char":       func(string) string { return "string" },
		"varchar":    func(string) string { return "string" },
		"tinytext":   func(string) string { return "string" },
		"mediumtext": func(string) string { return "string" },
		"longtext":   func(string) string { return "string" },
		"binary":     func(string) string { return "[]byte" },
		"varbinary":  func(string) string { return "[]byte" },
		"tinyblob":   func(string) string { return "[]byte" },
		"blob":       func(string) string { return "[]byte" },
		"mediumblob": func(string) string { return "[]byte" },
		"longblob":   func(string) string { return "[]byte" },
		"text":       func(string) string { return "string" },
		"json":       func(string) string { return "string" },
		"enum":       func(string) string { return "string" },
		"time":       func(string) string { return "time.Time" },
		"date":       func(string) string { return "time.Time" },
		"datetime":   func(string) string { return "time.Time" },
		"timestamp":  func(string) string { return "time.Time" },
		"year":       func(string) string { return "int32" },
		"bit":        func(string) string { return "[]uint8" },
		"boolean":    func(string) string { return "bool" },
		"tinyint": func(detailType string) string {
			if strings.HasPrefix(strings.TrimSpace(detailType), "tinyint(1)") {
				return "bool"
			}
			return "int32"
		},
	}
)

type dataTypeMapping func(detailType string) (finalType string)

type dataTypeMap map[string]dataTypeMapping

func (m dataTypeMap) Get(dataType, detailType string) string {
	if convert, ok := m[strings.ToLower(dataType)]; ok {
		return convert(detailType)
	}
	return defaultDataType
}

// GetDataType returns the corresponding Go data type for a given SQL data type string.
// It uses the dataType mapping to look up the Go type associated with the provided SQL data type.
// If the SQL data type is not found in the mapping, it returns an empty string.
//
// Parameters:
//   - sqlDataType: the SQL data type as a string.
//
// Returns:
//   - The Go data type as a string, or an empty string if not found.
func GetDataType(sqlDataType string) string {
	return dataType.Get(sqlDataType, "")
}

// SetDataType registers a mapping function for a specific database type.
// It associates the provided dbType string with the getTypeFunc, which defines
// how to map the database type to the application's data type.
//
// Parameters:
//   - dbType: the name of the database type to map.
//   - getTypeFunc: a function that returns the application's data type for the given dbType.
func SetDataType(dbType string, getTypeFunc dataTypeMapping) {
	dataType[dbType] = getTypeFunc
}

// Field user input structures
type Field struct {
	Name             string
	Type             string
	ColumnName       string
	ColumnComment    string
	MultilineComment bool
	Tag              field.Tag
	GORMTag          field.GormTag
	CustomGenType    string
	Relation         *field.Relation
}

// Tags ...
func (m *Field) Tags() string {
	if _, ok := m.Tag[field.TagKeyGorm]; ok {
		return m.Tag.Build()
	}

	if gormTag := strings.TrimSpace(m.GORMTag.Build()); gormTag != "" {
		m.Tag.Set(field.TagKeyGorm, gormTag)
	}
	return m.Tag.Build()
}

// IsRelation ...
func (m *Field) IsRelation() bool { return m.Relation != nil }

// GenType ...
func (m *Field) GenType() string {
	if m.IsRelation() {
		return m.Type
	}
	if m.CustomGenType != "" {
		return m.CustomGenType
	}
	typ := strings.TrimLeft(m.Type, "*")
	switch typ {
	case "string", "bytes":
		return strings.Title(typ)
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return strings.Title(typ)
	case "float64", "float32":
		return strings.Title(typ)
	case "bool":
		return strings.Title(typ)
	case "time.Time":
		return "Time"
	case "json.RawMessage", "[]byte":
		return "Bytes"
	case "serializer":
		return "Serializer"
	default:
		return "Field"
	}
}

// EscapeKeyword escape keyword
func (m *Field) EscapeKeyword() *Field {
	return m.EscapeKeywordFor(GormKeywords)
}

// EscapeKeywordFor escape for specified keyword
func (m *Field) EscapeKeywordFor(keywords KeyWord) *Field {
	if keywords.FullMatch(m.Name) {
		m.Name += "_"
	}
	return m
}

// SQLBuffer sql buffer
type SQLBuffer struct{ bytes.Buffer }

// WriteSQL ...
func (s *SQLBuffer) WriteSQL(b byte) {
	switch b {
	case '\n', '\t', ' ':
		if s.Len() == 0 || s.Bytes()[s.Len()-1] != ' ' {
			_ = s.WriteByte(' ')
		}
	default:
		_ = s.WriteByte(b)
	}
}

// Dump ...
func (s *SQLBuffer) Dump() string {
	defer s.Reset()
	return s.String()
}
