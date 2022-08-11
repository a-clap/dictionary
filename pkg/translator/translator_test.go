//  Copyright 2022 a-clap. All rights reserved.
//  Use of this source code is governed by a MIT-style
//  license that can be found in the LICENSE file.

package translator_test

import (
	"github.com/a-clap/dictionary/internal/deepl"
	"github.com/a-clap/dictionary/internal/logger"
	"github.com/a-clap/dictionary/internal/merriamw/dictionary"
	"github.com/a-clap/dictionary/pkg/translator"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestTranslator_Get(t *testing.T) {
	deeplKey, ok := os.LookupEnv("DEEPL_KEY")
	if !ok {
		t.Fatal("DEEPL_KEY not found in ENV")
		return
	}
	dictKey, ok := os.LookupEnv("MW_DICT_KEY")
	if !ok {
		t.Fatal("MW_DICT_KEY not found in ENV")
		return
	}
	thKey, ok := os.LookupEnv("MW_TH_KEY")
	if !ok {
		t.Fatal("MW_TH_KEY not found in ENV")
		return
	}

	type fields struct {
		translate translator.Translate
	}
	type args struct {
		text string
		from deepl.SourceLang
		to   deepl.TargetLang
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *translator.Translation
		wantErr bool
	}{
		{
			name: "TarPolish doesn't contain neither dictionary neither thesaurus",
			fields: fields{
				translate: translator.NewStandard(deeplKey, dictKey, thKey, logger.NewDummy()),
			},
			args: args{
				text: "brain",
				to:   deepl.TarPolish,
			},
			want: &translator.Translation{
				Deepl: []translator.DeeplTranslate{
					{
						Text: "mózg",
					},
				},
				Dictionary: nil,
				Thesaurus:  nil,
			},
			wantErr: false,
		},
		{
			name: "dict json",
			fields: fields{
				translate: translator.NewStandard(deeplKey, dictKey, thKey, logger.NewDummy()),
			},
			args: args{
				text: "mózg",
				to:   deepl.TarEnglishBritish,
			},
			want: &translator.Translation{
				Deepl: []translator.DeeplTranslate{
					{
						Text: "brain",
					},
				},
				Dictionary: &translator.DictionaryTranslate{
					Defs: []translator.Definition{
						{
							Offensive: false,
							Function:  "noun",
							Examples:  []string{},
							Definition: []string{
								"the portion of the vertebrate central nervous system enclosed in the skull and continuous with the spinal cord through the foramen magnum that is composed of neurons and supporting and nutritive structures (such as glia) and that integrates sensory information from inside and outside the body in controlling autonomic function (such as heartbeat and respiration), in coordinating and directing correlated motor responses, and in the process of learning",
								"a nervous center in invertebrates comparable in position and function to the vertebrate brain",
								"intellect, mind",
							},
							Audio: []dictionary.Pronunciation{
								{
									PhoneticNotation: "ˈbrān",
									Url:              "https://media.merriam-webster.com/audio/prons/en/us/mp3/b/brain001.mp3",
								},
							},
						},
						{
							Offensive: false,
							Function:  "verb",
							Examples:  []string{},
							Definition: []string{
								"to kill by smashing the skull",
								"to hit on the head",
							},
							Audio: []dictionary.Pronunciation{},
						},
					},
					Synonyms: []string{
						"brain attack",
						"brain coral",
						"brain cramp",
						"brain-dead",
						"brain death",
						"brain drain",
						"brain dump",
						"brain fog",
					},
				},
				Thesaurus: []translator.ThesaurusTranslate{
					{
						Text: "brain",
						Synonyms: []string{
							"brainiac",
							"genius",
							"intellect",
							"thinker",
							"whiz",
							"wiz",
							"wizard",
						},
						Antonyms: []string{
							"blockhead",
							"dodo",
							"dolt",
							"dope",
							"dumbbell",
							"dummy",
							"dunce",
							"fathead",
							"goon",
							"half-wit",
							"hammerhead",
							"idiot",
							"imbecile",
							"knucklehead",
							"moron",
							"nitwit",
							"numskull",
							"pinhead",
						},
						Offensive: false,
						Function:  "noun",
						Definition: []string{
							"a very smart person",
							"the ability to learn and understand or to deal with problems",
							"the part of a person that feels, thinks, perceives, wills, and especially reasons",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := require.New(t)
			translate := tt.fields.translate

			got, err := translate.Get(tt.args.text, tt.args.from, tt.args.to)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				req.Nil(got)
				return
			}

			req.Equal(got.Deepl, tt.want.Deepl, "got = %#v, want %#v", got.Deepl, tt.want.Deepl)
			req.Equal(got.Dictionary, tt.want.Dictionary, "got = %#v, want %#v", got.Dictionary, tt.want.Dictionary)
			req.Equal(got.Thesaurus, tt.want.Thesaurus, "got = %#v, want %#v", got.Thesaurus, tt.want.Thesaurus)
		})
	}
}
