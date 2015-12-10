package cereal

import "testing"
import "reflect"
import _ "fmt"

func TestUnmarshalBool1(t *testing.T) {
    v, _ := Unmarshal([]byte("b:0;"))
    if v.(bool) != false {
        t.Error("Expected false")
    }
}

func TestUnmarshalBool2(t *testing.T) {
    v, _ := Unmarshal([]byte("b:1;"))
    if v.(bool) != true {
        t.Error("Expected true")
    }
}

func TestUnmarshalInt(t *testing.T) {
    v, _ := Unmarshal([]byte("i:500;"))
    if v.(int) != 500 {
        t.Error("Expected intgere value 500")
    }
}

func TestUnmarshalStr(t *testing.T) {
    v, _ := Unmarshal([]byte("s:4:\"test\";"))
    if v.(string) != "test" {
        t.Error("Expected string value test")
    }
}

func TestUnmarshalArray1(t *testing.T) {
    v, _ := Unmarshal([]byte("a:3:{i:0;i:10;i:1;i:11;i:2;i:12;}"))
    if ! reflect.DeepEqual(v, map[interface{}]interface{}{0: 10, 1: 11, 2: 12}) {
        t.Error("Error in parsing array")
    }
}

func TestUnmarshalArray2(t *testing.T) {
    v, _ := Unmarshal([]byte("a:2:{s:3:\"foo\";i:4;s:3:\"bar\";i:2;}"))
    if ! reflect.DeepEqual(v, map[interface{}]interface{}{"foo": 4, "bar": 2}) {
        t.Error("Error in parsing array")
    }
}

func TestUnmarshalObj(t * testing.T) {
    expObj := PHPObject{ClassName: "Test", Properties: map[string]Property{"public": Property{PropType:Public, PropValue:1}}}
    v, _ := Unmarshal([]byte("O:4:\"Test\":1:{s:6:\"public\";i:1;}"))
    if !reflect.DeepEqual(v, expObj) {
        t.Error("Error in unmarshal object")
    }
}

// func TestUnmarshalObj(t *testing.T) {
//     v, _ := Unmarshal([]byte("O:4:\"Test\":3:{s:6:\"public\";i:1;s:9:\"protected\";i:2;s:7:\"private\";i:3;}\""))
//     if ! reflect.DeepEqual(v, map[interface{}]interface{}{Property{key:"public", propType:Public}: 1,
//      Property{key:"protected", propType:Public}: 2, Property{key:"private", propType:Public}: 3}) {
//         t.Error("Error in parsing object")
//     }
// }