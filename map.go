package xmlrpc

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	XMLName    struct{} `xml:"methodCall" json:"-"`
	MethodName string   `xml:"methodName" json:"method_name"`
	Data       any      `xml:"params>param>value" json:"data,omitempty"`
}

func (r *Request) MarshalXML(u *xml.Encoder, start xml.StartElement) (err error) {
	type TempRequest struct {
		XMLName    struct{} `xml:"methodCall" json:"-"`
		MethodName string   `xml:"methodName" json:"method_name"`
		Data       *Value   `xml:"params>param>value" json:"data,omitempty"`
	}

	tempValue := &TempRequest{
		MethodName: r.MethodName,
		Data:       &Value{Value: r.Data},
	}
	return u.Encode(tempValue)
}
func (r *Request) UnmarshalXML(u *xml.Decoder, start xml.StartElement) (err error) {
	type TempRequest struct {
		XMLName    struct{} `xml:"methodCall" json:"-"`
		MethodName string   `xml:"methodName" json:"method_name"`
		Data       Value    `xml:"params>param>value" json:"data,omitempty"`
		Error      *Error   `xml:"fault>value" json:"error,omitempty"`
	}
	tempValue := &TempRequest{}
	err = u.DecodeElement(tempValue, &start)
	if err != nil {
		return err
	}
	r.MethodName = tempValue.MethodName
	if tempValue.Data.Value != nil {
		r.Data = tempValue.Data.Value
	}

	return
}

type Response struct {
	XMLName struct{} `xml:"methodResponse" json:"-"`
	Data    any      `xml:"params>param>value" json:"data,omitempty"`
	Error   *Error   `xml:"fault>value" json:"error,omitempty"`
}

func (r *Response) MarshalXML(u *xml.Encoder, start xml.StartElement) (err error) {
	type TempResponse struct {
		XMLName struct{} `xml:"methodResponse" json:"-"`
		Data    *Value   `xml:"params>param>value" json:"data,omitempty"`
		Error   *Error   `xml:"fault>value" json:"error,omitempty"`
	}

	tempValue := &TempResponse{}
	if r.Error != nil {
		tempValue.Error = r.Error
	}
	if r.Data != nil {
		tempValue.Data = &Value{Value: r.Data}
	}

	return u.Encode(tempValue)
}
func (r *Response) UnmarshalXML(u *xml.Decoder, start xml.StartElement) (err error) {
	type TempResponse struct {
		XMLName struct{} `xml:"methodResponse" json:"-"`
		Data    Value    `xml:"params>param>value" json:"data,omitempty"`
		Error   *Error   `xml:"fault>value" json:"error,omitempty"`
	}
	tempValue := &TempResponse{}
	err = u.DecodeElement(tempValue, &start)
	if err != nil {
		return err
	}

	if tempValue.Data.Value != nil {
		r.Data = tempValue.Data.Value
	}
	r.Error = tempValue.Error
	return
}

// Value xml-rpc value
type Value struct {
	Value any `json:"value,omitempty" xml:"value,omitempty"`
}

func (r *Value) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.Value)
}

func (r *Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Value)
}

func (r *Value) MarshalXML(u *xml.Encoder, start xml.StartElement) (err error) {

	type Temp struct {
		XMLName struct{} `xml:"value"`
		Value   struct {
			XMLName xml.Name
			Value   any `xml:",innerxml"`
		} `xml:",any"`
	}

	structType, structValue, err := MarshalType(r.Value)
	if err != nil {

		return fmt.Errorf("MarshalXML error: %w", err)
	}
	data := Temp{}
	data.Value.XMLName.Local = structType
	if structType != "nil" {
		data.Value.Value = structValue
	}

	err = u.EncodeElement(data, start)

	return
}
func (r *Value) UnmarshalXML(u *xml.Decoder, start xml.StartElement) (err error) {
	type Temp struct {
		XMLName struct{} `xml:"value"`
		Value   struct {
			XMLName xml.Name
			Value   string `xml:",innerxml"`
		} `xml:",any"`
	}

	var temp Temp
	err = u.DecodeElement(&temp, &start)
	if err != nil {
		return fmt.Errorf("value xml decoder error: %w", err)
	}

	r.Value, err = UnmarshalType(temp.Value.XMLName.Local, temp.Value.Value)
	return
}

