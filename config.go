package log

const (
	// LogLevel configuration path.
	//
	// value: <root>.log.level
	LogLevel = "log.level"

	// LogFile configuration path.
	//
	// value: <root>.log.file
	LogFile = "log.file"

	// LogFileEnableFileLogging configuration path.
	//
	// value: <root>.log.file.enable
	LogFileEnableFileLogging = LogFile + ".enable"

	// LogFileDirectory configuration path.
	//
	// value: <root>.log.file.directory
	LogFileDirectory = LogFile + ".directory"

	// LogFileName configuration path.
	//
	// value: <root>.log.file.name
	LogFileName = LogFile + ".name"

	// LogFileSize configuration path.
	//
	// value: <root>.log.file.size
	LogFileSize = LogFile + ".size"

	// LogFileRetentionAge configuration path.
	//
	// value: <root>.log.file.retention.age
	LogFileRetentionAge = LogFile + ".retention.age"

	// LogFileRetentionBackups configuration path.
	//
	// value: <root>.log.file.retention.backups
	LogFileRetentionBackups = LogFile + ".retention.backups"
)
