package ordmap

func newIntInt() NodeBuiltin[int, int] {
	return NewBuiltin[int, int]()
}
func newIntOrder() NodeBuiltin[int, *order] {
	return NewBuiltin[int, *order]()
}
func newCompInt() *Node[comp, int] {
	return New[comp, int]()
}
func newCompOrder() *Node[comp, *order] {
	return New[comp, *order]()
}
