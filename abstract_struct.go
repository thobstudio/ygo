package ynotgo

type AbstractStruct struct {
	id     *ID
	length uint32
}

func newAbstractStruct(id *ID, length uint32) *AbstractStruct {
	return &AbstractStruct{
		id:     id,
		length: length,
	}
}

func (as *AbstractStruct) Id() *ID {
	return as.id
}

func (as *AbstractStruct) Length() uint32 {
	return as.length
}

func (as *AbstractStruct) State() uint32 {
	return as.id.clock + as.length
}

func (as *AbstractStruct) MergeWith() (bool, error) {
	return false, nil
}

func (as *AbstractStruct) Integrate(tx *Transaction, offset uint32) error {
	panic("to be implemented by the composite type")
}

func (as *AbstractStruct) Write(encoder UpdateEncoder, offset uint32) error {
	panic("to be implemented by the composite type")
}

func (as *AbstractStruct) Deleted() bool {
	panic("to be implemented by the composite type")
}

func (as *AbstractStruct) GetMissing() (uint32, bool) {
	panic("to be implemented by the composite type")
}
