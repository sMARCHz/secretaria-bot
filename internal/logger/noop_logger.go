package logger

// noopLogger is a minimal Logger implementation that does nothing.
// It's used as the default to make it safe to call logger.* helpers
// before a concrete logger is initialized (e.g. in tests).
type noopLogger struct{}

func (n *noopLogger) Debug(args ...interface{}) {}
func (n *noopLogger) Info(args ...interface{})  {}
func (n *noopLogger) Warn(args ...interface{})  {}
func (n *noopLogger) Error(args ...interface{}) {}
func (n *noopLogger) Fatal(args ...interface{}) {}

func (n *noopLogger) Debugf(format string, args ...interface{}) {}
func (n *noopLogger) Infof(format string, args ...interface{})  {}
func (n *noopLogger) Warnf(format string, args ...interface{})  {}
func (n *noopLogger) Errorf(format string, args ...interface{}) {}
func (n *noopLogger) Fatalf(format string, args ...interface{}) {}

// SetLogger allows replacing the global logger (useful in tests).
func SetLogger(l Logger) {
	if l == nil {
		logger = &noopLogger{}
		return
	}
	logger = l
}
