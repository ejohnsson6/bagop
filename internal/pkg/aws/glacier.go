package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	l "github.com/swexbe/bagop/internal/pkg/logging"
	"github.com/swexbe/bagop/internal/pkg/utility"
)

// UploadFile uploads a file to an AWS glacier vault
// The vault is determined by the S3_VAULT_NAME env variable
func UploadFile(fileLocation string, timestamp string) (*glacier.ArchiveCreationOutput, error) {

	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	file, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create Glacier client in default region
	svc := glacier.New(sess)

	vaultName := os.Getenv(utility.ENVVault)

	// start snippet

	description := fmt.Sprintf("Archive created by bagop %s", timestamp)

	l.Logger.Infof("Uploading file %s to %s", file.Name(), vaultName)

	result, err := svc.UploadArchive(&glacier.UploadArchiveInput{
		VaultName:          &vaultName,
		ArchiveDescription: &description,
		Body:               file, // 2 MB buffer
	})
	if err != nil {
		return nil, err
	}

	l.Logger.Infof("Upload succeeded with id %s", *result.ArchiveId)
	return result, nil
}
