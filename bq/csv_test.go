package bq

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestReplaceTimeUnit(t *testing.T) {
	in := "9分26秒"
	want := "9m26s"
	out := replaceTimeUnit(in)
	if out != want {
		err := errors.Errorf("Fail.\n[out ]: %s\n[want]: %s", out, want)
		t.Error(err)
	}
}

func TestCreateCsvDetails(t *testing.T) {
	tweetInfo := TweetInfo{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl:  "https://example.com",
	}
	in_lines := []string{
		"本日の運動結果", "test", "R", "画面を撮影する",
		"リングコン押しこみ",
		"611回(3558回)",
		"リングコン下押しこみ",
		"9回(47回)",
		"アームツイスト",
		"282回(611回)",
		"リングコン引っぱり",
		"6回(66回)",
		"バンザイコシフリ",
		"235回(376回)",
		"おなか押しこみスクワット",
		"1回(30)",
		"モモアゲアゲ",
		"108回(1019回)",
		"ジョギング",
		"215m(4479m)",
		"ねじり体側のポーズ",
		"96回(96回)",
		"ダッシュ",
		"160m(9175m)",
		"ワイドスクワット",
		"88回(198回)",
		"モモあげ",
		"131m(534m)",
		"バンザイモーニング",
		"66回(251回)",
		"スクワットキープ",
		"10秒(158秒)",
		"英雄2のポーズ",
		"リングコン下押しこみキープ",
		"リングコン引っぱりキープ",
		"60回(156回)",
		"4秒(52秒)",
		"モモデプッシュ",
		"22回(266回)",
		"4秒(376秒)",
		"おなか押しこみ",
		"20回(182回)",
		"カッコ内はプレイ開始からの累計値です", "とじる", " ",
	}

	_, err := tweetInfo.createCsvDetails(in_lines)
	if err != nil {
		t.Error(err)
	}
}

func TestSetSummaryB(t *testing.T) {
	tweetInfo := TweetInfo{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl:  "https://example.com",
	}
	in_lines := []string{
		"本日の運動結果",
		"R 画面を撮影する",
		"test",
		"12分13秒",
		"合計活動時間",
		"10.",
		"合計消費カロリー",
		"11kcal",
		"0.14km",
		"合計走行距離",
		"次へ",
		" ",
	}

	want := Summary{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl: "https://example.com",
		TotalTimeExcercising: time.Duration(12*time.Minute + 13*time.Second),
		TotalCaloriesBurned: 10.11,
		TotalDistanceRun: 0.14,
	}
	s, err := tweetInfo.setSummary(in_lines, 3)
	if err != nil {
		t.Fatal(err)
	}
	if s[0].TotalTimeExcercising != want.TotalTimeExcercising {
		t.Errorf("act:%v, except: %v", s[0].TotalTimeExcercising, want.TotalTimeExcercising)
	}
	if s[0].TotalCaloriesBurned != want.TotalCaloriesBurned {
		t.Errorf("act:%f, except: %f", s[0].TotalCaloriesBurned, want.TotalCaloriesBurned)
	}
	if s[0].TotalDistanceRun != want.TotalDistanceRun {
		t.Errorf("act:%f, except: %f", s[0].TotalDistanceRun, want.TotalDistanceRun)
	}
}

func TestSetSummaryC(t *testing.T) {
	tweetInfo := TweetInfo{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl:  "https://example.com",
	}
	in_lines := []string{
		"本日の運動結果",
		"R 画面を撮影する",
		"test",
		"12分13秒",
		"合計活動時間",
		"10.e", // 消費カロリーが整数部と小数部で分割
		"合計消費カロリー",
		".11kcal",
		"0.14km",
		"合計走行距離",
		"次へ",
		" ",
	}

	want := Summary{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl: "https://example.com",
		TotalTimeExcercising: time.Duration(12*time.Minute + 13*time.Second),
		TotalCaloriesBurned: 10.11,
		TotalDistanceRun: 0.14,
	}
	s, err := tweetInfo.setSummary(in_lines, 3)
	if err != nil {
		t.Fatal(err)
	}
	if s[0].TotalTimeExcercising != want.TotalTimeExcercising {
		t.Errorf("act:%v, except: %v", s[0].TotalTimeExcercising, want.TotalTimeExcercising)
	}
	if s[0].TotalCaloriesBurned != want.TotalCaloriesBurned {
		t.Errorf("act:%f, except: %f", s[0].TotalCaloriesBurned, want.TotalCaloriesBurned)
	}
	if s[0].TotalDistanceRun != want.TotalDistanceRun {
		t.Errorf("act:%f, except: %f", s[0].TotalDistanceRun, want.TotalDistanceRun)
	}
}

