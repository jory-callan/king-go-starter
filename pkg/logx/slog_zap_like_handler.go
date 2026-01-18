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
	w      io.Writer
	opts   slog.HandlerOptions
	level  slog.Leveler
	color  bool
	groups []string
}

func newZapLikeHandler(w io.Writer, opts *slog.HandlerOptions) *ZapLikeHandler {
	//var buf [10]uintptr
	//n := runtime.Callers(0, buf[:])
	//frames := runtime.CallersFrames(buf[:n])
	//for i := 0; i < n; i++ {
	//	frame, _ := frames.Next()
	//	fmt.Printf("[%d] %s:%d\n", i, filepath.Base(frame.File), frame.Line)
	//}
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	h := &ZapLikeHandler{
		w:     w,
		opts:  *opts,
		level: opts.Level,
		color: isTerminal(w),
	}
	return h
}

// Enabled implements slog.Handler.
func (h *ZapLikeHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle implements slog.Handler.
func (h *ZapLikeHandler) Handle(_ context.Context, r slog.Record) error {
	//var buf [10]uintptr
	//n := runtime.Callers(0, buf[:])
	//frames := runtime.CallersFrames(buf[:n])
	//for i := 0; i < n; i++ {
	//	frame, _ := frames.Next()
	//	fmt.Printf("[%d] %s:%d\n", i, filepath.Base(frame.File), frame.Line)
	//}

	// Time
	timeStr := r.Time.Format("2006-01-02 15:04:05")

	// Level with color
	levelStr := h.formatLevel(r.Level)

	// 3. Logger name
	loggerName := strings.Join(h.groups, ".")

	// Message
	msg := r.Message

	var parts []string
	parts = append(parts, timeStr, levelStr)
	if loggerName != "" {
		parts = append(parts, loggerName)
	}

	parts = append(parts, msg)

	// Build prefix: time level [caller] msg
	prefix := strings.Join(parts, "")

	prefix = fmt.Sprintf("%s  %s  %s", timeStr, levelStr, msg)

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
	if name == "" {
		return h
	}
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name
	return &ZapLikeHandler{
		w:      h.w,
		opts:   h.opts,
		level:  h.level,
		color:  h.color,
		groups: newGroups,
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

func (h *ZapLikeHandler) detectCallerSkip() int {
	// find the first caller, the one is this file:lineNumber
	//var buf [10]uintptr
	//n := runtime.Callers(0, buf[:])
	//frames := runtime.CallersFrames(buf[:n])
	//for i := 0; i < n; i++ {
	//	frame, _ := frames.Next()
	//	fmt.Printf("[%d] %s:%d\n", i, filepath.Base(frame.File), frame.Line)
	//}

	// Record a test log and capture its PC
	var testPC uintptr
	func() {
		var pcs [1]uintptr
		runtime.Callers(1, pcs[:]) // skip this func
		testPC = pcs[0]
	}()
	// Now simulate a log call and see what skip gives us testPC
	for skip := 0; skip <= 10; skip++ {
		var pcs [1]uintptr
		n := runtime.Callers(skip, pcs[:])
		if n > 0 && pcs[0] == testPC {
			return skip
		}
	}
	return 5 // fallback
}
