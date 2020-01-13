package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/profanity"
	"github.com/blend/go-sdk/ref"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	flagRulesFile            *string
	flagInclude, flagExclude *[]string
	flagVerbose              *bool
	flagDebug                *bool
	flagFailFast             *bool
)

type config struct {
	profanity.Config `yaml:",inline"`
}

// Resolve resolves the config.
func (c *config) Resolve() error {
	return configutil.AnyError(
		configutil.SetBool(&c.Verbose, configutil.Bool(flagVerbose), configutil.Bool(c.Verbose), configutil.Bool(ref.Bool(false))),
		configutil.SetBool(&c.FailFast, configutil.Bool(flagDebug), configutil.Bool(c.Debug), configutil.Bool(ref.Bool(false))),
		configutil.SetBool(&c.FailFast, configutil.Bool(flagFailFast), configutil.Bool(c.FailFast), configutil.Bool(ref.Bool(false))),
		configutil.SetString(&c.RulesFile, configutil.String(*flagRulesFile), configutil.String(c.RulesFile), configutil.String(profanity.DefaultRulesFile)),
		configutil.SetStrings(&c.Include, configutil.Strings(*flagInclude), configutil.Strings(c.Include)),
		configutil.SetStrings(&c.Exclude, configutil.Strings(*flagExclude), configutil.Strings(c.Exclude)),
	)
}

func command() *cobra.Command {
	root := &cobra.Command{
		Use:   "profanity",
		Short: "Enforce profanity rules in a directory tree.",
		Long:  "Enforce profanity rules in a directory tree with inherited rules for each child directory.",
		Example: `
# Run a basic rules set
profanity --rules=PROFANITY_RULES

# Run a basic rules set with excluded files by glob
profanity --rules=PROFANITY_RULES --exclude="*_test.go"

# Run a basic rules set with included and excluded files by glob
profanity --rules=PROFANITY_RULES --include="*.go" --exclude="*_test.go"

# An example rule file looks like

""" yaml
CONTAINS_EXAMPLE: # id is meant to be a de-duplicating identifier
  description: "please use 'foo.Bar', not a concrete type reference" # description should include remediation steps
  contains:
  - "foo.Bar"

EXCLUDES_EXAMPLE:
  description: "please dont use HerpDerp"
  contains: [ "HerpDerp" ]
  excludeFiles: [ "*_test.go" ]

IMPORTS_EXAMPLE: # you can assert a go AST doesnt contains a given import by glob
  description: "dont include command stuff"
  importsContain: [ "github.com/blend/go-sdk/cmd/*" ]
"""

For more example rule files, see https://github.com/blend/go-sdk/tree/master/PROFANITY_RULES.yml
`,
	}

	flagRulesFile = root.Flags().StringP("rules", "r", profanity.DefaultRulesFile, "The rules file to search for in each valid directory")
	flagInclude = root.Flags().StringArrayP("include", "i", nil, "Files to include in glob matching format; can be a csv.")
	flagExclude = root.Flags().StringArrayP("exclude", "e", nil, "Files to exclude in glob matching format; can be a csv.")
	flagVerbose = root.Flags().BoolP("verbose", "v", false, "If we should show verbose output.")
	flagDebug = root.Flags().BoolP("debug", "d", false, "If we should show debug output.")
	flagFailFast = root.Flags().Bool("fail-fast", false, "If we should fail the run after the first error.")
	return root
}

func main() {
	cmd := command()
	cmd.Run = func(parent *cobra.Command, args []string) {
		var cfg config
		if err := cfg.Resolve(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		engine := profanity.New(profanity.OptConfig(cfg.Config))
		engine.Stdout = os.Stdout
		engine.Stderr = os.Stderr

		if err := engine.Process(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
			return
		}
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}
	os.Exit(0)
}
