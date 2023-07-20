package media

const (
	Archive    Type = "archive"    // (.ZIP, .RAR, tar, gzip)
	Audio      Type = "audio"      // (MP3, MPEG, WAV, .ogg, .opus)
	Image      Type = "image"      // (.JPEG, .PNG, .GIF, .BMP, .TIFF, .SVG)
	MarkupCode Type = "markupcode" // (.CSS, .HTML, .PHP, .C, .CPP, .H, .HPP, .JS, .java, .py)
	Text       Type = "text"       // (.TXT)
	Video      Type = "video"      // (WebM, .MPEG4, .3GPP, .MOV, .AVI, .MPEGPS, .WMV, .FLV, .ogg)

	Other   Type = "other"
	Unknown Type = ""
)
