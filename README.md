USSD Web Server
---------------

USSD Web Server is a simple web server that can interact with GSM cellular
telephones and mobile applications. It aims to see the capacity of [USSD](https://en.wikipedia.org/wiki/Unstructured_Supplementary_Service_Data)
to host data exchange between mobile applications that only has GSM signal.

USSD Web Server is the backend counterpart of the project and works in partner with
[USSD POC (Android)](https://github.com/eapesa/ussdpoc-android) which connects to
the provisioned USSD code for this POC.


## Dependencies

*NOTE: To receive data sent to your provisioned USSD code, ensure that the server
       that will host this application can exchange data between USSD gateway
       specified for the project.*

- Go (assured to run in `go1.10`)
  - Properly configured `GOROOT` and `GOPATH` *to be automated*
- MySQL


## Setup

- Clone repo `$> git clone https://github.com/eapesa/ussdweb.git`

- Configure `GOROOT` and `GOPATH` in your `.bash_profile` or `.bashrc`

*NOTE: `GOROOT` is where your golang libraries are installed (e.g.
       /usr/local/opt/go/libexec for Mac OSX users). `GOPATH` is where your
       development directory is located.*

- Execute `$> make install`

- Run application through `$> make run`

*NOTE: To run the application in background, you may use `screen`*
