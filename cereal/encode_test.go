package cereal

import "testing"
import _ "reflect"
import _ "fmt"

func TestMarshalBool1(t *testing.T) {
    b, _ := Marshal(true)
    if string(b) != "b:1;" {
        t.Error("marshal failed for bool")
    }
}

func TestMarshalBool2(t *testing.T) {
    b, _ := Marshal(false)
    if string(b) != "b:0;" {
        t.Error("marshal failed for bool")
    }
}

func TestMarshalInt(t *testing.T) {
    b, _ := Marshal(5)
    if string(b) != "i:5;" {
        t.Error("marshal failed for int")
    }
}

func TestMarshalString(t *testing.T) {
    b, _ := Marshal("test")
    if string(b) != "s:4:\"test\";" {
        t.Error("marshal failed for string")
    }
}

func TestMarshalArray(t *testing.T) {
    b, _ := Marshal(map[interface{}]interface{}{0: 10, 1: 11, 2: 12})
    if string(b) != "a:3:{i:0;i:10;i:1;i:11;i:2;i:12;}" {
        t.Error("marshal failed for Array")
    }
}

func TestMarshalObject(t *testing.T) {
    obj := PHPObject{ClassName: "Test", Properties: map[string]Property{"public": Property{PropType:Public, PropValue:1}}}
    b, _ := Marshal(obj)
    if string(b) != "O:4:\"Test\":1:{s:6:\"public\";i:1;}" {
        t.Error("marshal failed for Array")
    }
}