package main

import (
  "net/http"
  "net/url"

  "encoding/json"

  "strings"
  "fmt"
)

func ussdCodeHandler(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    fmt.Printf("ERROR: Invalid method\n")
    return
  }

  sender := req.URL.Query().Get("sender")
  ds     := req.URL.Query().Get("dialstring")

  process_payload(sender, ds)
  res.Header().Set("Content-Type", "application/xml")
  // res.Write(sendCustomMsg(multiplyString("+123456789", 15)));
  // res.Write(sendInfoMsg(TESTMESSAGE))
  res.Write(sendCustomMsg(CustomMessage{
      Sender: "639473371382",
      Dialstring: "07",
      Recipient: "639473371382",
      TextPrompt: "Test Custom Message",
    }))
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
    Pubkey: serverPubkey,
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

/******************************************************************************
 * HELPER FUNCTIONS
 ******************************************************************************/
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

// For easier testing of maximum character count.
func multiplyString(text string, count int) string {
  textAll := ""
  for i := 0; i < count; i++ {
    textAll += text
  }
  return textAll
}
