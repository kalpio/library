package application

type TestLogWriter interface {
	Write(p []byte) (n int, err error)
	GetMessage() string
}

type logWriter struct {
	message string
}

func NewTestLogWriter() TestLogWriter {
	return &logWriter{}
}

func (t *logWriter) Write(p []byte) (n int, err error) {
	t.message = string(p)

	return len(p), nil
}

func (t *logWriter) GetMessage() string {
	return t.message
}
