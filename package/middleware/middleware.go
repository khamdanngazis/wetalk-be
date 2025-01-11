package middleware

import (
	"bytes"
	"chat-be/package/logging"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique request ID using UUID.
		requestID := uuid.New().String()
		startTime := time.Now()

		// Add the request ID to the request context.
		ctx := context.WithValue(r.Context(), logging.RequestIDKey, requestID)

		// Set the updated context in the request.
		r = r.WithContext(ctx)

		// Log the request details with the generated request ID using logrus.
		field := logrus.Fields{
			"request_id": requestID,
			"method":     r.Method,
			"uri":        r.RequestURI,
			"proto":      r.Proto,
			"body":       extractRequestBody(r),
		}
		logging.LogCustomField(logrus.InfoLevel, field, "Incoming request")

		// Create a custom response writer to capture the response status and body
		responseRecorder := NewResponseRecorder(w)

		// Call the next handler in the chain.
		next.ServeHTTP(responseRecorder, r)

		// Log the response details
		field = logrus.Fields{
			"request_id": requestID,
			"status":     responseRecorder.Status(),
			"body":       responseRecorder.Body(),
			"duration":   fmt.Sprintf("%v", time.Since(startTime)),
		}

		logging.LogCustomField(logrus.InfoLevel, field, "Outgoing response")
	})
}

type ResponseRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func extractRequestBody(r *http.Request) string {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "Error reading request body"
	}
	defer r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(body)
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{w, http.StatusOK, bytes.Buffer{}}
}

// WriteHeader captures the status code.
func (r *ResponseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write captures the response body.
func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Status returns the captured status code.
func (r *ResponseRecorder) Status() int {
	return r.status
}

// Body returns the captured response body.
func (r *ResponseRecorder) Body() string {
	return r.body.String()
}
