package meta

import "bytes"

type FileDescriptor struct {
	Structs []*Descriptor

	StructByName map[string]*Descriptor

	Enums []*Descriptor

	EnumByName map[string]*Descriptor

	unknownFields []*lazyField
}

func (self *FileDescriptor) resolveAll() error {

	for _, v := range self.unknownFields {
		if _, err := v.resolve(2); err != nil {
			return err
		}
	}

	return nil
}

func (self *FileDescriptor) String() string {

	var bf bytes.Buffer

	bf.WriteString("Enum:")
	for _, st := range self.Enums {
		bf.WriteString(st.String())
		bf.WriteString("\n")
	}

	bf.WriteString("\n")

	bf.WriteString("Structs:")
	for _, st := range self.Structs {
		bf.WriteString(st.String())
		bf.WriteString("\n")
	}

	return bf.String()
}

func (self *FileDescriptor) NameExists(name string) bool {
	if _, ok := self.StructByName[name]; ok {
		return true
	}

	if _, ok := self.EnumByName[name]; ok {
		return true
	}

	return false
}

func (self *FileDescriptor) addStruct(d *Descriptor) {
	self.Structs = append(self.Structs, d)
	self.StructByName[d.Name] = d
}

func (self *FileDescriptor) addEnum(d *Descriptor) {
	self.Enums = append(self.Enums, d)
	self.EnumByName[d.Name] = d
}

func NewFileDescriptor() *FileDescriptor {

	return &FileDescriptor{
		StructByName: make(map[string]*Descriptor),
		EnumByName:   make(map[string]*Descriptor),
	}

}
