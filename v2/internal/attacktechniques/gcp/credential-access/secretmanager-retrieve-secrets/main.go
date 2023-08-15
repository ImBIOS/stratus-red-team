/*
Service: secretmanager.googleapis.com
*/
package gcp

import (
	"context"
	_ "embed"
	"errors"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/iterator"

	"github.com/datadog/stratus-red-team/v2/pkg/stratus"
	"github.com/datadog/stratus-red-team/v2/pkg/stratus/mitreattack"
)

//go:embed main.tf
var tf []byte

func init() {
	const codeBlock = "```"
	stratus.GetRegistry().RegisterAttackTechnique(&stratus.AttackTechnique{
		ID:				"gcp.credential-access.secretmanager-retrieve-secrets",
		FriendlyName:	"Retrieve a High Number of Secrets Manager secrets",
		Description:`
Retrieves a high number of Secret Manager secrets

Warm-up:

- Create multiple secrets in Secret Manager.

Detonation:

- Enumerate the secrets
- Retrieve each secret value, one by one.
`,
		Detection:		"under construction",
		Platform:					stratus.GCP,
		IsIdempotent:				true,
		MitreAttackTactics:			[]mitreattack.Tactic{mitreattack.CredentialAccess},
		PrerequisitesTerraformCode:	tf,
		Detonate:					detonate,
	})
}

func detonate(params map[string]string, providers stratus.CloudProviders) error {
	ctx	:= context.Background()
	client, err := secretmanager.NewClient(ctx)
	defer client.Close()

	if err != nil {
		return errors.New("unable to create Secret Manager client: " + err.Error())
	}

	// build the request
	req := &secretmanagerpb.ListSecretsRequest { }

	// call the API
	it := client.ListSecrets(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return errors.New("unable to retrieve secret value: " + err.Error())
		}

		secret := resp.Name
		log.Println("Retrieving value of secret " + secret)
	}

	return nil
}