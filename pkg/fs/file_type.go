package fs

type Type string

// FileType returns the type associated with the specified filename,
// and UnknownType if it could not be matched.
func FileType(fileName string) Type {
	if t, found := Extensions[LowerExt(fileName)]; found {
		return t
	}

	return UnknownType
}
