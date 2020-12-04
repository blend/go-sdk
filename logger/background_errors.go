package logger

// BackgroundErrors reads errors from a channel and logs them as errors.
//
// You should call this method with it's own goroutine:
//
//    go logger.BackgroundErrors(log, flushErrors)
func BackgroundErrors(log ErrorReceiver, errors <-chan error) {
	if !IsLoggerSet(log) {
		return
	}
	var err error
	for {
		err = <-errors
		if err != nil {
			log.Error(err)
		}
	}
}
