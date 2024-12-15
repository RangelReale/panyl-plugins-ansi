package output

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RangelReale/panyl/v2"
	"github.com/fatih/color"
)

var _ panyl.ProcessResult = (*AnsiOutput)(nil)

type AnsiOutputSprintfFunc func(format string, a ...interface{}) string

type AnsiOutput struct {
	ColorInformation, ColorWarning, ColorError, ColorInternalError, ColorUnknown AnsiOutputSprintfFunc
}

func NewAnsiOutput(ansi bool) *AnsiOutput {
	ret := &AnsiOutput{
		ColorError:         fmt.Sprintf,
		ColorWarning:       fmt.Sprintf,
		ColorInformation:   fmt.Sprintf,
		ColorInternalError: fmt.Sprintf,
		ColorUnknown:       fmt.Sprintf,
	}
	if ansi {
		ret.ColorError = color.New(color.FgRed).SprintfFunc()
		ret.ColorWarning = color.New(color.FgYellow).SprintfFunc()
		ret.ColorInformation = color.New(color.FgGreen).SprintfFunc()
		ret.ColorInternalError = color.New(color.FgHiRed).SprintfFunc()
		ret.ColorUnknown = color.New(color.FgMagenta).SprintfFunc()
	}
	return ret
}

func (o *AnsiOutput) OnResult(ctx context.Context, item *panyl.Item) (cont bool) {
	var out bytes.Buffer

	// level
	var levelColor AnsiOutputSprintfFunc
	level := item.Metadata.StringValue(panyl.MetadataLevel)
	switch level {
	case panyl.MetadataLevelDEBUG, panyl.MetadataLevelINFO:
		levelColor = o.ColorInformation
	case panyl.MetadataLevelWARNING:
		levelColor = o.ColorWarning
	case panyl.MetadataLevelERROR:
		levelColor = o.ColorError
	default:
		level = "unknown"
		levelColor = o.ColorUnknown
	}

	// timestamp
	if ts, ok := item.Metadata[panyl.MetadataTimestamp]; ok {
		out.WriteString(fmt.Sprintf("%s ", ts.(time.Time).Local().Format("2006-01-02 15:04:05.000")))
	}

	// application
	if application := item.Metadata.StringValue(panyl.MetadataApplication); application != "" {
		out.WriteString(fmt.Sprintf("| %s | ", application))
	}

	// level
	if level != "" {
		out.WriteString(fmt.Sprintf("[%s] ", level))
	}

	// format
	if format := item.Metadata.StringValue(panyl.MetadataFormat); format != "" {
		out.WriteString(fmt.Sprintf("(%s) ", format))
	}

	// category
	if category := item.Metadata.StringValue(panyl.MetadataCategory); category != "" {
		out.WriteString(fmt.Sprintf("{{%s}} ", category))
	}

	// message
	if msg := item.Metadata.StringValue(panyl.MetadataMessage); msg != "" {
		out.WriteString(msg)
	} else if item.Line != "" {
		out.WriteString(item.Line)
	} else if len(item.Data) > 0 {
		dt, err := json.Marshal(item.Data)
		if err != nil {
			fmt.Println(o.ColorInternalError("Error marshaling data to json: %s", err.Error()))
			return
		}
		out.WriteString(fmt.Sprintf("| %s", string(dt)))
	}

	fmt.Println(levelColor(out.String()))

	return true
}

func (o *AnsiOutput) OnFlush(ctx context.Context) {}

func (o *AnsiOutput) OnClose(ctx context.Context) {}
