package main

import(
  "net/http"
  "fmt"
)

func main() {
  initializeServer();
}

func initializeServer() {
  port := ":8000"
  // ussdcode := 118
  fmt.Printf("Initializing server at port%s\n", port)
  http.HandleFunc("/118/handler", ussdCodeHandler)
  http.ListenAndServe(port, nil)
}

func ussdCodeHandler(res http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" {
    fmt.Printf("ERROR: Invalid method")
    return
  }

  sender     := req.URL.Query().Get("sender")
  dialstring := req.URL.Query().Get("dialstring")
  fmt.Printf("SENDER=%v || DIALSTRING=%v\n", sender, dialstring)

  res.Write(OK);
}
