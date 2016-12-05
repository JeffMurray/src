package ezjson

import(
	"fmt"
	"errs"
)

func SetBool(j_map map[string]interface{}, val bool, path ...string) *errs.ClnErr {
	rval := setInterface(j_map, val, path...)
	return errs.TraceClnErrIfErr(rval, "uw99ce", "Setting bool")
}
func GetBool(j_map map[string]interface{}, path ...string) (rval bool, e *errs.ClnErr) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return false, err.Traced("cdqmrr","Getting bool.")
	}
	bool_val, bool_val_ok := iface.(bool)
	if !bool_val_ok {
		return false, errs.NewClnErr("z44u2v", fmt.Sprintf("%s is not a bool", path[len(path)-1]))
	}
	return bool_val, nil
}









