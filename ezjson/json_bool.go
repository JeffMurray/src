package ezjson

import(
	"fmt"
	"errs"
)

func SetBool(j_map map[string]interface{}, val bool, path ...string) (*errs.ClientError) {
	return setInterface(j_map, val, path...)
}
func GetBool(j_map map[string]interface{}, path ...string) (rval bool, e *errs.ClientError) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return false, err
	}
	bool_val, bool_val_ok := iface.(bool)
	if !bool_val_ok {
		return false, errs.NewClientError("z44u2v", fmt.Sprintf("%s is not a bool", path[len(path)-1]))
	}
	return bool_val, nil
}
