package docker

import (
	"bytes"
	"io"
	"sort"
	"strings"
	"sync"
)

const (
	// maskPlaceholder replaces every masked value.
	maskPlaceholder = "[masked]"
	// minMaskedValueLen is the shortest value that is masked; the replaced
	// text of a shorter value is noise. The limit is part of the promise of
	// [driven.DevcontainerUpOptions.Progress].
	minMaskedValueLen = 4
)

// maskingWriter forwards what it receives to dst and replaces every
// occurrence of a value by [maskPlaceholder]. It holds back the last
// maxLen-1 bytes of the stream until more input or [maskingWriter.Flush]
// decides them, so a value that arrives in two Write calls is replaced as
// well. Values are matched byte for byte, longest first.
type maskingWriter struct {
	mu     sync.Mutex
	dst    io.Writer
	values [][]byte
	maxLen int
	held   []byte
}

// newMaskingWriter masks the given values (empty and too short ones are
// ignored) in what is written to dst. Without maskable values it forwards
// unchanged.
func newMaskingWriter(dst io.Writer, values []string) *maskingWriter {
	w := &maskingWriter{dst: dst}
	for _, v := range values {
		if len(v) >= minMaskedValueLen {
			w.values = append(w.values, []byte(v))
			w.maxLen = max(w.maxLen, len(v))
		}
	}
	sort.SliceStable(w.values, func(i, j int) bool { return len(w.values[i]) > len(w.values[j]) })
	return w
}

// maskString returns s with the values replaced as a maskingWriter replaces
// them in a stream.
func maskString(s string, values []string) string {
	var out strings.Builder
	w := newMaskingWriter(&out, values)
	_, _ = w.Write([]byte(s))
	_ = w.Flush()
	return out.String()
}

// Write implements [io.Writer].
func (w *maskingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.values) == 0 {
		return w.dst.Write(p)
	}
	buf := append(w.held, p...)
	consumed, out := w.mask(buf, false)
	w.held = append([]byte(nil), buf[consumed:]...)
	if _, err := w.dst.Write(out); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Flush writes the held bytes; call it when the stream has ended.
func (w *maskingWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.held) == 0 {
		return nil
	}
	_, out := w.mask(w.held, true)
	w.held = nil
	_, err := w.dst.Write(out)
	return err
}

// mask replaces the values in buf. Unless final, it stops maxLen-1 bytes
// before the end, where a value may still be incomplete, and returns the
// number of bytes consumed.
func (w *maskingWriter) mask(buf []byte, final bool) (int, []byte) {
	var out []byte
	start, i := 0, 0
	for i < len(buf) && (final || len(buf)-i >= w.maxLen) {
		if n := w.matchAt(buf[i:]); n > 0 {
			out = append(out, buf[start:i]...)
			out = append(out, maskPlaceholder...)
			i += n
			start = i
			continue
		}
		i++
	}
	return i, append(out, buf[start:i]...)
}

// matchAt returns the length of the longest value at the start of b, or 0.
func (w *maskingWriter) matchAt(b []byte) int {
	for _, v := range w.values {
		if bytes.HasPrefix(b, v) {
			return len(v)
		}
	}
	return 0
}
