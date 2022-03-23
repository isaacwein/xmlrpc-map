package xmlrpc

import (
	"bytes"
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
)

var (
	//go:embed examples
	reqFile       embed.FS
	fileReq       = "examples/req.xml"
	fileResp      = "examples/resp.xml"
	fileArrayResp = "examples/resp_array.xml"
	fileFault     = "examples/fault.xml"
)

func TestReqDecoder(t *testing.T) {
	reqData := Request{}
	file, _ := reqFile.Open(fileReq)
	err := xml.NewDecoder(file).Decode(&reqData)
	if err != nil {
		t.Fatal(err)
		return
	}
	j, _ := json.Marshal(reqData)
	t.Logf("%s", j)
}

func TestReqEncoder(t *testing.T) {
	resData := Request{
		MethodName: "some-method-name",
		Data: &Value{
			Value: Struct{
				"i_account":  3,
				"i_account2": 36,
			},
		},
	}
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	enc.Indent("\t", "	")
	err := enc.Encode(&resData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(buf.String())
}

func TestRespDecoder(t *testing.T) {
	respData := Response{}

	file, err := reqFile.Open(fileResp)
	if err != nil {
		t.Fatal(err)
	}

	err = xml.NewDecoder(file).Decode(&respData)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf(ToJsonString(respData))

}

func TestRespEncoder(t *testing.T) {
	resData := Response{
		Data: &Value{
			Value: Struct{
				"i_account_2": 3,
				"i_account_f": 56,
				"i_array":     Array{"a", "b", "c"},
			},
		},
	}
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	enc.Indent("\t", "	")
	err := enc.Encode(&resData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(buf.String())
}

func TestRespArrayDecoder(t *testing.T) {
	respData := Response{}

	file, err := reqFile.Open(fileArrayResp)
	if err != nil {
		t.Fatal(err)
	}

	err = xml.NewDecoder(file).Decode(&respData)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf(ToJsonString(respData))

}

func TestRespFault(t *testing.T) {
	respData := Response{}

	file, err := reqFile.Open(fileFault)
	if err != nil {
		t.Fatal(err)
	}

	err = xml.NewDecoder(file).Decode(&respData)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf(ToJsonString(respData))

}

func ToJsonString(d interface{}) string {
	res, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("error marshal data: %#+v", d)
	}
	return string(res)
}
