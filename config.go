package log

const (
	// LLevel configuration path.
	//
	// value: <root>.log.level
	LLevel = "log.level"

	// File configuration path.
	//
	// value: <root>.log.file
	File = "log.file"

	// FileEnableFileLogging configuration path.
	//
	// value: <root>.log.file.enable
	FileEnableFileLogging = File + ".enable"

	// FileDirectory configuration path.
	//
	// value: <root>.log.file.directory
	FileDirectory = File + ".directory"

	// FileName configuration path.
	//
	// value: <root>.log.file.name
	FileName = File + ".name"

	// FileSize configuration path.
	//
	// value: <root>.log.file.size
	FileSize = File + ".size"

	// FileRetentionAge configuration path.
	//
	// value: <root>.log.file.retention.age
	FileRetentionAge = File + ".retention.age"

	// FileRetentionBackups configuration path.
	//
	// value: <root>.log.file.retention.backups
	FileRetentionBackups = File + ".retention.backups"
)
