package filemanager

import (
	"os"
)

type FileEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

func ListDirectory(path string) ([]FileEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileEntry
	for _, entry := range entries {
		files = append(files, FileEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		})
	}

	return files, nil
}
