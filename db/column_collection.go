package db

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"github.com/blend/go-sdk/stringutil"
)

var (
	metaCacheLock sync.Mutex
	metaCache     map[string]*ColumnCollection
)

// --------------------------------------------------------------------------------
// Utility
// --------------------------------------------------------------------------------

// ColumnNamesCSV returns a csv of column names.
func ColumnNamesCSV(object DatabaseMapped) string {
	return CachedColumnCollectionFromInstance(object).ColumnNamesCSV()
}

// Columns returns the cached column metadata for an object.
func Columns(object DatabaseMapped) *ColumnCollection {
	return CachedColumnCollectionFromInstance(object)
}

// --------------------------------------------------------------------------------
// Column Collection
// --------------------------------------------------------------------------------

// NewColumnCollection returns a new empty column collection.
func NewColumnCollection() *ColumnCollection { return &ColumnCollection{lookup: map[string]*Column{}} }

// NewColumnCollectionWithPrefix makes a new column collection with a column prefix.
func NewColumnCollectionWithPrefix(columnPrefix string) *ColumnCollection {
	return &ColumnCollection{lookup: map[string]*Column{}, columnPrefix: columnPrefix}
}

// ColumnCollection represents the column metadata for a given struct.
type ColumnCollection struct {
	columns      []Column
	lookup       map[string]*Column
	columnPrefix string

	autos          *ColumnCollection
	notAutos       *ColumnCollection
	readOnly       *ColumnCollection
	notReadOnly    *ColumnCollection
	primaryKeys    *ColumnCollection
	notPrimaryKeys *ColumnCollection
	uniqueKeys     *ColumnCollection
	notUniqueKeys  *ColumnCollection
	writeColumns   *ColumnCollection
	updateColumns  *ColumnCollection
}

// Len returns the number of columns.
func (cc *ColumnCollection) Len() int {
	if cc == nil {
		return 0
	}
	return len(cc.columns)
}

// Add adds a column.
func (cc *ColumnCollection) Add(c Column) {
	cc.columns = append(cc.columns, c)
	cc.lookup[c.ColumnName] = &c
}

// Remove removes a column (by column name) from the collection.
func (cc *ColumnCollection) Remove(columnName string) {
	var newColumns []Column
	for _, c := range cc.columns {
		if c.ColumnName != columnName {
			newColumns = append(newColumns, c)
		}
	}
	cc.columns = newColumns
	delete(cc.lookup, columnName)
}

// HasColumn returns if a column name is present in the collection.
func (cc *ColumnCollection) HasColumn(columnName string) bool {
	_, hasColumn := cc.lookup[columnName]
	return hasColumn
}

// Copy creates a new column collection instance and carries over an existing column prefix.
func (cc ColumnCollection) Copy() *ColumnCollection {
	return newColumnCollectionWithPrefixFromColumns(cc.columnPrefix, cc.columns)
}

// CopyWithColumnPrefix applies a column prefix to column names and returns a new column collection.
func (cc ColumnCollection) CopyWithColumnPrefix(prefix string) *ColumnCollection {
	return newColumnCollectionWithPrefixFromColumns(prefix, cc.columns)
}

// WriteColumns are non-auto, non-primary key, non-readonly columns.
func (cc *ColumnCollection) WriteColumns() *ColumnCollection {
	if cc.writeColumns != nil {
		return cc.writeColumns
	}

	cc.writeColumns = cc.NotReadOnly().NotAutos()
	return cc.writeColumns
}

// UpdateColumns are non-readonly, non-serial columns.
func (cc *ColumnCollection) UpdateColumns() *ColumnCollection {
	if cc.updateColumns != nil {
		return cc.updateColumns
	}

	cc.updateColumns = cc.NotReadOnly()
	return cc.updateColumns
}

// PrimaryKeys are columns we use as where predicates and can't update.
func (cc *ColumnCollection) PrimaryKeys() *ColumnCollection {
	if cc.primaryKeys != nil {
		return cc.primaryKeys
	}

	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	for _, c := range cc.columns {
		if c.IsPrimaryKey {
			newCC.Add(c)
		}
	}

	cc.primaryKeys = newCC
	return cc.primaryKeys
}

// NotPrimaryKeys are columns we can update.
func (cc *ColumnCollection) NotPrimaryKeys() *ColumnCollection {
	if cc.notPrimaryKeys != nil {
		return cc.notPrimaryKeys
	}

	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)

	for _, c := range cc.columns {
		if !c.IsPrimaryKey {
			newCC.Add(c)
		}
	}

	cc.notPrimaryKeys = newCC
	return cc.notPrimaryKeys
}

