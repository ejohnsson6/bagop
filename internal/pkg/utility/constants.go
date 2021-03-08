package utility

const (
	// BackupLocation is the base location where are stored locally
	// This is removed every run
	BackupLocation = "/var/bagop/backups"
	// BackupDBLocation is the place where databases are dumped
	BackupDBLocation = BackupLocation + "/db"
	// ExtraLocation is where extra data to be backed up is pulled from
	ExtraLocation = "/extra"
	// ArchiveIDLocation is where the archive IDs are stored for persistance
	ArchiveIDLocation = "/var/bagop/ids.log"
	// Version is the version of bagop
	Version = "1.1.0"
	// ENVVault is the Environment Variable used for specifying the Glacier Vault
	ENVVault = "BAGOP_VAULT_NAME"
)
