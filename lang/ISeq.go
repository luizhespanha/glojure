package lang

type ISeq interface {
	IPersistentCollection

	First() interface{}
	Next() ISeq
	More() ISeq
	// Cons(interface{}) ISeq
}
