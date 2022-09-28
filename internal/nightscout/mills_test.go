package nightscout

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestMills_MarshalJSON(t *testing.T) {
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
		wantErr bool
	}{
		{"0", fields{unix0}, []byte("0"), false},
		{"now", fields{now}, []byte(nowStr), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mills{
				Time: tt.fields.Time,
			}
			got, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMills_UnmarshalJSON(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)

	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Mills
		wantErr bool
	}{
		{"now", args{[]byte(strconv.Itoa(int(now.UnixMilli())))}, Mills{now}, false},
		{"error", args{[]byte("a")}, Mills{time.Time{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Mills
			if err := m.UnmarshalJSON(tt.args.bytes); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(m, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", m, tt.want)
			}
		})
	}
}
