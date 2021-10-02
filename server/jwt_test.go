package server

import (
	_ "crypto/sha256"
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"parse0", args{[]byte(`{"alg":"sha1","typ":"jwt"}`)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parse(tt.args.b)
		})
	}
}

func TestToken_Sign(t *testing.T) {
	type fields struct {
		Header *Header
		Claims *Claims
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "sign0",
			fields: fields{
				Header: &Header{"SHA256", "jwt"},
				Claims: &Claims{1111, "Leng0"},
			},
			want:    []byte("xxx"),
			wantErr: false,
		},
		{
			name: "sign1",
			fields: fields{
				Header: &Header{"SHA256", "jwt"},
				Claims: &Claims{2222, "Leng1"},
			},
			want:    []byte("xxx"),
			wantErr: false,
		},
		{
			name: "sign2",
			fields: fields{
				Header: &Header{"SHA256", "jwt"},
				Claims: &Claims{3333, "Leng2"},
			},
			want:    []byte("xxx"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Token{
				Header: tt.fields.Header,
				Claims: tt.fields.Claims,
			}
			got, err := tr.Sign()
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Token.Sign() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestVarify(t *testing.T) {
	type args struct {
		token []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Claims
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "varify0",
			args:    args{[]byte("e30.eyJVaWQiOjExMTEsIk5hbWUiOiJMZW5nMCJ9.FUuVfNLqiDJ4952WL6_-U50IXT8dC1lisoac92vQY8E")},
			want:    &Claims{1111, "Leng0"},
			wantErr: false,
		},
		{
			name:    "varify1",
			args:    args{[]byte("e30.eyJVaWQiOjIyMjIsIk5hbWUiOiJMZW5nMSJ9.9MQdQ4Tm4LxgG13M9ySSQxHyNxPJwyFYdZ9LODsMNUg")},
			want:    &Claims{2222, "Leng1"},
			wantErr: false,
		},
		{
			name:    "varify2",
			args:    args{[]byte("e30.eyJVaWQiOjMzMzMsIk5hbWUiOiJMZW5nMiJ9.ugzV3hXLqrpF4hTGViEorC34l6SmnWT4i0dbKPrt88U")},
			want:    &Claims{3333, "Leng2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Varify(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Varify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Varify() = %v, want %v", got, tt.want)
			}
			t.Logf("%v", got)
		})
	}
}