func TestSetSummaryD(t *testing.T) {
	tweetInfo := TweetInfo{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl:  "https://example.com",
	}
	in_lines := []string{
		"本日の運動結果",
		"R 画面を撮影する",
		"test",
		"12 13", // 時間単位の欠落
		"合計活動時間",
		"10.11kcal",
		"合計消費カロリー",
		"0.14km",
		"合計走行距離",
		"次へ",
		" ",
	}

	want := Summary{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl: "https://example.com",
		TotalTimeExcercising: time.Duration(12*time.Minute + 13*time.Second),
		TotalCaloriesBurned: 10.11,
		TotalDistanceRun: 0.14,
	}
	s, err := tweetInfo.setSummary(in_lines, 3)
	if err != nil {
		t.Fatal(err)
	}
	if s[0].TotalTimeExcercising != want.TotalTimeExcercising {
		t.Errorf("act:%v, except: %v", s[0].TotalTimeExcercising, want.TotalTimeExcercising)
	}
	if s[0].TotalCaloriesBurned != want.TotalCaloriesBurned {
		t.Errorf("act:%f, except: %f", s[0].TotalCaloriesBurned, want.TotalCaloriesBurned)
	}
	if s[0].TotalDistanceRun != want.TotalDistanceRun {
		t.Errorf("act:%f, except: %f", s[0].TotalDistanceRun, want.TotalDistanceRun)
	}
}

func TestSetSummaryE(t *testing.T) {
	tweetInfo := TweetInfo{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl:  "https://example.com",
	}
	in_lines := []string{
		"本日の運動結果",
		"R 画面を撮影する",
		"test",
		"1時間12分13秒",
		"合計活動時間",
		"10.11kcal",
		"合計消費カロリー",
		"0.14km",
		"合計走行距離",
		"次へ",
		" ",
	}

	want := Summary{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl: "https://example.com",
		TotalTimeExcercising: time.Duration(1*time.Hour + 12*time.Minute + 13*time.Second),
		TotalCaloriesBurned: 10.11,
		TotalDistanceRun: 0.14,
	}
	s, err := tweetInfo.setSummary(in_lines, 3)
	if err != nil {
		t.Fatal(err)
	}
	if s[0].TotalTimeExcercising != want.TotalTimeExcercising {
		t.Errorf("act:%v, except: %v", s[0].TotalTimeExcercising, want.TotalTimeExcercising)
	}
	if s[0].TotalCaloriesBurned != want.TotalCaloriesBurned {
		t.Errorf("act:%f, except: %f", s[0].TotalCaloriesBurned, want.TotalCaloriesBurned)
	}
	if s[0].TotalDistanceRun != want.TotalDistanceRun {
		t.Errorf("act:%f, except: %f", s[0].TotalDistanceRun, want.TotalDistanceRun)
	}
}

func TestSetSummaryZ(t *testing.T) {
	tweetInfo := TweetInfo{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl:  "https://example.com",
	}
	in_lines := []string{
		"本日の運動結果",
		"R 画面を撮影する",
		"te7st", // 名前に数字
		"12 13", // 時間単位の欠落
		"合計活動時間",
		"分",
		"10.",
		"合計消費カロリー",
		".11kcal",
		"0.14km",
		"合計走行距離",
		"次へ",
		" ",
	}

	want := Summary{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl: "https://example.com",
		TotalTimeExcercising: time.Duration(12*time.Minute + 13*time.Second),
		TotalCaloriesBurned: 10.11,
		TotalDistanceRun: 0.14,
	}
	s, err := tweetInfo.setSummary(in_lines, 3)
	if err != nil {
		t.Fatal(err)
	}
	if s[0].TotalTimeExcercising != want.TotalTimeExcercising {
		t.Errorf("act:%v, except: %v", s[0].TotalTimeExcercising, want.TotalTimeExcercising)
	}
	if s[0].TotalCaloriesBurned != want.TotalCaloriesBurned {
		t.Errorf("act:%f, except: %f", s[0].TotalCaloriesBurned, want.TotalCaloriesBurned)
	}
	if s[0].TotalDistanceRun != want.TotalDistanceRun {
		t.Errorf("act:%f, except: %f", s[0].TotalDistanceRun, want.TotalDistanceRun)
	}
}

