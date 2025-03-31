package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessResponseSerialization(t *testing.T) {
	type Payload struct {
		Message string `json:"message"`
	}

	original := APIResponse[Payload]{
		Status: 200,
		Data:   &Payload{Message: "hello"},
		Error:  nil,
	}

	// Marshal to JSON
	bytes, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var decoded APIResponse[Payload]
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	require.Equal(t, original.Status, decoded.Status)
	require.NotNil(t, decoded.Data)
	require.Equal(t, "hello", decoded.Data.Message)
	require.Nil(t, decoded.Error)
}

func TestErrorResponseSerialization(t *testing.T) {
	original := APIResponse[any]{
		Status: 400,
		Data:   nil,
		Error: &APIError{
			Code:    "BadRequest",
			Message: "Something went wrong",
		},
	}

	// Marshal to JSON
	bytes, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var decoded APIResponse[any]
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	require.Equal(t, original.Status, decoded.Status)
	require.Nil(t, decoded.Data)
	require.NotNil(t, decoded.Error)
	require.Equal(t, "BadRequest", decoded.Error.Code)
	require.Equal(t, "Something went wrong", decoded.Error.Message)
}

func TestParseAPIResponse_SuccessAndError(t *testing.T) {
	type Payload struct {
		Name string `json:"name"`
	}

	// Mock successful JSON response
	success := `{"status":200,"data":{"name":"Alice"}}`
	resp := httptest.NewRecorder()
	resp.WriteString(success)

	data, apiErr, err := ParseAPIResponse[Payload](resp.Result())
	require.NoError(t, err)
	require.Nil(t, apiErr)
	require.Equal(t, "Alice", data.Name)

	// Mock error JSON response
	errorResp := `{"status":404,"error":{"code":"NotFound","message":"User missing"}}`
	resp = httptest.NewRecorder()
	resp.WriteString(errorResp)

	data, apiErr, err = ParseAPIResponse[Payload](resp.Result())
	require.NoError(t, err)
	require.Nil(t, data)
	require.Equal(t, "NotFound", apiErr.Code)
}

func TestParseAPIResponse_MalformedJSON(t *testing.T) {
	type Payload struct {
		Name string `json:"name"`
	}

	// 🚨 Intentionally broken JSON (missing closing brace)
	malformed := `{"status":200,"data":{"name":"Alice"`

	rec := httptest.NewRecorder()
	_, _ = rec.WriteString(malformed)

	data, apiErr, err := ParseAPIResponse[Payload](rec.Result())

	require.Nil(t, data)
	require.Nil(t, apiErr)
	require.Error(t, err)
	t.Logf("✅ Caught decoding error: %v", err)
}

func TestParseAPIResponse_UnexpectedDataType(t *testing.T) {
	type ExpectedPayload struct {
		Name string `json:"name"`
	}

	// 🔁 Valid JSON — but 'data' is a number, not an object
	jsonMismatch := `{
		"status": 200,
		"data": 12345
	}`

	rec := httptest.NewRecorder()
	_, _ = rec.WriteString(jsonMismatch)

	data, apiErr, err := ParseAPIResponse[ExpectedPayload](rec.Result())

	require.Nil(t, data)
	require.Nil(t, apiErr)
	require.Error(t, err)
	t.Logf("✅ Caught mismatched type error: %v", err)
}

func TestParseAPIResponse_EmptyBody(t *testing.T) {
	type Dummy struct {
		Value string `json:"value"`
	}

	rec := httptest.NewRecorder()
	// 👇 Do not write anything to body
	rec.WriteHeader(http.StatusOK)

	data, apiErr, err := ParseAPIResponse[Dummy](rec.Result())

	require.Nil(t, data)
	require.Nil(t, apiErr)
	require.Error(t, err)
	t.Logf("✅ Caught empty body decode error: %v", err)
}

func TestParseAPIResponse_EmptyJSONStruct(t *testing.T) {
	type Dummy struct {
		Value string `json:"value"`
	}

	body := `{}`

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusOK)
	_, _ = rec.Write([]byte(body))

	data, apiErr, err := ParseAPIResponse[Dummy](rec.Result())

	require.NoError(t, err)
	require.Nil(t, data)
	require.Nil(t, apiErr)

	t.Log("✅ Valid empty JSON parsed without error, no data or apiErr present")
}

func TestParseAPIResponse_ErrorMissingCodeField(t *testing.T) {
	body := `{
		"error": {
			"message": "Missing code field"
		}
	}`

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusBadRequest)
	_, _ = rec.Write([]byte(body))

	type Dummy struct{}
	data, apiErr, err := ParseAPIResponse[Dummy](rec.Result())

	require.NoError(t, err)
	require.Nil(t, data)
	require.NotNil(t, apiErr)
	require.Equal(t, "", apiErr.Code)
	require.Equal(t, "Missing code field", apiErr.Message)

	t.Log("✅ Parsed API error with missing code field")
}

func TestParseAPIResponse_ErrorMissingMessageField(t *testing.T) {
	body := `{
		"error": {
			"code": "ErrOnlyCode"
		}
	}`

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusBadRequest)
	_, _ = rec.Write([]byte(body))

	type Dummy struct{}
	data, apiErr, err := ParseAPIResponse[Dummy](rec.Result())

	require.NoError(t, err)
	require.Nil(t, data)
	require.NotNil(t, apiErr)
	require.Equal(t, "ErrOnlyCode", apiErr.Code)
	require.Equal(t, "", apiErr.Message)

	t.Log("✅ Parsed API error with missing message field")
}

func TestParseAPIResponse_EmptyErrorObject(t *testing.T) {
	body := `{
		"error": {}
	}`

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusBadRequest)
	_, _ = rec.Write([]byte(body))

	type Dummy struct{}
	data, apiErr, err := ParseAPIResponse[Dummy](rec.Result())

	require.NoError(t, err)
	require.Nil(t, data)
	require.NotNil(t, apiErr)
	require.Equal(t, "", apiErr.Code)
	require.Equal(t, "", apiErr.Message)

	t.Log("✅ Parsed empty error object correctly")
}
