package counter

import (
	"reflect"
	"testing"
)

func TestNewMemoryStorage(t *testing.T) {
	want := &MemoryStorage{counters: map[string]*Counter{}}

	got := NewMemoryStorage()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %+v, got: %+v", want, got)
	}
}

func TestMemoryStorage_Set(t *testing.T) {
	counter := &Counter{ID: "id", Value: 1}
	s := &MemoryStorage{counters: map[string]*Counter{}}

	err := s.Set(counter)

	if err != nil {
		t.Errorf("want: <nil>, got: %v", err)
	}
	if c, ok := s.counters[counter.ID]; !ok || !reflect.DeepEqual(c, counter) {
		t.Errorf("want: %+v, got: %+v", counter, c)
	}
}

func TestMemoryStorage_Get(t *testing.T) {
	for name, tt := range map[string]struct {
		s           *MemoryStorage
		id          string
		wantCounter *Counter
		wantErr     error
	}{
		"OK": {
			s: &MemoryStorage{
				counters: map[string]*Counter{
					"id": {ID: "id", Value: 1},
				},
			},
			id:          "id",
			wantCounter: &Counter{ID: "id", Value: 1},
			wantErr:     nil,
		},
		"ErrNotFound": {
			s:           &MemoryStorage{counters: map[string]*Counter{}},
			id:          "id",
			wantCounter: nil,
			wantErr:     ErrNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, err := tt.s.Get(tt.id)

			if !reflect.DeepEqual(c, tt.wantCounter) {
				t.Errorf("want: %+v, got: %+v", tt.wantCounter, c)
			}
			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestMemoryStorage_Delete(t *testing.T) {
	for name, tt := range map[string]struct {
		s       *MemoryStorage
		id      string
		wantErr error
	}{
		"OK": {
			s: &MemoryStorage{
				counters: map[string]*Counter{
					"id": {ID: "id", Value: 1},
				},
			},
			id:      "id",
			wantErr: nil,
		},
		"ErrNotFound": {
			s:       &MemoryStorage{counters: map[string]*Counter{}},
			id:      "id",
			wantErr: ErrNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := tt.s.Delete(tt.id)

			if err != tt.wantErr {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
