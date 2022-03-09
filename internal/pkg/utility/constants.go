package utility

const (
	// BackupLocation is the base location where are stored locally
	// This is removed every run
	BackupLocation = "/var/bagop/backups"
	// BackupDBLocation is the place where databases are dumped
	BackupDBLocation = BackupLocation + "/db"
	// ExtraLocation is where extra data to be backed up is pulled from
	ExtraLocation = "/extra"
	// ArchiveIDLocation is where the archive IDs are stored for persistence
	ArchiveIDLocation = "/var/bagop/ids.log"
	// Version is the version of bagop
	Version = "1.1.2"
	// ENVVault is the Environment Variable used for specifying the Glacier Vault
	ENVVault = "BAGOP_VAULT_NAME"
	// ENVCron is the Environment Variable used for specifying the regular backup schedule
	ENVCron = "CRON"
	// ENVLTCron is the Environment Variable used for specifying the long-term backup schedule
	ENVLTCron = "LT_CRON"
	// ENVTTL is the Environment Variable used for specifying the time to live for regular backups
	ENVTTL = "BAGOP_TTL"
	// ENVLTTTL is the Environment Variable used for specifying the time to live for long-term backups
	ENVLTTTL = "BAGOP_LT_TTL"
	// ENVColor is the Environment Variable used for making forcing bagop to use color output
	ENVColor = "BAGOP_COLOR"
	// ENVVerbose is the Environment Variable used for making bagop's output more verbose
	ENVVerbose = "BAGOP_VERBOSE"
)
