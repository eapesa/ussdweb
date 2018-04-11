package main

import(
  "net/http"
  "net/url"
  "strings"
  "fmt"
)

func main() {
  initializeServer();
}

func initializeServer() {
  port := ":8000"
  fmt.Printf("Initializing server at port%s\n", port)
  http.HandleFunc("/", ussdCodeHandler)
  http.ListenAndServe(port, nil)
}

func ussdCodeHandler(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    fmt.Printf("ERROR: Invalid method")
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
