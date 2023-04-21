package github

import (
	"reflect"
	"testing"
)

func TestPkg_LastTag(t *testing.T) {
	type fields struct {
		Owner    string
		Repo     string
		UseProxy bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Tag
		wantErr bool
	}{
		{
			name: "v0.3.15",
			fields: fields{
				Owner:    "ergoapi",
				Repo:     "util",
				UseProxy: false,
			},
			want: &Tag{
				Name: "v0.3.15",
				Commit: Commit{
					SHA: "85178fd6fafa3d7d2cf351f0a450d380b53433e0",
					URL: "https://api.github.com/repos/ergoapi/util/commits/85178fd6fafa3d7d2cf351f0a450d380b53433e0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pkg{
				Owner: tt.fields.Owner,
				Repo:  tt.fields.Repo,
			}
			got, err := p.LastTag()
			if (err != nil) != tt.wantErr {
				t.Errorf("LastTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LastTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}
