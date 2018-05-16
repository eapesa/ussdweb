package main

import(
  "fmt"
  "strings"

  "encoding/json"
)

var rawMenu []byte
var dsMenu Menu

type Menu map[string]interface{}
func LoadMenu() {
  // TODO: Hardcoded menu for now
  menu := []byte(`{
    "Menu-I": "Menu I",
    "Menu-II": {
      "Menu-I": "Menu II > A"
    },
    "Menu-III": {
      "Menu-I": {
        "Menu-I": "Menu III > B > i"
      }
    }
  }`)
  rawMenu = menu
  m := make(Menu)
  err := json.Unmarshal(menu, &m)
  if err != nil {
    fmt.Printf("Error loading menu\n")
  }
  fmt.Printf("Menu loaded successfully: %v\n", m)
  dsMenu = m
}

func Translate(key string) string {
  switch key {
  case "1": return "Menu-I"
  case "2": return "Menu-II"
  case "3": return "Menu-III"
  }
  return ""
}

func FindMenuValue(data Menu, keys []string) string {
  key := Translate(keys[0])
  value, ok := data[key].(string)
  if ok && len(keys) == 1 {
    // NOTE: May execute functions here...
    return "text:You traversed " + value
  }

  if !ok && len(keys) == 1 {
    subkeys := "menu:"
    for k := range data[key].(map[string]interface{}) {
      subkeys = subkeys + k + "|"
    }
    return subkeys[:len(subkeys) - 1]
  }

  if !ok && len(keys) > 1 {
    return FindMenuValue(data[key].(map[string]interface{}), keys[1:])
  }

  return INVALID_DIALSTRING
}

type DSPayload struct {
  Dialstring string
  Mode       string
  Rest       []string
}

func (ds *DSPayload) Parse() {
  syntax := strings.Split(ds.Dialstring, "*")
  ds.Mode = syntax[1]
  if ds.Rest = make([]string, 0); len(syntax) > 2 {
    ds.Rest = syntax[2:]
  }
}

func (ds *DSPayload) GenerateResponse() string {
  if ds.Mode == "0" {
    menu := GenerateMenu(dsMenu)
    result := "menu:" + menu[:len(menu) - 1]
    return result
  } else if ds.Mode == "1" {
    if len(ds.Rest) == 0 {
      subkeys := "menu:"
      for k := range dsMenu {
        subkeys = subkeys + k + "|"
      }
      subkeys = subkeys[:len(subkeys) - 1]
      return subkeys
    } else {
      return FindMenuValue(dsMenu, ds.Rest)
    }
  } else {
    return UNKNOWN_MODE
  }
}

func GenerateMenu(data Menu) string {
  menu := ""
  for key := range data {
    fmt.Printf("KEY IS %v || MENU IS %v\n", key, data)
    value, ok := data[key].(map[string]interface{})
    if ok {
      menu = menu + key + ">" + GenerateMenu(value)
    } else {
      menu = menu + key + "|"
    }
  }
  return menu
}
