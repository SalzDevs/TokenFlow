package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"os"

	tree_sitter            "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tree_sitter_python     "github.com/tree-sitter/tree-sitter-python/bindings/go"
	tree_sitter_go         "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_java       "github.com/tree-sitter/tree-sitter-java/bindings/go"
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

func ParserWrapper(filePath string) (*tree_sitter.Tree, error) {
	lang, err := LangDetector(filePath)
	if err != nil {
		// If the file type is unsupported, we can choose to skip it without returning an error
		fmt.Printf("Skipping unsupported file type: %s\n", filePath)
		return nil, nil
	}
	
	parser := tree_sitter.NewParser()
	defer parser.Close()	
	
	switch lang {
	case "javascript":
		parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_javascript.Language()))
	case "python":
		parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_python.Language()))
	case "go":
		parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language()))
	case "java":
	 	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_java.Language()))
	default:
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	code, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	tree := parser.Parse(code, nil)
	return tree, nil
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
		tree, err := ParserWrapper(file)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", file, err)
			continue
		}
		if tree != nil {
			fmt.Printf("Successfully parsed file: %s\n", file)
			tree.Close() // Don't forget to close the tree after use
		}
	}
}
