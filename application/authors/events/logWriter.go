package events

type logWriter struct {
	message string
}

func (t *logWriter) Write(p []byte) (n int, err error) {
	t.message = string(p)

	return len(p), nil
}
