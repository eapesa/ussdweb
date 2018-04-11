package main

import (
  "encoding/xml"
)

const (
  HELLOWORLD = "Hello World!",
  GETBALANCE = "Your current balance is "
)

type VxmlPrompt struct {
  XMLName   xml.Name    `xml:"vxml"`
  Form      *Form       `xml:"form"`
}

type Form struct {
  Id        string    `xml:"id,attr"`
  Name      string    `xml:"name,attr"`
  Block     string    `xml:"block,omitempty"`
  Field     *Field    `xml:"field,omitempty"`
  Filled    *Filled   `xml:"filled,omitempty"`
  Catch     *Catch    `xml:"catch,omitempty"`
  Property  *Property `xml:"property,omitempty"`
}

type Value struct {
  Expr  string  `xml:"expr,attr,omitempty"`
}

type Field struct {
  FieldName  string  `xml:"name,attr"`
  Prompt     string  `xml:"prompt"`
}

type Filled struct {
  Assign    *Assign `xml:"assign"`
  Goto      *Goto   `xml:"goto"`
}

type Assign struct {
  AssignName  string  `xml:"name,attr"`
  Expr        string  `xml:"expr,attr"`
}

type Goto struct {
  Next    string  `xml:"next,attr"`
}

type Catch struct {
  Event   string  `xml:"event,attr"`
  Prompt  string  `xml:"prompt"`
  Goto    *Goto   `xml:"goto"`
}

type Property struct {
  PropertyName  string  `xml:"name,attr"`
  PropertyValue string  `xml:"value,attr"`
}

func prompt_reply(textPrompt string) []byte {
  vxml := &VxmlPrompt{
    Form: &Form{
      Id: "Output",
      Name: "Output",
      Field: &Field{
        FieldName: "oc_Output",
        Prompt: textPrompt,
      },
      Filled: &Filled{
        Assign: &Assign{
          AssignName: "",
          Expr: "oc_Output",
        },
        Goto: &Goto{
          Next: "",
        },
      },
      Catch: &Catch{
        Event: "nomatch",
        Prompt: "Invalid input",
        Goto: &Goto{
          Next: "#Output",
        },
      },
      Property: &Property{
        PropertyName: "oc_bIsFinal",
        PropertyValue: "1",
      },
    },
  }

  output, _ := xml.MarshalIndent(vxml, "  ", "    ")
  output = []byte(xml.Header + string(output))
  return output
}
