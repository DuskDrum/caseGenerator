package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// findModFile 获取mod文件地址
func findModFile() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}
	for {
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return modPath
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	log.Fatal("go.mod file not found")
	return ""
}

// GetModulePath 获取系统的mod名称
func GetModulePath() string {
	modPath := findModFile()
	content, err := os.ReadFile(modPath)
	if err != nil {
		panic("failed to read go.mod file: " + err.Error())
	}

	modFile, err := modfile.Parse(modPath, content, nil)
	if err != nil {
		panic("failed to parse go.mod file: " + err.Error())
	}

	modulePath := modFile.Module.Mod.Path
	fmt.Printf("Module path: %s\n", modulePath)
	return modulePath
}
