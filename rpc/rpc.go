package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

// IDArg is the most useful argument type, so define it here
type IDArg struct {
	ID uint64 `json:"id,string"`
}

// Result is the API result type
type Result struct {
	Succeeded bool        `json:"succeeded"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// NewResult creates a new API result
func NewResult(succeeded bool, message string, data ...interface{}) *Result {
	r := &Result{
		Succeeded: succeeded,
		Message:   message,
	}

	if l := len(data); l > 1 {
		r.Data = data
	} else if l == 1 {
		r.Data = data[0]
	}

	return r
}

type route struct {
	handler reflect.Value
}

func (r *route) newArg() interface{} {
	t := r.handler.Type()
	if t.NumIn() < 2 {
		return nil
	}
	return reflect.New(t.In(1).Elem()).Interface()
}

func (r *route) callHandler(req *http.Request, arg interface{}) (interface{}, error) {
	in := make([]reflect.Value, 0, 2)
	in = append(in, reflect.ValueOf(req))
	if r.handler.Type().NumIn() > 1 {
		in = append(in, reflect.ValueOf(arg))
	}

	out := r.handler.Call(in)

	var res, e interface{}
	if len(out) == 1 {
		e = out[0].Interface()
	} else {
		res = out[0].Interface()
		e = out[1].Interface()
	}

	if e != nil {
		return res, e.(error)
	}
	return res, nil
}

const contentType = "application/json; charset=utf-8"

var routes = make(map[string]route, 64)

// Add register a API handler to route map.
// the function should only be called at initialization time,
// otherwise there can be a race condition
// 'handler' should be a function with proto type
//     func(r *http.Request, args *TypeXXX) (interface{}, error)
// or
//     func(r *http.Request) (interface{}, error)
// or
//     func(r *http.Request, args *TypeXXX) error
// or
//     func(r *http.Request) error
func Add(name string, handler interface{}) {
	if _, ok := routes[name]; ok {
		panic(fmt.Errorf("route '%v' already registered", name))
	}

	t := reflect.TypeOf(handler)
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("handler of route '%v' is not a function", name))
	}

	e := fmt.Errorf("handler proto type of route '%v' is wrong", name)
	if num := t.NumIn(); num < 1 || num > 2 {
		panic(e)
	}

	if t.In(0) != reflect.TypeOf((*http.Request)(nil)) {
		panic(e)
	}

	isErr := func(t reflect.Type) bool {
		return t.Implements(reflect.TypeOf((*error)(nil)).Elem())
	}

	if num := t.NumOut(); num == 1 && isErr(t.Out(0)) {
	} else if num == 2 && isErr(t.Out(1)) {
	} else {
		panic(e)
	}

	routes[name] = route{handler: reflect.ValueOf(handler)}
}

// ServeHTTP handles HTTP request
func ServeHTTP(urlPrefix string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Path[len(urlPrefix):]
	route, ok := routes[name]
	if !ok {
		w.WriteHeader(http.StatusFound)
		return
	}

	arg := route.newArg()
	if arg != nil {
		if e := json.NewDecoder(r.Body).Decode(arg); e != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	resp := Result{Succeeded: true}
	res, e := route.callHandler(r, arg)
	if e != nil {
		resp.Succeeded = false
		resp.Message = e.Error()
	}

	switch v := res.(type) {
	case int64:
		resp.Data = strconv.FormatInt(v, 10)
	case uint64:
		resp.Data = strconv.FormatUint(v, 10)
	default:
		resp.Data = res
	}

	w.Header().Add("Content-Type", contentType)
	json.NewEncoder(w).Encode(&resp)
}

// Call calls the specified API
func Call(url string, arg, result interface{}) error {
	var data []byte

	if arg != nil {
		if d, e := json.Marshal(arg); e != nil {
			return e
		} else {
			data = d
		}
	}

	resp, e := http.DefaultClient.Post(url, contentType, bytes.NewReader(data))
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	r := Result{Data: result}
	if e := json.NewDecoder(resp.Body).Decode(&r); e != nil {
		return e
	} else if !r.Succeeded {
		return errors.New(r.Message)
	}
	return nil
}