// UniqueKeys are columns we use as where predicates and can't update.
func (cc *ColumnCollection) UniqueKeys() *ColumnCollection {
	if cc.uniqueKeys != nil {
		return cc.uniqueKeys
	}

	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	for _, c := range cc.columns {
		if c.IsUniqueKey {
			newCC.Add(c)
		}
	}

	cc.uniqueKeys = newCC
	return cc.uniqueKeys
}

// NotUniqueKeys are columns we can update.
func (cc *ColumnCollection) NotUniqueKeys() *ColumnCollection {
	if cc.notUniqueKeys != nil {
		return cc.notUniqueKeys
	}

	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	for _, c := range cc.columns {
		if !c.IsUniqueKey {
			newCC.Add(c)
		}
	}

	cc.notUniqueKeys = newCC
	return cc.notUniqueKeys
}

// Autos are columns we have to return the id of.
func (cc *ColumnCollection) Autos() *ColumnCollection {
	if cc.autos != nil {
		return cc.autos
	}

	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	for _, c := range cc.columns {
		if c.IsAuto {
			newCC.Add(c)
		}
	}

	cc.autos = newCC
	return cc.autos
}

// NotAutos are columns we don't have to return the id of.
func (cc *ColumnCollection) NotAutos() *ColumnCollection {
	if cc.notAutos != nil {
		return cc.notAutos
	}

	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	for _, c := range cc.columns {
		if !c.IsAuto {
			newCC.Add(c)
		}
	}
	cc.notAutos = newCC
	return cc.notAutos
}

// ReadOnly are columns that we don't have to insert upon Create().
func (cc *ColumnCollection) ReadOnly() *ColumnCollection {
	if cc.readOnly != nil {
		return cc.readOnly
	}

	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	for _, c := range cc.columns {
		if c.IsReadOnly {
			newCC.Add(c)
		}
	}

	cc.readOnly = newCC
	return cc.readOnly
}

// NotReadOnly are columns that we have to insert upon Create().
func (cc *ColumnCollection) NotReadOnly() *ColumnCollection {
	if cc.notReadOnly != nil {
		return cc.notReadOnly
	}

	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	for _, c := range cc.columns {
		if !c.IsReadOnly {
			newCC.Add(c)
		}
	}

	cc.notReadOnly = newCC
	return cc.notReadOnly
}

// ColumnNames returns the string names for all the columns in the collection.
func (cc *ColumnCollection) ColumnNames() []string {
	if cc == nil {
		return nil
	}
	names := make([]string, len(cc.columns))
	for x := 0; x < len(cc.columns); x++ {
		c := cc.columns[x]
		if len(cc.columnPrefix) != 0 {
			names[x] = fmt.Sprintf("%s%s", cc.columnPrefix, c.ColumnName)
		} else {
			names[x] = c.ColumnName
		}
	}
	return names
}

// Columns returns the colummns
func (cc *ColumnCollection) Columns() []Column {
	return cc.columns
}

// Lookup gets the column name lookup.
func (cc *ColumnCollection) Lookup() map[string]*Column {
	if len(cc.columnPrefix) != 0 {
		lookup := map[string]*Column{}
		for key, value := range cc.lookup {
			lookup[fmt.Sprintf("%s%s", cc.columnPrefix, key)] = value
		}
		return lookup
	}
	return cc.lookup
}

// ColumnNamesFromAlias returns the string names for all the columns in the collection.
func (cc *ColumnCollection) ColumnNamesFromAlias(tableAlias string) []string {
	names := make([]string, len(cc.columns))
	for x := 0; x < len(cc.columns); x++ {
		c := cc.columns[x]
		if len(cc.columnPrefix) != 0 {
			names[x] = fmt.Sprintf("%s.%s as %s%s", tableAlias, c.ColumnName, cc.columnPrefix, c.ColumnName)
		} else {
			names[x] = fmt.Sprintf("%s.%s", tableAlias, c.ColumnName)
		}
	}
	return names
}

// ColumnNamesCSVFromAlias returns the string names for all the columns in the collection.
func (cc *ColumnCollection) ColumnNamesCSVFromAlias(tableAlias string) string {
	names := make([]string, len(cc.columns))
	for x := 0; x < len(cc.columns); x++ {
		c := cc.columns[x]
		if len(cc.columnPrefix) != 0 {
			names[x] = fmt.Sprintf("%s.%s as %s%s", tableAlias, c.ColumnName, cc.columnPrefix, c.ColumnName)
		} else {
			names[x] = fmt.Sprintf("%s.%s", tableAlias, c.ColumnName)
		}
	}
	return stringutil.CSV(names)
}

