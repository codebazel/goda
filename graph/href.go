package graph

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/loov/goda/pkggraph"
	exec "golang.org/x/sys/execabs"
)

func (ctx *Dot) genHrefs(graph *pkggraph.Graph) {
	hrefs := make([]string, 0, len(graph.Sorted))
	for _, n := range graph.Sorted {
		href := ctx.docs + n.ID
		hrefs = append(hrefs, href)
	}
	ctx.hrefs = ctx.transHrefs(hrefs)
}

func (ctx *Dot) transHrefs(hrefs []string) map[string]string {
	tool := os.Getenv("GODAHREFDRIVER")
	if tool == "" {
		var err error
		tool, err = exec.LookPath("godahrefdriver")
		if err != nil {
			return nil
		}
	}

	buf := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := exec.Command(tool)
	cmd.Stdin = strings.NewReader(strings.Join(hrefs, "\n"))
	cmd.Stdout = buf
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(ctx.err, "%v: %v: %s", tool, err, cmd.Stderr)
		return nil
	}

	rhs := make(map[string]string, len(hrefs))
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(ctx.err, "%v: %v", tool, err)
			break
		}
		// space splited href mapping.
		ss := strings.SplitN(strings.TrimSpace(line), " ", 2)
		if len(ss) == 2 {
			r, h := strings.TrimSpace(ss[0]), strings.TrimSpace(ss[1])
			if h != "" {
				rhs[r] = h
			}
		}
	}
	return rhs
}
