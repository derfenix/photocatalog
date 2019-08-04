package metadata

import (
	"reflect"
	"testing"
	"time"
)

func TestDefaultExtractor_Extract(t *testing.T) {
	local, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatal(err)
	}
	time.Local = local

	type args struct {
		fp string
	}
	tests := []struct {
		name    string
		args    args
		want    Metadata
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{"20190321_120325.jpg"},
			want: Metadata{
				Time: func() time.Time {
					tm, err := time.Parse(time.RFC3339, "2019-03-21T12:03:25+03:00")
					if err != nil {
						t.Fatal(err)
					}
					return tm
				}(),
			},
			wantErr: false,
		},
		{
			name: "Name with tilda",
			args: args{"20190321_120325~2.jpg"},
			want: Metadata{
				Time: func() time.Time {
					tm, err := time.Parse(time.RFC3339, "2019-03-21T12:03:25+03:00")
					if err != nil {
						t.Fatal(err)
					}
					return tm
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDefaultExtractor()
			if err != nil {
				t.Fatal(err)
			}
			got, err := d.Extract(tt.args.fp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract() got = %v, want %v", got, tt.want)
			}
		})
	}
}
