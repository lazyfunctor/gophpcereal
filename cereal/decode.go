package cereal

import (
    "bytes"
    "strconv"
    "io"
    "errors"
    "fmt"
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

const (
    colon = byte(':')
    semicol = byte(';')
    quote = byte('"')
    leftCurly = byte('{')
    rightCurly = byte('}')
    asterisk = byte('*')
)

type PropertyType int

const (
    Public PropertyType = iota
    Private PropertyType = iota
    Protected PropertyType = iota
)

type Property struct {
    key string
    propType PropertyType
}

//used defer/recovery for error handling
func Unmarshal(inp []byte) (v interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            if _, ok := r.(runtime.Error); ok {
                panic(r)
            }
            err = r.(error)
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
        fmt.Println(string(delimiter))
        fmt.Println(string(byte_))
        fmt.Println(d.offset)
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

func (d decodeState) readObject(val map[interface {}]interface{}) {
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
            p := Property{propType: Protected, key: key}
            val[p] = v
        } else if strings.HasPrefix(k, string(nulledClassName)) {
            key := k[length+2:]
            p := Property{propType: Private, key: key}
            val[p] = v
        } else {
            p := Property{propType: Public, key: k}
            val[p] = v
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
            case byte('n'):
                return None{}
            case byte('i'):
                var ival int
                d.readInt(&ival)
                return ival
            case byte('b'):
                var bval bool
                d.readBool(&bval)
                return bval
            case byte('d'):
                var dval float64
                d.readFloat(&dval)
                return dval
            case byte('s'):
                var sval string
                d.readString(&sval)
                return sval
            case byte('a'):
                aval := make(map[interface {}]interface{})
                d.readArray(aval)
                return aval
            case byte('O'):
                oval := make(map[interface {}]interface{})
                d.readObject(oval)
                return oval

        }
    }
    return nil
}