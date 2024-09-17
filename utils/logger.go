package utils

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/lmittmann/tint"
)

type SourceHandler struct {
	handler     slog.Handler
	source      string
	parentAttrs []slog.Attr
}

func NewSourceHandler(w io.Writer, opts *tint.Options) *SourceHandler {
	return &SourceHandler{
		handler: tint.NewHandler(w, opts),
	}
}

func (h *SourceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

const (
	FgYellow = "\x1b[33m"
	FgRed    = "\x1b[31m"
	FgCyan   = "\x1b[36m"
	FgWhite  = "\x1b[37m"
	FgReset  = "\x1b[0m"
)

func (h *SourceHandler) Handle(ctx context.Context, r slog.Record) error {
	// Use the source from the handler if available, otherwise try to extract from the record
	source := h.source
	if source == "" {
		r.Attrs(func(a slog.Attr) bool {
			if a.Key == "service" {
				source = a.Value.String()
				return false
			}
			return true
		})
	}

	// Modify the message to include the source if available
	var newMessage string
	if source != "" {
		coloredSource := fmt.Sprint(FgCyan + "[" + strings.ToUpper(source) + "]" + FgReset)
		newMessage = coloredSource + " " + r.Message
	} else {
		newMessage = r.Message
	}

	// Create a new record with the modified message
	newRecord := slog.NewRecord(r.Time, r.Level, newMessage, r.PC)

	// Add all attributes except 'service' to the new record
	for _, attr := range h.parentAttrs {
		if attr.Key != "service" {
			newRecord.AddAttrs(attr)
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != "service" {
			newRecord.AddAttrs(a)
		}
		return true
	})

	// Handle the modified record
	return h.handler.Handle(ctx, newRecord)
}

func (h *SourceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newParentAttrs := append([]slog.Attr{}, h.parentAttrs...)
	newSource := h.source

	// Check if 'service' is in the new attributes
	for _, attr := range attrs {
		if attr.Key == "service" {
			newSource = attr.Value.String()
			// remove the 'service' attribute from the new attributes

			break
		} else {
			newParentAttrs = append(newParentAttrs, attr)
		}

	}

	return &SourceHandler{
		handler:     h.handler.WithAttrs(newParentAttrs),
		source:      newSource,
		parentAttrs: newParentAttrs,
	}
}

func (h *SourceHandler) WithGroup(name string) slog.Handler {
	return &SourceHandler{handler: h.handler.WithGroup(name), parentAttrs: h.parentAttrs, source: h.source}
}
