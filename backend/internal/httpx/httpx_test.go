package httpx_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vbsontakke/boilerplate/internal/apperror"
	"github.com/vbsontakke/boilerplate/internal/httpx"
	"github.com/vbsontakke/boilerplate/internal/testutil"
)

func TestStatusFor_MapsEveryKind(t *testing.T) {
	cases := map[apperror.Kind]int{
		apperror.KindValidation:      http.StatusBadRequest,
		apperror.KindUnauthorized:    http.StatusUnauthorized,
		apperror.KindForbidden:       http.StatusForbidden,
		apperror.KindNotFound:        http.StatusNotFound,
		apperror.KindConflict:        http.StatusConflict,
		apperror.KindUnprocessable:   http.StatusUnprocessableEntity,
		apperror.KindTooManyRequests: http.StatusTooManyRequests,
		apperror.KindInternal:        http.StatusInternalServerError,
	}
	for kind, want := range cases {
		assert.Equal(t, want, httpx.StatusFor(kind), "kind %d", kind)
	}
}

func TestError_HidesInternalDetails(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)

	httpx.Error(rec, req, errors.New("pq: password authentication failed for user secret"))

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "secret")
	body := testutil.DecodeJSON[httpx.ErrorBody](t, rec)
	assert.Equal(t, "INTERNAL_ERROR", body.Error.Code)
}

func TestError_WrappedAppErrorKeepsKind(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)

	err := apperror.Conflict("INVALID_TRANSITION", "cannot do that").Wrap(errors.New("db detail"))
	httpx.Error(rec, req, err)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.JSONEq(t, `{"error":{"code":"INVALID_TRANSITION","message":"cannot do that"}}`, rec.Body.String())
}

func TestSentinelComparison(t *testing.T) {
	err := apperror.NotFound("item not found")
	assert.True(t, errors.Is(err, apperror.ErrNotFound))
	assert.False(t, errors.Is(err, apperror.ErrForbidden))
}

type createReq struct {
	Title string `json:"title" validate:"required,min=5,max=100"`
	Price int64  `json:"price" validate:"required,gt=0"`
	Email string `json:"email" validate:"omitempty,email"`
}

func decode(body string) error {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(body))
	var dst createReq
	return httpx.DecodeAndValidate(rec, req, &dst)
}

func TestDecodeAndValidate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantErr    bool
		wantFields []string
	}{
		{"valid", `{"title":"Hello world","price":100}`, false, nil},
		{"empty body", ``, true, nil},
		{"malformed", `{"title":`, true, nil},
		{"unknown field", `{"title":"Hello world","price":1,"hack":true}`, true, []string{"hack"}},
		{"wrong type", `{"title":"Hello world","price":"abc"}`, true, []string{"price"}},
		{"two objects", `{"title":"Hello world","price":1}{}`, true, nil},
		{"trailing garbage", `{"title":"Hello world","price":1}]`, true, nil},
		{"trailing whitespace is fine", "{\"title\":\"Hello world\",\"price\":1}\n", false, nil},
		{"validation fails", `{"title":"Hi","price":0,"email":"nope"}`, true, []string{"title", "price", "email"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := decode(tt.body)
			if !tt.wantErr {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			ae := apperror.From(err)
			assert.Equal(t, apperror.KindValidation, ae.Kind)

			var got []string
			for _, f := range ae.Fields {
				got = append(got, f.Field)
			}
			for _, f := range tt.wantFields {
				assert.Contains(t, got, f)
			}
		})
	}
}

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		wantLimit int
		wantErr   bool
	}{
		{"default when absent", "", httpx.DefaultLimit, false},
		{"explicit limit", "?limit=5", 5, false},
		{"max allowed", "?limit=100", 100, false},
		{"above max", "?limit=101", 0, true},
		{"zero", "?limit=0", 0, true},
		{"negative", "?limit=-1", 0, true},
		{"not a number", "?limit=abc", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/x"+tt.query, nil)
			p, err := httpx.ParsePagination(r)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantLimit, p.Limit)
		})
	}
}

func TestNewPage_EmptyIsArrayNotNull(t *testing.T) {
	var items []string // nil slice, as a repository returns for "no rows"
	raw, err := json.Marshal(httpx.NewPage(items, "", 20))
	require.NoError(t, err)
	assert.JSONEq(t, `{"data":[],"next_cursor":"","limit":20}`, string(raw))
}

func TestCursorRoundTrip(t *testing.T) {
	type pos struct {
		ID   string `json:"id"`
		Rank int    `json:"rank"`
	}
	c, err := httpx.EncodeCursor(pos{ID: "abc", Rank: 7})
	require.NoError(t, err)

	var got pos
	require.NoError(t, httpx.DecodeCursor(c, &got))
	assert.Equal(t, pos{ID: "abc", Rank: 7}, got)

	assert.Error(t, httpx.DecodeCursor("%%%not-base64", &got))
}

func TestRequireUserID(t *testing.T) {
	const validID = "6f9619ff-8b86-4d11-b42d-00c04fc964ff"
	h := httpx.RequireUserID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := httpx.UserID(r.Context())
		require.True(t, ok)
		httpx.OK(w, map[string]string{"user_id": id.String()})
	}))

	tests := []struct {
		name     string
		header   string
		wantCode int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"not a uuid", "not-a-uuid", http.StatusUnauthorized},
		{"valid uuid", validID, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := map[string]string{}
			if tt.header != "" {
				headers[httpx.UserIDHeader] = tt.header
			}
			rec := testutil.Do(t, h, testutil.Request{Method: http.MethodGet, Path: "/", Headers: headers})
			require.Equal(t, tt.wantCode, rec.Code)
			if tt.wantCode == http.StatusUnauthorized {
				body := testutil.DecodeJSON[httpx.ErrorBody](t, rec)
				assert.Equal(t, "UNAUTHORIZED", body.Error.Code)
			}
		})
	}
}
