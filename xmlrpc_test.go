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
		Data: Struct{
			"i_account":  3,
			"i_account2": 36,
			"nil-value":  nil,
			"i_array":    Array{"a", "b", "c", nil},
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
		Data: map[string]any{
			"i_account_2": 3,
			"i_account_f": 56,
			"i_array":     []any{"a", "b", "c", nil},
			"nil-value":   nil,
		},
	}
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	//enc.Indent("\t", "	")
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
		Data: &Value{
			Value: Array{
				Struct{
					"methodName": "wp.getUsersBlogs1",
					"params":     Array{"{{ Your Username }}", "{{ Your Password }}"},
				},
				Struct{
					"methodName": "wp.getUsersBlogs2",
					"params":     Array{Array{"{{ Your Username }}", "{{ Your Password }}"}},
				},
				Struct{
					"methodName": "wp.getUsersBlogs3",
					"params":     Array{Array{"{{ Your Username }}", "{{ Your Password }}"}},
				},
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

	reqData := Request{}
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
