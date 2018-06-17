package main

func LogError(v ...interface{}) string {
	if log == nil {
		return ""
	}

	return log.Error(v...)
}

func LogWarning(v ...interface{}) string {
	if log == nil {
		return ""
	}

	return log.Warning(v...)
}

func LogInfo(v ...interface{}) string {
	if log == nil {
		return ""
	}

	return log.Info(v...)
}

func LogTrace(v ...interface{}) string {
	if log == nil {
		return ""
	}

	return log.Trace(v...)
}

func LogDebug(v ...interface{}) string {
	if log == nil {
		return ""
	}

	return log.Debug(v...)
}
