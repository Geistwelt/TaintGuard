package logging

import (
	"bytes"
	"fmt"
	"reflect"
)

func toString(x interface{}) string {
	rv := reflect.ValueOf(x)

	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return fmt.Sprintf("%v", rv.Interface())
}

func normalize(kvs ...interface{}) string {
	if len(kvs)%2 != 0 {
		kvs = append(kvs, "ERROR_LOG")
	}
	buf := new(bytes.Buffer)
	buf.Reset()
	for i := 0; i < len(kvs)-1; i += 2 {
		fmt.Fprintf(buf, "%s=%s", toString(kvs[i]), toString(kvs[i+1]))
		if i < len(kvs)-2 {
			fmt.Fprint(buf, " ")
		}
	}
	return buf.String()
}
