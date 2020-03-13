package gobrake

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type gitLog struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Revision string    `json:"revision"`
	Time     time.Time `json:"time"`
}

type gitInfo struct {
	Repository   string
	Revision     string
	LastCheckout *gitLog
}

var (
	gitInfosMu sync.RWMutex
	gitInfos   = make(map[string]*gitInfo)
)

func getGitInfo(dir string) *gitInfo {
	gitInfosMu.RLock()
	info, ok := gitInfos[dir]
	gitInfosMu.RUnlock()

	if ok {
		return info
	}

	gitInfosMu.Lock()
	defer gitInfosMu.Unlock()

	info = new(gitInfo)
	gitInfos[dir] = info

	repo, err := gitRepository(dir)
	if err != nil {
		logger.Printf("gitRepository dir=%q failed: %s", dir, err)
	} else {
		info.Repository = repo
	}

	rev, err := gitRevision(dir)
	if err != nil {
		logger.Printf("gitRevision dir=%q failed: %s", dir, err)
	} else {
		info.Revision = rev
	}

	lastCheckout, err := gitLastCheckout(dir)
	if err != nil {
		logger.Printf("gitLastCheckout dir=%q failed: %s", dir, err)
	} else {
		info.LastCheckout = lastCheckout
	}

	return info
}

func gitRepository(dir string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(trimnl(out)), nil
}

// findGitDir returns first directory containing .git file checking the dir and parent dirs.
func findGitDir(dir string) (string, bool) {
	dir, err := filepath.Abs(dir)
	if err != nil || !exists(dir) {
		return "", false
	}

	for i := 0; i < 10; i++ {
		path := filepath.Join(dir, ".git")
		if exists(path) {
			return dir, true
		}

		if dir == "." || dir == "/" {
			return "", false
		}

		dir = filepath.Dir(dir)
	}

	return "", false
}

func gitRevision(dir string) (string, error) {
	head, err := gitHead(dir)
	if err != nil {
		return "", err
	}

	prefix := []byte("ref: ")
	if !bytes.HasPrefix(head, prefix) {
		return string(head), nil
	}
	head = head[len(prefix):]

	refFile := filepath.Join(dir, ".git", string(head))
	rev, err := ioutil.ReadFile(refFile)
	if err == nil {
		return string(trimnl(rev)), nil
	}

	refsFile := filepath.Join(dir, ".git", "packed-refs")
	fd, err := os.Open(refsFile)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) == 0 || b[0] == '#' || b[0] == '^' {
			continue
		}

		bs := bytes.Split(b, []byte{' '})
		if len(bs) != 2 {
			continue
		}

		if bytes.Equal(bs[1], head) {
			return string(bs[0]), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("git revision for ref=%q not found", head)
}

func gitHead(dir string) ([]byte, error) {
	headFile := filepath.Join(dir, ".git", "HEAD")
	b, err := ioutil.ReadFile(headFile)
	if err != nil {
		return nil, err
	}
	return trimnl(b), nil
}

func trimnl(b []byte) []byte {
	for _, c := range []byte{'\n', '\r'} {
		if len(b) > 0 && b[len(b)-1] == c {
			b = b[:len(b)-1]
		} else {
			break
		}
	}
	return b
}

func gitLastCheckout(dir string) (*gitLog, error) {
	headFile := filepath.Join(dir, ".git", "logs", "HEAD")
	line, err := lastCheckoutLine(headFile)
	if err != nil {
		return nil, err
	}

	ind := strings.IndexByte(line, '\t')
	if ind == -1 {
		return nil, fmt.Errorf("tab not found")
	}
	line = line[:ind]

	parts := strings.Split(line, " ")
	if len(parts) < 5 {
		return nil, fmt.Errorf("can't parse %q", line)
	}
	author := parts[2 : len(parts)-2]

	utime, err := strconv.ParseInt(parts[len(parts)-2], 10, 64)
	if err != nil {
		return nil, err
	}

	info := &gitLog{
		Revision: parts[1],
		Time:     time.Unix(utime, 0),
	}
	if email := cleanEmail(author[len(author)-1]); email != "" {
		info.Email = email
		author = author[:len(author)-1]
	}
	info.Username = strings.Join(author, " ")

	return info, nil
}

func lastCheckoutLine(filename string) (string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(fd)
	var lastCheckout string
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "\tclone: ") ||
			strings.Contains(s, "\tpull: ") ||
			strings.Contains(s, "\tcheckout: ") {
			lastCheckout = s
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	if lastCheckout == "" {
		return "", fmt.Errorf("no clone, pull, or checkout entries")
	}
	return lastCheckout, nil
}

func cleanEmail(s string) string {
	if s == "" {
		return ""
	}
	if s[0] == '<' && s[len(s)-1] == '>' {
		return s[1 : len(s)-1]
	}
	return ""
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
