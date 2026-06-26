package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GreenStarMatter/zenzore/internal/zyztem"
	"github.com/stretchr/testify/assert"
)

func TestPort_ReturnsEnvVarWhenSet(t *testing.T) {
	t.Setenv(PortEnvVar, "9090")

	port, err := Port()
	assert.NoError(t, err)
	assert.Equal(t, "9090", port)
}

func TestPort_ErrorsWhenUnset(t *testing.T) {
	t.Setenv(PortEnvVar, "")

	_, err := Port()
	assert.Error(t, err)
}

func TestGenerateID_ReturnsNonEmptyUniqueValues(t *testing.T) {
	id1, err := generateID()
	assert.NoError(t, err)
	assert.NotEmpty(t, id1)

	id2, err := generateID()
	assert.NoError(t, err)
	assert.NotEqual(t, id1, id2)
}

func TestRegistry_AddAndGet(t *testing.T) {
	s := NewServer()

	z := &zyztem.Zyztem{ID: "abc"}
	s.reg.add(z)

	got, ok := s.reg.get("abc")
	assert.True(t, ok)
	assert.Equal(t, z, got)
}

func TestRegistry_GetMissingReturnsFalse(t *testing.T) {
	s := NewServer()

	_, ok := s.reg.get("does-not-exist")
	assert.False(t, ok)
}

func TestRegistry_List(t *testing.T) {
	s := NewServer()

	s.reg.add(&zyztem.Zyztem{ID: "a"})
	s.reg.add(&zyztem.Zyztem{ID: "b"})

	got := s.reg.list()
	assert.Equal(t, 2, len(got))
}

func TestCreateZyztem(t *testing.T) {
	s := NewServer()

	req := httptest.NewRequest(http.MethodPost, "/zyztems/create", nil)
	rec := httptest.NewRecorder()

	s.createZyztem(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got zyztem.Zyztem
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.NotEmpty(t, got.ID)

	assert.Equal(t, 1, len(s.reg.list()), "expected created zyztem to be stored in the registry")
}

func TestCreateZyztem_GeneratesUniqueIDsAcrossCalls(t *testing.T) {
	s := NewServer()

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/zyztems/create", nil)
		rec := httptest.NewRecorder()
		s.createZyztem(rec, req)
	}

	assert.Equal(t, 3, len(s.reg.list()), "expected each create call to produce a distinct, stored zyztem")
}

func TestListZyztems(t *testing.T) {
	s := NewServer()
	s.reg.add(&zyztem.Zyztem{ID: "a"})
	s.reg.add(&zyztem.Zyztem{ID: "b"})

	req := httptest.NewRequest(http.MethodGet, "/zyztems/list", nil)
	rec := httptest.NewRecorder()

	s.listZyztems(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got []*zyztem.Zyztem
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(got))
}

func TestListZyztems_Empty(t *testing.T) {
	s := NewServer()

	req := httptest.NewRequest(http.MethodGet, "/zyztems/list", nil)
	rec := httptest.NewRecorder()

	s.listZyztems(rec, req)

	var got []*zyztem.Zyztem
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(got))
}

func TestAugmentZyztem_NotImplemented(t *testing.T) {
	s := NewServer()

	req := httptest.NewRequest(http.MethodPost, "/zyztems/augment", nil)
	rec := httptest.NewRecorder()

	s.augmentZyztem(rec, req)

	assert.Equal(t, http.StatusNotImplemented, rec.Code)
}

func TestCreateAndRemoveZyztem(t *testing.T) {
	s := NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /zyztems", s.createZyztem)
	mux.HandleFunc("DELETE /zyztems/{id}", s.removeZyztem)

	// create
	createReq := httptest.NewRequest(http.MethodPost, "/zyztems", nil)
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	var created zyztem.Zyztem
	assert.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &created))

	// delete
	removeReq := httptest.NewRequest(http.MethodDelete, "/zyztems/"+created.ID, nil)
	removeRec := httptest.NewRecorder()
	mux.ServeHTTP(removeRec, removeReq)

	assert.Equal(t, http.StatusNoContent, removeRec.Code)

	// delete again
	removeAgainReq := httptest.NewRequest(http.MethodDelete, "/zyztems/"+created.ID, nil)
	removeAgainRec := httptest.NewRecorder()
	mux.ServeHTTP(removeAgainRec, removeAgainReq)

	assert.Equal(t, http.StatusNotFound, removeAgainRec.Code)
}
