package responses

import "time"

// MapSlice converts each item with fn. It always returns a non-nil slice so
// the frontend receives [] instead of null.
func MapSlice[D any, R any](items []D, fn func(D) R) []R {
	out := make([]R, 0, len(items))
	for _, it := range items {
		out = append(out, fn(it))
	}
	return out
}

// FormatTime formats t as RFC3339 for the frontend.
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatTimePtr returns nil for a nil time, otherwise the RFC3339 string.
func FormatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := FormatTime(*t)
	return &s
}
