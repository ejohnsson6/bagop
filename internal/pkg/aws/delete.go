package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
)

// DeleteArchive deletes an archive from S3 Glacier
// The vault is determined by the S3_VAULT_NAME env variable
func DeleteArchive(vaultName string, archiveID string) (*glacier.DeleteArchiveOutput, error) {

	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create Glacier client in default region
	svc := glacier.New(sess)

	result, err := svc.DeleteArchive(&glacier.DeleteArchiveInput{
		VaultName: &vaultName,
		ArchiveId: &archiveID,
	})
	if err != nil {
		return nil, err
	}

	return result, nil

}
