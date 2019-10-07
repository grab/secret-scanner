package signatures

const (
	TypeSimple  = "simple"
	TypePattern = "pattern"

	PartExtension = "extension"
	PartFilename  = "filename"
	PartPath      = "path"
	PartContent   = "content"
)

var skippableExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".psd", ".xcf"}
var skippablePathIndicators = []string{"node_modules/", "vendor/bundle", "vendor/cache"}
