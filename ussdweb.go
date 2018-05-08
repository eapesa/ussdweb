package main

import (
  "net/http"
  "fmt"
)

type Pubkey struct {
  Pubkey  string  `json:"pubkey,omitempty"`
  UserId  string  `json:"user_id,omitempty"`
  Code    string  `json:"code,omitempty"`
  Message string  `json:"message,omitempty"`
}

var db *UssDB
var serverPubkey string
func main() {
  var err error
  db, err = initializeDb()
  if err != nil {
    fmt.Printf("ERROR in initializeDb: %v\n", err)
    return
  }
  defer db.PkInsertQuery.Close()
  defer db.DbObj.Close()

  // keyByte, err := loadKeys()
  // if err != nil {
  //   fmt.Printf("Error in loadKeys: %v\n", err)
  //   return
  // }
  // serverPubkey = string(keyByte)

  // Place this in environment variables
  port := ":8000"
  fmt.Printf("Initializing server at port%s\n", port)
  http.HandleFunc("/api/pubkey", pubkeyHandler)
  http.HandleFunc("/", ussdCodeHandler)
  http.ListenAndServe(port, nil)
}