// Struct is xml-rpc struct
type Struct map[string]any

func (r Struct) MarshalXML(u *xml.Encoder, start xml.StartElement) (err error) {
	type TempStructMember struct {
		Name  string `xml:"name"`
		Value Value  `xml:"value"`
	}
	var data []TempStructMember

	for key, value := range r {
		tm := TempStructMember{
			Name: key,
		}
		if value != nil {
			tm.Value.Value = value
		}

		data = append(data, tm)
	}

	start.Name.Local = "member"

	return u.EncodeElement(data, start)
}

func (r Struct) UnmarshalXML(u *xml.Decoder, start xml.StartElement) (err error) {
	type TempStructMember struct {
		Name  string `xml:"name"`
		Value Value  `xml:"value"`
	}

	type Temp struct {
		Members []TempStructMember `xml:"member"`
	}

	var temp Temp

	err = u.DecodeElement(&temp, &start)
	if err != nil {
		return fmt.Errorf("struct xml decoder error: %w", err)
	}
	if r == nil {
		err = fmt.Errorf("struct is nil")
		return
	}
	for _, m := range temp.Members {
		(r)[m.Name] = m.Value.Value
	}
	return nil
}

// Array xml-rpc array
type Array []any

func (r Array) MarshalXML(u *xml.Encoder, start xml.StartElement) (err error) {
	type Temp struct {
		XMLName xml.Name `xml:"data"`
		Values  []*Value `xml:"value"`
	}

	data := Temp{
		Values: make([]*Value, len(r)),
	}

	for i, v := range r {
		data.Values[i] = &Value{v}
	}

	return u.Encode(&data)

}
func (r *Array) UnmarshalXML(u *xml.Decoder, start xml.StartElement) (err error) {
	type Temp struct {
		XMLName xml.Name `xml:"array"`
		Values  []*Value `xml:"data>value"`
	}

	var temp Temp

	start.Name.Local = "array"
	err = u.DecodeElement(&temp, &start)

	if err != nil {
		return fmt.Errorf("array xml decoder error: %w", err)
	}

	for _, v := range temp.Values {
		*r = append(*r, v.Value)
	}
	return
}

var xmlRpcTime = "2006-01-02T15:04:05-0700"

func UnmarshalType(Type, data string) (res any, err error) {
	switch Type {
	case "int":
		res, err = strconv.ParseInt(data, 10, 64)
	case "double":
		res, err = strconv.ParseFloat(data, 64)
	case "boolean":
		res, err = strconv.ParseBool(data)
	case "dateTime.iso8601":
		res, err = time.Parse(data, xmlRpcTime)
	case "base64":
		res = []byte(data)
	case "string":
		res = data
	case "array":
		re := &Array{}
		data = fmt.Sprintf("<array>%s</array>", data)
		err = xml.NewDecoder(strings.NewReader(data)).Decode(re)
		if err != nil {
			err = fmt.Errorf("parsing array error: %w, %s ", err, data)
		}
		res = re
	case "nil":
		return nil, nil
	case "struct", "Struct":
		re := &Struct{}
		data = fmt.Sprintf("<string>%s</string>", data)
		err = xml.NewDecoder(strings.NewReader(data)).Decode(&re)
		if err != nil {
			err = fmt.Errorf("parsing struct error: %w - %s", err, data)
		}
		res = re
	default:
		fmt.Printf("uknown type: %s\n", Type)
	}
	return
}

