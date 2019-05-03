package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/cover"
)

// linker metadata block
// this block must be present
// it is used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	star             = "*"
	defaultFileFlags = 0644
)

var reportOutputPath = flag.String("output", "coverage.html", "the path to write the full html coverage report")
var update = flag.Bool("update", false, "if we should write the current coverage to `COVERAGE` files")
var enforce = flag.Bool("enforce", false, "if we should enforce coverage minimums defined in `COVERAGE` files")
var include = flag.String("include", "", "the include file filter in glob form, can be a csv.")
var exclude = flag.String("exclude", "", "the exclude file filter in glob form, can be a csv.")
var timeout = flag.String("timeout", "", "the timeout to pass to the package tests.")
var covermode = flag.String("covermode", "set", "the go test covermode.")
var coverprofile = flag.String("coverprofile", "coverage.cov", "the intermediate cover profile.")
var keepCoverageOut = flag.Bool("keep-coverage-out", false, "if we should keep coverage.out")
var v = flag.Bool("v", false, "show verbose output")

func verbose() bool {
	if v != nil && *v {
		return true
	}
	return false
}

func vf(format string, args ...interface{}) {
	if verbose() {
		fmt.Fprintf(os.Stdout, "coverage :: "+format+"\n", args...)
	}
}

func verrf(format string, args ...interface{}) {
	if verbose() {
		fmt.Fprintf(os.Stderr, "coverage :: err :: "+format+"\n", args...)
	}
}

func main() {
	flag.Parse()

	pwd, err := os.Getwd()
	maybeFatal(err)

	fmt.Fprintln(os.Stdout, "coverage starting")
	fmt.Fprintf(os.Stdout, "using covermode: %s\n", *covermode)
	fmt.Fprintf(os.Stdout, "using coverprofile: %s\n", *coverprofile)
	fullCoverageData, err := removeAndOpen(*coverprofile)
	if err != nil {
		maybeFatal(err)
	}
	fmt.Fprintln(fullCoverageData, "mode: set")

	// fileTotals is a map from the "package" file path to it's total line count
	fileTotals := map[string]int{}
	var fileName string
	maybeFatal(filepath.Walk("./", func(currentPath string, info os.FileInfo, err error) error {

		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		fileName = info.Name()

		if fileName == ".git" {
			vf("`%s` skipping dir; .git", currentPath)
			return filepath.SkipDir
		}
		if strings.HasPrefix(fileName, "_") {
			vf("`%s` skipping dir; '_' prefix", currentPath)
			return filepath.SkipDir
		}
		if fileName == "vendor" {
			vf("`%s` skipping dir; vendor", currentPath)
			return filepath.SkipDir
		}

		if !info.IsDir() {
			if strings.HasSuffix(fileName, ".go") {
				vf("`%s` counting file lines", currentPath)
				fileTotal, err := countFileLines(currentPath)
				if err != nil {
					return err
				}
				fileTotals[packageFilename(pwd, currentPath)] = fileTotal
			}
			return nil
		}

		if !dirHasGlob(currentPath, "*.go") {
			vf("`%s` skipping dir; no *.go files", currentPath)
			return nil
		}

		if len(*include) > 0 {
			if matches := globAnyMatch(*include, currentPath); !matches {
				vf("`%s` skipping dir; include no match: %s", currentPath, *include)
				return nil
			}
		}

		if len(*exclude) > 0 {
			if matches := globAnyMatch(*exclude, currentPath); matches {
				vf("`%s` skipping dir; exclude match: %s", currentPath, *exclude)
				return nil
			}
		}

		packageCoverReport := filepath.Join(currentPath, "profile.cov")
		err = removeIfExists(packageCoverReport)
		if err != nil {
			return err
		}

		var output []byte
		output, err = execCoverage(currentPath)
		if err != nil {
			verrf("error running coverage")
			fmt.Fprintln(os.Stderr, string(output))
			return err
		}

		coverage := extractCoverage(string(output))
		fmt.Fprintf(os.Stdout, "%s: %v%%\n", currentPath, colorCoverage(parseCoverage(coverage)))

		if enforce != nil && *enforce {
			vf("enforcing coverage minimums")
			err = enforceCoverage(currentPath, coverage)
			if err != nil {
				return err
			}
		}

		if update != nil && *update {
			fmt.Fprintf(os.Stdout, "%s updating coverage\n", currentPath)
			err = writeCoverage(currentPath, coverage)
			if err != nil {
				return err
			}
		}

		err = mergeCoverageOutput(packageCoverReport, fullCoverageData)
		if err != nil {
			return err
		}

		err = removeIfExists(packageCoverReport)
		if err != nil {
			return err
		}

		return nil
	}))

	// close the coverage data handle
	maybeFatal(fullCoverageData.Close())

	// complete summary steps
	covered, total, err := parseFullCoverProfile(pwd, *coverprofile, fileTotals)
	maybeFatal(err)
	finalCoverage := (float64(covered) / float64(total)) * 100
	maybeFatal(writeCoverage(pwd, formatCoverage(finalCoverage)))

	fmt.Fprintf(os.Stdout, "final coverage: %s%%\n", colorCoverage(finalCoverage))
	fmt.Fprintf(os.Stdout, "compiling coverage report: %s\n", *reportOutputPath)

	// compile coverage.html
	maybeFatal(execCoverageReportCompile())

	if !*keepCoverageOut {
		maybeFatal(removeIfExists(*coverprofile))
	}

	fmt.Fprintln(os.Stdout, "coverage complete")
}

