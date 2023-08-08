package app

import (
	"crypto/rand"
	"fmt"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/cryptoutil"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"github.com/urfave/cli/v2"
)

// ErrFailedToGenerateKey is returned when the key could not be generated.
const ErrFailedToGenerateKey xerrors.Error = "failed to generate API key"

// GenerateKeyAction is the action for the generate-key command.
func GenerateKeyAction(ctx *cli.Context) error {
	key, err := cryptoutil.GenerateAPIKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGenerateKey, err)
	}

	fmt.Fprintf(ctx.App.Writer, "%s\n", key)

	return nil
}
