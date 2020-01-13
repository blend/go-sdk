package template

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/uuid"
)

func TestTemplateFromFile(t *testing.T) {
	assert := assert.New(t)

	temp, err := NewFromFile("testdata/test.template.yml")
	assert.Nil(err)

	temp = temp.
		WithVar("name", "test-service").
		WithVar("container-name", "nginx").
		WithVar("container-image", "nginx:1.7.9").
		WithVar("container-port", "disabled")

	buffer := bytes.NewBuffer(nil)
	err = temp.Process(buffer)
	assert.Nil(err)

	result := buffer.String()
	assert.True(strings.Contains(result, "name: test-service"))
	assert.True(strings.Contains(result, "replicas: 2"))
	assert.False(strings.Contains(result, "containerPort:"))

	temp = temp.WithVar("container-port", 80)
	err = temp.Process(buffer)
	assert.Nil(err)
	result = buffer.String()
	assert.True(strings.Contains(result, "port: 80"))
}

func TestTemplateTemplate(t *testing.T) {
	assert := assert.New(t)

	main := `{{ template "test" . }}`
	test := `{{ define "test" }}{{ .Var "foo" }}{{end}}`

	buffer := bytes.NewBuffer(nil)
	err := New().WithBody(main).WithInclude(test).WithVar("foo", "bar").Process(buffer)
	assert.Nil(err)
	assert.Equal("bar", buffer.String())

	buffer = bytes.NewBuffer(nil)
	err = New().WithBody(main).WithVar("foo", "bar").Process(buffer)
	assert.NotNil(err)
}

func TestTemplateVar(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" }}`
	temp := New().WithBody(test).WithVar("foo", "bar")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("bar", buffer.String())
}

func TestTemplateVarMissing(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "baz" }}`
	temp := New().WithBody(test).WithVar("foo", "bar")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.NotNil(err)
}

func TestTemplateEnv(t *testing.T) {
	assert := assert.New(t)

	varName := uuid.V4().String()
	env.Env().Set(varName, "bar")
	defer env.Restore()

	test := fmt.Sprintf(`{{ .Env "%s" }}`, varName)
	temp := New().WithBody(test)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err, fmt.Sprintf("%+v", err))
	assert.Equal("bar", buffer.String())
}

func TestTemplateHasEnv(t *testing.T) {
	assert := assert.New(t)

	varName := uuid.V4().String()
	env.Env().Set(varName, "bar")
	defer env.Restore()

	test := fmt.Sprintf(`{{ if .HasEnv "%s" }}yep{{else}}nope{{end}}`, varName)
	temp := New().WithBody(test)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("yep", buffer.String())
}

func TestTemplateHasEnvFalsey(t *testing.T) {
	assert := assert.New(t)

	varName := uuid.V4().String()

	test := fmt.Sprintf(`{{ if .HasEnv "%s" }}yep{{else}}nope{{end}}`, varName)
	temp := New().WithBody(test)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("nope", buffer.String())
}

func TestTemplateEnvMissing(t *testing.T) {
	assert := assert.New(t)

	varName := uuid.V4().String()

	test := fmt.Sprintf(`{{ .Env "%s" }}`, varName)
	temp := New().WithBody(test)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.NotNil(err)
}

func TestTemplateExpandEnv(t *testing.T) {
	assert := assert.New(t)

	varName := uuid.V4().String()

	envVars := env.Env()
	envVars.Set(varName, "bar")

	templateBody := fmt.Sprintf(`{{ .ExpandEnv "${%s}.foo" }}`, varName)
	temp := New().WithBody(templateBody).WithEnvVars(envVars)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("bar.foo", buffer.String())
}

func TestTemplateReadFile(t *testing.T) {
	assert := assert.New(t)

	test := `{{ read_file "testdata/inline_file" }}`
	temp := New().WithBody(test)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("this is a test", buffer.String())
}

func TestTemplateFileExists(t *testing.T) {
	assert := assert.New(t)

	test := `{{ if file_exists "testdata/inline_file" }}yep{{else}}nope{{end}}`
	temp := New().WithBody(test)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("yep", buffer.String())
}

func TestTemplateFileExistsFalsey(t *testing.T) {
	assert := assert.New(t)

	fileName := uuid.V4().String()

	test := fmt.Sprintf(`{{ if file_exists "testdata/%s" }}yep{{else}}nope{{end}}`, fileName)
	temp := New().WithBody(test)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("nope", buffer.String())
}

func TestTemplateViewFuncTimeUnix(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "now" | unix }}`
	temp := New().WithBody(test).WithVar("now", time.Date(2017, 05, 20, 21, 00, 00, 00, time.UTC))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("1495314000", buffer.String())
}

