package meal

import (
	"reflect"
	"testing"
)

func TestParseText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    *Meal
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				text: `
				photo: https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Afghan_Palo.jpg/280px-Afghan_Palo.jpg
				===
				name: Pilav
				===
				instructions:
				1. Place the basmati rice in a large bowl and cover with hot water. Set aside.
				2. Wash the garlic head. Cut about a quarter inch off the top to expose the cloves. Set aside.
				===
				description: Afghan Palov`,
			},
			want: &Meal{
				ID:       0,
				Name:     "Pilav",
				PhotoURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Afghan_Palo.jpg/280px-Afghan_Palo.jpg",
				Instructions: `1. Place the basmati rice in a large bowl and cover with hot water. Set aside.
				2. Wash the garlic head. Cut about a quarter inch off the top to expose the cloves. Set aside.`,
				Description: "Afghan Palov",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseText(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseText() = %v, want %v", got, tt.want)
			}
		})
	}
}
