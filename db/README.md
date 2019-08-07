db
======

This is a very bare bones database interface for golang. It abstracts away a bunch of boilerplate so that the developer can concentrate on writing their app.

It does not abstract away actual sql, however.

# Gotchas & General Notes #

- Stuct to database table / column mapping is done through field tags.
- There are a bunch of helpers for common operations (Get, GetAll, Create, CreateMany, Update, Delete, etc.).
	- These will write sql for you, and generally simplify basic operations.
- Like other packages in this repo, the `db` package leverages the options pattern extensively. Most things can
be configured with functions that start with `Opt...`.
- We leverage statement caching aggressively. What this means is that if a query or exec is assigned a label, we will save the returned query plan for later increasing throughput.
	- The pre-built (`Get`, `All`, `Create`, `CreateMany` etc.) methods create query labels for you.
	- Your statements will not be cached if you don't set a query label.
	- You set the query label by:
		`conn.Invoke(db.OptCachedPlanKey("my_label")).[Query(...)|Exec(...)]`

# Mapping Structs Using `go-sdk/db` #

A sample database mapped type:
```go
type MyTable struct {
	Id int `db:"id,auto,pk"`
	Name string
	Excluded `db:"-"`
}
// note; if we don't do this, the table name will be inferred from the type name.
func (mt MyTable) TableName() string {
	return "my_table"
}
```

Two important things: we define `tags` on struct members to tell the orm how to interact with the db. We also define a method that returns
the table name for the struct as it is mapped in the db.

Tags are laid out in the following format `db:"<column_name>,<options>,..."`, where after the `column_name` there can be multiple `options`. An example above is `id,auto,pk`, which translates to a column name `id` and options `auto,pk` respectively.

Options include:
- `auto`,`serial` : denotes a column that will be read back on `Create` (there can be many of these).
- `pk` : deontes a column that consitutes a primary key. Will be used when creating SQL where clauses.
- `readonly` : denotes a column that is only read, not written to the db.

# Managing Connections and Aliases #

The next step in running a database driven app is to tell the app how to connect to the db.

We can create a connection with a configuration:
```golang
conn, err := db.New(db.OptConfig(cfg))
...
err = conn.Open()
...
```

The above snipped creates a connection, and opens it (establishing the connection). We can then pass this connection around to other things like controllers.

# ORM Actions: Create, Update, Delete, Get, GetAll

To create an object that has been mapped to a table, simply call:

```golang
obj := MyObj{...}
err := conn.Create(&obj) //note the reference! this is incase we have to write back a auto id.
```

Then we can get the object with:

```golang
var obj MyObj
found, err := conn.Get(&obj, "foo") // "foo" here is an imaginary primary key value.
```

To udpate an object:
```golang
var obj  MyObj
found, err := conn.Get(&obj, objID) //note, there can be multiple params (!!) if there are multiple pks
// .. handle the err
obj.Property = "new_value"
found, err = conn.Update(obj) //note we don't need a reference for this, as it's read only.
```

To delete an object:
```golang
var obj MyObj
err := conn.Get(&obj, objID) //note, there can be multiple params (!!) if there are multiple pks
// .. handle the err
err = conn.Delete(obj) //note we don't need a reference for this, as it's read only.
```

# Complex queries; using raw sql

