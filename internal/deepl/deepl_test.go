//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package deepl_test

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/deepl"
	"github.com/a-clap/dictionary/internal/logger"
	"log"
	"os"
	"reflect"
	"testing"
)

type fakeAccess struct {
	translation string
	expected    deepl.Word
	generateErr bool
}

func (f fakeAccess) Query(_ string, sourceLang deepl.SourceLang, _ deepl.TargetLang) ([]byte, error) {
	if f.generateErr {
		return nil, fmt.Errorf("generating random err")
	}
	d := deepl.Translations{
		DetectedSourceLanguage: string(sourceLang),
		Text:                   f.translation,
	}

	f.expected.Translations = append(f.expected.Translations, d)
	return json.Marshal(f.expected)
}

func TestDeepL_Translate(t *testing.T) {
	api, ok := os.LookupEnv("DEEPL_KEY")
	if !ok {
		log.Fatalf("DEEPL_KEY Not found")
	}

	type in struct {
		acc  deepl.Access
		text string
		src  deepl.SourceLang
		dst  deepl.TargetLang
	}
	type out struct {
		w        *deepl.Word
		err      bool
		expected string
	}
	type args struct {
		in  in
		out out
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "handle error",
			args: args{
				in: in{
					acc: fakeAccess{
						translation: "",
						expected:    deepl.Word{},
						generateErr: true,
					},
					text: "brain",
					src:  deepl.English,
					dst:  deepl.Polish,
				},
				out: out{
					w:        nil,
					err:      true,
					expected: "",
				},
			},
		},
		{
			name: "handle translation",
			args: args{
				in: in{
					acc: fakeAccess{
						translation: "mózg",
						expected:    deepl.Word{},
						generateErr: false,
					},
					text: "brain",
					src:  deepl.English,
					dst:  deepl.Polish,
				},
				out: out{
					w: &deepl.Word{Translations: []deepl.Translations{
						{
							DetectedSourceLanguage: string(deepl.English),
							Text:                   "mózg",
						},
					},
					},
					err:      false,
					expected: "mózg",
				},
			},
		},
		{
			name: "direct call EN -> PL",
			args: args{
				in: in{
					acc:  deepl.NewAccessDefault(api, logger.NewDummy()),
					text: "brain",
					src:  deepl.English,
					dst:  deepl.Polish,
				},
				out: out{
					w: &deepl.Word{Translations: []deepl.Translations{
						{
							DetectedSourceLanguage: string(deepl.English),
							Text:                   "mózg",
						},
					}},
					err: false,
				},
			},
		},
		{
			name: "direct call PL -> EN",
			args: args{
				in: in{
					acc:  deepl.NewAccessDefault(api, logger.NewDummy()),
					text: "mózg",
					src:  deepl.Polish,
					dst:  deepl.English,
				},
				out: out{
					w: &deepl.Word{Translations: []deepl.Translations{
						{
							DetectedSourceLanguage: string(deepl.Polish),
							Text:                   "brain",
						},
					}},
					err: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := deepl.NewDeepL(tt.args.in.acc, logger.NewDummy())
			w, err := d.Translate(tt.args.in.text, tt.args.in.src, tt.args.in.dst)

			if (err != nil) != tt.args.out.err {
				t.Errorf("%s: unexpected error %#v", t.Name(), err)
			}

			if !reflect.DeepEqual(w, tt.args.out.w) {
				t.Errorf("%s: got %#v, want %#v", t.Name(), w, &tt.args.out.w)
			}
		})
	}
}
