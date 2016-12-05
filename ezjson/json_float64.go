package ezjson

import (
	"errs"
	"encoding/json"
	"fmt"
	"reflect"
)

func SetFloat64(j_map map[string]interface{}, val float64, path ...string) *errs.ClnErr {
	rval := setInterface(j_map, val, path...)
	return errs.TraceClnErrIfErr(rval,"","Setting float64")
}
func GetFloat64(j_map map[string]interface{}, path ...string) (float64, *errs.ClnErr) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return 0.0, err.Traced("gr78zp","Getting float64")
	}
	rval, rval_err := IfaceToFloat64(iface)
	if rval_err != nil {
		return 0.0, rval_err.Traced("fe7kxs",fmt.Sprintf("%s is not a float64.",path[len(path)-1]))
	}
	return rval, nil
}
func IfaceToFloat64(iface interface{})(float64, *errs.ClnErr) {
	switch v := iface.(type) {
		case json.Number:
			val, err := v.Float64()
			if err != nil {
				return 0, errs.NewClnErr("x3n2zu","json.Number is not a float64.")
			} 
			return val, nil
		case float64:
			return float64(v), nil
		default:
			return 0, errs.NewClnErr("ov1s8k",fmt.Sprintf("value is type %v, not an float64.", reflect.TypeOf(v)))
	}
}