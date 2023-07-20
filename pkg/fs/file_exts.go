package fs

import (
	"path/filepath"
	"strings"
)

// FileExtensions maps file extensions to standard formats
type FileExtensions map[string]Type

// Extensions contains the filename extensions of file formats known to PhotoPrism.
var Extensions = FileExtensions{
	".zip":  ArchiveZip,
	".rar":  ArchiveRar,
	".tar":  ArchiveTar,
	".gzip": ArchiveGzip,

	".mp3":  AudioMp3,
	".mpeg": AudioMpeg,
	".wav":  AudioWav,
	".ogg":  AudioOgg,
	".opus": AudioOpus,

	".jpeg": ImageJpeg,
	".png":  ImagePng,
	".gif":  ImageGif,
	".bmp":  ImageBmp,
	".tiff": ImageTiff,
	".svg":  ImageSvg,

	".txt": Text,

	".webm":  VideoWebm,
	".mpeg4": VideoMpeg4,
	".3gpp":  Video3gpp,
	".mov":   VideoMov,
	".avi":   VideoAvi,
	".wmv":   VideoWmv,
}

// Known tests if the file extension is known (supported).
func (m FileExtensions) Known(name string) bool {
	if name == "" {
		return false
	}

	ext := strings.ToLower(filepath.Ext(name))

	if ext == "" {
		return false
	}

	if _, ok := m[ext]; ok {
		return true
	}

	return false
}
