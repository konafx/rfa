package vision_texts

import (
	"strings"

	"github.com/tosh223/rfa/firestore"
)

func ReplaceFalse(text string, replacers []firestore.Replacer) (trueText string) {
	trueText = text
	for _, replacer := range replacers {
		if len(replacer.Before) == 0 && len(replacer.After) == 0 {
			continue
		}
		trueText = strings.ReplaceAll(trueText, replacer.Before, replacer.After)
	}

	return
}
