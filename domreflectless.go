// +build reflectless

package vecty

func panicInvalidType(v interface{}) {
	panic("vecty: internal error (unexpected type)")
}

type ComponentCopier interface {
	Copier
	CopyProps(src Component)
	SameType(other Component) bool
}

func copyProps(src, dst Component) {
	if src == dst {
		return
	}
	dst.CopyProps(src)
}

func copyComponent(c Component) Component {
	if c == nil {
		panic("vecty: internal error (cannot copy nil Component)")
	}

	cpy := c.Copy()
	if cpy == c {
		panic("vecty: Component.Copy illegally returned an identical *MyComponent pointer")
	}
	return cpy
}

func sameType(first, second ComponentOrHTML) bool {
	if _, ok := first.(*HTML); ok {
		_, ok := second.(*HTML)
		return ok
	}

	return first.(Component).SameType(second.(Component))
}
