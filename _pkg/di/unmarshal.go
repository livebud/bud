package di

import (
	"fmt"
	"reflect"

	"github.com/livebud/bud/pkg/di/internal/reflector"
)

func Unmarshal(in Injector, ctx any) (err error) {
	if err := unmarshal(in, ctx); err != nil {
		return err
	}
	return nil
}

func unmarshal(in Injector, ctx any) error {
	v := reflect.ValueOf(ctx)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("ctx should be a pointer")
	}
	// If it is a pointer to pointer, check if it is nil and allocate a new struct
	if v.Elem().Kind() == reflect.Ptr && v.Elem().Elem().Kind() != reflect.Struct {
		if v.Elem().IsNil() {
			v.Elem().Set(reflect.New(v.Elem().Type().Elem()))
		}
		v = v.Elem() // Update v to point to the actual struct
	}
	// At this point, v should be a pointer to a struct
	if v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("ctx should be a pointer to struct or pointer to pointer to struct")
	}
	stct := v.Elem()
	for i := 0; i < stct.NumField(); i++ {
		fieldValue := stct.Field(i)
		fieldType := fieldValue.Type()
		typeID, err := reflector.ID(fieldType)
		if err != nil {
			return err
		}
		value, err := in.load(typeID)
		if err != nil {
			return err
		}
		fieldValue.Set(value)
	}
	return nil
}
