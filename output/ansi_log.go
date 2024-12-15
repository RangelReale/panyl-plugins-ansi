package output

import (
	"bytes"
	"context"
	"fmt"

	"github.com/RangelReale/panyl/v2"
	"github.com/RangelReale/panyl/v2/util"
	"github.com/fatih/color"
)

var _ panyl.DebugLog = (*AnsiLog)(nil)

type AnsiLog struct {
	ShowSource bool
}

func (l AnsiLog) LogSourceLine(ctx context.Context, n int, line, rawLine string) {
	red := color.New(color.FgRed)

	red.Printf("@@@ SOURCE LINE [%d]: '%s' @@@\n", n, line)
}

func (l AnsiLog) LogItem(ctx context.Context, item *panyl.Item) {
	green := color.New(color.FgGreen)

	var lineno string
	if item.LineCount > 1 {
		lineno = fmt.Sprintf("[%d-%d]", item.LineNo, item.LineNo+item.LineCount-1)
	} else {
		lineno = fmt.Sprintf("[%d]", item.LineNo)
	}

	var buf bytes.Buffer

	if len(item.Metadata) > 0 {
		_, _ = buf.WriteString(fmt.Sprintf("Metadata: %+v", item.Metadata))
	}
	if len(item.Data) > 0 {
		if buf.Len() > 0 {
			_, _ = buf.WriteString(" - ")
		}
		_, _ = buf.WriteString(fmt.Sprintf("Data: %+v", item.Data))
	}

	if len(item.Line) > 0 {
		if buf.Len() > 0 {
			_, _ = buf.WriteString(" - ")
		}
		_, _ = buf.WriteString(fmt.Sprintf("Line: \"%s\"", item.Line))
	}

	if l.ShowSource && len(item.Source) > 0 {
		if buf.Len() > 0 {
			_, _ = buf.WriteString(" - ")
		}
		_, _ = buf.WriteString(fmt.Sprintf("Source: \"%s\"", util.DoAnsiEscapeString(item.Source)))
	}

	green.Printf("*** PROCESS LINE %s: %s\n", lineno, buf.String())
}
