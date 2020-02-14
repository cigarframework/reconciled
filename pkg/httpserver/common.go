package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/storage"
)

var eol = []byte{'\n'}

func parseKey(key string) (group, kind, name string) {
	parts := strings.Split(key, "/")
	if len(parts) > 0 {
		group = parts[0]
	}
	if len(parts) > 1 {
		kind = parts[1]
	}
	if len(parts) > 2 {
		name = parts[2]
	}
	return
}

func httpError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError

	if errors.Is(err, api.ErrBadData) {
		code = http.StatusBadRequest
	} else if errors.Is(err, api.ErrNotExist) {
		code = http.StatusNotFound
	} else if errors.Is(err, api.ErrExist) {
		code = http.StatusConflict
	}

	http.Error(w, jsonError(err.Error()), code)
}

func jsonError(err string) string {
	return fmt.Sprintf("{\"error\":\"%s\"}", err)
}

func readState(reader io.Reader) (storage.State, error) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, api.ErrBadData
	}

	state := &storage.JSON{}
	if err := json.Unmarshal(body, state); err != nil {
		return nil, api.ErrBadData
	}

	return state, nil
}

func writeState(w http.ResponseWriter, state storage.State) {
	body, _ := json.Marshal(state)
	w.WriteHeader(200)
	w.Write(body)
}
