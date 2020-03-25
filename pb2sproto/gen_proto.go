package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/bobwong8975789757/pbmeta"
	pbprotos "github.com/bobwong8975789757/pbmeta/proto"
)

const codeTemplate = `# Generated by github.com/bobwong8975789757/gosproto/pb2sproto
# Source: {{.FileName}}

{{range .Enums}}
{{.Comment}}
enum {{.Name}} {
	{{range .Fields}}	
	{{.Name}} {{.Tag}}
	{{end}}
}
{{end}}

{{range .Structs}}
{{.Comment}}
.{{.Name}} {
	{{range .Fields}}	
	{{.Name}} {{.Tag}} : {{.TypeString}} {{.Comment}}
	{{end}}
}
{{end}}
`

type spEnumFieldModel struct {
	*pbmeta.EnumValueDescriptor
}

func (self *spEnumFieldModel) Tag() int32 {
	return self.Define.GetNumber()
}

type spFieldModel struct {
	*pbmeta.FieldDescriptor
}

func (self *spFieldModel) Tag() int32 {
	return self.Define.GetNumber()
}

func (self *spFieldModel) TypeString() (ret string) {

	if self.IsRepeated() {
		ret = "*"
	}

	switch self.Type() {
	case pbprotos.FieldDescriptorProto_TYPE_INT64:
		ret += "int64"
	case pbprotos.FieldDescriptorProto_TYPE_UINT64:
		ret += "uint64"
	case pbprotos.FieldDescriptorProto_TYPE_INT32,
		pbprotos.FieldDescriptorProto_TYPE_FLOAT, // 浮点数默认转为int32
		pbprotos.FieldDescriptorProto_TYPE_DOUBLE:
		ret += "int32"
	case pbprotos.FieldDescriptorProto_TYPE_UINT32:
		ret += "uint32"
	case pbprotos.FieldDescriptorProto_TYPE_BOOL:
		ret += "boolean"
	case pbprotos.FieldDescriptorProto_TYPE_STRING:
		ret += "string"
	case pbprotos.FieldDescriptorProto_TYPE_MESSAGE:
		ret += self.MessageDesc().Name()
	case pbprotos.FieldDescriptorProto_TYPE_ENUM:
		ret += self.EnumDesc().Name()
	}

	return
}

func addCommentSignAtEachLine(sign, comment string) string {

	if comment == "" {
		return ""
	}
	var out bytes.Buffer

	scanner := bufio.NewScanner(strings.NewReader(comment))

	var index int
	for scanner.Scan() {

		if index > 0 {
			out.WriteString("\n")
		}

		out.WriteString(sign)
		out.WriteString(" ")
		out.WriteString(scanner.Text())

		index++
	}

	return out.String()
}

func (self *spFieldModel) Comment() string {

	return addCommentSignAtEachLine("#", self.CommentMeta.TrailingComment())

}

type spStructModel struct {
	*pbmeta.Descriptor

	Fields []spFieldModel
}

func (self *spStructModel) Comment() string {

	return addCommentSignAtEachLine("#", self.CommentMeta.LeadingComment())
}

type spEnumModel struct {
	*pbmeta.EnumDescriptor

	Fields []spEnumFieldModel
}

func (self *spEnumModel) Comment() string {

	return addCommentSignAtEachLine("#", self.CommentMeta.LeadingComment())
}

type spFileModel struct {
	*pbmeta.FileDescriptor

	Structs []*spStructModel

	Enums []*spEnumModel
}

func gen_proto(fileD *pbmeta.FileDescriptor, outputDir string) {

	tpl, err := template.New("pb->sp").Parse(codeTemplate)
	if err != nil {
		fmt.Println("template error ", err.Error())
		os.Exit(1)
	}

	fm := &spFileModel{
		FileDescriptor: fileD,
	}

	for structIndex := 0; structIndex < fileD.MessageCount(); structIndex++ {
		st := fileD.Message(structIndex)

		stModel := &spStructModel{
			Descriptor: st,
		}

		for fieldIndex := 0; fieldIndex < st.FieldCount(); fieldIndex++ {
			fd := st.Field(fieldIndex)

			fdModel := spFieldModel{
				FieldDescriptor: fd,
			}

			stModel.Fields = append(stModel.Fields, fdModel)
		}

		fm.Structs = append(fm.Structs, stModel)
	}

	for enumIndex := 0; enumIndex < fileD.EnumCount(); enumIndex++ {
		st := fileD.Enum(enumIndex)

		stModel := &spEnumModel{
			EnumDescriptor: st,
		}

		for fieldIndex := 0; fieldIndex < st.ValueCount(); fieldIndex++ {
			fd := st.Value(fieldIndex)

			fdModel := spEnumFieldModel{
				EnumValueDescriptor: fd,
			}

			stModel.Fields = append(stModel.Fields, fdModel)
		}

		fm.Enums = append(fm.Enums, stModel)
	}

	var bf bytes.Buffer

	err = tpl.Execute(&bf, &fm)
	if err != nil {
		fmt.Println("template error ", err.Error())
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("format error ", err.Error())
	}

	final := path.Join(outputDir, changeExt(fileD.FileName(), ".sp"))

	if fileErr := ioutil.WriteFile(final, bf.Bytes(), 666); fileErr != nil {
		fmt.Println("write file error ", fileErr.Error())
		os.Exit(1)
	}
}

// newExt = .xxx
func changeExt(name, newExt string) string {
	ext := path.Ext(name)
	name = name[0 : len(name)-len(ext)]
	return name + newExt
}
