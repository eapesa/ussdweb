package main

import (
  "encoding/xml"
)

type VxmlPrompt struct {
  XMLName   xml.Name    `xml:"vxml"`
  Form      Form        `xml:"form"`
}

type Form struct {
  Id        string    `xml:"id,attr"`
  Name      string    `xml:"name,attr"`
  Field     Field     `xml:"field"`
  Filled    Filled    `xml:"filled"`
  Catch     Catch     `xml:"catch"`
  Property  Property  `xml:"property"`
}

type Field struct {
  FieldName  string  `xml:"name,attr"`
  Prompt     string  `xml:"prompt>field"`
}

type Filled struct {
  Assign    Assign  `xml:"assign"`
  Goto      Goto    `xml:"goto"`
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
  Prompt  string  `xml:"prompt>catch"`
  Goto    Goto    `xml:"goto"`
}

type Property struct {
  PropertyName  string  `xml:"name,attr"`
  PropertyValue string  `xml:"value,attr"`
}

func prompt_reply() []byte {
  vxml := &VxmlPrompt{
    Form: Form{
      Id: "Output",
      Name: "Output",
      Field: Field{
        FieldName: "oc_Output",
        Prompt: "Your transaction is being processed.",
      },
      Filled: Filled{
        Assign: Assign{
          AssignName: "",
          Expr: "",
        },
        Goto: Goto{
          Next: "",
        },
      },
      Catch: Catch{
        Event: "nomatch",
        Prompt: "Invalid choice. Try again.",
        Goto: Goto{
          Next: "#Output",
        },
      },
      Property: Property{
        PropertyName: "oc_bIsFinal",
        PropertyValue: "1",
      },
    },
  }

  output, _ := xml.MarshalIndent(vxml, "  ", "    ")
  output = []byte(xml.Header + string(output))
  return output
}