func TestTemplateViewFuncSince(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "ts" | since }}`
	temp := New().WithBody(test).
		WithVar("ts", time.Date(2018, 05, 20, 21, 00, 00, 00, time.UTC))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.NotEmpty(buffer.String())
}

func TestTemplateViewFuncSinceUTC(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "ts" | since_utc }}`
	temp := New().WithBody(test).
		WithVar("ts", time.Date(2018, 05, 20, 21, 00, 00, 00, time.UTC))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.NotEmpty(buffer.String())
}

func TestTemplateGenerateKey(t *testing.T) {
	assert := assert.New(t)

	test := `{{ generate_key 64 }}`
	temp := New().WithBody(test)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)

	assert.True(len(buffer.String()) > 64)
}

func TestTemplateViewFuncAsString(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | as_string }}`
	temp := New().WithBody(test).WithVar("foo", 123)

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("123", buffer.String())
}

func TestTemplateViewFuncAsBytes(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | as_bytes }}`
	temp := New().WithBody(test).WithVar("foo", "123")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("[49 50 51]", buffer.String())
}

func TestTemplateViewFuncParseTime(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_time "2006/01/02" | day }}`
	temp := New().WithBody(test).WithVar("foo", "2017/05/30")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("30", buffer.String())
}

func TestTemplateViewFuncTimeFromUnix(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_unix | year }}`
	temp := New().WithBody(test).WithVar("foo", time.Date(2017, 05, 20, 21, 00, 00, 00, time.UTC).Unix())

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("2017", buffer.String())
}

func TestTemplateViewFuncTimeMonth(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | month }}`
	temp := New().WithBody(test).WithVar("foo", time.Date(2017, 05, 20, 21, 00, 00, 00, time.UTC))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("5", buffer.String())
}

func TestTemplateViewFuncTimeDay(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | day }}`
	temp := New().WithBody(test).WithVar("foo", time.Date(2017, 05, 20, 21, 00, 00, 00, time.UTC))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("20", buffer.String())
}

func TestTemplateViewFuncTimeHour(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | hour }}`
	temp := New().WithBody(test).WithVar("foo", time.Date(2017, 05, 20, 21, 00, 00, 00, time.UTC))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("21", buffer.String())
}

func TestTemplateViewFuncTimeMinute(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | minute }}`
	temp := New().WithBody(test).WithVar("foo", time.Date(2017, 05, 20, 21, 16, 17, 18, time.UTC))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("16", buffer.String())
}

func TestTemplateViewFuncTimeSecond(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | second }}`
	temp := New().WithBody(test).WithVar("foo", time.Date(2017, 05, 20, 21, 16, 17, 18, time.UTC))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("17", buffer.String())
}

func TestTemplateViewFuncTimeFromUnixString(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | as_string | parse_int64 | parse_unix | year }}`
	temp := New().WithBody(test).WithVar("foo", time.Date(2017, 05, 20, 21, 00, 00, 00, time.UTC).Unix())

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("2017", buffer.String())
}

func TestTemplateViewFuncParseBool(t *testing.T) {
	assert := assert.New(t)

	test := `{{ if .Var "foo" | parse_bool }}yep{{end}}`
	temp := New().WithBody(test).WithVar("foo", "true")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("yep", buffer.String())
}

func TestTemplateViewFuncParseInt(t *testing.T) {
	assert := assert.New(t)

	test := `{{ if .Var "foo" | parse_int | eq 123 }}yep{{end}}`
	temp := New().WithBody(test).WithVar("foo", "123")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("yep", buffer.String())
}

func TestTemplateViewFuncParseInt64(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_int64 }}`
	temp := New().WithBody(test).WithVar("foo", fmt.Sprintf("%d", (1<<33)))

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("8589934592", buffer.String())
}

func TestTemplateViewFuncParseFloat64(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_float64 }}`
	temp := New().WithBody(test).WithVar("foo", "3.14")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("3.14", buffer.String())
}

func TestTemplateViewFuncFormatMoney(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_float64 | format_money }}`
	temp := New().WithBody(test).WithVar("foo", "3.00")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("$3.00", buffer.String())
}

func TestTemplateViewFuncFormatPct(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_float64 | format_pct }}`
	temp := New().WithBody(test).WithVar("foo", "0.24")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("24.00%", buffer.String())
}

