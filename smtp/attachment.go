package smtp

// TODO 附件
type attachment struct {
	filename    string
	contentType string
	data        []byte
}
