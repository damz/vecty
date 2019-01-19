// +build !reflectless

package vecty

import (
	"reflect"
)

func panicInvalidType(v interface{}) {
	panic("vecty: internal error (unexpected type " + reflect.TypeOf(v).String() + ")")
}

type ComponentCopier interface {
}

// copyProps copies all struct fields from src to dst that are tagged with
// `vecty:"prop"`.
//
// If src and dst are different types or non-pointers, copyProps panics.
func copyProps(src, dst Component) {
	if src == dst {
		return
	}
	s := reflect.ValueOf(src)
	d := reflect.ValueOf(dst)
	if s.Type() != d.Type() {
		panic("vecty: internal error (attempted to copy properties of incompatible structs)")
	}
	if s.Kind() != reflect.Ptr || d.Kind() != reflect.Ptr {
		panic("vecty: internal error (attempted to copy properties of non-pointer)")
	}
	for i := 0; i < s.Elem().NumField(); i++ {
		sf := s.Elem().Field(i)
		if s.Elem().Type().Field(i).Tag.Get("vecty") == "prop" {
			df := d.Elem().Field(i)
			if sf.Type() != df.Type() {
				panic("vecty: internal error (should never be possible, struct types are identical)")
			}
			df.Set(sf)
		}
	}
}

// copyComponent makes a copy of the given component.
func copyComponent(c Component) Component {
	if c == nil {
		panic("vecty: internal error (cannot copy nil Component)")
	}

	// If the Component implements the Copier interface, then use that to
	// perform the copy.
	if copier, ok := c.(Copier); ok {
		cpy := copier.Copy()
		if cpy == c {
			panic("vecty: Component.Copy illegally returned an identical *MyComponent pointer")
		}
		return cpy
	}

	// Component does not implement the Copier interface, so perform a shallow
	// copy.
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("vecty: Component must be pointer to struct, found " + reflect.TypeOf(c).String())
	}
	cpy := reflect.New(v.Elem().Type())
	cpy.Elem().Set(v.Elem())
	return cpy.Interface().(Component)
}

// sameType returns whether first and second ComponentOrHTML are of the same
// underlying type.
func sameType(first, second ComponentOrHTML) bool {
	return reflect.TypeOf(first) == reflect.TypeOf(second)
}
