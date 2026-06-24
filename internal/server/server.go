package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/GreenStarMatter/zenzore/internal/message"
	"github.com/GreenStarMatter/zenzore/internal/zyztem"
)

const PortEnvVar = "ZENZORE_PORT"

type Server struct {
	reg *registry
}

// Port resolves the port the root server should listen on.
func Port() (string, error) {
	port := os.Getenv(PortEnvVar)
	if port == "" {
		return "", fmt.Errorf("%s not set", PortEnvVar)
	}
	return port, nil
}

// generateID returns a random, collision-resistant hex string for
// use as a Zyztem's unique ID.
func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// registry holds every zyztem the server has created, in memory,
// for the lifetime of the process.
type registry struct {
	mu      sync.Mutex
	zyztems map[string]*zyztem.Zyztem
}

func (reg *registry) add(z *zyztem.Zyztem) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.zyztems[z.ID] = z
}

func (reg *registry) get(id string) (*zyztem.Zyztem, bool) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	z, ok := reg.zyztems[id]
	return z, ok
}

func (reg *registry) list() []*zyztem.Zyztem {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	out := make([]*zyztem.Zyztem, 0, len(reg.zyztems))
	for _, z := range reg.zyztems {
		out = append(out, z)
	}
	return out
}

func (reg *registry) remove(id string) error {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if _, exists := reg.zyztems[id]; !exists {
		return fmt.Errorf("zyztem %q does not exist", id)
	}
	delete(reg.zyztems, id)
	return nil
}

func (s *Server) SendAllZyztems(topicName string) error {
	for _, z := range s.reg.list() {
		psm := message.New()

		if err := psm.CreatePubSubClient(); err != nil {
			return fmt.Errorf("creating pubsub client for zyztem %q: %w", z.ID, err)
		}

		data, err := json.Marshal(z)
		if err != nil {
			psm.Client.Close()
			return fmt.Errorf("marshaling zyztem %q: %w", z.ID, err)
		}
		psm.AcceptGenericJson(data)

		err = psm.SendMessageToPubSub(topicName)
		psm.Client.Close()
		if err != nil {
			return fmt.Errorf("sending zyztem %q: %w", z.ID, err)
		}
	}
	return nil
}

// NewServer builds a Server with a fresh, empty registry.
func NewServer() *Server {
	return &Server{reg: &registry{zyztems: make(map[string]*zyztem.Zyztem)}}
}

func (s *Server) createZyztem(w http.ResponseWriter, r *http.Request) {
	id, err := generateID()
	if err != nil {
		http.Error(w, fmt.Sprintf("generating id: %v", err), http.StatusInternalServerError)
		return
	}

	z := zyztem.New()
	z.ID = id
	s.reg.add(z)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(z); err != nil {
		http.Error(w, fmt.Sprintf("encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) removeZyztem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id parameter", http.StatusBadRequest)
		return
	}

	if err := s.reg.remove(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listZyztems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.reg.list()); err != nil {
		http.Error(w, fmt.Sprintf("encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Server) augmentZyztem(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) sendAllZyztems(w http.ResponseWriter, r *http.Request) {
	topicName := os.Getenv(message.TOPIC_ID_ENV_VAR)
	if topicName == "" {
		http.Error(w, fmt.Sprintf("%s not set", message.TOPIC_ID_ENV_VAR), http.StatusInternalServerError)
		return
	}

	if err := s.SendAllZyztems(topicName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Run starts the root HTTP server and blocks until it is shut down,
// either via SIGINT/SIGTERM or an internal server error.
func (s *Server) Run() error {
	port, err := Port()
	if err != nil {
		return fmt.Errorf("resolving port: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/zyztems/create", s.createZyztem)
	mux.HandleFunc("/zyztems/list", s.listZyztems)
	mux.HandleFunc("/zyztems/augment", s.augmentZyztem)
	mux.HandleFunc("/zyztems/remove", s.removeZyztem)
	mux.HandleFunc("/zyztems/send", s.sendAllZyztems)
	// ... register other routes

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
