package opc

// DefaultContentTypes maps file extensions to their default OPC content
// types, as defined by the OPC specification. These are used to determine
// whether a part's content type can be expressed as a Default mapping in
// [Content_Types].xml or requires an Override.
var DefaultContentTypes = []struct {
	Ext         string
	ContentType string
}{
	{"bin", CT_PML_PRINTER_SETTINGS},
	{"bin", CT_SML_PRINTER_SETTINGS},
	{"bin", CT_WML_PRINTER_SETTINGS},
	{"bmp", CT_BMP},
	{"emf", CT_X_EMF},
	{"fntdata", CT_X_FONTDATA},
	{"gif", CT_GIF},
	{"jpe", CT_JPEG},
	{"jpeg", CT_JPEG},
	{"jpg", CT_JPEG},
	{"png", CT_PNG},
	{"rels", CT_OPC_RELATIONSHIPS},
	{"tif", CT_TIFF},
	{"tiff", CT_TIFF},
	{"wdp", CT_MS_PHOTO},
	{"wmf", CT_X_WMF},
	{"xlsx", CT_SML_SHEET},
	{"xml", CT_XML},
}

// IsDefaultContentType reports whether the given extension and content type
// pair matches one of the default content type entries in the OPC spec.
func IsDefaultContentType(ext, contentType string) bool {
	for _, d := range DefaultContentTypes {
		if d.Ext == ext && d.ContentType == contentType {
			return true
		}
	}
	return false
}
