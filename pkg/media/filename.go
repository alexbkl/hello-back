package media

import "github.com/Hello-Storage/hello-back/pkg/fs"

// FromName returns the content type matching the file extension.
func FromName(fileName string) Type {
	if fileName == "" {
		return Unknown
	}

	// Find media type based on the file type.
	if result, found := Formats[fs.FileType(fileName)]; found {
		return result
	}

	// Default.
	return Other
}
