package ynotgo

type SharedStruct interface {
	Id() *ID
	Length() uint32
	MergeWith(right SharedStruct) (bool, error)
	Integrate(tx *Transaction, offset uint32) error
	Deleted() bool
	Write(encoder UpdateEncoder, offset uint32) error
}
