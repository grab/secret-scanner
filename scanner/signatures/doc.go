package signatures

const (
	// TypeSimple ...
	TypeSimple = "simple"

	// TypePattern ...
	TypePattern = "pattern"

	// PartExtension ...
	PartExtension = "extension"

	// PartFilename ...
	PartFilename = "filename"

	// PartPath ...
	PartPath = "path"

	// PartContent ...
	PartContent = "content"
)

var skippableExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".psd", ".xcf"}
var skippablePathIndicators = []string{"node_modules/", "vendor/bundle", "vendor/cache"}
