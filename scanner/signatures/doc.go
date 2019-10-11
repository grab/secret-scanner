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

var skippableExtensions = []string{".exe", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".psd", ".xcf", ".zip", ".tar.gz", ".ttf", ".lock"}
var skippablePathIndicators = []string{"node_modules/", "vendor/"}
var skippableTestContexts = []string{"test", "_spec", "fixture", "mock", "stub", "fake", "demo", "sample"}
