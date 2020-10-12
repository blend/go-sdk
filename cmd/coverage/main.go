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
	expand           = "/..."
)

var reportOutputPath = flag.String("output", "coverage.html", "the path to write the full html coverage report")
var update = flag.Bool("update", false, "if we should write the current coverage to `COVERAGE` files")
var enforce = flag.Bool("enforce", false, "if we should enforce coverage minimums defined in `COVERAGE` files")
var timeout = flag.String("timeout", "", "the timeout to pass to the package tests.")
var race = flag.Bool("race", false, "if we should add -race to test invocations")
var covermode = flag.String("covermode", "atomic", "the go test covermode.")
var coverprofile = flag.String("coverprofile", "coverage.cov", "the intermediate cover profile.")
var keepCoverageOut = flag.Bool("keep-coverage-out", false, "if we should keep coverage.out")
var v = flag.Bool("v", false, "show verbose output")
var disableEarlyExitOnCoverageFailure = flag.Bool("include-coverage-report-failures", false, "enable to produce coverage reports when coverage fails")

var (
	includes Paths
	excludes Paths
)

func main() {
	flag.Var(&includes, "include", "glob patterns to include explicitly")
	flag.Var(&excludes, "exclude", "glob patterns to exclude explicitly")
	flag.Parse()

	pwd, err := os.Getwd()
	maybeFatal(err)

	fmt.Fprintln(os.Stdout, "coverage starting")
	fmt.Fprintf(os.Stdout, "using covermode: %s\n", *covermode)
	fmt.Fprintf(os.Stdout, "using coverprofile: %s\n", *coverprofile)
	if *timeout != "" {
		fmt.Fprintf(os.Stdout, "using timeout: %s\n", *timeout)
	}
	if len(includes) > 0 {
		fmt.Fprintf(os.Stdout, "using includes: %s\n", strings.Join(includes, ", "))
	}
	if len(excludes) > 0 {
		fmt.Fprintf(os.Stdout, "using excludes: %s\n", strings.Join(excludes, ", "))
	}
	if *race {
		fmt.Fprintln(os.Stdout, "using race detection")
	}

	//
	// start
	//

	fullCoverageData, err := removeAndOpen(*coverprofile)
	if err != nil {
		maybeFatal(err)
	}
	fmt.Fprintf(fullCoverageData, "mode: %s\n", *covermode)

	paths := flag.Args()

	if len(paths) == 0 {
		paths = []string{"./..."}
	}

	var allPathCoverageErrors []error
	for _, path := range paths {
		fmt.Fprintf(os.Stdout, "walking path: %s\n", path)
		if coverageErrors := walkPath(path, fullCoverageData); len(coverageErrors) > 0 {
			allPathCoverageErrors = append(allPathCoverageErrors, coverageErrors...)
		}
	}

	// close the coverage data handle
	maybeFatal(fullCoverageData.Close())

	// complete summary steps
	covered, total, err := parseFullCoverProfile(pwd, *coverprofile)
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

	if len(allPathCoverageErrors) > 0 {
		fmt.Fprintln(os.Stderr, "coverage thresholds not met")
		for _, coverageError := range allPathCoverageErrors {
			fmt.Fprintf(os.Stderr, "%+v\n", coverageError)
		}
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "coverage complete")
}

func walkPath(walkedPath string, fullCoverageData *os.File) []error {
	recursive := strings.HasSuffix(walkedPath, expand)
	rootPath := filepath.Dir(walkedPath)
	var coverageErrors []error

	maybeFatal(filepath.Walk(rootPath, func(currentPath string, info os.FileInfo, fileErr error) error {
		packageCoverReport, err := getPackageCoverage(currentPath, info, fileErr)

		if err != nil && !*disableEarlyExitOnCoverageFailure {
			return err
		}

		if err != nil && len(packageCoverReport) > 0 {
			coverageErrors = append(coverageErrors, err)
		}

		if len(packageCoverReport) == 0 {
			return nil
		}

		err = mergeCoverageOutput(packageCoverReport, fullCoverageData)
		if err != nil {
			return err
		}

		err = removeIfExists(packageCoverReport)
		if err != nil {
			return err
		}

		if !recursive && info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}))
	return coverageErrors
}

