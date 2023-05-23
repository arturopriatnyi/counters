package counter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestNewManager(t *testing.T) {
	s := NewMockStorage(gomock.NewController(t))

	m := NewManager(s)

	if !reflect.DeepEqual(m, &Manager{s: s}) {
		t.Errorf("want: %+v, got: %+v", &Manager{s: s}, m)
	}
}

func TestManager_Add(t *testing.T) {
	errUnexpected := errors.New("unexpected error")

	for name, tt := range map[string]struct {
		id      string
		s       func(*gomock.Controller) Storage
		wantErr error
	}{
		"OK": {
			id: "id",
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Get("id").
					Return(nil, ErrNotFound)
				s.
					EXPECT().
					Set(&Counter{ID: "id"}).
					Return(nil)

				return s
			},
			wantErr: nil,
		},
		"ErrExists": {
			id: "id",
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Get("id").
					Return(&Counter{ID: "id"}, nil)

				return s
			},
			wantErr: ErrExists,
		},
		"ErrUnexpected": {
			id: "id",
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Get("id").
					Return(nil, errUnexpected)

				return s
			},
			wantErr: errUnexpected,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{s: tt.s(gomock.NewController(t))}

			err := m.Add(tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestManager_Get(t *testing.T) {
	for name, tt := range map[string]struct {
		s           func(*gomock.Controller) Storage
		id          string
		wantCounter *Counter
		wantErr     error
	}{
		"OK": {
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Get("id").
					Return(&Counter{ID: "id", Value: 1}, nil)

				return s
			},
			id:          "id",
			wantCounter: &Counter{ID: "id", Value: 1},
			wantErr:     nil,
		},
		"ErrNotFound": {
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Get("id").
					Return(nil, ErrNotFound)

				return s
			},
			id:          "id",
			wantCounter: nil,
			wantErr:     ErrNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{s: tt.s(gomock.NewController(t))}

			c, err := m.Get(tt.id)

			if !reflect.DeepEqual(c, tt.wantCounter) {
				t.Errorf("want: %+v, got: %+v", tt.wantCounter, c)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestManager_Inc(t *testing.T) {
	errUnexpected := errors.New("unexpected error")

	for name, tt := range map[string]struct {
		id      string
		s       func(*gomock.Controller) Storage
		wantErr error
	}{
		"OK": {
			id: "id",
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Get("id").
					Return(&Counter{ID: "id", Value: 1}, nil)
				s.
					EXPECT().
					Set(&Counter{ID: "id", Value: 2}).
					Return(nil)

				return s
			},
			wantErr: nil,
		},
		"ErrNotFound": {
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Get("id").
					Return(nil, ErrNotFound)

				return s
			},
			id:      "id",
			wantErr: ErrNotFound,
		},
		"ErrUnexpected": {
			id: "id",
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Get("id").
					Return(&Counter{ID: "id", Value: 1}, nil)
				s.
					EXPECT().
					Set(&Counter{ID: "id", Value: 2}).
					Return(errUnexpected)

				return s
			},
			wantErr: errUnexpected,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{s: tt.s(gomock.NewController(t))}

			err := m.Inc(tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestManager_Delete(t *testing.T) {
	for name, tt := range map[string]struct {
		id      string
		s       func(*gomock.Controller) Storage
		wantErr error
	}{
		"OK": {
			id: "id",
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Delete("id").
					Return(nil)

				return s
			},
			wantErr: nil,
		},
		"ErrNotFound": {
			id: "id",
			s: func(c *gomock.Controller) Storage {
				s := NewMockStorage(c)

				s.
					EXPECT().
					Delete("id").
					Return(ErrNotFound)

				return s
			},
			wantErr: ErrNotFound,
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := &Manager{s: tt.s(gomock.NewController(t))}

			err := m.Delete(tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("want: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