To use sql directly, we need to use either an `Exec` (when we don't need to return results) or a `Query` (when we do want the results).

## Query

There are a couple options / paths we can take to actually running a query, and it's important to understand when to use each path.

- We need to run a query against the database without a transaction or a statement cache label:
```golang
conn.Query(<SQL Statement>, <Args...>).<Out(...)|OutMany(...)|Each(...)|First(...)|Scan(...)|Any()|None()>
```
- We need to run a query against the database with statement cache label:
```golang
conn.Invoke(db.OptCachedPlanKey("cached_statement")).Query(<SQL Statement>, <Args...>).<Out(...)|OutMany(...)|Each(...)|First(...)|Scan(...)|Any()|None()>
```
- We need to run a query against the database using a transaction, with a cache label:
```golang
conn.Invoke(db.OptTx(tx), db.OptCachedPlanKey("cached_statement")).Query(<SQL Statement>, <Args...>).<Out(...)|OutMany(...)|Each(...)|First(...)|Scan(...)|Any()|None()>
```

Note that after the `Query(...)` function itself, there can be various collectors. Each collector serves a different purpose:
- `Out(&obj)`: take the first result and automatically populate it against the object reference `obj` (must be passed by addr `&`)
- `OutMany(&objs)`: take all results and automatically populate them against the object array reference `objs` (must be passed by addr `&`)
- `Each(func(*sql.Rows) error)`: run the given handler for each result. This is useful if you need to read nested objects.
- `First(func(*sql.Rows) error)`: run the given handler for the first result. This is useful if you need to read a single complicated object.
- `Scan(<Args...>)`: read the first result into a given set of references. Useful for scalar return values.
- `Any`, `None`: return if there are results present, or conversely no results present.

## Exec

Executes have very similar preambles to queries:

- We need to execute a sql statement:
```golang
conn.Exec(<SQL Statement>, <Args...>)
```

The only difference is the lack of a collector, `Exec` will only return the `sql.Result` and the `error` for the statement execution.

# Common Patterns / Advanced Usage

## Nested objects

Lets say you have to model the following:

```golang
type Parent struct {
	ID int `db:"id,pk,serial"`
	TimestampUTC time.Time `db:"timestamp_utc"`
	Children []Child `db:"-"` // not we don't actually map this to the db.
}

type Child struct {
	ID int `db:"id,pk,serial"`
	ParentID int `db:"parent_id"`
	Name string `db:"name"`
}
```

What would the best way be to read all the `Parent` objects out with a given query?

We would want to query the parent objects, and while we're doing so, create a way to modify the parents as we read all the children.

To do this we use a map as a lookup, and some careful handling of pointers.

```golang
type Manager struct {
	DB *db.Connection
}

func (m Manager) GetAllParents() (parents []Parent, err error) {
	parentLookup := map[int]*Parent{} // note the pointer! this is so we can modify it.

	if err = m.DB.Query("select * from parent").Each(func(r *sql.Rows) error {
		var parent Parent
		// populate by name is a helper to set an object from a given row result
		// it is used internally by `OutMany` on `Query`.
		if err := db.PopulateByName(&parent, r, db.Columns(parent)); err != nil {
			return err
		}
		parents = append(parents, parent)
		parentLookup[parent.ID] = &parent
		return nil
	}); err != nil {
		return
	}

	// now we need to do a second query to get all the children.
	if err = m.DB.Query("select * from children").Each(func(r *sql.Rows) error {
		var child Child
		if err := db.PopulateByName(&child, r. db.Columns(child)); err != nil {
			return err
		}
		// here is the key part, we're looking up the parent to add the children.
		// because we're modifying references, the changes propagate to the original instances.
		if parent, hasParent := parentLookup[child.ParentID]; hasParent {
			parent.Children = append(parent.Children, child)
		}
		return nil
	}); err != nil {
		return
	}
	return
}

```

# Performance #

Generally it's pretty good. There is a comparison test in `spiffy_test.go` if you want to see for yourself. It creates 5000 objects with 5 properties each, then reads them out using the orm or manual scanning.

The results were:

| manual  |   orm    |
|---------|----------|
|17.11ms  | 38.08ms  |

This would be shabby in .net land, but for Go where we can't dynamically emit populate methods, and we're stuck doing in-order population or by name population, it's ok.

If you implement `Populatable`, performance improves dramatically.

| manual  |   orm (Populatable)    |
|---------|------------------------|
|14.33ms  | 16.95ms                |

The strategy then is to impelement populate on your "hot read" objects, and let the orm figure out the other ones.
