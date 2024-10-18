package epub

import (
	"errors"
	"os"
	"runtime"
	"slices"
)

// Replaces invalid characters in directory name
// For *nix this simply replace nul characters
// For windows all control chars are invalid as well as some punctuation
func sanitizeDirName(dir string) string {
	var invalidDirChars []rune
	if runtime.GOOS == "windows" {

		invalidDirChars = []rune{
			0,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, '|',
		}
	} else {
		invalidDirChars = []rune{0}
	}

	dirRunes := []rune(dir)

	for i, r := range dirRunes {
		if slices.Contains(invalidDirChars, r) {
			dirRunes[i] = '_'
		}
	}

	return string(dirRunes)
}

// Replaces invalid characters in file name
// For *nix this simply replace nul characters and forward slashes
// For windows all control chars are invalid as well as some punctuation
func sanitizeFileName(filename string) string {
	var invalidFileChars []rune
	if runtime.GOOS == "windows" {
		invalidFileChars = []rune{
			0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31,
			':', '*', '?', '\\', '/', '"', '<', '>', '|',
		}
	} else {
		invalidFileChars = []rune{0, '/'}
	}

	filenameRunes := []rune(filename)

	for i, r := range filenameRunes {
		if slices.Contains(invalidFileChars, r) {
			filenameRunes[i] = '_'
		}
	}
	return string(filenameRunes)
}

func fileExists(filepath string) bool {
	_, error := os.Stat(filepath)
	return !errors.Is(error, os.ErrNotExist)
}
