package logx

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

// ZapLikeHandler mimics zap's console output style for slog.
type ZapLikeHandler struct {
	w          io.Writer
	opts       slog.HandlerOptions
	level      slog.Leveler
	color      bool
	callerSkip int
}

func newZapLikeHandler(w io.Writer, opts *slog.HandlerOptions) *ZapLikeHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	h := &ZapLikeHandler{
		w:          w,
		opts:       *opts,
		level:      opts.Level,
		color:      isTerminal(w),
		callerSkip: 6,
	}
	return h
}

// Enabled implements slog.Handler.
func (h *ZapLikeHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle implements slog.Handler.
func (h *ZapLikeHandler) Handle(_ context.Context, r slog.Record) error {

	// Time
	timeStr := r.Time.Format("2006-01-02 15:04:05")

	// Level with color
	levelStr := h.formatLevel(r.Level)

	// Message
	msg := r.Message

	// Caller (if AddSource is enabled)
	callerStr := ""
	if h.opts.AddSource && r.PC != 0 {
		var pcs [1]uintptr
		n := runtime.Callers(h.callerSkip, pcs[:])
		if n > 0 {
			frames := runtime.CallersFrames(pcs[:n])
			frame, _ := frames.Next()
			if frame.File != "" {
				// cut filename from program directory
				d, _ := os.Getwd()
				parts := strings.Split(frame.File, d)
				filename := parts[len(parts)-1]
				filename = strings.TrimPrefix(filename, "/")
				callerStr = fmt.Sprintf("%s:%d", filename, frame.Line)
			}
		}
	}

	var parts []string
	parts = append(parts, timeStr, levelStr)
	if callerStr != "" {
		parts = append(parts, callerStr)
	}

	parts = append(parts, msg)

	prefix := strings.Join(parts, "  ")

	// Append attributes (key=value)
	attrs := make([]string, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, h.formatAttr(a))
		return true
	})

	// Final output
	var out string
	if len(attrs) > 0 {
		out = prefix + "  " + strings.Join(attrs, "  ")
	} else {
		out = prefix
	}

	_, err := fmt.Fprintln(h.w, out)
	return err
}

// WithAttrs implements slog.Handler.
func (h *ZapLikeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we ignore context attrs in this demo.
	// In production, you'd store them and prepend to each record.
	return h
}

// WithGroup implements slog.Handler.
func (h *ZapLikeHandler) WithGroup(name string) slog.Handler {
	return &ZapLikeHandler{
		w:     h.w,
		opts:  h.opts,
		level: h.level,
		color: h.color,
	}
}

// formatLevel returns colored or plain level string.
func (h *ZapLikeHandler) formatLevel(level slog.Level) string {
	s := level.String()
	if !h.color {
		return s
	}
	switch level {
	case slog.LevelDebug:
		return "\033[36mDEBUG\033[0m" // Cyan
	case slog.LevelInfo:
		return "\033[32mINFO\033[0m" // Green
	case slog.LevelWarn:
		return "\033[33mWARN\033[0m" // Yellow
	case slog.LevelError:
		return "\033[31mERROR\033[0m" // Red
	default:
		return s
	}
}

// formatAttr formats a single attribute.
func (h *ZapLikeHandler) formatAttr(a slog.Attr) string {
	// Quote string values, leave others as-is
	if v, ok := a.Value.Any().(string); ok {
		return fmt.Sprintf("%s=%q", a.Key, v)
	}
	return fmt.Sprintf("%s=%v", a.Key, a.Value.Any())
}

// isTerminal checks if w is a terminal (for color support).
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return isatty(f.Fd())
	}
	return false
}

// isatty is a minimal implementation (you can use github.com/mattn/go-isatty)
func isatty(fd uintptr) bool {
	// Simplified: assume stdout/stderr are terminals
	// For production, use: go get github.com/mattn/go-isatty
	return fd == 1 || fd == 2 // stdout/stderr
}
