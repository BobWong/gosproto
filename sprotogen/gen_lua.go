package main

import (
	"github.com/bobwong89757/gosproto/meta"
)

const luaCodeTemplate = `-- Generated by github.com/bobwong8975789757/gosproto/sprotogen
-- DO NOT EDIT!
{{if .EnumValueGroup}}
ResultToString = function ( result )
	if result == 0 then
		return "OK"
	end

	local str = ResultByID[result]
	if str == nil then
		return string.format("unknown result: %d", result )
	end

	return str
end

ResultByID = {
	{{range $a, $enumObj := .Enums}} {{if .IsResultEnum}} {{range .Fields}} {{if ne .TagNumber 0}}
	[{{.TagNumber}}] = "{{$enumObj.Name}}.{{.Name}}", {{end}} {{end}} {{end}} {{end}}
}

{{end}}

Enum = {
{{range $a, $enumObj := .Enums}}
	{{$enumObj.Name}} = { {{range .Fields}}
		{{.Name}} = {{.TagNumber}}, {{end}}
	},
	{{end}}
}

local sproto = {
	Schema = [[
{{range .Structs}}
.{{.Name}} {	{{range .StFields}}	
	{{.Name}} {{.TagNumber}} : {{.CompatibleTypeString}} {{end}}
}
{{end}}
	]],

	NameByID = { {{range .Structs}}
		[{{.MsgID}}] = "{{.Name}}",{{end}}
	},
	
	IDByName = {},

	ResetByID = { {{range .Structs}}
		[{{.MsgID}}] = function( obj ) -- {{.Name}}
			if obj == nil then return end {{range .StFields}}
			obj.{{.Name}} = {{.LuaDefaultValueString}} {{end}}
		end, {{end}}
	},
}

local t = sproto.IDByName
for k, v in pairs(sproto.NameByID) do
	t[v] = k
end

return sproto

`

func (self *fieldModel) LuaDefaultValueString() string {

	if self.Repeatd {
		return "nil"
	}

	switch self.Type {
	case meta.FieldType_Bool:
		return "false"
	case meta.FieldType_Int32,
		meta.FieldType_Int64,
		meta.FieldType_UInt32,
		meta.FieldType_UInt64,
		meta.FieldType_Integer,
		meta.FieldType_Float32,
		meta.FieldType_Float64,
		meta.FieldType_Enum:
		return "0"
	case meta.FieldType_String:
		return "\"\""
	case meta.FieldType_Struct,
		meta.FieldType_Bytes:
		return "nil"
	}

	return "unknown type" + self.Type.String()
}

func gen_lua(fm *fileModel, filename string) {

	addData(fm, "lua")

	generateCode("sp->lua", luaCodeTemplate, filename, fm, nil)

}
