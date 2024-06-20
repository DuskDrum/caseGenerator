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
		log.Fatalf("failed to get working directory: %v", err)
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
		log.Fatalf("failed to read go.mod file: %v", err)
	}

	modFile, err := modfile.Parse(modPath, content, nil)
	if err != nil {
		log.Fatalf("failed to parse go.mod file: %v", err)
	}

	modulePath := modFile.Module.Mod.Path
	fmt.Printf("Module path: %s\n", modulePath)
	return modulePath
}