func TestTemplateViewFuncBase64(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | base64 }}`
	temp := New().WithBody(test).WithVar("foo", "bar")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("YmFy", buffer.String())
}

func TestTemplateViewFuncBase64Decode(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | base64 | base64decode }}`
	temp := New().WithBody(test).WithVar("foo", "bar")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("bar", buffer.String())
}

func TestTemplateViewFuncSplit(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | split "," }}`
	temp := New().WithBody(test).WithVar("foo", "bar,baz,biz")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("[bar baz biz]", buffer.String())
}

func TestTemplateViewFuncFirst(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | split "," | first }}`
	temp := New().WithBody(test).WithVar("foo", "bar,baz,biz")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("bar", buffer.String())
}
func TestTemplateViewFuncIndex(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | split "," | at_index 1 }}`
	temp := New().WithBody(test).WithVar("foo", "bar,baz,biz")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err, fmt.Sprintf("%+v", err))
	assert.Equal("baz", buffer.String())
}

func TestTemplateViewFuncLast(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | split "," | last }}`
	temp := New().WithBody(test).WithVar("foo", "bar,baz,biz")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("biz", buffer.String())
}

func TestTemplateViewFuncSlice(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | split "," | slice 1 3 }}`
	temp := New().WithBody(test).WithVar("foo", "bar,baz,biz,boof")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("[baz biz]", buffer.String())
}
func TestTemplateViewFuncJoin(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | split "," | slice 1 3 | join "/" }}`
	temp := New().WithBody(test).WithVar("foo", "bar,baz,biz,boof")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("baz/biz", buffer.String())
}

func TestTemplateViewFuncHasPrefix(t *testing.T) {
	assert := assert.New(t)

	test := `{{ if .Var "foo" | has_prefix "http" }}yep{{end}}`
	temp := New().WithBody(test).WithVar("foo", "http://foo.bar.com")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("yep", buffer.String())
}

func TestTemplateViewFuncHasSuffix(t *testing.T) {
	assert := assert.New(t)

	test := `{{ if .Var "foo" | has_suffix "com" }}yep{{end}}`
	temp := New().WithBody(test).WithVar("foo", "http://foo.bar.com")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("yep", buffer.String())
}

func TestViewfuncConcat(t *testing.T) {
	assert := assert.New(t)

	test := `{{ concat "foo" "." "bar" "." "baz" }}`
	temp := New().WithBody(test)
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("foo.bar.baz", buffer.String())
}

func TestViewfuncPrefix(t *testing.T) {
	assert := assert.New(t)

	test := `{{ "foo" | prefix "bar." }}`
	temp := New().WithBody(test)
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("bar.foo", buffer.String())
}

func TestViewfuncSuffix(t *testing.T) {
	assert := assert.New(t)

	test := `{{ "foo" | suffix ".bar" }}`
	temp := New().WithBody(test)
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("foo.bar", buffer.String())
}
func TestTemplateViewFuncContains(t *testing.T) {
	assert := assert.New(t)

	test := `{{ if .Var "foo" | contains "bar" }}yep{{end}}`
	temp := New().WithBody(test).WithVar("foo", "http://foo.bar.com")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("yep", buffer.String())
}

func TestTemplateViewFuncMatches(t *testing.T) {
	assert := assert.New(t)

	test := `{{ if .Var "foo" | matches "^[a-z]+$" }}yep{{else}}nope{{end}}`
	temp := New().WithBody(test).WithVar("foo", "http://foo.bar.com")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("nope", buffer.String())
}

func TestTemplateViewFuncSha1(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | sha1 }}`
	temp := New().WithBody(test).WithVar("foo", "this is only a test")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("e7ee879d16c08f616c32e5bbe2253bdba18cf003", buffer.String())
}

func TestTemplateViewFuncMD5(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | md5 }}`
	temp := New().WithBody(test).WithVar("foo", "this is only a test")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("e668034188ba397a9b6ff95d2a8e7203", buffer.String())
}

func TestTemplateViewFuncSha256(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | sha256 }}`
	temp := New().WithBody(test).WithVar("foo", "this is only a test")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("1661186b8e38e79f434e4549a2d53f84716cfff7c45d334bbc67c9d41d1e3be6", buffer.String())
}

func TestTemplateViewFuncSha512(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | sha512 }}`
	temp := New().WithBody(test).WithVar("foo", "this is only a test")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("819bf8f4c3c5508a061d0637d09858cf098ef8ef7cafa312d07fca8480703eccf1a00b24b8915e24f926a8106331d7bc064e63c04262dbed65e05b28e208e53e", buffer.String())
}

