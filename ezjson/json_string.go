package ezjson

import (
	"errs"
)

func SetString(j_map map[string]interface{}, val string, path ...string) *errs.ClnErr {
	rval := setInterface(j_map, val, path...)
	return errs.TraceClnErrIfErr(rval, "ihntfb", "Setting string")
}
func GetString(j_map map[string]interface{}, path ...string) (string, *errs.ClnErr) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return "", err.Traced("l1752e","Getting string")
	}
	rval, ok := iface.(string)
	if !ok {
		return "", errs.NewClnErr("mu3n13",path[len(path)-1] + " is not a string.")
	}
	return rval, nil
}
