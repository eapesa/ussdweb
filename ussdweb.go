package main

import(
  "net/http"
  "net/url"
  "strings"
  "encoding/json"
  "fmt"
)

type Pubkey struct {
  Pubkey  string  `json:"pubkey,omitempty"`
  UserId  string  `json:"user_id,omitempty"`
  Code    string  `json:"code,omitempty"`
  Message string  `json:"message,omitempty"`
}

var db *UssDB
func main() {
  var err error
  db, err = initializeDb()
  if err != nil {
    fmt.Printf("ERROR in initializeDb: %v\n", err)
  }

  defer db.DbObj.Close()
  defer db.PkInsertQuery.Close()
  
  initializeServer()
}

func initializeServer() {
  port := ":8000"
  fmt.Printf("Initializing server at port%s\n", port)
  http.HandleFunc("/api/pubkey", pubkeyHandler)
  http.HandleFunc("/", ussdCodeHandler)
  http.ListenAndServe(port, nil)
}

func ussdCodeHandler(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    fmt.Printf("ERROR: Invalid method\n")
    return
  }

  sender := req.URL.Query().Get("sender")
  ds     := req.URL.Query().Get("dialstring")

  process_payload(sender, ds)
  res.Header().Set("Content-Type", "application/xml")
  res.Write(prompt_reply(GETBALANCE + "P10000."));
}

func process_payload(sender, ds string) {
  if sender == "" {
    fmt.Printf("ERROR: Invalid inputs (`sender`)\n")
    return
  }

  if ds  == "" {
    fmt.Printf("ERROR: Invalid inputs (`dialstring`)\n")
    return
  }

  dsdecoded, err := url.QueryUnescape(ds)
  if err != nil {
    fmt.Printf("ERROR: Dialstring can't be decoded\n")
  }

  fmt.Printf("TODO: Process dialstring %s\n", dsdecoded)
}

func parse(dialstring string) []string {
  dsdecoded, err := url.QueryUnescape(dialstring)
  if err != nil {
    fmt.Printf("ERROR: Encountered error decoding query string\n")
  }

  return strings.Split(dsdecoded, "*#")
}

func pubkeyHandler(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    fmt.Printf("ERROR: Invalid method\n")
    return
  }

  input := json.NewDecoder(req.Body)
  defer req.Body.Close()

  var pubkey Pubkey
  err := input.Decode(&pubkey)
  if err != nil {
    fmt.Printf("ERROR: Error decoding post body\n")
    return
  }

  insertPubkey(db.PkInsertQuery, pubkey.UserId, pubkey.Pubkey)

  output, err2 := json.Marshal(&Pubkey{
    Pubkey: "some-pubkey",
    Code: "OK",
    Message: "Operation completed successfully",
  })
  if err2 != nil {
    fmt.Printf("ERROR: Error encoding response body\n")
    return
  }

  res.Header().Set("Content-Type", "application/json")
  res.Write(output)
}

/*
  NOTES:
  Use goroutine to start the server and create a channel for checking if there's a signal for exit
 */
