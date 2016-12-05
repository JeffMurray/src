package ezjson

import(
	"fmt"
	"errs"
	"time"
)

func SetDate(format string, j_map map[string]interface{}, val time.Time, path ...string) *errs.ClnErr {
	time_format, time_format_err := Str_format_to_go_parse_format(format)
	if time_format_err != nil {
		return time_format_err.Traced("a4j4vh","Setting date")
	}
	rval := setInterface(j_map, val.Format(time_format), path...)
	return errs.TraceClnErrIfErr(rval,"gt76w5","Setting date")
}
func GetDate(format string, j_map map[string]interface{}, path ...string) (time.Time, *errs.ClnErr) {
	iface, err := GetInterface(j_map, path...)
	if err != nil {
		return time.Time{}, err.Traced("nhynro","Getting date")
	}
	time_val, time_val_err := Str_to_date(format, iface)
	if time_val_err != nil {
		return time.Time{}, time_val_err.Traced("j1ogfn", "Key is " + path[len(path)-1])
	}
	return time_val, nil
}
func Str_format_to_go_parse_format( format string ) (string, *errs.ClnErr) {
	switch format {
        case "ANSIC": return time.ANSIC, nil //"Mon Jan _2 15:04:05 2006"
        case "UnixDate": return time.UnixDate, nil //""Mon Jan _2 15:04:05 MST 2006"
        case "RubyDate": return time.RubyDate, nil //""Mon Jan 02 15:04:05 -0700 2006"
        case "RFC822": return time.RFC822, nil //""02 Jan 06 15:04 MST"
        case "RFC822Z": return time.RFC822Z, nil //""02 Jan 06 15:04 -0700" // RFC822 with numeric zone
        case "RFC850": return time.RFC850, nil //""Monday, 02-Jan-06 15:04:05 MST"
        case "RFC1123": return time.RFC1123, nil //""Mon, 02 Jan 2006 15:04:05 MST"
        case "RFC1123Z": return time.RFC1123Z, nil //""Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
        case "RFC3339": return time.RFC3339, nil //""2006-01-02T15:04:05Z07:00"
        case "RFC3339Nano": return time.RFC3339Nano, nil //""2006-01-02T15:04:05.999999999Z07:00"
        //Kitchen     = "3:04PM"
        // Handy time stamps.
        case "Stamp": return time.Stamp, nil //""Jan _2 15:04:05"
        case "StampMilli": return time.StampMilli, nil //""Jan _2 15:04:05.000"
        case "StampMicro": return time.StampMicro, nil //""Jan _2 15:04:05.000000"
        case "StampNano": return time.StampNano, nil //""Jan _2 15:04:05.000000000"
		default: return "Error", errs.NewClnErr("hzhdzh", fmt.Sprintf("Unrecognized time format: %s",format))
	}
}
func Str_to_date(format string, date_str_i interface{}) (time.Time, *errs.ClnErr) {
	date_str, ok := date_str_i.(string)
	if !ok {
		return time.Now(), errs.NewClnErr("mfmftw","Date is not in a string.")
	}
	go_format, go_format_err := Str_format_to_go_parse_format(format)
	if go_format_err != nil {
		return time.Now(), go_format_err.Traced("wif2o6","Traced.")
	}
	rval, rval_err := time.Parse(go_format, date_str)
	if rval_err != nil {
		return time.Now(), errs.NewClnErr("err2nc",fmt.Sprintf("%s cannot parse as %s",date_str,go_format))
	}
	return rval, nil
}