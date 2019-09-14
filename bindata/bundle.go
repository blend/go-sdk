package bindata

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/blend/go-sdk/ex"
)

// Bundle is an assets bundle with associated options.
type Bundle struct {
	PackageName string
	Ignores     []*regexp.Regexp
}

// PackageNameOrDefault returns the package name or a default.
func (b *Bundle) PackageNameOrDefault() string {
	if b.PackageName != "" {
		return b.PackageName
	}
	return "static"
}

// Start writes the file preamble.
func (b *Bundle) Start(output io.Writer) error {
	return b.anyError(
		b.writeHeader(output),
		b.writeHelpers(output),
		b.writeTypeFile(output),
		b.writeAssetsHeader(output),
	)
}

// ProcessPath processes a path with a given config.
func (b *Bundle) ProcessPath(output io.Writer, pc PathConfig) error {
	if !pc.Recursive {
		f, err := b.readFile(pc.Path)
		if err != nil {
			return err
		}
		if err := b.writeAssetsFile(output, f); err != nil {
			return err
		}
	} else {
		if err := b.findFiles(pc.Path, b.Ignores, func(f *File) error {
			return b.writeAssetsFile(output, f)
		}); err != nil {
			return err
		}
	}
	return nil
}

// Finish closes out the output file
func (b *Bundle) Finish(output io.Writer) error {
	return b.writeAssetsFooter(output)
}

func (b *Bundle) writeHeader(output io.Writer) error {
	return b.anyError(
		b.writeLines(output,
			"// Code generated by bindata.",
			"// DO NOT EDIT!",
			"",
		),
		b.writeLines(output,
			"package "+b.PackageNameOrDefault(),
			"",
			"import (",
			"\t\"bytes\"",
			"\t\"compress/gzip\"",
			"\t\"io/ioutil\"",
			"\t\"os\"",
			"\t\"path/filepath\"",
			")",
			"",
		),
	)
}

func (b *Bundle) writeHelpers(output io.Writer) error {
	return b.writeLines(output,
		"// GetBinaryAsset returns a binary asset file or",
		"// os.ErrNotExist if it is not found.",
		"func GetBinaryAsset(path string) (*BinaryFile, error) {",
		"\tfile, ok := BinaryAssets[filepath.Clean(path)]",
		"\tif !ok {",
		"\t\treturn nil, os.ErrNotExist",
		"\t}",
		"\treturn file, nil",
		"}",
		"",
	)
}

func (b *Bundle) writeTypeFile(output io.Writer) error {
	return b.writeLines(output,
		"// BinaryFile represents a statically managed binary asset.",
		"type BinaryFile struct {",
		"\tName               string",
		"\tModTime            int64",
		"\tMD5                []byte",
		"\tCompressedContents []byte",
		"}",
		"",
		"// Contents returns the raw uncompressed content bytes",
		"func (bf *BinaryFile) Contents() ([]byte, error) {",
		"\tgzr, err := gzip.NewReader(bytes.NewReader(bf.CompressedContents))",
		"\tif err != nil {",
		"\t\treturn nil, err",
		"\t}",
		"\treturn ioutil.ReadAll(gzr)",
		"}",
		"",
		"// Decompress returns a decompression stream.",
		"func (bf *BinaryFile) Decompress() (*gzip.Reader, error) {",
		"\treturn gzip.NewReader(bytes.NewReader(bf.CompressedContents))",
		"}",
		"",
	)
}

func (b *Bundle) writeAssetsHeader(output io.Writer) error {
	return b.writeLines(output,
		"// BinaryAssets are a map from relative filepath to the binary file contents.",
		"// The binary file contents include the file name, md5, modtime, and binary contents.",
		"var BinaryAssets = map[string]*BinaryFile{",
	)
}

func (b *Bundle) writeAssetsFile(output io.Writer, file *File) error {
	return b.anyError(
		b.ignoreCount(io.WriteString(output, "\t\""+file.Name+"\": ")),
		b.writeAssetsFileContents(output, file),
	)
}

func (b *Bundle) writeAssetsFileContents(output io.Writer, file *File) error {
	return b.anyError(
		b.writeLines(output,
			"&BinaryFile{",
			"\t\tName:    \""+file.Name+"\",",
			"\t\tModTime: "+strconv.FormatInt(file.Modtime.Unix(), 10)+",",
			"\t\tMD5: []byte{",
		),
		b.ignoreCount((&ByteWriter{Writer: output, Indent: []byte("\t\t\t")}).Write(file.Contents.MD5.Sum(nil))),
		b.writeLines(output,
			"",
			"\t\t},",
			"\t\tCompressedContents: []byte{",
		),
		b.ignoreBigCount(file.Contents.WriteTo(&ByteWriter{Writer: output, Indent: []byte("\t\t\t")})),
		b.writeLines(output,
			"",
			"\t\t},",
			"\t},",
		),
	)
}

func (b *Bundle) writeAssetsFooter(output io.Writer) error {
	return b.writeLines(output,
		"}",
		"",
	)
}

func (b *Bundle) anyError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bundle) write(output io.Writer, text string) error {
	return b.ignoreCount(io.WriteString(output, text))
}

func (b *Bundle) writeLines(output io.Writer, lines ...string) error {
	var err error
	for _, line := range lines {
		err = b.ignoreCount(io.WriteString(output, line+"\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bundle) ignoreBigCount(_ int64, err error) error {
	return err
}

func (b *Bundle) ignoreCount(_ int, err error) error {
	return err
}

// FindFiles traverses a root recursively, ignoring files that match the optional ignore
// expression, and calls the handler when it finds a file.
func (b *Bundle) findFiles(root string, ignores []*regexp.Regexp, handler func(*File) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, ignore := range ignores {
			if ignore.MatchString(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if info.IsDir() {
			return nil
		}
		f, err := b.readFile(path)
		if err != nil {
			return err
		}
		if err = handler(f); err != nil {
			return err
		}
		return nil
	})
}

func (b Bundle) readFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, ex.New(err)
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, ex.New(err)
	}

	return &File{
		Name:     path,
		Modtime:  stat.ModTime(),
		Contents: NewFileCompressor(f),
	}, nil
}
