package dbassert

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

const (
    tableName = "dbassert_test"
)

func startUp() (*sql.DB, error) {
    host := os.Getenv("POSTGRESQL_HOST")
    port := os.Getenv("POSTGRESQL_PORT")
    user := os.Getenv("POSTGRESQL_USER")
    password := os.Getenv("POSTGRESQL_PASSWORD")
    database := os.Getenv("POSTGRESQL_DATABASE")

    conn, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database))
    if err != nil {
        return nil, err
    }
    err = shutDown(conn)
    if err != nil {
        return nil, err
    }
    _, err = conn.Exec(fmt.Sprintf(`
CREATE TABLE %s
(
    id    serial,
    name  character varying(255) not null
);`, tableName))
    if err != nil {
        return nil, err
    }

    return conn, err
}

func shutDown(conn *sql.DB) error {
    _, err := conn.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName))
    return err
}

func TestPublicExists(t *testing.T) {
    conn, err := startUp()
    if err != nil {
        t.Error(err.Error())
        t.Fail()
        return
    }
    defer shutDown(conn)
    defer conn.Close()
    s := "random string"
    _, err = conn.Exec(fmt.Sprintf(`INSERT INTO %s(name) VALUES('%s')`, tableName, s))
    assert.Nil(t, err)
    dbassert := New(conn)
    result := dbassert.ExistsInDatabase(t, tableName, map[string]any{
        "name": s,
    })
    assert.True(t, result)
}

func TestPublicNotExists(t *testing.T) {
    conn, err := startUp()
    if err != nil {
        t.Error(err.Error())
        t.Fail()
        return
    }
    defer shutDown(conn)
    defer conn.Close()
    s := "random string"
    _, err = conn.Exec(fmt.Sprintf(`INSERT INTO %s(name) VALUES('%s')`, tableName, s))
    assert.Nil(t, err)
    dbassert := New(conn)
    result := dbassert.NotExistsInDatabase(t, tableName, map[string]any{
        "name": "some string",
    })
    assert.True(t, result)
}

func TestPrivateExists(t *testing.T) {
    conn, err := startUp()
    if err != nil {
        t.Error(err.Error())
        t.Fail()
        return
    }
    defer shutDown(conn)
    defer conn.Close()
    s := "random string"
    _, err = conn.Exec(fmt.Sprintf(`INSERT INTO %s(name) VALUES('%s')`, tableName, s))
    assert.Nil(t, err)
    dbassert := New(conn)
    result, err := dbassert.existsInDatabase(tableName, map[string]any{
        "name": s,
    })
    assert.Nil(t, err)
    assert.True(t, result)
}

func TestPrivateNotExists(t *testing.T) {
    conn, err := startUp()
    if err != nil {
        t.Error(err.Error())
        t.Fail()
        return
    }
    defer shutDown(conn)
    defer conn.Close()
    s := "random string"
    _, err = conn.Exec(fmt.Sprintf(`INSERT INTO %s(name) VALUES('%s')`, tableName, s))
    assert.Nil(t, err)
    dbassert := New(conn)
    result, err := dbassert.existsInDatabase(tableName, map[string]any{
        "name": "some string",
    })
    assert.Nil(t, err)
    assert.False(t, result)
}

func TestPrivateEmptyAttributes(t *testing.T) {
    conn, err := startUp()
    if err != nil {
        t.Error(err.Error())
        t.Fail()
        return
    }
    defer shutDown(conn)
    defer conn.Close()
    s := "random string"
    _, err = conn.Exec(fmt.Sprintf(`INSERT INTO %s(name) VALUES('%s')`, tableName, s))
    assert.Nil(t, err)
    dbassert := New(conn)
    result, err := dbassert.existsInDatabase(tableName, map[string]any{})
    assert.NotNil(t, err)
    assert.False(t, result)
}
