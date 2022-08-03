package translator_test

import (
	"github.com/a-clap/dictionary/pkg/translator"
	"testing"
)

func TestTranslate(t *testing.T) {
	type args struct {
		text string
		from translator.Language
	}
	tests := []struct {
		name string
		args args
		want translator.Translation
	}{
		{
			name: "basic test",
			args: args{
				text: "brain",
				from: translator.English,
			},
			want: translator.Translation{
				Text:        "brain",
				Translation: "mózg",
				Antonyms:    nil,
				Synonyms:    nil,
				IsOffensive: false,
				Audio:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translator.Translate(tt.args.text, tt.args.from)
			if got.Text != tt.want.Text || got.Translation != tt.want.Translation {
				t.Errorf("Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}
