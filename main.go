package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

func LangDetector(filePath string) (string, error) {
	ext := filepath.Ext(filePath)
	
	switch ext {
	case ".js":
		return "javascript", nil
	case ".ts":
		return "typescript", nil
	case ".py":
		return "python", nil
	case ".go":
		return "go", nil
	case ".java":
		return "java", nil
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

func FileWalker(dir string) ([]string, error) {
    var files []string

    err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            files = append(files, path)
        }

				// Skip hidden directories
				if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
					return filepath.SkipDir
				}

        return nil
    })

    if err != nil {
        return nil, err
    }

    return files, nil
}

func main(){

	code := []byte("const foo = 1 + 2")

	parser := tree_sitter.NewParser()
	defer parser.Close()
	

	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_javascript.Language()))

	tree := parser.Parse(code,nil)
	defer tree.Close()
	
	files, err := FileWalker("./")
	if err != nil {
		fmt.Println("Error walking the directory:", err)
		return
	}

	for _, file := range files {
		fmt.Println(file)
	}

}
