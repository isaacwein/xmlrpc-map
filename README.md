# xmlrpc-map
this is a good package for proxies server that doesn't need to parse to struct just passing to json

###Getting started
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
resData := Request{
    MethodName: "some-method-name",
    Data: &Value{
		Value: Struct{
			"i_account":  3,
			"i_account2": 36,
		},
	},
}
// creating reader
buf := &bytes.Buffer{}

// Marshaling XML
err := xml.NewEncoder(buf).Encode(&resData)

if err != nil {
	// handle error
}
```