package trackingmiddleware

import (
	"net/http"

	"github.com/arquivei/foundationkit/request"
	"github.com/arquivei/foundationkit/trace"
)

func New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = request.WithRequestID(ctx)
		ctx = trace.WithTrace(ctx, trace.GetTraceFromHTTRequest(r))

		w.Header().Set("X-REQUESTID", request.GetRequestIDFromContext(ctx).String())
		trace.SetTraceInHTTPResponse(trace.GetTraceFromContext(ctx), w)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}