func TestSetSummaryY(t *testing.T) {
	tweetInfo := TweetInfo{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl:  "https://example.com",
	}
	in_lines := []string{
		"R 画面を撮影する",
		"本日の運動結果",
		"test",
		"9分1秒",
		"合計活動時間",
		"48.12", // 消費カロリーが
		"合計消費力ロリー",
		"12kcal", // 重複している
		"0.9km",
		"合計走行距離",
		"次へ",
		" ",
	}

	want := Summary{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl: "https://example.com",
		TotalTimeExcercising: time.Duration(9*time.Minute + 1*time.Second),
		TotalCaloriesBurned: 48.12,
		TotalDistanceRun: 0.9,
	}
	s, err := tweetInfo.setSummary(in_lines, 3)
	if err != nil {
		t.Fatal(err)
	}
	if s[0].TotalTimeExcercising != want.TotalTimeExcercising {
		t.Errorf("act:%v, except: %v", s[0].TotalTimeExcercising, want.TotalTimeExcercising)
	}
	if s[0].TotalCaloriesBurned != want.TotalCaloriesBurned {
		t.Errorf("act:%f, except: %f", s[0].TotalCaloriesBurned, want.TotalCaloriesBurned)
	}
	if s[0].TotalDistanceRun != want.TotalDistanceRun {
		t.Errorf("act:%f, except: %f", s[0].TotalDistanceRun, want.TotalDistanceRun)
	}
}

func TestSetSummary20220414(t *testing.T) {
	tweetInfo := TweetInfo{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl:  "https://example.com",
	}
	in_lines := []string{
		"R 画面を撮影する",
		"本日の運動結果",
		"年中無休でスワイショウ",
		"カイト",
		"22 5»",
		"分",
		"合計活動時間",
		"73.21kcal",
		"合計消費力ロリー",
		"合計走行距離",
		"1.26km",
		"次へ",
		"1.",
		" ",
	}

	want := Summary{
		TwitterId: "test",
		CreatedAt: time.Time{},
		ImageUrl: "https://example.com",
		TotalTimeExcercising: time.Duration(22*time.Minute + 5*time.Second),
		TotalCaloriesBurned: 73.21,
		TotalDistanceRun: 1.26,
	}
	s, err := tweetInfo.setSummary(in_lines, 3)
	if err != nil {
		t.Fatal(err)
	}
	if s[0].TotalTimeExcercising != want.TotalTimeExcercising {
		t.Errorf("act:%v, except: %v", s[0].TotalTimeExcercising, want.TotalTimeExcercising)
	}
	if s[0].TotalCaloriesBurned != want.TotalCaloriesBurned {
		t.Errorf("act:%f, except: %f", s[0].TotalCaloriesBurned, want.TotalCaloriesBurned)
	}
	if s[0].TotalDistanceRun != want.TotalDistanceRun {
		t.Errorf("act:%f, except: %f", s[0].TotalDistanceRun, want.TotalDistanceRun)
	}
}

func TestClassificateDetectedText(t *testing.T) {
	if class := classificateDetectedText("次へ"); class != SummaryText {
		t.Errorf("act:%v, except:%v", class, SummaryText)
	}
	if class := classificateDetectedText("Next"); class != SummaryText {
		t.Errorf("act:%v, except:%v", class, SummaryText)
	}
	if class := classificateDetectedText("とじる"); class != DetailsText {
		t.Errorf("act:%v, except:%v", class, SummaryText)
	}
	if class := classificateDetectedText("Close"); class != DetailsText {
		t.Errorf("act:%v, except:%v", class, SummaryText)
	}
	if class := classificateDetectedText("next"); class != UndefinedText {
		t.Errorf("act:%v, except:%v", class, SummaryText)
	}
}
