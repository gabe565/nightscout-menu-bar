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
