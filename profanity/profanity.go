package profanity

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/yaml"
)

// New creates a new profanity engine with a given set of config options.
func New(options ...ConfigOption) *Profanity {
	var cfg Config
	for _, option := range options {
		option(&cfg)
	}
	return &Profanity{
		Config: cfg,
	}
}

// Profanity parses rules from the filesystem and applies them to a given root path.
// Creating a full rules set.
type Profanity struct {
	Config Config
	Stdout io.Writer
	Stderr io.Writer
}

// Printf writes to the output stream.
func (p *Profanity) Printf(format string, args ...interface{}) {
	if p.Stdout != nil {
		fmt.Fprintf(p.Stdout, format, args...)
	}
}

// Errorf writes to the error output stream.
func (p *Profanity) Errorf(format string, args ...interface{}) {
	if p.Stdout != nil {
		fmt.Fprintf(p.Stderr, format, args...)
	}
}

// Process processes the profanity rules.
func (p *Profanity) Process() error {
	if p.Config.VerboseOrDefault() {
		p.Printf("using rules file: %s\n", p.Config.RulesFileOrDefault())
	}

	if p.Config.VerboseOrDefault() {
		if len(p.Config.Include) > 0 {
			p.Printf("using include filter: %s\n", p.Config.Include)
		}
		if len(p.Config.Exclude) > 0 {
			p.Printf("using exclude filter: %s\n", p.Config.Exclude)
		}
	}

	var didError bool

	// rule cache is shared between files and directories during the full walk.
	ruleCache := make(map[string]Rules)
	// make sure the root rules are initialized if they exist.
	if _, err := os.Stat("./" + p.Config.RulesFileOrDefault()); err == nil {
		_, err = p.RulesForPathOrCached(ruleCache, ".")
		if err != nil {
			return err
		}
	}

	var fileBase string
	if err := filepath.Walk(".", func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasSuffix(file, ".git") { // don't ever process git directories
			if p.Config.VerboseOrDefault() {
				p.Printf("%s ... skipping (is .git dir)\n", ansi.LightWhite(file))
			}
			return filepath.SkipDir
		}
		if info.IsDir() {
			if p.Config.VerboseOrDefault() {
				p.Printf("%s ... skipping (is dir)\n", ansi.LightWhite(file))
			}
			return nil
		}

		fileBase = filepath.Base(file)

		if len(p.Config.Include) > 0 {
			if matches := GlobAnyMatch(p.Config.Include, file); !matches {
				if p.Config.VerboseOrDefault() {
					p.Printf("%s ... skipping (does not match include filter)\n", ansi.LightWhite(file))
				}
				return nil
			}
		}

		if len(p.Config.Exclude) > 0 {
			if matches := GlobAnyMatch(p.Config.Exclude, file); matches {
				if p.Config.VerboseOrDefault() {
					p.Printf("%s ... skipping (matches exclude filter)\n", ansi.LightWhite(file))
				}
				return nil
			}
		}

		if matches, err := filepath.Match(p.Config.RulesFileOrDefault(), fileBase); err != nil {
			return ex.New(err)
		} else if matches {
			if p.Config.VerboseOrDefault() {
				p.Printf("%s ... skipping (is rules `%s` file)\n", ansi.LightWhite(file), p.Config.RulesFileOrDefault())
			}
			return nil
		}

		fullPath := filepath.Dir(file)
		rules, err := p.RulesForPathOrCached(ruleCache, fullPath)
		if err != nil {
			return err
		}

		contents, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		for _, rule := range rules {
			if matches := rule.ShouldInclude(file); !matches {
				if p.Config.VerboseOrDefault() {
					p.Printf("%s ... skipping rule %s (fails include)\n", ansi.LightWhite(file), rule.ID)
				}
				continue
			}

			if matches := rule.ShouldExclude(file); matches {
				if p.Config.VerboseOrDefault() {
					p.Printf("%s ... skipping rule %s (fails exclude)\n", ansi.LightWhite(file), rule.ID)
				}
				continue
			}

			if p.Config.VerboseOrDefault() {
				p.Printf("%s ... checking rule %s\n", ansi.LightWhite(file), rule.ID)
			}
			if res := rule.Apply(file, contents); !res.OK {
				didError = true

				// check if there was an error with the rule ...
				if res.Err != nil {
					return res.Err
				}

				// handle the failure
				failure := res.Failure(rule)
				p.Errorf("%v\n", failure)
				if p.Config.FailFastOrDefault() {
					return failure
				}
			}
		}

		if p.Config.VerboseOrDefault() {
			p.Printf("%s ... %s!\n", ansi.LightWhite(file), ansi.Green("ok"))
		}

		return nil
	}); err != nil {
		return err
	}
	if didError {
		p.Printf("profanity %s!\n", ansi.Red("failed"))
		return ErrFailure
	}
	p.Printf("profanity %s!\n", ansi.Green("ok"))
	return nil
}

