package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WorkDir see makefile
var WorkDir string

func main() {
	if WorkDir != "" {
		if err := os.Chdir(WorkDir); err != nil {
			panic(err)
		}

		fmt.Printf("Working directory set to: %s\n", WorkDir)
	}

	srcDir := filepath.Join("Frontend", "dist")
	dstDir := "static"

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		panic(err)
	}

	if _, err := os.Stat(dstDir); err == nil {
		fmt.Printf("Directory '%s' already exists. Overwrite? [Y/N] ", dstDir)

		var response string

		_, err := fmt.Scanln(&response)
		if err != nil || (strings.ToLower(response) != "y" && response != "") {
			panic("Aborting...")
		}

		if err := os.RemoveAll(dstDir); err != nil {
			panic(err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(dstDir), 0o750); err != nil {
		panic(err)
	}

	if err := os.Rename(srcDir, dstDir); err != nil {
		fmt.Printf("Cross-device move detected, using copy and remove...\n")

		if err := copyDir(srcDir, dstDir); err != nil {
			panic(err)
		}

		if err := os.RemoveAll(srcDir); err != nil {
			fmt.Printf("Warning: failed to remove source directory: %v\n", err)
		}
	}

	fmt.Printf("Successfully moved '%s' to '%s'\n", srcDir, dstDir)
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o750); err != nil {
		return fmt.Errorf("creating destination directory: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("reading source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return fmt.Errorf("copying directory %s: %w", entry.Name(), err)
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("copying file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src) //nolint:gosec
	if err != nil {
		return fmt.Errorf("reading source file: %w", err)
	}

	if err := os.WriteFile(dst, data, 0o600); err != nil {
		return fmt.Errorf("writing destination file: %w", err)
	}

	return nil
}
