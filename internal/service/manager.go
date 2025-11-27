package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	"task/internal/model"
	"task/internal/storage"
)

type Manager struct {
	storage storage.Storage

	mu    sync.RWMutex
	state model.State

	jobs chan int64
}

func NewManager(st storage.Storage) (*Manager, error) {
	s, err := st.Load()
	if err != nil {
		return nil, err
	}
	if s.Batches == nil {
		s.Batches = make(map[int64]model.Batch)
	}

	return &Manager{
		storage: st,
		state:   s,
		jobs:    make(chan int64, 100),
	}, nil
}

func (m *Manager) StartWorkers(ctx context.Context, workers int) {
	for i := 0; i < workers; i++ {
		go m.worker(ctx)
	}

	m.mu.RLock()
	for id, b := range m.state.Batches {
		for _, st := range b.Status {
			if st == model.StatusPending {
				m.jobs <- id
				break
			}
		}
	}
	m.mu.RUnlock()
}

func (m *Manager) worker(ctx context.Context) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			return
		case batchID := <-m.jobs:
			m.processBatch(ctx, client, batchID)
		}
	}
}

func (m *Manager) processBatch(ctx context.Context, client *http.Client, id int64) {
	m.mu.RLock()
	batch, ok := m.state.Batches[id]
	m.mu.RUnlock()
	if !ok {
		return
	}

	for _, link := range batch.Links {
		select {
		case <-ctx.Done():
			return
		default:
		}
		batch.Status[link] = checkLink(ctx, client, link)
	}

	m.mu.Lock()
	m.state.Batches[id] = batch
	if err := m.storage.Save(m.state); err != nil {

	}
	m.mu.Unlock()
}

func (m *Manager) CreateBatch(links []string) (model.Batch, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.state.NextID
	m.state.NextID++

	statuses := make(map[string]model.LinkStatus, len(links))
	for _, l := range links {
		statuses[l] = model.StatusPending
	}

	b := model.Batch{
		ID:     id,
		Links:  links,
		Status: statuses,
	}

	m.state.Batches[id] = b

	if err := m.storage.Save(m.state); err != nil {
		return model.Batch{}, err
	}

	m.jobs <- id

	return b, nil
}

func (m *Manager) GetBatches(ids []int64) []model.Batch {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]model.Batch, 0, len(ids))
	for _, id := range ids {
		if b, ok := m.state.Batches[id]; ok {
			out = append(out, b)
		}
	}
	return out
}

func checkLink(ctx context.Context, client *http.Client, link string) model.LinkStatus {
	url := normalizeURL(link)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.StatusNotAvailable
	}

	resp, err := client.Do(req)
	if err != nil {
		return model.StatusNotAvailable
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return model.StatusAvailable
	}
	return model.StatusNotAvailable
}

func normalizeURL(s string) string {
	if len(s) > 8 && (s[:7] == "http://" || s[:8] == "https://") {
		return s
	}
	return "https://" + s
}
