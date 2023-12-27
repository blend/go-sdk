/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"reflect"
	"strings"
	"sync"

	"github.com/blend/go-sdk/stringutil"
)

var (
	metaCacheMu	sync.RWMutex
	metaCache	= make(map[string]*ColumnCollection)
)

// --------------------------------------------------------------------------------
// Common helpers
// --------------------------------------------------------------------------------

// Columns returns the cached column metadata for an object.
func Columns(object DatabaseMapped) *ColumnCollection {
	objectType := reflect.TypeOf(object)
	return ColumnsFromType(newColumnCacheKey(objectType), objectType)
}

// ColumnsFromType reflects a reflect.Type into a column collection.
// The results of this are cached for speed.
func ColumnsFromType(identifier string, t reflect.Type) *ColumnCollection {
	// check with read lock ...
	metaCacheMu.RLock()
	if value, ok := metaCache[identifier]; ok {
		metaCacheMu.RUnlock()
		return value
	}
	metaCacheMu.RUnlock()

	// grab write lock ...
	metaCacheMu.Lock()
	defer metaCacheMu.Unlock()

	// double checked lock
	if value, ok := metaCache[identifier]; ok {
		return value
	}

	metadata := NewColumnCollection(generateColumnsForType(nil, t)...)
	metaCache[identifier] = metadata
	return metadata
}

// --------------------------------------------------------------------------------
// Utility
// --------------------------------------------------------------------------------

// ColumnNamesCSV returns a csv of column names.
func ColumnNamesCSV(object DatabaseMapped) string {
	return Columns(object).ColumnNamesCSV()
}

// --------------------------------------------------------------------------------
// Column Collection
// --------------------------------------------------------------------------------

// NewColumnCollection returns a new empty column collection.
func NewColumnCollection(columns ...Column) *ColumnCollection {
	cc := ColumnCollection{
		columns: columns,
	}
	lookup := make(map[string]*Column)
	for i := 0; i < len(columns); i++ {
		col := &columns[i]
		lookup[col.ColumnName] = col
	}
	cc.lookup = lookup
	return &cc
}

// NewColumnCollectionWithPrefix makes a new column collection with a column prefix.
func NewColumnCollectionWithPrefix(columnPrefix string, columns ...Column) *ColumnCollection {
	cc := ColumnCollection{
		columns: columns,
	}
	lookup := make(map[string]*Column)
	for i := 0; i < len(columns); i++ {
		col := &columns[i]
		lookup[col.ColumnName] = col
	}
	cc.lookup = lookup
	cc.columnPrefix = columnPrefix
	return &cc
}

// ColumnCollection represents the column metadata for a given struct.
type ColumnCollection struct {
	columns		[]Column
	lookup		map[string]*Column
	columnPrefix	string

	autos		*ColumnCollection
	notAutos	*ColumnCollection
	readOnly	*ColumnCollection
	notReadOnly	*ColumnCollection
	primaryKeys	*ColumnCollection
	notPrimaryKeys	*ColumnCollection
	uniqueKeys	*ColumnCollection
	notUniqueKeys	*ColumnCollection
	insertColumns	*ColumnCollection
	updateColumns	*ColumnCollection
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
	return NewColumnCollectionWithPrefix(cc.columnPrefix, cc.columns...)
}

// CopyWithColumnPrefix applies a column prefix to column names and returns a new column collection.
func (cc ColumnCollection) CopyWithColumnPrefix(prefix string) *ColumnCollection {
	return NewColumnCollectionWithPrefix(prefix, cc.columns...)
}

// InsertColumns are non-auto, non-readonly columns.
func (cc *ColumnCollection) InsertColumns() *ColumnCollection {
	if cc.insertColumns != nil {
		return cc.insertColumns
	}

	cc.insertColumns = cc.NotReadOnly().NotAutos()
	return cc.insertColumns
}

// UpdateColumns are non-primary key, non-readonly columns.
func (cc *ColumnCollection) UpdateColumns() *ColumnCollection {
	if cc.updateColumns != nil {
		return cc.updateColumns
	}

	cc.updateColumns = cc.NotReadOnly().NotPrimaryKeys()
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

// Zero returns unset fields on an instance that correspond to fields in the column collection.
func (cc *ColumnCollection) Zero(instance interface{}) *ColumnCollection {
	objValue := ReflectValue(instance)
	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	var fieldValue reflect.Value
	for _, c := range cc.columns {
		fieldValue = objValue.Field(c.Index)
		if fieldValue.IsZero() {
			newCC.Add(c)
		}
	}
	return newCC
}

// NotZero returns set fields on an instance that correspond to fields in the column collection.
func (cc *ColumnCollection) NotZero(instance interface{}) *ColumnCollection {
	objValue := ReflectValue(instance)
	newCC := NewColumnCollectionWithPrefix(cc.columnPrefix)
	var fieldValue reflect.Value
	for _, c := range cc.columns {
		fieldValue = objValue.Field(c.Index)
		if !fieldValue.IsZero() {
			newCC.Add(c)
		}
	}
	return newCC
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
			names[x] = cc.columnPrefix + c.ColumnName
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
			lookup[cc.columnPrefix+key] = value
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
		if cc.columnPrefix != "" {
			names[x] = tableAlias + "." + c.ColumnName + " as " + cc.columnPrefix + c.ColumnName
		} else {
			names[x] = tableAlias + "." + c.ColumnName
		}
	}
	return names
}

// ColumnNamesCSVFromAlias returns the string names for all the columns in the collection.
func (cc *ColumnCollection) ColumnNamesCSVFromAlias(tableAlias string) string {
	return stringutil.CSV(cc.ColumnNamesFromAlias(tableAlias))
}

// ColumnValues returns the reflected value for all the columns on a given instance.
func (cc *ColumnCollection) ColumnValues(instance interface{}) []interface{} {
	value := ReflectValue(instance)

	values := make([]interface{}, len(cc.columns))
	for x := 0; x < len(cc.columns); x++ {
		c := cc.columns[x]
		valueField := value.FieldByName(c.FieldName)
		if c.IsJSON {
			values[x] = JSON(valueField.Interface())
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
	return NewColumnCollection(total...)
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
			if col.Inline && field.Anonymous {	// if it's not anonymous, whatchu doin
				cols = append(cols, generateColumnsForType(col, col.FieldType)...)
			} else if !field.Anonymous {
				cols = append(cols, *col)
			}
		}
	}

	return cols
}
