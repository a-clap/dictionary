//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package deepl_test

import (
	"encoding/json"
	"fmt"
	"github.com/a-clap/dictionary/internal/deepl"
	"github.com/a-clap/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
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

func init() {
	logger.Init(logger.NewDefaultZap(zapcore.DebugLevel))
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
		acc  deepl.Deepler
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
					src:  deepl.SrcEnglish,
					dst:  deepl.TarPolish,
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
					src:  deepl.SrcEnglish,
					dst:  deepl.TarPolish,
				},
				out: out{
					w: &deepl.Word{Translations: []deepl.Translations{
						{
							DetectedSourceLanguage: string(deepl.SrcEnglish),
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
					acc:  deepl.NewDeeplerDefault(api),
					text: "brain",
					src:  deepl.SrcEnglish,
					dst:  deepl.TarPolish,
				},
				out: out{
					w: &deepl.Word{Translations: []deepl.Translations{
						{
							DetectedSourceLanguage: string(deepl.SrcEnglish),
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
					acc:  deepl.NewDeeplerDefault(api),
					text: "mózg",
					src:  deepl.SrcPolish,
					dst:  deepl.TarEnglishBritish,
				},
				out: out{
					w: &deepl.Word{Translations: []deepl.Translations{
						{
							DetectedSourceLanguage: string(deepl.SrcPolish),
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
			d := deepl.NewDeepL(tt.args.in.acc)
			require.NotNil(t, d)
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
