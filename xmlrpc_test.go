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
	reqFile          embed.FS
	fileReq          = "examples/req.xml"
	fileResp         = "examples/resp.xml"
	fileArrayResp    = "examples/resp_array.xml"
	fileFault        = "examples/fault.xml"
	fileMultiCallReq = "examples/multicall_req.xml"
	fileMultiCallRes = "examples/multicall_res.xml"
)

func TestReqDecoder(t *testing.T) {

	file, _ := reqFile.Open(fileReq)
	// reqData is a pointer to Request struct
	reqData, err := NewDecoder(file).DecodeRequest()
	if err != nil {
		t.Fatal(err)
		return
	}
	j, _ := json.Marshal(reqData)
	t.Logf("%s", j)
}

func TestReqEncoder(t *testing.T) {
	reqData := Struct{
		"i_account":  3,
		"i_account2": 36,
		"nil-value":  nil,
		"i_struct": &Struct{
			"i_account":  3,
			"i_struct_2": (*Struct)(nil),
		},
		"i_array": Array{"a", "b", "c", nil},
	}
	buf := &bytes.Buffer{}

	enc := NewEncoder(buf)
	enc.Indent("\t", "	")
	err := enc.EncodeRequest("some-method-name", "some-string-param", 555, &reqData, &reqData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(buf.String())
}

func TestRespDecoder(t *testing.T) {

	file, err := reqFile.Open(fileResp)
	if err != nil {
		t.Fatal(err)
	}
	respData, err := NewDecoder(file).DecodeResponse()

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf(ToJsonString(respData))

}

func TestRespEncoder(t *testing.T) {
	resData := map[string]any{
		"i_account_2": 3,
		"i_account_f": 56,
		"i_array":     []any{"a", "b", "c", nil},
		"nil-value":   nil,
	}

	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	//enc.Indent("\t", "	")
	err := enc.EncodeResponse(resData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(buf.String())
}

func TestRespErrorEncoder(t *testing.T) {
	resData := Error{
		FaultCode:   1,
		FaultString: "some error",
	}

	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	enc.Indent("\t", "	")
	err := enc.EncodeError(&resData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(buf.String())
}

func TestRespArrayDecoder(t *testing.T) {

	file, err := reqFile.Open(fileArrayResp)
	if err != nil {
		t.Fatal(err)
	}
	respData, err := NewDecoder(file).DecodeResponse()

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf(ToJsonString(respData))

}

func TestRespFault(t *testing.T) {

	file, err := reqFile.Open(fileFault)
	if err != nil {
		t.Fatal(err)
	}

	respData, err := NewDecoder(file).DecodeResponse()
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf(ToJsonString(respData))

}

func TestMultiCallReqDecode(t *testing.T) {
	reqData := Request{}
	file, _ := reqFile.Open(fileMultiCallReq)
	err := xml.NewDecoder(file).Decode(&reqData)
	if err != nil {
		t.Fatal(err)
		return
	}
	j, _ := json.Marshal(reqData)
	t.Logf("%s", j)

}

func TestMultiCallReqEncode(t *testing.T) {
	resData := Request{
		MethodName: "system.multicall",
		Data: Array{
			&Struct{
				"methodName": "wp.getUsersBlogs1",
				"params":     Array{"{{ Your Username1 }}", "{{ Your Password1 }}"},
			},
			Struct{
				"methodName": "wp.getUsersBlogs2",
				"params":     Array{Array{"{{ Your Username2 }}", "{{ Your Password2 }}"}},
			},
			Struct{
				"methodName": "wp.getUsersBlogs3",
				"params":     Array{Array{"{{ Your Username3 }}", "{{ Your Password3 }}"}},
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
func TestMultiCallResDecode(t *testing.T) {

	reqData := Response{}
	file, _ := reqFile.Open(fileMultiCallRes)
	err := xml.NewDecoder(file).Decode(&reqData)
	if err != nil {
		t.Fatal(err)
		return
	}
	j, _ := json.Marshal(reqData)
	t.Logf("%s", j)
}

func ToJsonString(d any) string {
	res, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("error marshal data: %#+v", d)
	}
	return string(res)
}
