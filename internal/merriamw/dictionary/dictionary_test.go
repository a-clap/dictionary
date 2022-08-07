//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package dictionary_test

import (
	"fmt"
	"github.com/a-clap/dictionary/internal/logger"
	"github.com/a-clap/dictionary/internal/merriamw/dictionary"
	"os"
	"reflect"
	"testing"
)

type errDefinitioner struct {
}

func (e errDefinitioner) Get(_ string) ([]byte, error) {
	return nil, fmt.Errorf("handle me, please")
}

func TestDictionary_Definition(t *testing.T) {
	dictKey, ok := os.LookupEnv("MW_DICT_KEY")
	if !ok {
		t.Fatal("MW_DICT_KEY not defined in ENV")
	}

	type fields struct {
		Definitioner dictionary.Definitioner
		Logger       logger.Logger
	}
	type args struct {
		text string
	}
	type wants struct {
		text        string
		definition  []string
		suggestions *dictionary.Suggestions
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
				Definitioner: errDefinitioner{},
				Logger:       logger.NewDummy(),
			},
			args: args{
				text: "",
			},
			wantData: wants{
				text:        "",
				definition:  nil,
				suggestions: nil,
			},
			wantErr: true,
		},
		{
			name: "test some obvious word \"world\"",
			fields: fields{
				Definitioner: dictionary.NewDefaultGetDefinition(dictKey),
				Logger:       logger.NewDummy(),
			},
			args: args{
				text: "world",
			},
			wantData: wants{
				text:        "world",
				definition:  []string{"the earthly state of human existence", "life after death —used with a qualifier", "the earth with its inhabitants and all things upon it"},
				suggestions: nil,
			},
			wantErr: false,
		},
		{
			name: "test typo \"warld\"",
			fields: fields{
				Definitioner: dictionary.NewDefaultGetDefinition(dictKey),
				Logger:       logger.NewDummy(),
			},
			args: args{
				text: "warld",
			},
			wantData: wants{
				text:       "",
				definition: nil,
				suggestions: &dictionary.Suggestions{Suggestions: []string{"world", "wared", "varlet", "warbled", "worlds", "warlord", "worldly", "warble", "warded", "warily", "warmed", "warmly",
					"warned", "warlords", "wards",
					"wares", "warms", "warns", "warps", "warts"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dictionary.NewDictionary(tt.fields.Definitioner, tt.fields.Logger)
			gotData, suggestions, err := d.Definition(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s: Definition() error = %v, wantErr %v", t.Name(), err, tt.wantErr)
			}
			// If we want error, then rest of testing may fail because of "index out of range"
			if tt.wantErr {
				if gotData != nil {
					t.Fatalf("%s: Definition() gotData = %#v, want %v", t.Name(), gotData, nil)
				}
				if suggestions != nil {
					t.Fatalf("%s: Definition() suggestions = %#v, want %v", t.Name(), suggestions, nil)
				}
				return
			}

			if !reflect.DeepEqual(suggestions, tt.wantData.suggestions) {
				t.Fatalf("%s: Definition() suggestions = %#v, want %v", t.Name(), suggestions, tt.wantData.suggestions)
			}

			if !reflect.DeepEqual(gotData[0].Definition(), tt.wantData.definition) {
				t.Fatalf("%s: Definition() gotData.Definition() = %#v, want %v", t.Name(), gotData[0].Definition(), tt.wantData.definition)
			}

			if !reflect.DeepEqual(gotData[0].Text(), tt.wantData.text) {
				t.Fatalf("%s: Definition() gotData.Text() = %#v, want %v", t.Name(), gotData[0].Text(), tt.wantData.text)
			}
		})
	}
}
