package common

import "reflect"

func HasNonZeroField(s interface{}) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := v.Field(i)
		switch sf.Type.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			if !fv.IsNil() {
				return true
			}
		case reflect.Struct:
			if HasNonZeroField(fv.Interface()) {
				return true
			}
		// case reflect.Array:
		// TODO: call recursively for array elements
		default:
			// if reflect.Zero(sf.Type).Interface() != fv.Interface() {
			// 	return true
			// }
			if !v.Field(i).IsZero() {
				return true
			}
		}
	}
	return false
}
