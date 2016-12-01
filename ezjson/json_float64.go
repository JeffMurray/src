package ezjson

import (
	"errs"
	"encoding/json"
	"fmt"
	"reflect"
)

func SetFloat64(j_map map[string]interface{}, val float64, path ...string) (*errs.ClientError) {
	return setInterface(j_map, val, path...)
}
func GetFloat64(j_map map[string]interface{}, path ...string) (float64, *errs.ClientError) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return 0.0, err
	}
	rval, rval_err := IfaceToFloat64(iface)
	if rval_err != nil {
		return 0.0, rval_err.Traced("fe7kxs",path[len(path)-1] + " is not a float64.")
	}
	return rval, nil
}
func IfaceToFloat64(iface interface{})(float64, *errs.ClientError) {
	switch v := iface.(type) {
		case json.Number:
			val, err := v.Float64()
			if err != nil {
				return 0, errs.NewClientError("x3n2zu","json.Number is not a float64.")
			} 
			return val, nil
		case float64:
			return float64(v), nil
		default:
			return 0, errs.NewClientError("ov1s8k",fmt.Sprintf("value is type %v, not an float64.", reflect.TypeOf(v)))
	}
}
