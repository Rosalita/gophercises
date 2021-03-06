package main

import "testing"
import "flag"
import "os"
import "errors"
import "github.com/stretchr/testify/assert"
import "fmt"

var globalAnswerNumber = 1

func TestMain(m *testing.M) {

	flag.Parse()
	exitCode := m.Run()
	os.Exit(exitCode)

}

func TestReadFromCsvErrors(t *testing.T) {
	var tests = []struct {
		input string
		err   error
	}{
		{"testdata/test.csv", nil},
		{"testdata/helloworld.txt", errors.New("Error: Questions and answers must be in a .csv file, received .txt")},
	}
	for _, test := range tests {
		_, e := readFromCsv(test.input)
		assert.Equal(t, test.err, e, "unexpected error returned")
	}
}

func TestReadFromCsvReturnsDataOfCorrectLength(t *testing.T) {
	var tests = []struct {
		input      string
		quizLength int
	}{
		{"testdata/test.csv", 2},
		{"testdata/missing.csv", 0},
		{"testdata/helloworld.txt", 0},
	}
	for _, test := range tests {
		quizData, _ := readFromCsv(test.input)
		assert.Equal(t, len(quizData), test.quizLength, "unexpected data returned")
	}
}

func TestCanReadQuizFromValidCsv(t *testing.T) {
	path := "testdata/test.csv"
	data, _ := readFromCsv(path)

	expected_row1 := []string{"1 + 1", "2"}
	expected_row2 := []string{"Why can't spaghetti code?", "Impasta syndrome"}
	expected := [][]string{expected_row1, expected_row2}

	for j, row := range data {
		for i, item := range row {
			assert.Equal(t, expected[j][i], item, "unexpected values read from CSV")
		}
	}
}

func TestCsvFileIsMissing(t *testing.T) {
	path := "testdata/missing.csv"
	_, err := readFromCsv(path)

	expected := "open testdata/missing.csv: no such file or directory"
	assert.Equal(t, expected, err.Error(), "unexpected error message")
}

func mockSomeIncorrectAnswers() string {
	answer := fmt.Sprintf("answer %d", globalAnswerNumber)
	globalAnswerNumber++
	return answer
}

func TestGetAnswers(t *testing.T) {
	data, _ := readFromCsv("testdata/test.csv")
	result := getAnswers(data, mockSomeIncorrectAnswers)
	expected := []string{"answer 1", "answer 2"}

	assert.Equal(t, expected[0], result[0], "unexpected answer received")
	assert.Equal(t, expected[1], result[1], "unexpected answer received")
}

func TestMarkQuiz(t *testing.T) {
	testQuestions, _ := readFromCsv("testdata/test.csv")
	incorrectAnswers := []string{"fish", "chips"}
	correctAnswers := []string{"2", "Impasta syndrome"}
	allCapsCorrectAnswers := []string{"2", "IMPASTA SYNDROME"}

	var tests = []struct {
		questions [][]string
		answers   []string
		score     int
	}{
		{testQuestions, incorrectAnswers, 0},
		{testQuestions, correctAnswers, 2},
		{testQuestions, allCapsCorrectAnswers, 2},
	}
	for _, test := range tests {
		score := markQuiz(testQuestions, test.answers)
		assert.Equal(t, len(test.questions), len(test.answers), "number of questions and number of answers do not match")
		assert.Equal(t, test.score, score, "unexpected number of correct answers")
	}

}
