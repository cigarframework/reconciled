package httpserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/storage"
	"go.uber.org/zap"
)

type Server struct {
	stateServer api.Server
	options     *options
	logger      *zap.Logger
}

func New(stateServer api.Server, logger *zap.Logger, opts ...optionFunc) *Server {
	opt := &options{}
	for _, o := range opts {
		opt = o(opt)
	}

	return &Server{
		stateServer: stateServer,
		options:     opt,
		logger:      logger,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := api.WithToken(r.Context(), r.Header.Get(api.AuthHeader))

	group, kind, name := parseKey(r.URL.Path[1:])
	switch r.Method {
	case http.MethodDelete:
		{
			if err := s.stateServer.Remove(ctx, group, kind, name); err != nil {
				httpError(w, err)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte("{}"))
			return
		}
	case http.MethodPost:
		{
			state, err := readState(r.Body)
			if err != nil {
				httpError(w, err)
				return
			}

			res, err := s.stateServer.Create(ctx, state)
			if err != nil {
				httpError(w, err)
				return
			}

			writeState(w, res)
			return
		}
	case http.MethodPut:
		{
			state, err := readState(r.Body)
			if err != nil {
				httpError(w, err)
				return
			}

			res, err := s.stateServer.Update(ctx, state)
			if err != nil {
				httpError(w, err)
				return
			}

			writeState(w, res)
			return
		}
	case http.MethodGet:
		{
			if group != "" && kind != "" && name != "" {
				res, err := s.stateServer.Get(ctx, group, kind, name)
				if err != nil {
					httpError(w, err)
					return
				}

				writeState(w, res)
				return
			}

			watch := r.URL.Query().Get("watch") == "true"
			expression, err := url.QueryUnescape(r.URL.Query().Get("expression"))
			if err != nil {
				httpError(w, err)
				return
			}

			var req *api.WatchOptions
			if watch {
				req = &api.WatchOptions{BufferSize: s.options.streamBufferSize}
			}

			listOptions := &api.ListOptions{
				Expression: expression,
				Group:      r.URL.Query().Get("group"),
				Kind:       r.URL.Query().Get("kind"),
				Name:       r.URL.Query().Get("name"),
			}

			list, ch, err := s.stateServer.List(ctx, listOptions, req)

			if err != nil {
				httpError(w, err)
				return
			}

			if !watch {
				body, _ := json.Marshal(list)
				w.WriteHeader(200)
				w.Write(body)
				return
			}

			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			for _, state := range list {
				body, _ := json.Marshal(&api.Notification{State: state})
				if _, err := w.Write(body); err != nil {
					s.logger.Error(err.Error())
					return
				}

				if _, err := w.Write(eol); err != nil {
					s.logger.Error(err.Error())
					return
				}
			}

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

			for n := range ch {
				body, _ := json.Marshal(n)
				if _, err := w.Write(body); err != nil {
					s.logger.Error(err.Error())
					return
				}

				if _, err := w.Write(eol); err != nil {
					s.logger.Error(err.Error())
					return
				}

				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
			return
		}
	case http.MethodPatch:
		{
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				httpError(w, api.ErrBadData)
				return
			}

			var ops []*api.Patch
			if err := json.Unmarshal(body, &ops); err != nil {
				httpError(w, api.ErrBadData)
				return
			}

			res, err := s.stateServer.Patch(ctx, storage.NewMetaState(group, kind, name), ops)
			if err != nil {
				httpError(w, err)
				return
			}

			writeState(w, res)
			return
		}
	default:
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
	}
}