func TestTemplateViewFuncSemverMajor(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_semver | semver_major }}`
	temp := New().WithBody(test).WithVar("foo", "1.2.3-beta1")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("1", buffer.String())
}

func TestTemplateViewFuncSemverMinor(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_semver | semver_minor }}`
	temp := New().WithBody(test).WithVar("foo", "1.2.3-beta1")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("2", buffer.String())
}

func TestTemplateViewFuncSemverPatch(t *testing.T) {
	assert := assert.New(t)

	test := `{{ .Var "foo" | parse_semver | semver_patch }}`
	temp := New().WithBody(test).WithVar("foo", "1.2.3-beta1")
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("3", buffer.String())
}

type label struct {
	Name string `yaml:"name"`
	Vaue string `yaml:"value"`
}

func TestTemplateViewFuncYAML(t *testing.T) {
	assert := assert.New(t)

	test := `
type: foo
meta:
	name:
	labels:
{{ .Var "labels" | to_yaml | indent_tabs 1 }}
`
	temp := New().WithBody(test).WithVar("labels", []label{
		{"foo", "bar"},
		{"bar", "baz"},
		{"moobar", "zoobar"},
	})
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.NotEmpty(buffer.String())
}

func TestTemplateViewFuncJSON(t *testing.T) {
	assert := assert.New(t)

	test := `
type: foo
meta:
	name:
	labels:
{{ .Var "labels" | to_json | indent_tabs 1 }}
`
	temp := New().WithBody(test).WithVar("labels", []label{
		{"foo", "bar"},
		{"bar", "baz"},
		{"moobar", "zoobar"},
	})
	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.NotEmpty(buffer.String())
}

func TestTemplateWithVars(t *testing.T) {
	assert := assert.New(t)

	temp := New().WithVars(map[string]interface{}{
		"foo": "baz",
		"bar": "buz",
	})
	assert.True(temp.HasVar("foo"))
	assert.True(temp.HasVar("bar"))
	assert.False(temp.HasVar("baz"))
	assert.False(temp.HasVar("buz"))
	val, err := temp.Var("foo")
	assert.Nil(err)
	assert.Equal("baz", val)
	val, err = temp.Var("bar")
	assert.Nil(err)
	assert.Equal("buz", val)
}

func TestTemplateWithDelims(t *testing.T) {
	assert := assert.New(t)

	curly := `{{ .Var "foo" }}`
	pointy := `<< .Var "foo" >>`
	temp := New().WithBody(curly+pointy).WithDelims("<<", ">>").WithVar("foo", "bar")

	buffer := bytes.NewBuffer(nil)
	err := temp.Process(buffer)
	assert.Nil(err)
	assert.Equal(curly+"bar", buffer.String())

	temp = New().WithBody(curly+pointy).WithDelims("", "").WithVar("foo", "bar")
	buffer.Reset()
	err = temp.Process(buffer)
	assert.Nil(err)
	assert.Equal("bar"+pointy, buffer.String())
}

func TestOrdinalNames(t *testing.T) {
	assert := assert.New(t)
	tmp := New().WithBody("{{ generate_ordinal_names \"cockroachdb-%d\" 5 | join \",\" }}")

	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("cockroachdb-0,cockroachdb-1,cockroachdb-2,cockroachdb-3,cockroachdb-4", buffer.String())
}

func TestViewfuncNow(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ now | unix | parse_int64 }}")
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.NotEmpty(buffer.String())
}

func TestViewfuncNowUTC(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ now_utc | unix | parse_int64 }}")
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.NotEmpty(buffer.String())
}

func TestViewFuncRFC3339(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | rfc3339 }}").WithVar("data", time.Date(2018, 10, 06, 12, 0, 0, 0, time.UTC))
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("2018-10-06T12:00:00Z", buffer.String())
}

func TestViewFuncTimeShort(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | time_short }}").WithVar("data", time.Date(2018, 10, 06, 12, 0, 0, 0, time.UTC))
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("10/06/2018 12:00:00 PM", buffer.String())
}

func TestViewFuncTimeMedium(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | time_medium }}").WithVar("data", time.Date(2018, 10, 06, 12, 0, 0, 0, time.UTC))
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("Oct 06, 2018 12:00:00 PM", buffer.String())
}

