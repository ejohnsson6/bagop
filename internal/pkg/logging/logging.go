package logging

import "github.com/sirupsen/logrus"

var (
	// Logger is the global logging instance
	Logger logrus.Logger = *logrus.New()
)
