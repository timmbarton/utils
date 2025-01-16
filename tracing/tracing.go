package tracing

import (
	"context"
	"runtime"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const slash = "/"

//goland:noinspection ALL
func NewSpan(c context.Context, caller ...string) (ctx context.Context, span trace.Span) {
	spanName := ""
	if len(caller) != 0 {
		spanName = caller[0]
	} else {
		pc, _, _, ok := runtime.Caller(1)

		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			detailsItems := strings.Split(details.Name(), slash)
			if len(detailsItems) != 0 {
				spanName = detailsItems[len(detailsItems)-1]
			} else {
				spanName = details.Name()
			}
		} else {
			spanName = "anonymous function"
		}
	}
	ctx, span = otel.Tracer("").Start(c, spanName)

	return ctx, span
}

//goland:noinspection ALL
func GetTraceID(span trace.Span) string {
	return span.SpanContext().TraceID().String()
}
