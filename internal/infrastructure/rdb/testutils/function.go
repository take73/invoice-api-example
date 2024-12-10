package testutils

import (
	"runtime"
	"strings"
)

func GetFuncName() string {
	pc, _, _, ok := runtime.Caller(1) // Caller(1)で呼び出し元の関数を取得
	if !ok {
		return "unknown"
	}
	// フルパスの関数名を取得
	fullFuncName := runtime.FuncForPC(pc).Name()

	// 最後のドット以降を取得
	parts := strings.Split(fullFuncName, ".")
	funcName := parts[len(parts)-1]

	return funcName
}
