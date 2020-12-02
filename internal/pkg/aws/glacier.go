package aws

import (
	"bytes"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
)

func Test() {
	// Initialize a session that the SDK uses to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and configuration from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create Glacier client in default region
	svc := glacier.New(sess)

	// start snippet
	vaultName := "YOUR_VAULT_NAME"

	result, err := svc.UploadArchive(&glacier.UploadArchiveInput{
		VaultName: &vaultName,
		Body:      bytes.NewReader(make([]byte, 2*1024*1024)), // 2 MB buffer
	})
	if err != nil {
		log.Println("Error uploading archive.", err)
		return
	}

	log.Println("Uploaded to archive", *result.ArchiveId)
}
