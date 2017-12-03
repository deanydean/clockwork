package utils

import "container/list"
import "io/ioutil"
import "log"
import "os"

// FilterFunc will return true if the provided string matches the filter, false
// if not
type FilterFunc func(string) bool

// PathExists return true if the provided path exists and false if not
func PathExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

// FilesInDir gets a List of files in the provided directory that match the
// provided FilterFunc
func FilesInDir(dir string, filter FilterFunc) *list.List {
	// Read the files in the directory
	files, err := ioutil.ReadDir(dir)

	result := list.New()

	// If there's an error, just result nothing
	if err != nil {
		log.Fatal(err)
		return result
	}

	// Return the files that match the filter
	for _, f := range files {
		if filter(f.Name()) {
			result.PushBack(f)
		}
	}
	return result
}

// GetFileAsString gets file contents as a string
func GetFileAsString(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	str := string(bytes)
	return str, nil
}
