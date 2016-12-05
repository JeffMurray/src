package ezjson

import (
	"fmt"
	"bytes"
	"errs"
	"encoding/json"
)

//package ezjson has get and set functions that make woring with 
//the native google go map[string]interface{} scheme easier
func Decode(json_str string) ( map[string]interface{}, *errs.SysErr ) {
	m := map[string]interface{}{}
	d := json.NewDecoder(bytes.NewBuffer([]byte(json_str)))
	d.UseNumber() //<----
	err := d.Decode(&m)
	if err != nil {
		return m, errs.NewSysErr(err.Error(), "ii6upj", json_str)
	}
	return m, nil
}
func Byteify(json_map map[string]interface{}, pretty bool) ([]byte, *errs.SysErr) {
	if !pretty {
		b, err := json.Marshal(json_map)
		if err != nil {
			return b, errs.NewSysErr(err.Error(), "y9xl0j", "")
		}
		return b, nil
	} else {
		b, err := json.MarshalIndent(json_map, "", "\t")
		if err != nil {
			return b, errs.NewSysErr(err.Error(), "awuw6r", "")
		}
		return b, nil
	}
}
func IsNotExistErr(token string) bool {
	return token == "d2v0h1" || token == "gfj3n4"
}
func GetMap(j_map map[string]interface{}, path ...string) (map[string]interface{}, *errs.ClnErr) {
	next := j_map
	ok := true
	for _, p:= range  path {
		if next[p] == nil {
			return next, errs.NewClnErr("d2v0h1",fmt.Sprintf("%s does not exist",p))
		}
		next, ok = next[p].(map[string]interface{})
		if !ok {
			return next, errs.NewClnErr("j3tprh",fmt.Sprintf("%s not a map",p))
		}
	}
	return next, nil
}
func GetInterface(j_map map[string]interface{}, path ...string)(interface{}, *errs.ClnErr){
	if len(path) == 0 {
		return "", errs.NewClnErr("v6v3ai","Path is required")
	}
	doc, err := GetMap(j_map, path[0:len(path)-1]...)
	if err != nil {
		return "", err.Traced("tmt0h6","Getting interface.")
	}
	if doc[path[len(path)-1]] == nil {	
		return "", errs.NewClnErr("gfj3n4", fmt.Sprintf("%s does not exist.",path[len(path)-1]))
	}
	return doc[path[len(path)-1]], nil
}
func setInterface(j_map map[string]interface{}, val interface{}, path ...string)(*errs.ClnErr){
	if len(path) == 0 {
		return errs.NewClnErr("nic0hc","Path is required")
	}
	built_path := []string{}
	last_doc := j_map
	for _, key := range path {
		built_path = append(built_path, key)
		doc, err := GetMap(j_map, path[0:len(path)-1]...)
		if err != nil {
			switch err.Token {
				case "d2v0h1": //Does not exist
					last_doc[key] = map[string]interface{}{}
					last_doc = last_doc[key].(map[string]interface{})
				default:
					return err.Traced("x62p64",fmt.Sprintf("Key is: %s", key))
			}
		} else {
			last_doc = doc
		}
	}
	last_doc[path[len(path)-1]] = val
	return nil
}
func GetArray(j_map map[string]interface{}, path ...string) ([]interface{}, *errs.ClnErr) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return []interface{}{}, err.Traced("y0ezzo","Getting array.")
	}
	rval, ok := iface.([]interface{})
	if !ok {
		return []interface{}{}, errs.NewClnErr("c0dxq4",path[len(path)-1] + " is not an array.")
	}
	return rval, nil	
}
func NonMatchingKeys(master_doc map[string]interface{}, doc_to_check map[string]interface{}) []string {
	list := []string{}
	for key, _ := range doc_to_check {
		if master_doc[key] == nil {
			list = append( list, key )
		}
	}
	return list
}
func KeysWithValue( doc map[string]interface{}, value interface{}, path ...string ) []string {
	list := []string{}
	for key, _ := range doc {
		newPath := append([]string{key}, path...)
		val, err := GetInterface( doc, newPath... )
		if err == nil && val == value {
			list = append( list, key )
		}
	}
	return list
}
func KeysToMap( keys ...string ) map[string]interface{} {
	m := map[string]interface{}{}
	for _, key := range keys {
		m[key] = true
	}
	return m
}
//Use to make sure numbers are json.Number if you have a raw 
//(not produced from the json decoder) map you want to use.
//See ezjson_test
func Cycle( doc map[string]interface{} ) ( map[string]interface{}, *errs.SysErr ) {
	bts, err := Byteify( doc, false )
	if err != nil {
		return doc, err.Traced("c2byil","Cycling")
	}
	rval, e := Decode(string(bts))
	return rval, errs.TraceSysErrIfErr(e, "c2byil", string(bts))
}