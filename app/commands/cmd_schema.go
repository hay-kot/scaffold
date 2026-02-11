package commands

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

//go:embed schema/scaffold.schema.json
var scaffoldSchema []byte

//go:embed schema/scaffoldrc.schema.json
var scaffoldrcSchema []byte

// FlagsSchema contains flags for the schema command
type FlagsSchema struct {
	Type string // "scaffold" or "scaffoldrc"
}

func (ctrl *Controller) Schema(_ context.Context, c *cli.Command) error {
	schemaType := c.String("type")

	var schema []byte
	switch schemaType {
	case "scaffold":
		schema = scaffoldSchema
	case "scaffoldrc":
		schema = scaffoldrcSchema
	default:
		return fmt.Errorf("unknown schema type: %s (use 'scaffold' or 'scaffoldrc')", schemaType)
	}

	_, err := os.Stdout.Write(schema)
	return err
}