// ColumnValues returns the reflected value for all the columns on a given instance.
func (cc *ColumnCollection) ColumnValues(instance interface{}) []interface{} {
	value := ReflectValue(instance)

	values := make([]interface{}, len(cc.columns))
	for x := 0; x < len(cc.columns); x++ {
		c := cc.columns[x]
		valueField := value.FieldByName(c.FieldName)
		if c.IsJSON {
			jsonBytes, _ := json.Marshal(valueField.Interface())
			if result := string(jsonBytes); result != "null" { // explicitly bad.
				values[x] = result
			} else {
				values[x] = nil
			}
		} else {
			values[x] = valueField.Interface()
		}
	}
	return values
}

// FirstOrDefault returns the first column in the collection or `nil` if the collection is empty.
func (cc *ColumnCollection) FirstOrDefault() *Column {
	if len(cc.columns) > 0 {
		return &cc.columns[0]
	}
	return nil
}

// ConcatWith merges a collection with another collection.
func (cc *ColumnCollection) ConcatWith(other *ColumnCollection) *ColumnCollection {
	total := make([]Column, len(cc.columns)+len(other.columns))
	var x int
	for ; x < len(cc.columns); x++ {
		total[x] = cc.columns[x]
	}
	for y := 0; y < len(other.columns); y++ {
		total[x+y] = other.columns[y]
	}
	return newColumnCollectionFromColumns(total)
}

func (cc *ColumnCollection) String() string {
	return strings.Join(cc.ColumnNames(), ", ")
}

// ColumnNamesCSV returns a csv of column names.
func (cc *ColumnCollection) ColumnNamesCSV() string {
	return stringutil.CSV(cc.ColumnNames())
}

//
// helpers
//

// newColumnCollectionFromColumns creates a column lookup for a slice of columns.
func newColumnCollectionFromColumns(columns []Column) *ColumnCollection {
	cc := ColumnCollection{columns: columns}
	lookup := make(map[string]*Column)
	for i := 0; i < len(columns); i++ {
		col := &columns[i]
		lookup[col.ColumnName] = col
	}
	cc.lookup = lookup
	return &cc
}

// newColumnCollectionWithPrefixFromColumns creates a column lookup for a slice of columns.
func newColumnCollectionWithPrefixFromColumns(prefix string, columns []Column) *ColumnCollection {
	cc := ColumnCollection{columns: columns, columnPrefix: prefix}
	lookup := make(map[string]*Column)
	for i := 0; i < len(columns); i++ {
		col := &columns[i]
		lookup[col.ColumnName] = col
	}
	cc.lookup = lookup
	return &cc
}

// newColumnCacheKey creates a cache key for a type.
func newColumnCacheKey(objectType reflect.Type) string {
	typeName := objectType.String()
	instance := reflect.New(objectType).Interface()
	if typed, ok := instance.(ColumnMetaCacheKeyProvider); ok {
		return typeName + "_" + typed.ColumnMetaCacheKey()
	}
	if typed, ok := instance.(TableNameProvider); ok {
		return typeName + "_" + typed.TableName()
	}
	return typeName
}

// CachedColumnCollectionFromInstance reflects an object instance into a new column collection.
func CachedColumnCollectionFromInstance(object interface{}) *ColumnCollection {
	objectType := reflect.TypeOf(object)
	return CachedColumnCollectionFromType(newColumnCacheKey(objectType), objectType)
}

// CachedColumnCollectionFromType reflects a reflect.Type into a column collection.
// The results of this are cached for speed.
func CachedColumnCollectionFromType(identifier string, t reflect.Type) *ColumnCollection {
	metaCacheLock.Lock()
	defer metaCacheLock.Unlock()

	if metaCache == nil {
		metaCache = map[string]*ColumnCollection{}
	}

	cachedMeta, ok := metaCache[identifier]
	if !ok {
		metadata := newColumnCollectionFromColumns(generateColumnsForType(nil, t))
		metaCache[identifier] = metadata
		return metadata
	}
	return cachedMeta
}

// generateColumnsForType generates a column list for a given type.
func generateColumnsForType(parent *Column, t reflect.Type) []Column {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var tableName string
	if parent != nil {
		tableName = parent.TableName
	} else {
		tableName = TableNameByType(t)
	}

	numFields := t.NumField()

	var cols []Column
	for index := 0; index < numFields; index++ {
		field := t.Field(index)
		col := NewColumnFromFieldTag(field)
		if col != nil {
			col.Parent = parent
			col.Index = index
			col.TableName = tableName
			if col.Inline && field.Anonymous { // if it's not anonymous, whatchu doin
				cols = append(cols, generateColumnsForType(col, col.FieldType)...)
			} else if !field.Anonymous {
				cols = append(cols, *col)
			}
		}
	}

	return cols
}
