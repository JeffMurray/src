package ezjson

import (
	"fmt"
	"reflect"
	"errs"
	"encoding/json"
)

func SetInt64(j_map map[string]interface{}, val int64, path ...string) (*errs.ClientError) {
	return setInterface(j_map, val, path...)
}
func GetInt64(j_map map[string]interface{}, path ...string) (int64, *errs.ClientError) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return 0, err
	}
	val, val_err := IfaceToInt64(iface)
	if val_err != nil {
		return 0, val_err.Traced("rtoumf", fmt.Sprintf("Key is %v", path[len(path)-1]))
	}
	return val, nil	
}
func IfaceToInt64(iface interface{})(int64, *errs.ClientError) {
	switch v := iface.(type) {
		case json.Number:
			val, err := v.Int64()
			if err != nil {
				return 0, errs.NewClientError("g2lcfi","json.Number is not a int64.")
			} 
			return val, nil
		case float64:
			if v  != float64(int64(v)) {
				return 0, errs.NewClientError("jcw4yo",fmt.Sprintf("json.Number is out of range: %v", v))
			}
			return int64(v), nil
		case int64:
			return int64(v), nil
		default:
			return 0, errs.NewClientError("wc7epu",fmt.Sprintf("value is type %v, not an 40.", reflect.TypeOf(v)))
	}
}