func TestViewFuncTimeKitchen(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | time_kitchen }}").WithVar("data", time.Date(2018, 10, 06, 12, 0, 0, 0, time.UTC))
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("12:00PM", buffer.String())
}

func TestViewFuncDateMonthDay(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | date_month_day }}").WithVar("data", time.Date(2018, 10, 06, 12, 0, 0, 0, time.UTC))
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("10/6", buffer.String())
}

func TestViewFuncTimeInLocation(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | in_loc \"UTC\" }}").WithVar("data", time.Date(2018, 10, 06, 12, 0, 0, 0, time.UTC))
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("2018-10-06 12:00:00 +0000 UTC", buffer.String())
}

func TestViewFuncRound(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | round 2 }}").WithVar("data", 12.34567)
	buffer := new(bytes.Buffer)
	err := tmp.Process(buffer)
	assert.Nil(err, fmt.Sprintf("%+v", err))
	assert.Equal("12.35", buffer.String())
}

func TestViewFuncCeil(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | ceil }}").WithVar("data", 12.34567)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("13", buffer.String())
}

func TestViewFuncFloor(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Var \"data\" | floor }}").WithVar("data", 12.34567)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("12", buffer.String())
}

func TestViewfuncUUID(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ uuidv4 | as_string | parse_uuid }}")
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.NotEmpty(buffer.String())
}

func TestViewfuncToUpper(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "foo" | to_upper }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("FOO", buffer.String())
}

func TestViewfuncToLower(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "FOO" | to_lower }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("foo", buffer.String())
}

func TestViewfuncTrimSpace(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "  foo  " | trim_space }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("foo", buffer.String())
}

func TestViewfuncSplitN(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "a:b:c:d" | split_n ":" 2 }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("[a b:c:d]", buffer.String())
}

func TestViewfuncRandomLetters(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ random_letters 10 }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Len(buffer.String(), 10)
}

func TestViewfuncRandomLettersWithNumbers(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ random_letters_with_numbers 10 }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Len(buffer.String(), 10)
}

func TestViewfuncRandomLettersWithNumbersAndSymbols(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ generate_password 10 }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Len(buffer.String(), 10)
}

func TestViewfuncURLEncode(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "foo bar" | urlencode }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal(`foo+bar`, buffer.String())
}

func TestViewfuncURLScheme(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "http://foo.com/bar?fuzz=buzz" | parse_url | url_scheme }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("http", buffer.String())
}

func TestViewfuncURLHost(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "http://foo.com/bar?fuzz=buzz" | parse_url | url_host }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("foo.com", buffer.String())
}

func TestViewfuncURLPort(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "http://foo.com:8080/bar?fuzz=buzz" | parse_url | url_port }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("8080", buffer.String())
}

func TestViewfuncURLPath(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "http://foo.com:8080/bar?fuzz=buzz" | parse_url | url_path }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("/bar", buffer.String())
}

func TestViewfuncURLQuery(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ "http://foo.com:8080/bar?fuzz=buzz&up=down" | parse_url | url_query "up" }}`)
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("down", buffer.String())
}

func TestTemplateProcess(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ read_file "testdata/process.yml" | process . }}`).
		WithVar("db.database", "postgres").
		WithVar("db.username", "root").
		WithVar("db.password", "password")

	buffer := new(bytes.Buffer)
	err := tmp.Process(buffer)
	assert.Nil(err, fmt.Sprintf("%+v", err))
	assert.NotEmpty(buffer.String())
}

func TestViewfuncCSV(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ .Var "things" | csv}}`).WithVar("things", []string{"a", "b", "c"})
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("a,b,c", buffer.String())
}

func TestViewfuncTSV(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody(`{{ .Var "things" | tsv}}`).WithVar("things", []string{"a", "b", "c"})
	buffer := new(bytes.Buffer)
	assert.Nil(tmp.Process(buffer))
	assert.Equal("a\tb\tc", buffer.String())
}

func TestViewfuncQuote(t *testing.T) {
	assert := assert.New(t)

	tmp := New().WithBody("{{ .Vars.foo | quote }}").WithVar("foo", "foo")
	assert.Equal(`"foo"`, tmp.MustProcessString())

	tmp = New().WithBody("{{ .Vars.foo | quote }}").WithVar("foo", "\"foo\"")
	assert.Equal(`"foo"`, tmp.MustProcessString())

	tmp = New().WithBody("{{ .Vars.foo | quote }}").WithVar("foo", "\n\"foo\"\t")
	assert.Equal(`"foo"`, tmp.MustProcessString())
}
