package main

const sprotoCodeTemplate = `# Generated by github.com/bobwong8975789757/gosproto/sprotogen
# DO NOT EDIT!

{{range .Structs}}
.{{.Name}} {
	{{range .StFields}}
	{{.Name}} {{.TagNumber}} : {{.CompatibleTypeString}}
	{{end}}
}
{{end}}

`

func gen_sproto(fm *fileModel, filename string) {

	addData(fm, "sproto")

	generateCode("sp->sproto", sprotoCodeTemplate, filename, fm, nil)

}
