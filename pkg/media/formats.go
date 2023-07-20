package media

import "github.com/Hello-Storage/hello-back/pkg/fs"

// Formats maps file formats to general media types.
var Formats = map[fs.Type]Type{
	fs.ArchiveZip:  Archive,
	fs.ArchiveRar:  Archive,
	fs.ArchiveTar:  Archive,
	fs.ArchiveGzip: Archive,

	fs.AudioMp3:  Audio,
	fs.AudioMpeg: Audio,
	fs.AudioWav:  Audio,
	fs.AudioOgg:  Audio,
	fs.AudioOpus: Audio,

	fs.ImageJpeg: Image,
	fs.ImagePng:  Image,
	fs.ImageGif:  Image,
	fs.ImageBmp:  Image,
	fs.ImageTiff: Image,
	fs.ImageSvg:  Image,

	fs.Text: Text,

	fs.VideoWebm:  Video,
	fs.VideoMpeg4: Video,
	fs.Video3gpp:  Video,
	fs.VideoMov:   Video,
	fs.VideoAvi:   Video,
	fs.VideoWmv:   Video,

	fs.UnknownType: Other,
}
