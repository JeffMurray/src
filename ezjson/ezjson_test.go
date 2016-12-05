package ezjson

import (
	"testing"
	"time"
	"fmt"
)

var test_specs_x = map[string]interface{}{
	"l5h39n9cqo6cz5sb1:1":	map[string]interface{}{
		"Spec": map[string]interface{}{
			"Key":	"l5h39n9cqo6cz5sb1",
			"Ver":	1,
			"Auth":	0,
		},
		"Doc":	map[string]interface{}{
			"Name": "Test spec",
			"KeySpecs":	map[string]interface{}{
				"Fantastic Int":	map[string]interface{}{
					"Type": "Int",
					"Required":	true,
					"IntRules": map[string]interface{}{
						"MaxVal":	1000,
						"MinVal":	0,
					},
				},
				"Title":	map[string]interface{}{
					"Type": "String",
					"Required":	false,
					"StringRules": map[string]interface{}{
						"MaxLen":	10,
						"MinLen":	5,
					},
				},
			},
		},
	},
}
func TestEzjson(t *testing.T) {
	test_specs, _ := Cycle(test_specs_x)
	title_map, err := GetMap(test_specs,"l5h39n9cqo6cz5sb1:1", "Doc","KeySpecs","Title")
	//**STRING
	if err != nil {
		t.Error(fmt.Sprintf("%v", err))
	}	
	str_1, _ := GetString(test_specs,"l5h39n9cqo6cz5sb1:1", "Doc","KeySpecs","Title", "Type")
	if str_1 != "String" {
		t.Error( "str_1 != String" )
	}
	str_2, _ := GetString(title_map,"Type")
	if str_1 != str_2 {
		t.Error( "str_1 != str_2" )
	}
	_ , err3 := GetString(test_specs,"l5h39n9cqo6cz5sb1:1", "Doc","KeySpecs","Fantastic Int","IntRules", "MaxVal")
	if err3.Token != "mu3n13" {
		t.Error(fmt.Sprintf("%v", err3))
	}
	//**INT
	var f64 float64 = 1000
	floatInt, _ := IfaceToInt(f64)
	if floatInt != 1000 {
		t.Error( "floatInt != 1000" )
	}
	IntInt, _ := IfaceToInt(int(1000))
	if IntInt != 1000 {
		t.Error( "IntInt != 1000" )
	}
	max_val, _ := GetInt(test_specs,"l5h39n9cqo6cz5sb1:1", "Doc","KeySpecs","Fantastic Int","IntRules", "MaxVal")
	if max_val != 1000 {
		t.Error( "max_val != 1000" )
	}
	_, err2 := GetInt(test_specs,"l5h39n9cqo6cz5sb1:1", "Doc","KeySpecs","Fantastic Int","IntRules")
	if err2 == nil {
		t.Error( err2.Error() )
	} else {
		if err2.Token != "pjex5i" {
			t.Error(fmt.Sprintf("%v", err2))
		}
	}
	byte_val, err4 := GetByte(test_specs,"l5h39n9cqo6cz5sb1:1", "Doc","KeySpecs","Title","StringRules", "MinLen")
	if byte_val != 5 {
		t.Error( "byte_val != 5, " + fmt.Sprintf("%v", err4) )
	}
	byte_val, err5 := GetByte(test_specs,"l5h39n9cqo6cz5sb1:1", "Doc","KeySpecs","Fantastic Int","IntRules", "MaxVal")
	if err5.Token != "md4byk" {
		t.Error( fmt.Sprintf("%v", err5) )
	}	
	m := map[string]interface{}{
		"Fantastic Int": 0,
		"Im out": 0,
	}
	master, _ := GetMap(test_specs, "l5h39n9cqo6cz5sb1:1", "Doc", "KeySpecs")
	nm_keys := NonMatchingKeys(master, m)
	if len(nm_keys) != 1 || nm_keys[0] != "Im out" {
	
		t.Error("!Im out")
	}
	req_keys := KeysWithValue(master, true, "Required")
	if len(req_keys) != 1 || req_keys[0] != "Fantastic Int" {
		t.Error("!Fantastic Int")
	}	
	keys := []string{"Bull 1", "Bull 2"}
	m2 := KeysToMap(keys...)
	if m2["Bull 1"] == false || m2["Bull 2"] == false {
		t.Error("KeysToMap")
	}
	nm_map := KeysToMap(req_keys...)
	keys = NonMatchingKeys(m, nm_map )
	if len(keys) != 0 {
		t.Error("NonMatchingKeys")
	}
	delete(m,"Fantastic Int")
	keys = NonMatchingKeys(m, nm_map )
	if len(keys) != 1 || keys[0] != "Fantastic Int" {
		t.Error("!Fantastic Int")
	}	
	test_insert := map[string]interface{}{}
	test_ins_err := setInterface(test_insert, "yo", "yo_field")
	if test_ins_err != nil {
		t.Error(test_ins_err.Error())
	}
	test_str, _ := GetString(test_insert, "yo_field")
	if test_str != "yo" {
		t.Error("yo")
	}
	test_ins_err = SetString(test_insert, "we dat", "yo_field", "whodat")
	if test_ins_err.Token != "j3tprh" {
		t.Error(test_ins_err.Error())
	}
	test_ins_err = SetString(test_insert, "we dat", "testdat", "whodat")
	if test_ins_err!= nil {		
		t.Error(test_ins_err.Error())
	}
	test_ins_str, test_ins_str_error := GetString(test_insert, "testdat","whodat")
	if test_ins_str != "we dat" {
		t.Error(test_ins_str_error.Error())
	}
	set_int64_err := SetInt64(test_insert, 1255353618542527051, "testdat", "countdat")
	if set_int64_err!= nil {
		t.Error(set_int64_err.Error())
	}
	test_int64, _ := GetInt64(test_insert, "testdat", "countdat")
	if test_int64 != 1255353618542527051 {
		t.Error(fmt.Sprintf("test_int64 = %v", test_int64))
	}
	set_int_err := SetInt(test_insert, 1255, "testdat", "countdat")
	if set_int_err!= nil {
		t.Error(set_int_err.Error())
	}
	test_int, _ := GetInt(test_insert, "testdat", "countdat")
	if test_int != 1255 {
		t.Error(fmt.Sprintf("test_int = %v", test_int))
	}	
	set_float64_err := SetFloat64(test_insert, 1255.1255, "testdat", "countdat")
	if set_float64_err!= nil {
		t.Error(set_float64_err)
	}
	test_float64, _ := GetFloat64(test_insert, "testdat", "countdat")
	if test_float64 != 1255.1255 {
		t.Error(fmt.Sprintf("test_float64 = %v", test_float64))
	}
	some_date := time.Now()
	set_date_err := SetDate("RFC3339Nano",test_insert, some_date, "testdat", "countdat")
	if set_date_err!= nil {
		t.Error(set_date_err)
	}
	test_date, _ := GetDate("RFC3339Nano", test_insert, "testdat", "countdat")
	if test_date != some_date {
		t.Error(fmt.Sprintf("some_date = %v, test_date = %v", some_date, test_date))
	}
	set_byte_err := SetByte(test_insert, 12, "testdat", "countdat")
	if set_byte_err!= nil {
		t.Error(set_byte_err)
	}
	test_byte, _ := GetByte(test_insert, "testdat", "countdat")
	if test_byte != 12 {
		t.Error(fmt.Sprintf("test_byte = %v", test_byte))
	}	
	set_bool_err := SetBool(test_insert, true, "testdat", "countdat")
	if set_bool_err!= nil {
		t.Error(set_bool_err)
	}
	test_bool, _ := GetBool(test_insert, "testdat", "countdat")
	if test_bool != true {
		t.Error(fmt.Sprintf("test_bool = %v", test_bool))
	}		
}