package ezjson

import (
	"errs"
)

func SetString(j_map map[string]interface{}, val string, path ...string) (*errs.ClientError) {
	return setInterface(j_map, val, path...)
}
func GetString(j_map map[string]interface{}, path ...string) (string, *errs.ClientError) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return "", err
	}
	rval, ok := iface.(string)
	if !ok {
		return "", errs.NewClientError("mu3n13",path[len(path)-1] + " is not a string.")
	}
	return rval, nil
}
