package vision_texts

import (
	"testing"

	"github.com/tosh223/rfa/firestore"
)

func TestReplace(t *testing.T) {
	text := "R 画面を撮影する\n本日の運動結果\n年中無休でスワイショウ\nカイト\n22 44\n合計活動時間\n74.80kcal\n合計消費力ロリー\n0.47km\n合計走行距離\n次へ\n"
	replacers := []firestore.Replacer{
		firestore.Replacer{
			Before: "Om(",
			After: "0m(",
		},
		firestore.Replacer{
			Before:"0(",
			After:"回(",
		},
		firestore.Replacer{
			Before: "押しにみ",
			After: "押しこみ",
		},
		firestore.Replacer{
			Before: " m",
			After: "m",
		},
		firestore.Replacer{
			Before: "Im(",
			After: "1m(",
		},
		firestore.Replacer{
			Before: "- ",
			After: "",
		},
	}

	trueText := ReplaceFalse(text, replacers)
	if text != trueText {
		t.Error(text, trueText)
	}
	return
}

