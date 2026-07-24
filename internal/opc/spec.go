package opc

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

func IsDefaultContentType(ext, contentType string) bool {
	for _, d := range DefaultContentTypes {
		if d.Ext == ext && d.ContentType == contentType {
			return true
		}
	}
	return false
}
