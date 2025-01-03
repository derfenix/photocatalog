package metadata_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	. "github.com/derfenix/photocatalog/v2/internal/metadata"
)

func TestDefault_Extract(t *testing.T) {
	t.Parallel()

	const timeFormat = "20060102_150405"

	tests := []struct {
		name        string
		filePath    string
		prefix      string
		timeFormat  string
		want        Metadata
		wantErr     bool
		wantErrType error
	}{
		{
			name:       "simple",
			filePath:   "20241108_160834.jpg",
			timeFormat: timeFormat,
			want: Metadata{
				Created: time.Date(2024, 11, 8, 16, 8, 34, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:       "simple in subdirectory",
			filePath:   "foo/20241108_160834.jpg",
			timeFormat: timeFormat,
			want: Metadata{
				Created: time.Date(2024, 11, 8, 16, 8, 34, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:       "simple with suffix",
			filePath:   "20241108_160834_(1).jpg",
			timeFormat: timeFormat,
			want: Metadata{
				Created: time.Date(2024, 11, 8, 16, 8, 34, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:       "prefix",
			filePath:   "foo_20241108_160834.jpg",
			timeFormat: timeFormat,
			prefix:     "foo_",
			want: Metadata{
				Created: time.Date(2024, 11, 8, 16, 8, 34, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:       "bad prefix",
			filePath:   "foo_20241108_160834.jpg",
			timeFormat: timeFormat,
			prefix:     "bar_",
			want: Metadata{
				Created: time.Time{},
			},
			wantErr:     true,
			wantErrType: ErrBadPrefix,
		},
		{
			name:       "no prefix",
			filePath:   "20241108_160834.jpg",
			timeFormat: timeFormat,
			prefix:     "bar_",
			want: Metadata{
				Created: time.Time{},
			},
			wantErr:     true,
			wantErrType: ErrBadPrefix,
		},
		{
			name:       "unexpected prefix",
			filePath:   "foo_20241108_160834.jpg",
			timeFormat: timeFormat,
			prefix:     "",
			want: Metadata{
				Created: time.Time{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &Default{
				TimeFormat: tt.timeFormat,
				Prefix:     tt.prefix,
			}

			got, err := d.Extract(tt.filePath, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantErrType != nil {
				if !errors.Is(err, tt.wantErrType) {
					t.Errorf("Extract() errorType = %v, wantErrType %v", err, tt.wantErrType)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract() got = %v, want %v", got, tt.want)
			}
		})
	}
}
