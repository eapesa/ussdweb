package main

import(
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
)

const (
  INSERT_PUBKEY_QUERY = `INSERT INTO users (user_id,pubkey,expires_at)
                          VALUES(?, ?, DATE_ADD(NOW(), INTERVAL 1 WEEK))
                          ON DUPLICATE KEY UPDATE
                            pubkey=?, expires_at=DATE_ADD(NOW(), INTERVAL 1 WEEK)`
)

type UssDB struct {
  DbObj *sql.DB
  PkInsertQuery *sql.Stmt
}

func initializeDb() (*UssDB, error) {
  var db *sql.DB
  var stmt *sql.Stmt
  var err error

  db, err = sql.Open("mysql", "SOMEUSER:SOMEPASS@/SOMEDB")
  if err != nil {
    fmt.Printf("ERROR in initializeDb: %v\n", err)
    return &UssDB{}, err
  }
  fmt.Printf("MySQL starting... OK\n")

  stmt, err = db.Prepare(INSERT_PUBKEY_QUERY)
  if err != nil {
    fmt.Printf("ERROR in initializeDb: %v\n", err)
    return &UssDB{}, err
  }

  return &UssDB{ DbObj: db, PkInsertQuery: stmt }, nil
}

func insertPubkey(stmt *sql.Stmt, user_id, pubkey string) {
  _, err := stmt.Exec(user_id, pubkey, pubkey)
  if err != nil {
    fmt.Printf("ERROR in insertPubkey: %v\n", err)
    return
  }
  fmt.Printf("Insert/Update pubkey... OK\n")
  return
}

/*
  NOTES:
  - initializer return the object and have the main check if there's no db object, do stuff
  - create maybe a struct for db objects and its prepared queries
 */
