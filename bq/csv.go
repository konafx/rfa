package bq

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

type Summary struct {
	TwitterId            string        `json:"twitter_id" csv:"twitter_id"`
	TotalTimeExcercising time.Duration `json:"total_time_excercising" csv:"total_time_excercising"`
	TotalCaloriesBurned  float64       `json:"total_calories_burned" csv:"total_calories_burned"`
	TotalDistanceRun     float64       `json:"total_distance_run" csv:"total_distance_run"`
	CreatedAt            time.Time     `json:"created_at" csv:"created_at"`
}

type Details struct {
	TwitterId     string    `json:"twitter_id" csv:"twitter_id"`
	ExerciseName  string    `json:"exercise_name" csv:"exercise_name"`
	Quantity      int       `json:"quantity" csv:"quantity"`
	TotalQuantity int       `json:"total_quantity" csv:"total_quantity"`
	CreatedAt     time.Time `json:"created_at" csv:"created_at"`
}

func CreateCsv(twitterId string, createdAtStr string, text string) string {
	var csvName string = ""

	createdAt, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", createdAtStr)
	lines := replaceLines(strings.Split(text, "\n"))
	lastWords := lines[len(lines)-2]

	switch {
	// summary
	case strings.HasPrefix(lastWords, "次へ"), strings.HasPrefix(lastWords, "Next"):
		csvName = createCsvSummary(twitterId, createdAt, lines)

	// details
	case strings.HasPrefix(lastWords, "とじる"), strings.HasPrefix(lastWords, "Close"):
		csvName = createCsvDetails(twitterId, createdAt, lines)
	}

	return csvName
}

func replaceLines(lines []string) []string {
	var rLines []string
	for _, line := range lines {
		rLine := strings.TrimSpace(strings.Trim(line, "*"))
		rLine = strings.Replace(rLine, "Om(", "0m(", 1)
		rLine = strings.Replace(rLine, "0(", "回(", 1)
		rLineSplited := strings.Split(rLine, " ")
		rLines = append(rLines, rLineSplited...)
	}
	return rLines
}

func createCsvSummary(twitterId string, createdAt time.Time, lines []string) string {
	var csvName string = "./csv/summary.csv"
	summary := setSummary(twitterId, createdAt, lines)

	_ = os.Remove(csvName)
	csvfile, _ := os.OpenFile(csvName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer csvfile.Close()

	gocsv.MarshalFile(&summary, csvfile)
	return csvName
}

func setSummary(twitterId string, createdAt time.Time, lines []string) []*Summary {
	var summary []*Summary
	rQuantity := regexp.MustCompile(`^[0-9.]+`)

	for i, line := range lines {
		if rQuantity.MatchString(line) {
			strTotalTime := strings.ReplaceAll(line, "時", "h")
			strTotalTime = strings.ReplaceAll(strTotalTime, "分", "m")
			strTotalTime = strings.ReplaceAll(strTotalTime, "秒", "s")
			totalTimeExcercising, _ := time.ParseDuration(strTotalTime)
			strTotalCalories := rQuantity.FindAllString(lines[i+2], 1)[0]
			totalCaloriesBurned, _ := strconv.ParseFloat(strTotalCalories, 64)
			totalDistanceRun, _ := strconv.ParseFloat(rQuantity.FindAllString(lines[i+4], 1)[0], 64)

			summary = append(summary, &Summary{
				TwitterId:            twitterId,
				TotalTimeExcercising: totalTimeExcercising,
				TotalCaloriesBurned:  totalCaloriesBurned,
				TotalDistanceRun:     totalDistanceRun,
				CreatedAt:            createdAt,
			})
			break
		}
	}
	return summary
}

func createCsvDetails(twitterId string, createdAt time.Time, lines []string) string {
	var csvName string = "./csv/details.csv"
	var isEven bool = (len(lines)%2 == 0)
	var isExercise bool = false
	rExercise := regexp.MustCompile(`^[^0-9]+`)
	details := []*Details{}

	for i, line := range lines {
		if strings.HasPrefix(line, "カッコ内はプレイ開始からの累計値です") {
			break
		} else if isExercise && !isEven &&
			rExercise.MatchString(line) &&
			rExercise.MatchString(lines[i+1]) {
			details = setDetails(details, twitterId, createdAt, line, lines[i+4])
			details = setDetails(details, twitterId, createdAt, lines[i+1], lines[i+3])
			details = setDetails(details, twitterId, createdAt, lines[i+2], lines[i+5])
			break
		} else if isExercise && rExercise.MatchString(line) {
			details = setDetails(details, twitterId, createdAt, line, lines[i+1])
		}

		if strings.HasPrefix(line, "R画面を撮影する") ||
			strings.HasPrefix(line, "画面を撮影する") {
			isExercise = true
		}
	}

	_ = os.Remove(csvName)
	csvfile, _ := os.OpenFile(csvName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer csvfile.Close()

	gocsv.MarshalFile(&details, csvfile)
	return csvName
}

func setDetails(details []*Details, twitterId string, createdAt time.Time, nameLine string, quantityLine string) []*Details {
	rQuantity := regexp.MustCompile(`^[0-9]+`)
	rTotalQuantity := regexp.MustCompile(`\([0-9]+`)

	quantity, _ := strconv.Atoi(rQuantity.FindAllString(quantityLine, 1)[0])
	strTotalQuantity := rTotalQuantity.FindAllString(quantityLine, 1)
	totalQuantity, _ := strconv.Atoi(strings.Trim(strTotalQuantity[0], "("))
	details = append(details, &Details{
		TwitterId:     twitterId,
		ExerciseName:  nameLine,
		Quantity:      quantity,
		TotalQuantity: totalQuantity,
		CreatedAt:     createdAt,
	})

	return details
}