package commander

import (
	"fmt"
	"strconv"
)

type Int struct {
	target *int
	value  *int
}

func (v *Int) Default(value int) {
	*v.value = value
}

func (v *Int) Optional() {

}

type intValue struct{ inner *Int }

func (v *intValue) Get() interface{} {
	return *v.inner.target
}

func (v *intValue) Set(val string) error {
	n, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*v.inner.target = n
	return nil
}

func (v *intValue) String() string {
	return fmt.Sprintf("%d", *v.inner.target)
}
