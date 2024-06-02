package constant

const (
	KiloBytes = 1024
	MegaBytes = 1024 * KiloBytes

	Image = "image"
	PDF   = "pdf"

	pngContentType                = "image/png"
	pdfContentType                = "application/pdf"
	SniffingFirstResponsibleBytes = 512
)

var (
	MaxFileSize = map[string]int{
		Image: 500 * KiloBytes,
		PDF:   2 * MegaBytes,
	}

	FileType = map[string]string{
		Image: pngContentType,
		PDF:   pdfContentType,
	}
)
