package nightscout

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMills_MarshalJSON(t *testing.T) {
	t.Parallel()
	unix0 := time.Unix(0, 0)

	now := time.Now()
	nowStr := strconv.Itoa(int(now.UnixMilli()))

	type fields struct {
		Time time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr require.ErrorAssertionFunc
	}{
		{"0", fields{unix0}, []byte("0"), require.NoError},
		{"now", fields{now}, []byte(nowStr), require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Mills{
				Time: tt.fields.Time,
			}
			got, err := m.MarshalJSON()
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMills_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	now := time.Now().Truncate(time.Millisecond)

	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Mills
		wantErr require.ErrorAssertionFunc
	}{
		{"now", args{[]byte(strconv.Itoa(int(now.UnixMilli())))}, Mills{now}, require.NoError},
		{"error", args{[]byte("a")}, Mills{time.Time{}}, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var m Mills
			tt.wantErr(t, m.UnmarshalJSON(tt.args.bytes))
			assert.Equal(t, tt.want, m)
		})
	}
}

func TestMills_Relative(t *testing.T) {
	t.Parallel()
	type fields struct {
		Time time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"0m", fields{time.Now()}, "0m"},
		{"59s", fields{time.Now().Add(-59 * time.Second)}, "0m"},
		{"1m", fields{time.Now().Add(-time.Minute)}, "1m"},
		{"1m30s", fields{time.Now().Add(-time.Minute - 30*time.Second)}, "1m"},
		{"2m35s", fields{time.Now().Add(-2*time.Minute - 35*time.Second)}, "2m"},
		{"4m15s", fields{time.Now().Add(-4*time.Minute - 15*time.Second)}, "4m"},
		{"5m1s", fields{time.Now().Add(-5*time.Minute - time.Second)}, "5m"},
		{"now", fields{time.Now()}, "0m"},
		{"unix 0", fields{time.Unix(0, 0)}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Mills{Time: tt.fields.Time}
			assert.Equal(t, tt.want, m.Relative())
		})
	}
}