// gets coverage for a directory and returns the path to the coverage file for that directory
func getPackageCoverage(currentPath string, info os.FileInfo, err error) (string, error) {
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	fileName := info.Name()

	if fileName == ".git" {
		vf("%q skipping dir; .git", currentPath)
		return "", filepath.SkipDir
	}
	if strings.HasPrefix(fileName, "_") {
		vf("%q skipping dir; '_' prefix", currentPath)
		return "", filepath.SkipDir
	}
	if fileName == "vendor" {
		vf("%q skipping dir; vendor", currentPath)
		return "", filepath.SkipDir
	}

	if !dirHasGlob(currentPath, "*.go") {
		vf("%q skipping dir; no *.go files", currentPath)
		return "", nil
	}

	for _, include := range includes {
		if matches := glob(include, currentPath); !matches { // note the !
			vf("%q skipping dir; include no match: %s", currentPath, include)
			return "", nil
		}
	}

	for _, exclude := range excludes {
		if matches := glob(exclude, currentPath); matches {
			vf("%q skipping dir; exclude match: %s", currentPath, exclude)
			return "", nil
		}
	}

	packageCoverReport := filepath.Join(currentPath, "profile.cov")
	err = removeIfExists(packageCoverReport)
	if err != nil {
		return "", err
	}

	var output []byte
	output, err = execCoverage(currentPath)
	if err != nil {
		verrf("error running coverage")
		fmt.Fprintln(os.Stderr, string(output))
		return "", err
	}

	coverage := extractCoverage(string(output))
	fmt.Fprintf(os.Stdout, "%s: %v%%\n", currentPath, colorCoverage(parseCoverage(coverage)))

	if enforce != nil && *enforce {
		vf("enforcing coverage minimums")
		err = enforceCoverage(currentPath, coverage)
		if err != nil {
			return packageCoverReport, err
		}
	}

	if update != nil && *update {
		fmt.Fprintf(os.Stdout, "%q updating coverage\n", currentPath)
		err = writeCoverage(currentPath, coverage)
		if err != nil {
			return "", err
		}
	}

	return packageCoverReport, nil
}

// --------------------------------------------------------------------------------
// utilities
// --------------------------------------------------------------------------------

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

func gopath() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath
	}
	return build.Default.GOPATH
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
	args := []string{
		"test",
		"-short",
		fmt.Sprintf("-covermode=%s", *covermode),
		"-coverprofile=profile.cov",
	}
	if *timeout != "" {
		args = append(args, fmt.Sprintf("-timeout=%s", *timeout))
	}
	if *race {
		args = append(args, "-race")
	}
	cmd := exec.Command(gobin(), args...)
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
func parseFullCoverProfile(pwd string, path string) (covered, total int, err error) {
	vf("parsing coverage profile: %q", path)
	files, err := cover.ParseProfiles(path)
	if err != nil {
		return
	}

	var fileCovered, numLines int

	for _, file := range files {
		fileCovered = 0

		for _, block := range file.Blocks {
			numLines = block.EndLine - block.StartLine

			total += numLines
			if block.Count != 0 {
				fileCovered += numLines
			}
		}

		vf("processing coverage profile: %q result: %s (%d/%d lines)", path, file.FileName, fileCovered, numLines)
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

// Paths are cli flag input paths.
type Paths []string

// String returns the param as a string.
func (p *Paths) String() string {
	return fmt.Sprint(*p)
}

// Set sets a value.
func (p *Paths) Set(value string) error {
	for _, val := range strings.Split(value, ",") {
		*p = append(*p, val)
	}
	return nil
}