func MarshalType(d any) (t string, v any, err error) {

	t = fmt.Sprintf("%T", d)
	v = fmt.Sprintf("%v", d)
	if isNil(d) {
		t = "nil"
		return
	}
	switch val := d.(type) {
	case int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64, uint,
		*uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64:
		t = "int"
	case float64, *float64, float32, *float32:
		t = "double"
	case bool, *bool:
		t = "boolean"
	case time.Time:
		t = "dateTime.iso8601"
		v = val.Format(xmlRpcTime)
	case *time.Time:
		t = "dateTime.iso8601"
		v = val.Format(xmlRpcTime)
	case []byte:
		t = "base64"
		v = base64.StdEncoding.EncodeToString(val)
	case string, *string:
		t = "string"
	case Array, *Array:
		t = "array"
		v = val
	case []any:
		t = "array"
		v = Array(val)
	case Struct, *Struct:
		t = "struct"
		v = val
	case map[string]any:
		t = "struct"
		v = Struct(val)
	case nil:
		t = "nil"
		v = val
	}

	return
}
func isNil(i any) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

type Error struct {
	FaultCode   int    `xml:"faultCode"`
	FaultString string `xml:"faultString"`
}

func (r *Error) Error() string {
	return fmt.Sprintf("faultCode: %d, faultString: %s", r.FaultCode, r.FaultString)
}

func (r *Error) UnmarshalXML(u *xml.Decoder, start xml.StartElement) (err error) {

	temp := Value{}
	err = u.DecodeElement(&temp, &start)
	if err != nil {
		return fmt.Errorf("decoding error-struct error: %w", err)
	}

	s, ok := temp.Value.(*Struct)
	if !ok {
		return fmt.Errorf("error is not a struct")
	}

	// finding faultCode
	faultCode, ok := (*s)["faultCode"]
	if !ok {
		return fmt.Errorf("error finding faultCode")
	}
	faultCodeI, ok := faultCode.(int64)
	if !ok {
		return fmt.Errorf("error converting faultCode to int64")
	}
	r.FaultCode = int(faultCodeI)

	// finding FaultString
	FaultString, ok := (*s)["faultString"]
	if !ok {
		return fmt.Errorf("error finding FaultString")
	}
	// converting FaultString to string
	r.FaultString, ok = FaultString.(string)
	if !ok {
		return fmt.Errorf("error converting FaultString")
	}
	return
}

func (r *Error) MarshalXML(u *xml.Encoder, start xml.StartElement) (err error) {
	temp := Value{&Struct{
		"faultCode":   r.FaultCode,
		"faultString": r.FaultString,
	}}
	err = u.EncodeElement(&temp, start)
	if err != nil {
		return fmt.Errorf("encoding error-struct error: %w", err)
	}
	return
}

type XmlRpcTypes interface {
	Value | Struct | Array | map[string]any | []any | *Value | *Struct | *Array
}

type Encoder struct {
	*xml.Encoder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{xml.NewEncoder(w)}
}
func (e *Encoder) Encode() error {
	return fmt.Errorf("not implemented use EncodeRequest or EncodeResponse")
}
func (e *Encoder) EncodeRequest(method string, data any) error {
	req := &Request{
		MethodName: method,
		Data:       data,
	}
	return e.Encoder.Encode(req)
}
func (e *Encoder) EncodeResponse(data any, err *Error) error {

	req := &Response{
		Data:  data,
		Error: err,
	}
	return e.Encoder.Encode(req)
}

type Decoder struct {
	*xml.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{xml.NewDecoder(r)}
}

func (d *Decoder) Decode() (err error) {
	err = fmt.Errorf("not implemented use DecodeRequest or DecodeResponse")
	return
}

func (d *Decoder) DecodeRequest() (res *Request, err error) {
	res = &Request{}
	err = d.Decoder.Decode(res)
	return
}

func (d *Decoder) DecodeResponse() (res *Response, err error) {
	res = &Response{}
	err = d.Decoder.Decode(res)
	return
}
