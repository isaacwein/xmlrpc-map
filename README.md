# xmlrpc-map

## its parsing xmlrpc to map[string]any like the default json package

this is a good package for proxies server that doesn't need to parse to struct just passing to json
it can be used

- `Struct` is a `map[string]any` that can be used to marshal and unmarshal xmlrpc data
- `Array` is a `[]any` that can be used to marshal and unmarshal xmlrpc data
- `Value` is a `any` that can be used to marshal and unmarshal xmlrpc data

if you want to use it with custom types you can marshel into a json and then unmarshal into your struct


### Getting started
```go get github/telebroad/xmlrpc-map```

### Demos

#### Unmarshal

```go
respData := Response{}

// resp.Body is the io.Reader from an api 
err = xml.NewDecoder(resp.Body).Decode(&respData)
if err != nil {
    // handle error
}

```

#### Marshal
```go
reqData := &Request{
	MethodName: "some-method-name",
	Data: Struct{
        "i_account":  3,
        "i_account2": 36,
        "nil-value":  nil,
        "i_array":    Array{"a", "b", "c", nil},
    },
}
// creating reader
buf := &bytes.Buffer{}

// Marshaling XML
err := xml.NewEncoder(buf).Encode(&reqData)

if err != nil {
	// handle error
}
```

encodiing Response example

```go
resData := Response{
	Data: map[string]any{
		"i_account_2": 3,
		"i_account_f": 56,
		"i_array":     []any{"a", "b", "c", nil},
		"nil-value":   nil,
	},
}
buf := &bytes.Buffer{}
err := xml.NewEncoder(buf).Encode(&resData)
if err != nil {
// handle error
}
```

output
```xml
<methodCall><params><param><value><struct><member><name>i_account_2</name><value><int>3</int></value></member><member><name>i_account_f</name><value><int>56</int></value></member><member><name>i_array</name><value><array><data><value><string>a</string></value><value><string>b</string></value><value><string>c</string></value><value><nil></nil></value></data></array></value></member><member><name>nil-value</name><value><nil></nil></value></member></struct></value></param></params></methodCall>
```