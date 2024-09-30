// SPDX-License-Identifier: MPL-2.0

package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*ReadCommand)(nil)
	_ cli.CommandAutocomplete = (*ReadCommand)(nil)
)

type ReadCommand struct {
	*BaseCommand

	testStdin io.Reader // for tests
}

func (c *ReadCommand) Synopsis() string {
	return "Read data and retrieves secrets"
}

func (c *ReadCommand) Help() string {
	helpText := `
Usage: bao read [options] PATH

  Reads data from Vault at the given path. This can be used to read secrets,
  generate dynamic credentials, get configuration details, and more.

  Read a secret from the static secrets engine:

      $ bao read secret/my-secret

  For a full list of examples and paths, please see the documentation that
  corresponds to the secrets engine in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ReadCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
}

func (c *ReadCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *ReadCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ReadCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args, ParseOptionAllowRawFormat(true)); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// client.ReadRaw* methods require a manual timeout override
	ctx, cancel := context.WithTimeout(context.Background(), client.ClientTimeout())
	defer cancel()

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}
	if c.flagNonInteractive {
		stdin = bytes.NewReader(nil)
	}

	path := sanitizePath(args[0])

	data, err := parseArgsDataStringLists(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	if Format(c.UI) != "raw" {
		secret, err := client.Logical().ReadWithDataWithContext(ctx, path, data)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error reading %s: %s", path, err))
			return 2
		}
		if secret == nil {
			c.UI.Error(fmt.Sprintf("No value found at %s", path))
			return 2
		}

		if c.flagField != "" {
			return PrintRawField(c.UI, secret, c.flagField)
		}

		return OutputSecret(c.UI, secret)
	}

	resp, err := client.Logical().ReadRawWithDataWithContext(ctx, path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading: %s: %s", path, err))
		return 2
	}
	if resp == nil || resp.Body == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return 2
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading: %s: %s", path, err))
		return 2
	}

	return OutputData(c.UI, contents)
}
