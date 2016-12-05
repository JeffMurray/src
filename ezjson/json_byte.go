package ezjson

import(
	"fmt"
	"reflect"
	"errs"
	"encoding/json"
)

func SetByte(j_map map[string]interface{}, val byte, path ...string) *errs.ClnErr {
	rval := setInterface(j_map, val, path...)
	return errs.TraceClnErrIfErr(rval,"cq1bqn","Setting byte")
}
func GetByte(j_map map[string]interface{}, path ...string) (rval byte, e *errs.ClnErr) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return 0, errs.TraceClnErrIfErr(err,"u8dw2u","Getting byte")
	}
	byte_val, byte_val_err := IfaceToByte(iface)
	if byte_val_err != nil {
		return 0, byte_val_err.Traced("ogpj6t", fmt.Sprintf("%s is not a byte %s", path[len(path)-1], byte_val_err.Error()))
	}
	return byte_val, nil
}

func IfaceToByte(byte_i interface{}) (byte, *errs.ClnErr) {
	switch v := byte_i.(type) {
		case json.Number:
			val, err := v.Int64()
			if err != nil {
				return 0, errs.NewClnErr("ehy2v2",fmt.Sprintf("json.Number is not an byte: %v", v))
			} 
			if val  != int64(byte(val)) {
				return 0, errs.NewClnErr("md4byk",fmt.Sprintf("json.Number is out of range: %v", v))
			} 
			return byte(val), nil
		case float64:
			if v  != float64(byte(v)) {
				return 0, errs.NewClnErr("yrbrbp",fmt.Sprintf("json.Number is out of range: %v", v))
			}
			return byte(v), nil
		case byte:
			return byte(v), nil
		default:
			return 0, errs.NewClnErr("ikni0f",fmt.Sprintf("value is type %v, not a byte.", reflect.TypeOf(v)))
	}
}