package main

import (
	"bytes"
	"fmt"
	"go/token"
	"sort"

	"github.com/davyxu/gosproto/meta"
)

const csharpCodeTemplate = `// Generated by github.com/davyxu/gosproto/sprotogen
// DO NOT EDIT!
using System;
using Sproto;
using System.Collections.Generic;

namespace {{.PackageName}}
{
{{range $a, $enumobj := .Enums}}
	public enum {{.Name}} {
		{{range .CSFields}}
		{{.Name}} = {{.Tag}},
		{{end}}
	}
{{end}}

{{range .Structs}}
	public class {{.Name}} : SprotoTypeBase {
		private static int max_field_count = {{.MaxFieldCount}};
		
		{{range .CSFields}}
		private {{.CSTypeString}} _{{.Name}}; // tag {{.Tag}}
		public {{.CSTypeString}} {{.Name}} {
			get{ return _{{.Name}}; }
			set{ base.has_field.set_field({{.FieldIndex}},true); _{{.Name}} = value; }
		}
		public bool Has{{.UpperFirstName}}{
			get { return base.has_field.has_field({{.FieldIndex}}); }
		}
		{{end}}
		
		public {{.Name}}() : base(max_field_count) {}
		
		public {{.Name}}(byte[] buffer) : base(max_field_count, buffer) {
			this.decode ();
		}
		
		protected override void decode () {
			int tag = -1;
			while (-1 != (tag = base.deserialize.read_tag ())) {
				switch (tag) {
				{{range .CSFields}}
				case {{.Tag}}:
					this.{{.Name}} = base.deserialize.{{.ReadFunc}}{{.CSTemplate}}({{.LamdaFunc}});
					break;
				{{end}}
				default:
					base.deserialize.read_unknow_data ();
					break;
				}
			}
		}
		
		public override int encode (SprotoStream stream) {
			base.serialize.open (stream);

			{{range .CSFields}}
			if (base.has_field.has_field ({{.FieldIndex}})) {
				base.serialize.{{.WriteFunc}}(this.{{.Name}}, {{.Tag}});
			}
			{{end}}

			return base.serialize.close ();
		}
	}
{{end}}

}
`

type csharpFieldModel struct {
	*meta.FieldDescriptor
	FieldIndex int
}

func (self *csharpFieldModel) UpperFirstName() string {
	return publicFieldName(self.Name)
}

func (self *csharpFieldModel) FieldName() string {
	pname := publicFieldName(self.Name)

	// 碰到关键字在尾部加_
	if token.Lookup(pname).IsKeyword() {
		return pname + "_"
	}

	return pname
}

func (self *csharpFieldModel) CSTemplate() string {

	var buf bytes.Buffer

	var needTemplate bool

	switch self.Type {
	case meta.FieldType_Struct,
		meta.FieldType_Enum:
		needTemplate = true
	}

	if needTemplate {
		buf.WriteString("<")
	}

	if self.MainIndex != nil {
		buf.WriteString(csharpTypeName(self.MainIndex))
		buf.WriteString(",")
	}

	if needTemplate {
		buf.WriteString(self.Complex.Name)
		buf.WriteString(">")
	}

	return buf.String()
}

func (self *csharpFieldModel) LamdaFunc() string {
	if self.MainIndex == nil {
		return ""
	}

	return fmt.Sprintf("v => v.%s", self.MainIndex.Name)
}

func (self *csharpFieldModel) WriteFunc() string {

	return "write_" + self.serializer()
}

func (self *csharpFieldModel) ReadFunc() string {

	funcName := "read_"

	if self.Repeatd {

		if self.MainIndex != nil {
			return funcName + "map"
		} else {
			return funcName + self.serializer() + "_list"
		}

	}

	return funcName + self.serializer()
}

func (self *csharpFieldModel) serializer() string {

	var baseName string

	switch self.Type {
	case meta.FieldType_Integer:
		baseName = "integer"
	case meta.FieldType_Int32:
		baseName = "int32"
	case meta.FieldType_Int64:
		baseName = "int64"
	case meta.FieldType_UInt32:
		baseName = "uint32"
	case meta.FieldType_UInt64:
		baseName = "uint64"
	case meta.FieldType_String:
		baseName = "string"
	case meta.FieldType_Bool:
		baseName = "boolean"
	case meta.FieldType_Struct:
		baseName = "obj"
	case meta.FieldType_Enum:
		baseName = "enum"
	default:
		baseName = "unknown"
	}

	return baseName
}

func (self *csharpFieldModel) CSTypeName() string {
	// 字段类型映射go的类型
	return csharpTypeName(self.FieldDescriptor)
}

func csharpTypeName(fd *meta.FieldDescriptor) string {
	switch fd.Type {
	case meta.FieldType_Integer:
		return "Int64"
	case meta.FieldType_Int32:
		return "Int32"
	case meta.FieldType_Int64:
		return "Int64"
	case meta.FieldType_UInt32:
		return "UInt32"
	case meta.FieldType_UInt64:
		return "UInt64"
	case meta.FieldType_String:
		return "string"
	case meta.FieldType_Bool:
		return "bool"
	case meta.FieldType_Struct,
		meta.FieldType_Enum:
		return fd.Complex.Name
	}
	return "unknown"
}

func (self *csharpFieldModel) CSTypeString() string {

	var b bytes.Buffer
	if self.Repeatd {

		if self.MainIndex != nil {
			b.WriteString("Dictionary<")

			b.WriteString(csharpTypeName(self.MainIndex))

			b.WriteString(",")

		} else {
			b.WriteString("List<")
		}

	}

	b.WriteString(self.CSTypeName())

	if self.Repeatd {
		b.WriteString(">")
	}

	return b.String()
}

type csharpStructModel struct {
	*meta.Descriptor

	CSFields []csharpFieldModel
}

func (self *csharpStructModel) FieldCount() int {
	return len(self.CSFields)
}

type csharpFileModel struct {
	*meta.FileDescriptor

	Structs []*csharpStructModel
	Enums   []*csharpStructModel

	PackageName string
}

func (self *csharpFileModel) Len() int {
	return len(self.Structs)
}

func (self *csharpFileModel) Swap(i, j int) {
	self.Structs[i], self.Structs[j] = self.Structs[j], self.Structs[i]
}

func (self *csharpFileModel) Less(i, j int) bool {

	a := self.Structs[i]
	b := self.Structs[j]

	return a.Name < b.Name
}

func addCSStruct(descs []*meta.Descriptor, callback func(*csharpStructModel)) {

	for _, st := range descs {

		stModel := &csharpStructModel{
			Descriptor: st,
		}

		for index, fd := range st.Fields {

			fdModel := csharpFieldModel{
				FieldDescriptor: fd,
				FieldIndex:      index,
			}

			stModel.CSFields = append(stModel.CSFields, fdModel)

		}

		callback(stModel)

	}
}

func gen_csharp(fileD *meta.FileDescriptor, packageName, filename string) {

	fm := &csharpFileModel{
		FileDescriptor: fileD,
		PackageName:    packageName,
	}

	addCSStruct(fileD.Structs, func(stModel *csharpStructModel) {
		fm.Structs = append(fm.Structs, stModel)
	})

	addCSStruct(fileD.Enums, func(stModel *csharpStructModel) {
		fm.Enums = append(fm.Enums, stModel)
	})

	sort.Sort(fm)

	generateCode("sp->cs", csharpCodeTemplate, filename, fm, nil)

}
