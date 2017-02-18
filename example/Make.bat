set CURR_DIR=%cd%

: Build generator
cd ..\..\..\..\..
set GOPATH=%cd%
go build -o %CURR_DIR%\sprotogen.exe github.com/davyxu/gosproto/sprotogen
cd %CURR_DIR%

: Generate go source file by sproto
sprotogen --type=go --out=addressbook.go --gopackage=example addressbook.sp

: Convert to standard sproto file
sprotogen --type=sproto --out=addressbook.sproto addressbook.sp

: Generate c# source file by sproto
sprotogen --type=csharp --out=addressbook.cs --gopackage=example addressbook.sp