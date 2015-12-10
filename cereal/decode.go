package cereal

import (
    "bytes"
    "strconv"
    "io"
    "errors"
    _ "fmt"
    "runtime"
    "strings"
    )


type decodeState struct {
    r *bytes.Reader
    err error
    offset int
}

type None struct{}

var errUnexpectedToken = errors.New("Unexpected token")
var errMissingToken = errors.New("Missing Token")

//used defer/recovery for error handling
func Unmarshal(inp []byte) (v interface{}, err error) {
    defer func() {
        if rec := recover(); rec != nil {
            if _, ok := rec.(runtime.Error); ok {
                panic(rec)
            }
            err = rec.(error)
        }
    }()    
    r := bytes.NewReader(inp)
    d := decodeState{r: r}
    v = d.Decode()
    return
}

func (d decodeState) error(err error) {
    panic(err)
}

func (d decodeState) readUntil(delimiter byte, data *[]byte) {
    for {
        byte_, err := d.r.ReadByte()
        if err != nil {
            if err == io.EOF {
                d.error(errMissingToken)
            } else {
                d.error(err)
            }
            return
        }
        d.offset += 1
        if byte_ == delimiter {
            return
        } else {
            *data = append(*data, byte_)
        }
    }
}

func (d decodeState) expect(delimiter byte) {
    var byte_ byte
    byte_, err := d.r.ReadByte()
    if err != nil {
        if err == io.EOF {
            d.error(errUnexpectedToken)
        } else {
            d.error(err)
        }
    }
    d.offset += 1
    if byte_ != delimiter {
        d.error(errUnexpectedToken)
    }        
}

func (d decodeState) readInt(val *int) {
    d.expect(colon)
    var data []byte
    d.readUntil(semicol, &data)
    inp := string(data)
    v, err := strconv.Atoi(inp)
    if err != nil {
        d.error(err)
    }
    *val = v
    return
}

func (d decodeState) readFloat(val *float64 ) {
    d.expect(colon)
    var data []byte
    d.readUntil(semicol, &data)
    inp := string(data)
    v, err := strconv.ParseFloat(inp, 64)
    if err != nil {
        d.error(err)
    }
    *val = v
    return
}

func (d decodeState) readBool(val * bool) {
    d.expect(colon)
    var data []byte
    d.readUntil(semicol, &data)
    if string(data) == "1" {
        *val = true
    } else if string(data) == "0" {
        *val = false
    } else {
        d.error(errUnexpectedToken)
    }
}

func (d decodeState) readString(val *string) {
    d.expect(colon)
    var lengthData []byte
    d.readUntil(colon, &lengthData)
    length, err := strconv.Atoi(string(lengthData))
    if err != nil {
        d.error(err)
        return
    }
    d.expect(quote)
    strBytes := make([]byte, length)
    n, err := d.r.Read(strBytes)
    if err != nil {
        d.error(err)
        return
    }
    d.offset += n
    d.expect(quote)
    *val = string(strBytes)
    return 
}

func (d decodeState) readArray(val map[interface {}]interface{}) {
    d.expect(colon)
    var lengthData []byte
    d.readUntil(colon, &lengthData)
    length, err := strconv.Atoi(string(lengthData))
    if err != nil {
        d.error(err)
        return
    }
    d.expect(leftCurly)
    for i := 1; i <= length; i += 1 {
        k := d.Decode()
        v := d.Decode()
        val[k] = v
    }
    d.expect(rightCurly)
    return
}

func (d decodeState) readObject(val *PHPObject) {
    d.expect(colon)
    var lengthData []byte
    d.readUntil(colon, &lengthData)
    length, err := strconv.Atoi(string(lengthData))
    if err != nil {
        d.error(err)
        return
    }
    d.expect(quote)
    className := make([]byte, length)
    n, err := d.r.Read(className)
    val.ClassName = string(className)
    val.Properties = make(map[string]Property)
    if err != nil {
        d.error(err)
        return
    }
    d.offset += n
    d.expect(quote)    
    d.expect(colon)
    var propLengthData []byte
    d.readUntil(colon, &propLengthData)
    propLength, err := strconv.Atoi(string(propLengthData))
    if err != nil {
        d.error(err)
        return
    }
    d.expect(leftCurly)
    for i := 1; i <= propLength; i += 1 {
        k := d.Decode().(string)
        v := d.Decode()
        nulledClassName := append([]byte{0}, className...)
        nulledClassName = append(nulledClassName, byte(0))
        if strings.HasPrefix(k, string([]byte{0, asterisk, 0})) {
            key := k[3:]
            p := Property{PropType: Protected, PropValue: v}
            val.Properties[key] = p
        } else if strings.HasPrefix(k, string(nulledClassName)) {
            key := k[length+2:]
            p := Property{PropType: Private, PropValue: v}
            val.Properties[key] = p
        } else {
            p := Property{PropType: Public, PropValue: v}
            val.Properties[k] = p
        }
    }
    d.expect(rightCurly)

}

func (d decodeState) Decode() interface{} {
    for {
        type_, err := d.r.ReadByte()
        d.offset += 1
        if err == io.EOF {
            break
        } else if err != nil {
            d.err = err
            return nil
        }        
        switch type_ {
            case nullMark:
                return None{}
            case intMark:
                var ival int
                d.readInt(&ival)
                return ival
            case boolMark:
                var bval bool
                d.readBool(&bval)
                return bval
            case floatMark:
                var dval float64
                d.readFloat(&dval)
                return dval
            case strMark:
                var sval string
                d.readString(&sval)
                return sval
            case arrMark:
                aval := make(map[interface {}]interface{})
                d.readArray(aval)
                return aval
            case objMark:
                oval := &PHPObject{}
                d.readObject(oval)
                return *oval

        }
    }
    return nil
}