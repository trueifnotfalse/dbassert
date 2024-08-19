package dbassert

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"
    "testing"
)

// DBAssert provides database assertion methods around the testing.TB interface.
type DBAssert struct {
    conn *sql.DB
}

// New creates a new DBAssert
func New(conn *sql.DB) *DBAssert {
    return &DBAssert{
        conn: conn,
    }
}

// ExistsInDatabase check if row exists in table
func (r *DBAssert) ExistsInDatabase(t testing.TB, tableName string, data map[string]any) bool {
    exists, err := r.existsInDatabase(tableName, data)
    if err != nil {
        t.Errorf(err.Error())
        t.Fail()
        return false
    }
    if !exists {
        d, err := json.Marshal(data)
        if err != nil {
            t.Errorf(err.Error())
            t.Fail()
            return false
        }
        t.Errorf("unable to find row in database table [%s] that matched attributes [%s]", tableName, d)
        t.Fail()
        return false
    }

    return true
}

// NotExistsInDatabase check if row exists in table
func (r *DBAssert) NotExistsInDatabase(t testing.TB, tableName string, data map[string]any) bool {
    exists, err := r.existsInDatabase(tableName, data)
    if err != nil {
        t.Errorf(err.Error())
        t.Fail()
        return false
    }
    if exists {
        d, err := json.Marshal(data)
        if err != nil {
            t.Errorf(err.Error())
            t.Fail()
            return false
        }
        t.Errorf("found unexpected records in database table [%s] that matched attributes [%s]", tableName, d)
        t.Fail()
        return false
    }

    return true
}

func (r *DBAssert) existsInDatabase(tableName string, data map[string]any) (bool, error) {
    if len(data) == 0 {
        return false, fmt.Errorf("no attributes to find")
    }
    var (
        fields        []string
        valueAsString string
        exists        bool
    )
    for column, value := range data {
        switch value.(type) {
        case string:
            valueAsString = "'" + value.(string) + "'"
        case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
            valueAsString = fmt.Sprintf("%d", value)
        case float32:
            valueAsString = strconv.FormatFloat(float64(value.(float32)), 'f', -1, 64)
        case float64:
            valueAsString = strconv.FormatFloat(value.(float64), 'f', -1, 64)
        }
        fields = append(fields, fmt.Sprintf("%s=%s", column, valueAsString))
    }
    s := fmt.Sprintf(`select exists (select * from %s where %s);`, tableName, strings.Join(fields, " and "))
    row := r.conn.QueryRow(s)
    err := row.Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("database error [%v]", err)
    }

    return exists, nil
}