// --------------------------------------------------------------------------------
// utilities
// --------------------------------------------------------------------------------

func gopath() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath
	}
	return build.Default.GOPATH
}

// globIncludeMatch tests if a file matches a (potentially) csv of glob filters.
func globAnyMatch(filter, file string) bool {
	parts := strings.Split(filter, ",")
	for _, part := range parts {
		if matches := glob(strings.TrimSpace(part), file); matches {
			return true
		}
	}
	return false
}

func glob(pattern, subj string) bool {
	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == star {
		return true
	}

	parts := strings.Split(pattern, star)

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, star)
	trailingGlob := strings.HasSuffix(pattern, star)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}

		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(subj, parts[end])
}

func enforceCoverage(path, actualCoverage string) error {
	actual, err := strconv.ParseFloat(actualCoverage, 64)
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadFile(filepath.Join(path, "COVERAGE"))
	if err != nil {
		return err
	}
	expected, err := strconv.ParseFloat(strings.TrimSpace(string(contents)), 64)
	if err != nil {
		return err
	}

	if expected == 0 {
		return nil
	}

	if actual < expected {
		return fmt.Errorf(
			"%s %s coverage: %0.2f%% vs. %0.2f%%",
			path, colorRed.Apply("fails"), expected, actual,
		)
	}
	return nil
}

func extractCoverage(corpus string) string {
	regex := `coverage: ([0-9,.]+)% of statements`
	expr := regexp.MustCompile(regex)

	results := expr.FindStringSubmatch(corpus)
	if len(results) > 1 {
		return results[1]
	}
	return "0"
}

func writeCoverage(path, coverage string) error {
	return ioutil.WriteFile(filepath.Join(path, "COVERAGE"), []byte(strings.TrimSpace(coverage)), defaultFileFlags)
}

func dirHasGlob(path, glob string) bool {
	files, _ := filepath.Glob(filepath.Join(path, glob))
	return len(files) > 0
}

func gobin() string {
	gobin, err := exec.LookPath("go")
	maybeFatal(err)
	return gobin
}

func execCoverage(path string) ([]byte, error) {
	var cmd *exec.Cmd
	if *timeout != "" {
		cmd = exec.Command(gobin(), "test", "-timeout", *timeout, "-short", fmt.Sprintf("-covermode=%s", *covermode), "-coverprofile=profile.cov")
	} else {
		cmd = exec.Command(gobin(), "test", "-short", fmt.Sprintf("-covermode=%s", *covermode), "-coverprofile=profile.cov")
	}
	cmd.Env = os.Environ()
	cmd.Dir = path
	return cmd.CombinedOutput()
}

func execCoverageReportCompile() error {
	cmd := exec.Command(gobin(), "tool", "cover", fmt.Sprintf("-html=%s", *coverprofile), fmt.Sprintf("-o=%s", *reportOutputPath))
	cmd.Env = os.Environ()
	return cmd.Run()
}

func mergeCoverageOutput(temp string, outFile *os.File) error {
	contents, err := ioutil.ReadFile(temp)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(contents))

	var skip int
	for scanner.Scan() {
		skip++
		if skip < 2 {
			continue
		}
		_, err = fmt.Fprintln(outFile, scanner.Text())
		if err != nil {
			return err
		}
	}
	return nil
}

func removeIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func removeAndOpen(path string) (*os.File, error) {
	if _, err := os.Stat(path); err == nil {
		if err = os.Remove(path); err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}

func countFileLines(path string) (lines int, err error) {
	if filepath.Ext(path) != ".go" {
		err = fmt.Errorf("count lines path must be a .go file")
		return
	}

	var contents []byte
	contents, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	lines = bytes.Count(contents, []byte{'\n'})
	return
}

// joinCoverPath takes a pwd, and a filename, and joins them
// overlaying parts of the suffix of the pwd, and the prefix
// of the filename that match.
// ex:
// - pwd: /foo/bar/baz, filename: bar/baz/buzz.go => /foo/bar/baz/buzz.go
func joinCoverPath(pwd, fileName string) string {
	pwdPath := lessEmpty(strings.Split(pwd, "/"))
	fileDirPath := lessEmpty(strings.Split(filepath.Dir(fileName), "/"))

	for index, dir := range pwdPath {
		if dir == first(fileDirPath) {
			pwdPath = pwdPath[:index]
			break
		}
	}

	return filepath.Join(maybePrefix(strings.Join(pwdPath, "/"), "/"), fileName)
}

// pacakgeFilename returns the github.com/foo/bar/baz.go form of the filename.
func packageFilename(pwd, relativePath string) string {
	fullPath := filepath.Join(pwd, relativePath)
	return strings.TrimPrefix(strings.TrimPrefix(fullPath, filepath.Join(gopath(), "src")), "/")
}

// parseFullCoverProfile parses the final / merged cover output.
func parseFullCoverProfile(pwd string, path string, fileTotals map[string]int) (covered, total int, err error) {
	vf("parsing coverage profile: %s", path)
	files, err := cover.ParseProfiles(path)
	if err != nil {
		return
	}

	var fileCovered int
	for _, fileTotal := range fileTotals {
		total += fileTotal
	}

	for _, file := range files {
		fileTotal := fileTotals[file.FileName]
		fileCovered = 0
		for _, block := range file.Blocks {
			fileCovered += (block.EndLine - block.StartLine) + 1
		}
		vf("processing coverage profile: %s result: %s (%d/%d lines)", path, file.FileName, fileCovered, fileTotal)
		covered += fileCovered
	}

	return
}

func lessEmpty(values []string) (output []string) {
	for _, value := range values {
		if len(value) > 0 {
			output = append(output, value)
		}
	}
	return
}

func first(values []string) (output string) {
	if len(values) == 0 {
		return
	}
	output = values[0]
	return
}

func maybePrefix(root, prefix string) string {
	if strings.HasPrefix(root, prefix) {
		return root
	}
	return prefix + root
}

// AnsiColor represents an ansi color code fragment.
type ansiColor string

func (acc ansiColor) escaped() string {
	return "\033[" + string(acc)
}

// Apply returns a string with the color code applied.
func (acc ansiColor) Apply(text string) string {
	return acc.escaped() + text + colorReset.escaped()
}

const (
	// ColorWhite is the posix escape code fragment for white.
	colorWhite ansiColor = "97m"
	// ColorBlack is the posix escape code fragment for black.
	colorBlack ansiColor = "30m"
	// ColorGray is the posix escape code fragment for black.
	colorGray ansiColor = "90m"
	// ColorRed is the posix escape code fragment for red.
	colorRed ansiColor = "31m"
	// ColorYellow is the posix escape code fragment for yellow.
	colorYellow ansiColor = "33m"
	// ColorGreen is the posix escape code fragment for green.
	colorGreen ansiColor = "32m"
	// ColorReset is the posix escape code fragment to reset all formatting.
	colorReset ansiColor = "0m"
)

func parseCoverage(coverage string) float64 {
	coverage = strings.TrimSpace(coverage)
	coverage = strings.TrimSuffix(coverage, "%")
	value, _ := strconv.ParseFloat(coverage, 64)
	return value
}

func colorCoverage(coverage float64) string {
	text := formatCoverage(coverage)
	if coverage > 80.0 {
		return colorGreen.Apply(text)
	} else if coverage > 70 {
		return colorYellow.Apply(text)
	} else if coverage == 0 {
		return colorGray.Apply(text)
	}
	return colorRed.Apply(text)
}

func formatCoverage(coverage float64) string {
	return fmt.Sprintf("%.2f", coverage)
}
