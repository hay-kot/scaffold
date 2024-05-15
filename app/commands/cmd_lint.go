package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func (ctrl *Controller) Lint(ctx *cli.Context) error {
	pfpath := ctx.Args().First()

	if pfpath == "" {
		return errors.New("no file provided")
	}

	file, err := os.OpenFile(pfpath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	pf, err := scaffold.ReadScaffoldFile(file)
	if err != nil {
		return err
	}

	errs := make([]error, 0)

	for _, q := range pf.Questions {
		// Template variables only allow alphanumeric characters and underscores.
		if !engine.IsValidIdentifier(q.Name) {
			errs = append(errs, fmt.Errorf("invalid template variable name: %s (only alphanumeric and underscore characters are supported)", q.Name))
		}

		types := [...]bool{
			q.Prompt.IsInput(),
			q.Prompt.IsConfirm(),
			q.Prompt.IsSelect(),
			q.Prompt.IsMultiSelect(),
			q.Prompt.IsInputLoop(),
			q.Prompt.IsTextInput(),
		}

		isAny := false
		for _, t := range types {
			if t {
				isAny = true
				break
			}
		}

		if !isAny {
			errs = append(errs, fmt.Errorf("unknown prompt type for question %s", q.Name))
		}
	}

	// Check Computed variable names are valid identifiers.
	for k := range pf.Computed {
		if !engine.IsValidIdentifier(k) {
			errs = append(errs, fmt.Errorf("invalid computed variable name: %s (only alphanumeric and underscore characters are supported)", k))
		}
	}

	// Validate skip patterns
	for _, skip := range pf.Skip {
		ok := doublestar.ValidatePathPattern(skip)
		if !ok {
			errs = append(errs, fmt.Errorf("invalid skip pattern: %s", skip))
		}
	}

	// Validate rewrites from fields exist
	scaffolddir := filepath.Dir(pfpath)
	for _, rewrite := range pf.Rewrites {
		abs, _ := filepath.Abs(filepath.Join(scaffolddir, rewrite.From))

		_, err := os.Stat(abs)
		if err != nil {
			errs = append(errs, fmt.Errorf("rewrite from path does not exist: %s", rewrite.From))
		}
	}

	// Validate injectjons
	for _, injection := range pf.Inject {
		if injection.Mode != "" {
			if injection.Mode != "before" && injection.Mode != "after" {
				errs = append(errs, fmt.Errorf("invalid injection mode: %s", injection.Mode))
			}
		}
	}

	for _, err := range errs {
		log.Error().Err(err).Msg("")
	}

	return nil
}
