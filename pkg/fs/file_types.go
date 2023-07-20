package fs

// Supported file types.
const (
	ArchiveZip  Type = "zip"  // zip Archive
	ArchiveRar  Type = "rar"  // rar Archive
	ArchiveTar  Type = "tar"  // tar Archive
	ArchiveGzip Type = "gzip" // gzip Archive

	AudioMp3  Type = "mp3"  // mp3 Audio
	AudioMpeg Type = "mpeg" // mpeg Audio
	AudioWav  Type = "wav"  // wag Audio
	AudioOgg  Type = "ogg"  // ogg Audio
	AudioOpus Type = "opus" // opus Audio

	ImageJpeg Type = "jpeg" // jpeg Image
	ImagePng  Type = "png"  // png Image
	ImageGif  Type = "gif"  // gif Image
	ImageBmp  Type = "bmp"  // bmp Image
	ImageTiff Type = "tiff" // tiff Image
	ImageSvg  Type = "svg"  // svg Image

	Text Type = "txt" // Text

	VideoWebm  Type = "webm"  // webm Video
	VideoMpeg4 Type = "mpeg4" // mpeg4 Video
	Video3gpp  Type = "3gpp"  // 3gpp Video
	VideoMov   Type = "mov"   // mov Video
	VideoAvi   Type = "avi"   // avi Video
	VideoWmv   Type = "wmv"   // wmv Video

	UnknownType Type = "" // Unknown file
)
