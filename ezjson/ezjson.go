package ezjson

import (
	"fmt"
	"bytes"
	"errs"
	"encoding/json"
)

//package ezjson has get and set functions that make woring with 
//the native google go map[string]interface{} scheme easier
func Decode(json_str string) ( map[string]interface{}, *errs.SysError ) {
	m := map[string]interface{}{}
    d := json.NewDecoder(bytes.NewBuffer([]byte(json_str)))
    d.UseNumber() //<----
	err := d.Decode(&m)
	return m, errs.IfSysError(err, "ii6upj", "ezjson.go", "")
}
func Byteify(json_map map[string]interface{}, pretty bool) ([]byte, *errs.SysError) {
	if !pretty {
		b, err := json.Marshal(json_map)
		return b, errs.IfSysError(err, "y9xl0j", "ezjson.go", "")
	} else {
		b, err := json.MarshalIndent(json_map, "", "\t")
		return b, errs.IfSysError(err, "awuw6r", "ezjson.go", "")
	}
}
func IsNotExistErr(token string) bool {
	return token == "d2v0h1" || token == "gfj3n4"
}
func GetMap(j_map map[string]interface{}, path ...string) (map[string]interface{}, *errs.ClientError) {
	next := j_map
	ok := true
	for i := 0; i < len(path); i++ {
		if next[path[i]] == nil {
			return next, errs.NewClientError("d2v0h1",path[i] + " does not exist")
		}
		next, ok = next[path[i]].(map[string]interface{})
		if !ok {
			return next, errs.NewClientError("j3tprh",path[i] + " is not a map")
		}
	}
	return next, nil
}
func GetInterface(j_map map[string]interface{}, path ...string)(interface{}, *errs.ClientError){
	if len(path) == 0 {
		return "", errs.NewClientError("v6v3ai","path is required")
	}
	doc, err := GetMap(j_map, path[0:len(path)-1]...)
	if err != nil {
		return "", err
	}
	if doc[path[len(path)-1]] == nil {	
		return "", errs.NewClientError("gfj3n4",path[len(path)-1] + " does not exist.")
	}
	return doc[path[len(path)-1]], nil
}
func setInterface(j_map map[string]interface{}, val interface{}, path ...string)(*errs.ClientError){
	if len(path) == 0 {
		return errs.NewClientError("nic0hc","path is required")
	}
	built_path := []string{}
	last_doc := j_map
	for i := 0; i < len(path)-1; i++ {
		built_path = append(built_path, path[i])
		doc, err := GetMap(j_map, path[0:len(path)-1]...)
		if err != nil {
			switch err.Token {
				case "d2v0h1": //Does not exist
					last_doc[path[i]] = map[string]interface{}{}
					last_doc = last_doc[path[i]].(map[string]interface{})
				default:
					return err.Traced("x62p64",fmt.Sprintf("Key is: %s", path[i]))
			}
		} else {
			last_doc = doc
		}
	}
	last_doc[path[len(path)-1]] = val
	return nil
}
func GetArray(j_map map[string]interface{}, path ...string) ([]interface{}, *errs.ClientError) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return []interface{}{}, err
	}
	rval, ok := iface.([]interface{})
	if !ok {
		return []interface{}{}, errs.NewClientError("c0dxq4",path[len(path)-1] + " is not an array.")
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

//mainly used for testing so that numbers are json.Number
func Cycle( doc map[string]interface{} ) ( map[string]interface{}, error ) {
	bts, err := Byteify( doc, false )
	if err != nil {
		return doc, err
	}
	return Decode(string(bts))
}