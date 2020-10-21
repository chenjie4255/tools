package structfix

import (
	"errors"
	"reflect"
)

type fixer interface {
	FixNilArray()
}

// FixNilArray 修复为nil的数组
func FixNilArray(obj interface{}, justTopLevel bool) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Slice {
		return errors.New("object must be a struct pointer or slice")
	}

	return fixNilArrayForReflectVal(v, justTopLevel)
}

func fixNilArrayForReflectVal(val reflect.Value, justTopLevel bool) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			fixNilArrayForReflectVal(item, justTopLevel)
		}
		return nil
	}

	if val.Kind() != reflect.Struct {
		return errors.New("object must be a struct pointer or slice")
	}

	fixerType := reflect.TypeOf(new(fixer)).Elem()

	valT := val.Type()
	if valT.Implements(fixerType) || reflect.PtrTo(valT).Implements(fixerType) {
		if val.CanSet() && val.CanAddr() {

			fixerObj, ok := val.Addr().Interface().(fixer)
			if !ok {
				panic("should implement fixer type")
			}

			fixerObj.FixNilArray()
			// return nil
		}
	}

	if justTopLevel {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldT := field.Type()
		// fmt.Printf("fixxing nil array for %s...\n", fieldT.Name())
		if fieldT.Implements(fixerType) || reflect.PtrTo(fieldT).Implements(fixerType) {
			if !field.CanSet() || !field.CanAddr() {
				// fmt.Printf("field(%s) cannot set or addr,%s, %v ,%v \n", val.Type().Field(i).Name, fieldT.Name(), field.CanSet(), field.CanAddr())
				continue
			}

			fixerObj, ok := field.Addr().Interface().(fixer)
			if !ok {
				panic("should implement fixer type")
			}

			// fmt.Printf("field(%s) fix nil array...%s\n", val.Type().Field(i).Name, fieldT.Name())

			fixerObj.FixNilArray()
		}
		if field.Kind() == reflect.Struct || field.Kind() == reflect.Slice || field.Kind() == reflect.Ptr {
			fixNilArrayForReflectVal(field, justTopLevel)
		}
	}
	return nil
}
