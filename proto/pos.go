package proto

import (
	"path"
	"runtime"
	"strings"
)

var gopath = func() string {
	_, file, _, _ := runtime.Caller(0)
	parts := strings.Split(file, "/src/")
	if len(parts) == 2 {
		return parts[0] + "/src/"
	}
	return ""
}()

type FilePos struct {
	Pkg  string `json:"pkg"`
	File string `json:"file"`
	Func string `json:"func"`
	Line int    `json:"line"`
}

func getFilePos(skip int) FilePos {
	pc, fullPath, line, _ := runtime.Caller(skip + 1)
	pkg, file := pkgFile(fullPath)
	return FilePos{
		Pkg:  pkg,
		File: file,
		Func: funcName(pc, pkg),
		Line: line,
	}
}

func pkgFile(fullPath string) (pkg, file string) {
	relativePath := strings.TrimPrefix(fullPath, gopath)
	pkg, file = path.Split(relativePath)
	pkg = strings.TrimSuffix(pkg, "/")
	return
}

func funcName(pc uintptr, pkg string) string {
	fn := runtime.FuncForPC(pc)
	s := strings.TrimPrefix(fn.Name(), pkg)
	i := strings.Index(s, ".")
	if i >= 0 {
		s = s[i+1:]
	}
	return s
}
