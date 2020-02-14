package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/zap/buffer"
	"golang.org/x/tools/cover"
)

func main() {

	tmplHTML := `
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<style>
			body {
				background: black;
				color: rgb(80, 80, 80);
			}
			body, pre, #legend span {
				font-family: Menlo, monospace;
				font-weight: bold;
			}
			#topbar {
				background: black;
				position: fixed;
				top: 0; left: 0; right: 0;
				height: 42px;
				border-bottom: 1px solid rgb(80, 80, 80);
			}
			#content {
				margin-top: 50px;
			}
			#nav, #legend {
				float: left;
				margin-left: 10px;
			}
			#legend {
				margin-top: 12px;
			}
			#nav {
				margin-top: 10px;
			}
			#legend span {
				margin: 0 5px;
			}
			{{colors}}
		</style>
	</head>
	<body>
		<div id="topbar">
			<div id="nav">
				<select id="files">
				{{range $i, $f := .Files}}
				<option value="file{{$i}}">{{$f.Name}} ({{printf "%.1f" $f.Coverage}}%)</option>
				{{end}}
				</select>
			</div>
			<div id="legend">
				<span>not tracked</span>
			{{if .Set}}
				<span class="cov0">not covered</span>
				<span class="cov8">covered</span>
			{{else}}
				<span class="cov0">no coverage</span>
				<span class="cov1">low coverage</span>
				<span class="cov2">*</span>
				<span class="cov3">*</span>
				<span class="cov4">*</span>
				<span class="cov5">*</span>
				<span class="cov6">*</span>
				<span class="cov7">*</span>
				<span class="cov8">*</span>
				<span class="cov9">*</span>
				<span class="cov10">high coverage</span>
			{{end}}
			</div>
		</div>
		<div id="content">
		{{range $i, $f := .Files}}
		<pre class="file" id="file{{$i}}" {{if $i}}style="display: none"{{end}}>{{$f.Body}}</pre>
		{{end}}
		</div>
	</body>
	<script>
	(function() {
		var files = document.getElementById('files');
		var visible = document.getElementById('file0');
		files.addEventListener('change', onChange, false);
		function onChange() {
			visible.style.display = 'none';
			visible = document.getElementById(files.value);
			visible.style.display = 'block';
			window.scrollTo(0, 0);
		}
	})();
	</script>
</html>
`

	die := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	htmlGen := func(w io.Writer, src []byte, boundaries []cover.Boundary) error {
		dst := bufio.NewWriter(w)
		for i := range src {
			for len(boundaries) > 0 && boundaries[0].Offset == i {
				b := boundaries[0]
				if b.Start {
					n := 0
					if b.Count > 0 {
						n = int(math.Floor(b.Norm*9)) + 1
					}
					fmt.Fprintf(dst, `<span class="cov%v" title="%v">`, n, b.Count)
				} else {
					dst.WriteString("</span>")
				}
				boundaries = boundaries[1:]
			}
			switch b := src[i]; b {
			case '>':
				dst.WriteString("&gt;")
			case '<':
				dst.WriteString("&lt;")
			case '&':
				dst.WriteString("&amp;")
			case '\t':
				dst.WriteString("        ")
			default:
				dst.WriteByte(b)
			}
		}
		return dst.Flush()
	}

	rgb := func(n int) string {
		if n == 0 {
			return "rgb(192, 0, 0)" // Red
		}
		// Gradient from gray to green.
		r := 128 - 12*(n-1)
		g := 128 + 12*(n-1)
		b := 128 + 3*(n-1)
		return fmt.Sprintf("rgb(%v, %v, %v)", r, g, b)
	}

	colors := func() template.CSS {
		var buf bytes.Buffer
		for i := 0; i < 11; i++ {
			fmt.Fprintf(&buf, ".cov%v { color: %v }\n", i, rgb(i))
		}
		return template.CSS(buf.String())
	}

	htmlTemplate := template.Must(template.New("html").Funcs(template.FuncMap{
		"colors": colors,
	}).Parse(tmplHTML))

	type templateFile struct {
		Name     string
		Body     template.HTML
		Coverage float64
	}

	type templateData struct {
		Files []*templateFile
		Set   bool
	}

	root := os.Getenv("ROOT")
	var files []string
	die(filepath.Walk(filepath.Join(root, "bazel-out", runtime.GOOS+"-fastbuild", "testlogs"), func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".dat") {
			files = append(files, path)
		}
		return err
	}))

	var output []string
	output = append(output, "mode: set")

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		die(err)
		lines := strings.Split(string(content), "\n")
		if len(lines) > 1 {
			for i := 1; i < len(lines); i++ {
				line := strings.TrimSpace(lines[i])
				if line == "" {
					continue
				}
				output = append(output, line)
			}
		}
	}

	die(ioutil.WriteFile(filepath.Join(root, "coverall.dat"), []byte(strings.Join(output, "\n")), 0644))
	profiles, err := cover.ParseProfiles(filepath.Join(root, "coverall.dat"))
	die(err)

	var d templateData
	d.Set = true

	jsonCoverageFiles := make(map[string][]int64)
	var overallCoverage float64

	for _, profile := range profiles {
		var (
			total   int64
			covered int64
		)
		for _, block := range profile.Blocks {
			total += int64(block.NumStmt)
			if block.Count > 0 {
				covered += int64(block.NumStmt)
			}
		}

		f := profile.FileName
		for _, char := range root {
			if rune(f[0]) == rune(char) {
				f = f[1:]
			}
		}

		src, err := ioutil.ReadFile(filepath.Join(root, f))
		die(err)

		coveragePercentage := float64(covered) / float64(total) * 100
		overallCoverage += coveragePercentage

		var buf bytes.Buffer
		die(htmlGen(&buf, src, profile.Boundaries(src)))
		d.Files = append(d.Files, &templateFile{
			Name:     profile.FileName,
			Body:     template.HTML(buf.String()),
			Coverage: coveragePercentage,
		})

		fmt.Println(profile.FileName, total, covered)
		jsonCoverageFiles[profile.FileName] = []int64{covered, total}
	}

	var outputBuf buffer.Buffer
	die(htmlTemplate.Execute(&outputBuf, d))
	die(ioutil.WriteFile(filepath.Join(root, "coverall.html"), outputBuf.Bytes(), 0644))

	jsonBytes, _ := json.Marshal(&map[string]interface{}{
		"overall": overallCoverage / float64(len(profiles)),
		"files":   jsonCoverageFiles,
	})
	die(ioutil.WriteFile(filepath.Join(root, "coverall.json"), jsonBytes, 0644))
}
