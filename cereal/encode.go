package cereal

import (
    "reflect"
    "bytes"
    "runtime"
    "strconv"
    _ "errors"
    "fmt"
    )



func Marshal(v interface{}) (data []byte, err error) {
    defer func() {
        if rec := recover(); rec != nil {
            fmt.Println(rec)
            if _, ok := rec.(runtime.Error); ok {
                panic(rec)
            }
            err = rec.(error)
        }
    }()    
    w := bytes.NewBuffer(data)
    e := encodeState{w: w}
    e.Encode(v)
    data = e.w.Bytes()
    return
}


type encodeState struct {
    w *bytes.Buffer
    err error
}

func (e *encodeState) error(err error) {
    panic(err)
}

func (e* encodeState) write(data []byte) {
    _, err := e.w.Write(data)
    if err != nil {
        e.error(err)
    }
}

func (e* encodeState) writeString(data string) {
    _, err := e.w.WriteString(data)
    if err != nil {
        e.error(err)
    }
}

func (e *encodeState) encodeBool(v bool) {
    if v {
        e.write([]byte("b:1;"))
    } else {
        e.write([]byte("b:0;"))
    }
}

func (e *encodeState) encodeInt(v int) {
    e.writeString("i:")
    e.writeString(strconv.Itoa(v))
    e.writeString(";")
}

func (e *encodeState) encodeString(v string) {
    e.writeString("s:")
    e.writeString(strconv.Itoa(len(v)))
    e.writeString(":\"")
    e.writeString(v)
    e.writeString("\";")
}

func (e *encodeState) encodeArray(v map[interface{}]interface{}) {
    n := len(v)
    e.writeString("a:")
    e.writeString(strconv.Itoa(n))
    e.writeString(":{")
    for key, val := range(v) {
        e.Encode(key)
        e.Encode(val)
    }
    e.writeString("}")
}

func (e *encodeState) encodeObject(v PHPObject) {
    lenClassName := len(v.ClassName)
    propCount := len(v.Properties)
    e.writeString("O:")
    e.writeString(strconv.Itoa(lenClassName)); e.writeString(":\"")
    e.writeString(v.ClassName)
    e.writeString("\":")
    e.writeString(strconv.Itoa(propCount))
    e.writeString(":{")
    for key, prop := range(v.Properties) {
        if prop.PropType == Protected {
            e.write([]byte{0, asterisk, 0})
        } else if prop.PropType == Private {
            e.write([]byte{0})
            e.writeString(v.ClassName)
            e.write([]byte{0})
        } 
        e.Encode(key)
        e.Encode(prop.PropValue)
    }
    e.writeString("}")
}


func (e *encodeState) Encode(v interface{}) {
    rv := reflect.ValueOf(v)
    if rv.Type().Kind() == reflect.Ptr {
        fmt.Println("pointer")
        rv = rv.Elem()
    }
    kind := rv.Type().Kind()
    switch kind {
        case reflect.Bool:
            e.encodeBool(v.(bool))
        case reflect.Int:
            e.encodeInt(v.(int))
        case reflect.String:
            e.encodeString(v.(string))
        case reflect.Map:
            e.encodeArray(v.(map[interface{}]interface{}))
        case reflect.Struct:
            e.encodeObject(v.(PHPObject))            
    }
}