package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// XML-RPC response success body example
var sucResp = `
<?xml version='1.0'?>
<methodResponse>
    <params>
        <param>
            <value>
                <struct>
                    <member>
                        <name>name</name>
                        <value>
                            <string>jack</string>
                        </value>
                    </member>
                </struct>
            </value>
        </param>
    </params>
</methodResponse>
`

// XML-RPC response error body example
var errorResp = `
<?xml version='1.0'?>
<methodResponse>
    <fault>
        <value>
            <struct>
                <member>
                    <name>faultCode</name>
                    <value>
                        <int>400</int>
                    </value>
                </member>
                <member>
                    <name>faultString</name>
                    <value>
                        <string>Account not found</string>
                    </value>
                </member>
            </struct>
        </value>
    </fault>
</methodResponse>
`

func main() {
	type SUC struct {
		Name string `xmlrpc:"name"`
	}

	fmt.Println("---- parsing on success body ----")
	err := XmlRpcDecoder(strings.NewReader(sucResp), &SUC{})
	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("---- parsing on error body ----")
	err = XmlRpcDecoder(strings.NewReader(errorResp), &SUC{})
	if err != nil {
		fmt.Println("err", err)
	}
}

// XmlRpcErr XML-RPC response error
type XmlRpcErr struct {
	FaultCode   int    `xmlrpc:"faultCode"`
	FaultString string `xmlrpc:"faultString"`
}

func (e *XmlRpcErr) Error() string {
	return fmt.Sprintf("xmlrpc error: %d - %s", e.FaultCode, e.FaultString)
}

// respData XML-RPC response data that holds both success and error data
type respData struct {
	Suc any        `xmlrpc:"suc"`
	Err *XmlRpcErr `xmlrpc:"err"`
}

// UnmarshalXML checking if the xml-rpc body is a success or error
func (r *respData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// decode inner elements

	for {
		t, err := d.Token()
		if err != nil {
			return fmt.Errorf("getting first element error: %w", err)
		}
		//var i any
		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "params":
				fmt.Printf("you decoding success data on r.Suc %T\n", r.Suc)
				return d.DecodeElement(r.Suc, &tt)
			case "fault":
				fmt.Printf("you decoding error data on r.Err %T\n", r.Err)
				return d.DecodeElement(r.Err, &tt)
			}

		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}

	}
}

// XmlRpcDecoder XML-RPC response decoder
func XmlRpcDecoder(body io.Reader, sucData any) error {

	resp := &respData{Suc: sucData}

	err := xml.NewDecoder(body).Decode(resp)
	if err != nil {
		return err
	}
	if resp.Err != nil {
		return resp.Err
	}
	return nil
}
