package main

import (
	"encoding/json"
	"strings"
)

type Something[T any] struct {
	ReturnType T
}

type customType struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type NoReturn struct{}

func main() {
	var s Something[customType]
	_ = lol(&s, strings.NewReader(`{"name":"gamet","age":17}`))

	var n Something[NoReturn]
	_ = lol(&n, strings.NewReader(`null`))

	println("dwqdqwqwdwqd")
}

func lol[T any](dst *Something[T], r *strings.Reader) error {
	var zero T
	if _, ok := any(zero).(NoReturn); ok {
		return nil
	}
	return json.NewDecoder(r).Decode(&dst.ReturnType)
}
