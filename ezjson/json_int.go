package ezjson

import (
	"errs"
	"encoding/json"
	"fmt"
	"reflect"
)

func SetInt(j_map map[string]interface{}, val int, path ...string) *errs.ClnErr {
	rval := setInterface(j_map, val, path...)
	return errs.TraceClnErrIfErr(rval, "nv22nj", "Setting int")
}
func GetInt(j_map map[string]interface{}, path ...string) (int, *errs.ClnErr) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return 0, err.Traced("c2z0js","Setting int")
	}
	int_val, int_val_err := IfaceToInt(iface)
	if int_val_err != nil {
		return 0, int_val_err.Traced("nhgx12", fmt.Sprintf("Key is %v", path[len(path)-1]))
	}
	return int_val, nil
}
func IfaceToInt(iface interface{})(int, *errs.ClnErr) {
	//fmt.Println(fmt.Sprintf("%v %v", reflect.TypeOf(iface), iface))
	switch v := iface.(type) {
		case json.Number:
			val, err := v.Int64()
			if err != nil {
				return 0, errs.NewClnErr("mmkurv",fmt.Sprintf("json.Number is not an int: %v", v))
			} 
			if val  != int64(int(val)) {
				return 0, errs.NewClnErr("xnnxo9",fmt.Sprintf("json.Number is out of range: %v", v))
			} 
			return int(val), nil
		case float64:
			//How better to do this?
			if v  != float64(int(v)) {
				return 0, errs.NewClnErr("mpsrtu",fmt.Sprintf("json.Number is out of range: %v", v))
			}
			return int(v), nil
		case int:
			return int(v), nil
		default:
			return 0, errs.NewClnErr("pjex5i",fmt.Sprintf("value is type %v, not an int.", reflect.TypeOf(v)))
	}
}