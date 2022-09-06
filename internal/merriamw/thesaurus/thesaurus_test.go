//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package thesaurus_test

import (
	"fmt"
	"github.com/a-clap/dictionary/internal/merriamw/thesaurus"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

func init() {
}

type errThesauruser struct {
}

func (e errThesauruser) Get(_ string) ([]byte, error) {
	return nil, fmt.Errorf("handle me, please")
}

func TestThesaurus_Translate(t1 *testing.T) {
	thKey, ok := os.LookupEnv("MW_TH_KEY")
	if !ok {
		t1.Fatal("MW_TH_KEY not defined in ENV")
	}

	type fields struct {
		Thesauruser thesaurus.Thesauruser
	}
	type args struct {
		text string
	}
	type wants struct {
		w thesaurus.Word
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData wants
		wantErr  bool
	}{
		{
			name: "handle error gracefully",
			fields: fields{
				Thesauruser: errThesauruser{},
			},
			args:     args{},
			wantData: wants{},
			wantErr:  true,
		},
		{
			name: "test some obvious word \"world\"",
			fields: fields{
				Thesauruser: thesaurus.NewDefaultThesauruser(thKey),
			},
			args: args{
				text: "world",
			},
			wantData: wants{
				w: thesaurus.Word{
					Meta: struct {
						Id        string     `json:"id"`
						Syns      [][]string `json:"syns"`
						Ants      [][]string `json:"ants"`
						Offensive bool       `json:"offensive"`
					}{
						Id: "world",
						Syns: [][]string{{"folks", "humanity", "humankind", "people", "public", "species"}, {"earth", "globe", "planet"}, {"cosmos", "creation", "macrocosm", "nature", "universe"},
							{"galaxy", "light-year"}},
						Ants:      [][]string{},
						Offensive: false,
					},
					Fl:       "noun",
					Shortdef: []string{"human beings in general", "the celestial body on which we live", "the whole body of things observed or assumed"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(tester *testing.T) {
			t := thesaurus.NewThesaurus(tt.fields.Thesauruser)
			gotWords, err := t.Translate(tt.args.text)
			if (err != nil) != tt.wantErr {
				tester.Fatalf("%s: Translate() error = %v, wantErr %v", tester.Name(), err, tt.wantErr)
			}
			// If we want error, then rest of the testing may fail, because of nil dereference
			if tt.wantErr {
				if gotWords != nil {
					tester.Fatalf("%s: Translate() gotWords = %v, want nil", tester.Name(), gotWords)
				}
				return
			}
			if len(gotWords) == 0 {
				tester.Fatalf("%s: Translate() expected to receive anything", tester.Name())
			}

			diff := cmp.Diff(gotWords[0], &tt.wantData.w)
			if diff != "" {
				tester.Fatalf("%s: Translate(), diff: %s", tt.name, diff)
			}
		})
	}
}