// RulesForPathOrCached returns rules cached or rules from disk.
// It prevents re-reading the full rules set for each file in a path.
func (p *Profanity) RulesForPathOrCached(packageRules map[string]Rules, path string) (Rules, error) {
	if rules, ok := packageRules[path]; ok {
		return rules, nil
	}

	rules, err := p.RulesForPath(packageRules, path)
	if err != nil {
		return nil, err
	}

	packageRules[path] = rules
	return rules, nil
}

// RulesForPath adds rules in a given path and child paths to an existing rule set.
// `workingSet` are the current working rules keyed on the path they
// came from, including '.' for the root rules.
func (p *Profanity) RulesForPath(workingSet map[string]Rules, path string) (Rules, error) {
	pathRules, err := p.ReadRules(path)
	if err != nil {
		return nil, err
	}

	for key, workingRules := range workingSet {
		if strings.HasPrefix(path, key) && key != path {
			if p.Config.VerboseOrDefault() {
				p.Printf("%s including inherited rules from %s", ansi.LightWhite(path), ansi.LightWhite(key))
			}
			pathRules = MergeRules(workingRules, pathRules)
		}
	}

	// always include rules from "." if they were set
	if rootRules, hasRootRules := workingSet[Root]; hasRootRules && path != Root {
		pathRules = MergeRules(rootRules, pathRules)
	}

	return pathRules, nil
}

// ReadRules reads rules at a given directory path.
// Path is meant to be the slash terminated dir, which will have the configured rule path appended to it.
func (p *Profanity) ReadRules(path string) (Rules, error) {
	if p.Config.DebugOrDefault() {
		p.Printf("checking for profanity file: %s/%s", ansi.LightWhite(path), p.Config.RulesFileOrDefault())
	}
	profanityPath := filepath.Join(path, p.Config.RulesFileOrDefault())
	if _, err := os.Stat(profanityPath); err != nil {
		if p.Config.VerboseOrDefault() {
			p.Printf("%s/ local rules file not found %s\n", ansi.LightWhite(path), p.Config.RulesFileOrDefault())
		}
		return nil, nil
	}
	rules, err := p.RulesFromPath(profanityPath)
	if err != nil {
		if p.Config.DebugOrDefault() {
			p.Errorf("error reading profanity file: %s/%s %v", ansi.LightWhite(path), p.Config.RulesFileOrDefault(), err)
		}
		return nil, err
	}
	return rules, nil
}

// RulesFromPath reads rules from a path
func (p *Profanity) RulesFromPath(path string) (rules Rules, err error) {
	contents, readErr := os.Open(path)
	if readErr != nil {
		err = ex.New(readErr, ex.OptMessagef("file: %s", path))
		return
	}
	defer contents.Close()
	rules, err = p.RulesFromReader(path, contents)
	return
}

// RulesFromReader reads rules from a reader.
func (p *Profanity) RulesFromReader(path string, reader io.Reader) (rules Rules, err error) {
	var fileRules Rules
	yamlErr := yaml.NewDecoder(reader).Decode(&fileRules)
	if yamlErr != nil {
		err = ex.New("cannot unmarshal rules file", ex.OptMessagef("file: %s", path), ex.OptInnerClass(yamlErr))
		return
	}
	rules = make(Rules)
	for id, fileRule := range fileRules {
		rule := fileRule
		rule.ID = id
		rule.File = path
		rules[id] = rule
	}
	return
}
