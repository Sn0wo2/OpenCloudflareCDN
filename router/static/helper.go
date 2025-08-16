package static

import (
	"os"
	"path/filepath"
	"strings"
)

func GetStaticFile(index, staticRoot, reqPath string) string {
	base := filepath.Join(staticRoot, strings.TrimPrefix(reqPath, "/"))
	defaultIndex := filepath.Join(staticRoot, index)

	if rel, err := filepath.Rel(staticRoot, base); err != nil || strings.HasPrefix(rel, "..") {
		return defaultIndex
	}

	if info, err := os.Stat(base); err == nil {
		switch mode := info.Mode(); {
		case mode.IsDir():
			idx := filepath.Join(base, index)
			if _, err = os.Stat(idx); err == nil {
				return idx
			}

			return defaultIndex

		case mode.IsRegular():
			return base
		}
	}

	return defaultIndex
}
