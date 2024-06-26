package httperror_test

import (
	"artk.dev/apperror"
	"artk.dev/assume"
	"artk.dev/broken"
	"artk.dev/httperror"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncodeToText_encodes_kind_into_status_code(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		t.Run(kind.String(), func(t *testing.T) {
			err := apperror.New(kind, errorMessage)

			w := httptest.NewRecorder()
			httperror.EncodeToText(w, err)

			expected := httperror.EncodeKind(kind)
			if got := w.Code; got != expected {
				t.Errorf("expected %v, got %v", expected, got)
			}
		})
	}
}

func TestEncodeToText_content_type_is_plain_text(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		t.Run(kind.String(), func(t *testing.T) {
			err := apperror.New(kind, errorMessage)

			w := httptest.NewRecorder()
			httperror.EncodeToText(w, err)

			const expected = "text/plain"
			contentType := w.Header().Get("Content-Type")
			mediaType, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				t.Error("failed:", err)
			}
			if mediaType != expected {
				t.Errorf(
					"expected %v, got %v",
					expected,
					mediaType,
				)
			}
		})
	}
}

func TestEncodeToText_contains_error_message_in_body(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		if kind == apperror.OK || kind == apperror.UnknownError {
			// Special cases.
			continue
		}

		t.Run(kind.String(), func(t *testing.T) {
			err := apperror.New(kind, errorMessage)
			w := httptest.NewRecorder()
			httperror.EncodeToText(w, err)

			expected := errorMessage
			got := strings.TrimSpace(w.Body.String())
			if got != expected {
				t.Errorf(
					`expected "%v", got "%v"`,
					errorMessage,
					got,
				)
			}
		})
	}
}

func TestEncodeToText_message_is_empty_if_kind_is_OK(t *testing.T) {
	var err error
	w := httptest.NewRecorder()
	httperror.EncodeToText(w, err)

	expected := ""
	if got := strings.TrimSuffix(w.Body.String(), "\n"); got != expected {
		t.Errorf(
			`expected "%v", got "%v"`,
			expected,
			got,
		)
	}
}

func TestEncodeToText_redacts_the_message_of_unknown_errors(t *testing.T) {
	err := errors.New("an unknown error")
	w := httptest.NewRecorder()
	httperror.EncodeToText(w, err)

	const expected = "Internal Server Error"
	if got := strings.TrimSuffix(w.Body.String(), "\n"); got != expected {
		t.Errorf(
			`expected "%v", got "%v"`,
			errorMessage,
			got,
		)
	}
}

func TestEncodeToText_panics_for_nil_response_writer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("missing expected panic")
		}
	}()

	err := errors.New(errorMessage)
	httperror.EncodeToText(nil, err)
}

func TestDecodeFromText_kind_encoding_is_reversible(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		t.Run(kind.String(), func(t *testing.T) {
			originalErr := apperror.New(kind, errorMessage)
			decodedErr := encodeAndDecode(originalErr)
			if got := apperror.KindOf(decodedErr); got != kind {
				t.Errorf("expected %v, got %v", kind, got)
			}
		})
	}
}

func TestDecodeFromText_message_encoding_is_reversible(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		if kind == apperror.UnknownError {
			// Special case: the message is redacted.
			continue
		}

		t.Run(kind.String(), func(t *testing.T) {
			originalErr := apperror.New(kind, errorMessage)
			decodedErr := encodeAndDecode(originalErr)
			assertEqualMessage(t, originalErr, decodedErr)
		})
	}
}

func TestDecodeFromText_return_unknown_error_on_failure(t *testing.T) {
	for _, kind := range apperror.KindValues() {
		if kind == apperror.OK {
			// The message is ignored for the OK kind.
			continue
		}

		t.Run(kind.String(), func(t *testing.T) {
			statusCode := httperror.EncodeKind(kind)
			reader := io.NopCloser(broken.Reader{})
			response := &http.Response{
				Status:     http.StatusText(statusCode),
				StatusCode: statusCode,
				Header: http.Header{
					"Content-Type": []string{"text/plain"},
				},
				Body: reader,
			}

			defer func() {
				assume.Success(reader.Close())
			}()

			err := httperror.DecodeFromText(response)
			expected := apperror.UnknownError
			if got := apperror.KindOf(err); got != expected {
				t.Errorf("expected %v, got %v", expected, got)
			}
		})
	}
}

func assertEqualMessage(t *testing.T, originalErr, decodedErr error) {
	t.Helper()

	// Ensure that either both are nil or none is nil.
	if originalErr == nil && decodedErr != nil {
		t.Fatal("expected nil, got", decodedErr)
	}
	if originalErr != nil && decodedErr == nil {
		t.Fatal("expected not nil, got nil")
	}

	// Cannot compare messages for nil errors.
	// This happens for apperror.OK.
	if originalErr == nil {
		return
	}

	// Handle all other error kinds.
	expected := originalErr.Error()
	got := decodedErr.Error()
	if got != expected {
		t.Errorf(`expected "%v", got "%v"`, expected, got)
	}
}

func encodeAndDecode(err error) error {
	w := httptest.NewRecorder()
	httperror.EncodeToText(w, err)

	response := w.Result()
	defer func() {
		assume.Success(response.Body.Close())
	}()

	return httperror.DecodeFromText(response)
}

const errorMessage = "test error"
