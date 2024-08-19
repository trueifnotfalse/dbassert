# dbassert

The `dbassert` package provides several useful functions to help you write integration tests.

# Example exists usage:

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/trueifnotfalse/dbassert"
)

func TestSomeDb(t *testing.T) {
	conn, err := sql.Open("postgres", "postgres://postgres:secret@localhost:%s?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
    _, err = conn.Exec(`INSERT INTO your_table(name) VALUES('John')`)
	dbassert := dbassert.New(conn)

    // check row exists in table
	dbassert.ExistsInDatabase(t, "your_table", map[string]any{
        "name": "John",
    }))

    // check row not exists in table
    dbassert.NotExistsInDatabase(t, "your_table", map[string]any{
        "name": "Boris",
    }))
}
```
