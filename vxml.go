package main

import (
  "encoding/xml"
)

const (
  HELLOWORLD = "Hello World!"
  GETBALANCE = "Your current balance is "
  TESTMESSAGE = "DS07 Triggered. This reply is an Info Message."
  INVALID_MESSAGE = "You have entered an invalid message."
  INVALID_DIALSTRING = "Invalid dialstring. No mapping found."
  CHOICE_I = "You clicked for Choice I"
  CHOICE_II_A = "You clicked for Choice II > A"
  CHOICE_II_B_i = "You clicked for Choice II > B > i"
  CHOICE_II_B_ii = "You clicked for Choice II > B > ii"
  CHOICE_III_A = "You clicked for Choice III > A"
)

type CustomMessage struct {
  Sender     string
  Recipient  string
  Dialstring string
  TextPrompt string
}

type VxmlPrompt struct {
  XMLName   xml.Name    `xml:"vxml"`
  Form      Forms       `xml:"form"`
}

type Form struct {
  Id        string    `xml:"id,attr"`
  Name      string    `xml:"name,attr"`
  Field     *Field    `xml:"field,omitempty"`
  Filled    *Filled   `xml:"filled,omitempty"`
  Catch     *Catch    `xml:"catch,omitempty"`
  Property  *Property `xml:"property,omitempty"`
  Var       *Var      `xml:"var,omitempty"`
  Block     Blocks    `xml:"block,omitempty"`
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
  AssignExpr  string  `xml:"expr,attr"`
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

type Block struct {
  BlockName  string  `xml:"name,attr"`
  Goto       *Goto   `xml:"goto"`
}

type Var struct {
  VarName  string  `xml:"name,attr"`
  VarExpr  string  `xml:"expr,attr"`
}

type Forms  []*Form
type Blocks []*Block

func sendInfoMsg(textPrompt string) []byte {
  vxml := &VxmlPrompt{
    Form: Forms{
      &Form{
        Id: "Output",
        Name: "Output",
        Field: &Field{
          FieldName: "oc_Output",
          Prompt: textPrompt,
        },
        Filled: &Filled{
          Assign: &Assign{
            AssignName: "",
            AssignExpr: "oc_Output",
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
    },
  }
  return send(vxml)
}

func sendCustomMsg(params CustomMessage) []byte {
  gotoNext := "http_client://HTTP_CHIKKA_MAGIC/118/handler.php?sender=" +
      params.Sender + "&dialstring=" + params.Dialstring +
      "&send_type=custom&recipient=" + params.Recipient + "&message=%var_custom_msg%"
  vxml := &VxmlPrompt{
    Form: Forms{
      &Form{
        Id: "form_custom_msg",
        Name: "form_custom_msg",
        Field: &Field{
          FieldName: "field_custom_msg",
          Prompt: params.TextPrompt,
        },
        Filled: &Filled{
          Assign: &Assign{
            AssignName: "var_custom_msg",
            AssignExpr: "field_custom_msg",
          },
          Goto: &Goto{
            Next: "#form_send_custom_msg",
          },
        },
        Catch: &Catch{
          Event: "nomatch",
          Prompt: "Invalid input",
          Goto: &Goto{
            Next: "#form_custom_msg",
          },
        },
      },
      &Form{
        Id: "form_send_custom_msg",
        Name: "form_send_custom_msg",
        Block: Blocks{
          &Block{
            BlockName: "oc_ActionUrl",
            Goto: &Goto{
              Next: gotoNext,
            },
          },
          &Block{
            BlockName: "oc_NextNodeUrl",
            Goto: &Goto{
              Next: "#End",
            },
          },
        },
      },
      &Form{
        Id: "End",
        Name: "End",
        Var: &Var{
          VarName: "strfinal",
          VarExpr: "#!%ResponseText%",
        },
        Block: Blocks{
          &Block{
            BlockName: "oc_NextNodeUrl",
            Goto: &Goto{
              Next: "var://strfinal",
            },
          },
        },
      },
    },
  }
  return send(vxml)
}

func send(vxml *VxmlPrompt) []byte {
  output, _ := xml.MarshalIndent(vxml, "  ", "    ")
  output = []byte(xml.Header + string(output))
  return output
}
