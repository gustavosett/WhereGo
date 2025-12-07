package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJSONSerializer_Serialize covers serialization scenarios using a table-driven approach.
func TestJSONSerializer_Serialize(t *testing.T) {
	e := echo.New()
	serializer := &JSONSerializer{}
	e.JSONSerializer = serializer

	tests := []struct {
		name         string
		input        interface{}
		indent       string
		expectedBody string
	}{
		{
			name:         "Default Serialization",
			input:        map[string]string{"hello": "world"},
			indent:       "",
			expectedBody: `{"hello":"world"}`,
		},
		{
			name:         "Indented Serialization",
			input:        map[string]string{"hello": "world"},
			indent:       "  ",
			expectedBody: "{\n  \"hello\": \"world\"\n}\n",
		},
		{
			name:         "Complex Struct",
			input:        struct{ ID int }{ID: 1},
			indent:       "",
			expectedBody: `{"ID":1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := serializer.Serialize(c, tc.input, tc.indent)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}

// TestJSONSerializer_Deserialize covers deserialization scenarios.
func TestJSONSerializer_Deserialize(t *testing.T) {
	e := echo.New()
	serializer := &JSONSerializer{}
	e.JSONSerializer = serializer

	t.Run("Success", func(t *testing.T) {
		jsonBody := `{"hello":"world"}`
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var result map[string]string
		err := serializer.Deserialize(c, &result)
		require.NoError(t, err)
		assert.Equal(t, "world", result["hello"])
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		jsonBody := `{"hello":` // Invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(jsonBody))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var result map[string]string
		err := serializer.Deserialize(c, &result)
		assert.Error(t, err)
	})
}

func TestNewServer(t *testing.T) {
	dbPath := "../../data/city.db"

	t.Run("Success Initialization", func(t *testing.T) {
		// Check if file exists before trying, to avoid failing in CI environments without the DB.
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Skipf("Skipping test: Database file not found at %s", dbPath)
		}

		e, svc, err := NewServer(dbPath)
		require.NoError(t, err)
		require.NotNil(t, e)
		require.NotNil(t, svc)

		// Clean up
		svcErr := svc.DB.Close()
		require.NoError(t, svcErr)

		// Verify Routes are Registered
		foundRoutes := 0
		for _, r := range e.Routes() {
			if r.Path == "/health" && r.Method == http.MethodGet {
				foundRoutes++
			}
			if r.Path == "/lookup/:ip" && r.Method == http.MethodGet {
				foundRoutes++
			}
		}
		assert.Equal(t, 2, foundRoutes, "Expected /health and /lookup/:ip routes to be registered")
	})

	t.Run("Failure Invalid Path", func(t *testing.T) {
		e, svc, err := NewServer("invalid/path/to/db.mmdb")
		assert.Error(t, err)
		assert.Nil(t, e)
		assert.Nil(t, svc)
	})
}

func TestHealthCheck_Integration(t *testing.T) {
	dbPath := "../../data/city.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Skip("Skipping integration test: Database not found")
	}

	e, svc, err := NewServer(dbPath)
	require.NoError(t, err)
	defer func() {
		closeErr := svc.DB.Close()
		require.NoError(t, closeErr)
	}()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Inject request into Echo
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "status")
}
