package setup

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

type SprintFunc func(a ...interface{}) string

type CustomLogHandler struct {
	withColor    bool
	minimumLevel slog.Level

	debugColor SprintFunc
	infoColor  SprintFunc
	warnColor  SprintFunc
	errorColor SprintFunc
	panicColor SprintFunc
}

func NewCustomLogHandler(withColor bool, level slog.Level) *CustomLogHandler {
	white := color.New(color.Bold, color.FgHiWhite).SprintFunc()
	green := color.New(color.Bold, color.FgGreen).SprintFunc()
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
	errorFunc := color.New(color.Bold, color.FgRed).SprintFunc()
	panicFunc := color.New(color.Bold, color.FgHiRed, color.BgWhite).SprintFunc()

	return &CustomLogHandler{
		withColor:    withColor,
		minimumLevel: level,
		debugColor:   green,
		infoColor:    white,
		warnColor:    yellow,
		errorColor:   errorFunc,
		panicColor:   panicFunc,
	}
}

func (h CustomLogHandler) Enabled(context context.Context, level slog.Level) bool {
	if level >= h.minimumLevel {
		return true
	} else {
		return false
	}
}

func (h CustomLogHandler) Handle(context context.Context, record slog.Record) error {
	message := record.Message

	// appends each attribute to the message
	// An attribute is of the form `<key>=<value>` and specified as in `slog.Error(<message>, <key>, <value>, ...)`.
	record.Attrs(func(attr slog.Attr) bool {
		message += fmt.Sprintf(" %v", attr)
		return true
	})

	// timestamp := record.Time.Format(time.RFC3339)
	timestamp := record.Time.Format("2006-01-02 15:04:05")

	if h.withColor {
		switch record.Level {
		case slog.LevelDebug:
			fmt.Fprintf(os.Stderr, "[%v] %v %v\n", h.debugColor(record.Level), timestamp, message)
		case slog.LevelInfo:
			fmt.Fprintf(os.Stderr, "[%v] %v %v\n", h.infoColor(record.Level), timestamp, message)
		case slog.LevelWarn:
			fmt.Fprintf(os.Stderr, "[%v] %v %v\n", h.warnColor(record.Level), timestamp, message)
		case slog.LevelError:
			fmt.Fprintf(os.Stderr, "[%v] %v %v\n", h.errorColor(record.Level), timestamp, message)
		default:
			panic("unreachable")
		}
	} else {
		switch record.Level {
		case slog.LevelDebug:
			fmt.Fprintf(os.Stderr, "[%v] %v %v\n", record.Level, timestamp, message)
		case slog.LevelInfo:
			fmt.Fprintf(os.Stderr, "[%v] %v %v\n", record.Level, timestamp, message)
		case slog.LevelWarn:
			fmt.Fprintf(os.Stderr, "[%v] %v %v\n", record.Level, timestamp, message)
		case slog.LevelError:
			fmt.Fprintf(os.Stderr, "[%v] %v %v\n", record.Level, timestamp, message)
		default:
			panic("unreachable")
		}
	}

	return nil
}

// for advanced users
func (h CustomLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("unimplemented")
}

// for advanced users
func (h CustomLogHandler) WithGroup(name string) slog.Handler {
	panic("unimplemented")
}
