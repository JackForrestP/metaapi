//Auto generated with MetaApi https://github.com/exyzzy/metaapi
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
    "strconv"
)

import	"time"


var testDb *sql.DB
var configdb map[string]interface{}
const testDbName = "testtodo"

// ======= helpers

//assumes a configlocaldb.json file as:
//{
//    "Host": "localhost",
//    "Port": "5432",
//    "User": "dbname",
//    "Pass": "dbname",
//    "Name": "dbname",
//    "SSLMode": "disable"
//}
func loadConfig() {
	fmt.Println("  loadConfig")
	file, err := os.Open("configlocaldb.json")
	if err != nil {
		log.Panicln("Cannot open configlocaldb file", err.Error())
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configdb)
	if err != nil {
		log.Panicln("Cannot get local configurationdb from file", err.Error())
	}
}

func createDb(db *sql.DB, dbName string, owner string) (err error) {
	ss := fmt.Sprintf("CREATE DATABASE %s OWNER %s", dbName, owner)
	fmt.Println("  " + ss)
	_, err = db.Exec(ss)
	return
}

func setTzDb(db *sql.DB) (err error) {
	ss := fmt.Sprintf("SET TIME ZONE UTC")
	fmt.Println("  " + ss)
	_, err = db.Exec(ss)
	return
}

func dropDb(db *sql.DB, dbName string) (err error) {
	ss := fmt.Sprintf("DROP DATABASE %s", dbName)
	fmt.Println("  " + ss)
	_, err = db.Exec(ss)
	return
}

func rowExists(db *sql.DB, query string, args ...interface{}) (exists bool, err error) {
	query = fmt.Sprintf("SELECT EXISTS (%s)", query)
	fmt.Println("  " + query)
	err = db.QueryRow(query, args...).Scan(&exists)
	return
}

func tableExists(db *sql.DB, table string) (valid bool, err error) {

	valid, err = rowExists(db, "SELECT 1 FROM pg_tables WHERE tablename = $1", table)
	return
}

func initTestDb() (err error) {
	loadConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		"sslmode=%s", configdb["Host"], configdb["Port"], configdb["User"], configdb["Pass"], configdb["SSLMode"])
	testDb, err = sql.Open("postgres", psqlInfo)
	return
}

func TestMain(m *testing.M) {
	//test setup
	err := initTestDb()
	if err != nil {
		log.Panicln("cannot initTestDb ", err.Error())
	}

	err = createDb(testDb, testDbName, configdb["User"].(string))
	if err != nil {
		log.Panicln("cannot CreateDb ", err.Error())
	}

	err = setTzDb(testDb)
	if err != nil {
		log.Panicln("cannot setTzDb ", err.Error())
	}

	//run tests
	exitVal := m.Run()

	//test teardown
	err = dropDb(testDb, testDbName)
	if err != nil {
		log.Panicln("cannot DropDb ", err.Error())
	}
	os.Exit(exitVal)
}

type compareType func(interface{}, interface{}) bool

func noCompare(result, expect interface{}) bool {
	fmt.Printf("  noCompare: %v, %v -  %T, %T \n", result, expect, result, expect)
	return (true)
}

func defaultCompare(result, expect interface{}) bool {
	fmt.Printf("  defaultCompare: %v, %v -  %T, %T \n", result, expect, result, expect)
	return (result == expect)
}

func jsonCompare(result, expect interface{}) bool {
	fmt.Printf("  jsonCompare: %v, %v -  %T, %T \n", result, expect, result, expect)

	//json fields can be any order after db return, so read into map[string]interface and look up
	resultMap := make(map[string]interface{})
	expectMap := make(map[string]interface{})

	if reflect.TypeOf(result).String() == "sql.NullString" {
		err := json.Unmarshal([]byte(result.(sql.NullString).String), &resultMap)
		if err != nil {
			log.Panic(err)
		}
		err = json.Unmarshal([]byte(expect.(sql.NullString).String), &expectMap)
		if err != nil {
			log.Panic(err)
		}
	} else {
		err := json.Unmarshal([]byte(result.(string)), &resultMap)
		if err != nil {
			log.Panic(err)
		}
		err = json.Unmarshal([]byte(expect.(string)), &expectMap)
		if err != nil {
			log.Panic(err)
		}
	}

	for k, v := range expectMap {
		if v != resultMap[k] {
			fmt.Printf("Key: %v, Result: %v, Expect: %v", k, resultMap[k], v)
			return false
		}
	}
	return true


	for k, v := range expectMap {
		if v != resultMap[k] {
			fmt.Printf("Key: %v, Result: %v, Expect: %v", k, resultMap[k], v)
			return false
		}
	}
	return true
}

func stringCompare(result, expect interface{}) bool {

	resultJson, err := json.Marshal(result)
	if err != nil {
		log.Panic(err)
	}
	expectJson, err := json.Marshal(expect)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("  stringCompare: %v, %v -  %T, %T \n", string(resultJson), string(expectJson), result, expect)
	return (strings.TrimSpace(string(resultJson)) == strings.TrimSpace(string(expectJson)))
}

//psgl truncs reals at 6 digits
func realCompare(result, expect interface{}) bool {

	fmt.Printf("  realCompare: %v, %v -  %T, %T \n", result, expect, result, expect)

	var resultStr string
	var expectStr string
	if reflect.TypeOf(result).String() == "sql.NullFloat64" {
		resultStr = strconv.FormatFloat(result.(sql.NullFloat64).Float64, 'f', 6, 32)
		expectStr = strconv.FormatFloat(expect.(sql.NullFloat64).Float64, 'f', 6, 32)
	} else {
		resultStr = strconv.FormatFloat(float64(result.(float32)), 'f', 6, 32)
		expectStr = strconv.FormatFloat(float64(expect.(float32)), 'f', 6, 32)
	}
	return (resultStr == expectStr)
}

//iterate through each field of struct and apply the compare function to each field based on compareType map
func equalField(result, expect interface{}, compMap map[string]compareType) error {

	u := reflect.ValueOf(expect)
	v := reflect.ValueOf(result)
	typeOfS := u.Type()

	for i := 0; i < u.NumField(); i++ {

		if !(compMap[typeOfS.Field(i).Name])(v.Field(i).Interface(), u.Field(i).Interface()) {
			return fmt.Errorf("Field: %s, Result: %v, Expect: %v", typeOfS.Field(i).Name, v.Field(i).Interface(), u.Field(i).Interface())
		}
	}
	return nil
}


//table specific 


const surveystableName = "surveys"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testSurvey = [2]Survey{  Id: 1, Uuid: "eLhnvKCpJ3nO9ASl", Version: "1SQvYkP6NnVjtBcm", Title: "jzJLrZtoZFo9yPqZ", Description: "DwtdMUmLletVw44L", Open: false, CreationDate: time.Now().UTC().Truncate(time.Microsecond),  Id: 2, Uuid: "m9ipnUc9bcMDkv7c", Version: "rm7Xh1SOsCN0yc7V", Title: "Rp5j5eHzb0DidLD2", Description: "82tLNUFRdtZJppi4", Open: true, CreationDate: time.Now().UTC().Truncate(time.Microsecond) }

var updateSurvey = Survey Id: 1, Uuid: "QaLqqJccZs7Fvjhn", Version: "m0ybMQGe9gyGhiTL", Title: "z1xIrI6Zc3j6MOMQ", Description: "PMwapA2vOy1HqJob", Open: false, CreationDate: time.Now().UTC().Truncate(time.Microsecond)

//compare functions
var compareSurveys = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"Version": defaultCompare,
	"Title": defaultCompare,
	"Description": defaultCompare,
	"Open": defaultCompare,
	"CreationDate": stringCompare,

}

func reverseSurveys(surveys []Survey) (result []Survey) {

	for i := len(surveys) - 1; i >= 0; i-- {
		result = append(result, surveys[i])
	}
	return
}

// ======= tests: Survey =======

func TestCreateTableSurveys(t *testing.T) {
	fmt.Println("==CreateTableSurveys")
	err := CreateTableSurveys(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableSurveys " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableSurveys")
	}
	exists, err := tableExists(testDb, "surveys")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(surveys) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateSurvey(t *testing.T) {
	fmt.Println("==CreateSurvey")
	result, err := testSurvey[0].CreateSurvey(testDb)
	if err != nil {
		t.Errorf("cannot CreateSurvey " + err.Error())
	} else {
		fmt.Println("  Done: CreateSurvey")
	}
	err = equalField(result, testSurvey[0], compareSurveys)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveSurvey(t *testing.T) {
	fmt.Println("==RetrieveSurvey")
	result, err := testSurvey[0].RetrieveSurvey(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveSurvey " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveSurvey")
	}
	err = equalField(result, testSurvey[0], compareSurveys)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllSurveys(t *testing.T) {
	fmt.Println("==RetrieveAllSurveys")
	_, err := testSurvey[1].CreateSurvey(testDb)
	if err != nil {
		t.Errorf("cannot CreateSurvey " + err.Error())
	} else {
		fmt.Println("  Done: CreateSurvey")
	}
	result, err := RetrieveAllSurveys(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllSurveys " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllSurveys")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseSurveys(testSurvey[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareSurveys)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateSurvey(t *testing.T) {
	fmt.Println("==UpdateSurvey")
	result, err := updateSurvey.UpdateSurvey(testDb)
	if err != nil {
		t.Errorf("cannot UpdateSurvey " + err.Error())
	} else {
		fmt.Println("  Done: UpdateSurvey")
	}
	err = equalField(result, updateSurvey, compareSurveys)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const questionstableName = "questions"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testQuestion = [2]Question{  Id: 1, Uuid: "v2SsM3zFc6lf1VYx", Question: "rLKyDb5waQKovXOF", Focus: "YQnInHaa8Eh8glHh", Unit: "Wg3aXgxiuCdmBpkE", BackImage: "6h92EH061muN6q7r", MinValue: 1162371545, MaxValue: 484807343, DefaultValue: 1687829751, WarpFunction: "1T7VfA2MwpcvjDOZ", CreationDate: time.Now().UTC().Truncate(time.Microsecond),  Id: 2, Uuid: "8KHDBpxa5vF0MWYq", Question: "AhztBdQxyDx8VtDW", Focus: "qydnLKTy5HLZ7WUv", Unit: "RW765MFIcIAPbwOe", BackImage: "jqycBhEqoicV06ej", MinValue: 1114666444, MaxValue: 782924087, DefaultValue: 1731407471, WarpFunction: "m5aBFkWfb293j1Ev", CreationDate: time.Now().UTC().Truncate(time.Microsecond) }

var updateQuestion = Question Id: 1, Uuid: "uKEvLpUQAzQrf0kk", Question: "lU7C9g593phgpjz5", Focus: "b6BdY5tSwPE8EaoW", Unit: "pNQIANpt2lM5Vm2r", BackImage: "fGjNZBpVGeqvXRBH", MinValue: 1028302608, MaxValue: 439085936, DefaultValue: 1300039253, WarpFunction: "p9gNoCnkRicO8paZ", CreationDate: time.Now().UTC().Truncate(time.Microsecond)

//compare functions
var compareQuestions = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"Question": defaultCompare,
	"Focus": defaultCompare,
	"Unit": defaultCompare,
	"BackImage": defaultCompare,
	"MinValue": defaultCompare,
	"MaxValue": defaultCompare,
	"DefaultValue": defaultCompare,
	"WarpFunction": defaultCompare,
	"CreationDate": stringCompare,

}

func reverseQuestions(questions []Question) (result []Question) {

	for i := len(questions) - 1; i >= 0; i-- {
		result = append(result, questions[i])
	}
	return
}

// ======= tests: Question =======

func TestCreateTableQuestions(t *testing.T) {
	fmt.Println("==CreateTableQuestions")
	err := CreateTableQuestions(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableQuestions " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableQuestions")
	}
	exists, err := tableExists(testDb, "questions")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(questions) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateQuestion(t *testing.T) {
	fmt.Println("==CreateQuestion")
	result, err := testQuestion[0].CreateQuestion(testDb)
	if err != nil {
		t.Errorf("cannot CreateQuestion " + err.Error())
	} else {
		fmt.Println("  Done: CreateQuestion")
	}
	err = equalField(result, testQuestion[0], compareQuestions)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveQuestion(t *testing.T) {
	fmt.Println("==RetrieveQuestion")
	result, err := testQuestion[0].RetrieveQuestion(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveQuestion " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveQuestion")
	}
	err = equalField(result, testQuestion[0], compareQuestions)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllQuestions(t *testing.T) {
	fmt.Println("==RetrieveAllQuestions")
	_, err := testQuestion[1].CreateQuestion(testDb)
	if err != nil {
		t.Errorf("cannot CreateQuestion " + err.Error())
	} else {
		fmt.Println("  Done: CreateQuestion")
	}
	result, err := RetrieveAllQuestions(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllQuestions " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllQuestions")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseQuestions(testQuestion[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareQuestions)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateQuestion(t *testing.T) {
	fmt.Println("==UpdateQuestion")
	result, err := updateQuestion.UpdateQuestion(testDb)
	if err != nil {
		t.Errorf("cannot UpdateQuestion " + err.Error())
	} else {
		fmt.Println("  Done: UpdateQuestion")
	}
	err = equalField(result, updateQuestion, compareQuestions)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const answerstableName = "answers"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testAnswer = [2]Answer{  Id: 1, Uuid: "xSKHQHIFrknEDVhb", AnswerColumn: "0fYa7Ha8xmVCqSRT", IntAnswer: sql.NullInt32{1717237539, true}, StringAnswer: sql.NullString{"CHyANEJAG7GUZCgr", true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond),  Id: 2, Uuid: "UGxpUwQHvSNUxrKx", AnswerColumn: "eHT4cSrcuyDBcu4Z", IntAnswer: sql.NullInt32{1322153763, true}, StringAnswer: sql.NullString{"bDqHxlAMbh9TyJDZ", true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond) }

var updateAnswer = Answer Id: 1, Uuid: "UAcARPjyFzKONGQQ", AnswerColumn: "qWOoiGH2jH3Z340C", IntAnswer: sql.NullInt32{1816880339, true}, StringAnswer: sql.NullString{"ZeqXiUgkCKq9lNfl", true}, CreationDate: time.Now().UTC().Truncate(time.Microsecond)

//compare functions
var compareAnswers = map[string]compareType{
	"Id": defaultCompare,
	"Uuid": defaultCompare,
	"AnswerColumn": defaultCompare,
	"IntAnswer": defaultCompare,
	"StringAnswer": defaultCompare,
	"CreationDate": stringCompare,

}

func reverseAnswers(answers []Answer) (result []Answer) {

	for i := len(answers) - 1; i >= 0; i-- {
		result = append(result, answers[i])
	}
	return
}

// ======= tests: Answer =======

func TestCreateTableAnswers(t *testing.T) {
	fmt.Println("==CreateTableAnswers")
	err := CreateTableAnswers(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableAnswers " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableAnswers")
	}
	exists, err := tableExists(testDb, "answers")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(answers) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateAnswer(t *testing.T) {
	fmt.Println("==CreateAnswer")
	result, err := testAnswer[0].CreateAnswer(testDb)
	if err != nil {
		t.Errorf("cannot CreateAnswer " + err.Error())
	} else {
		fmt.Println("  Done: CreateAnswer")
	}
	err = equalField(result, testAnswer[0], compareAnswers)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAnswer(t *testing.T) {
	fmt.Println("==RetrieveAnswer")
	result, err := testAnswer[0].RetrieveAnswer(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAnswer " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAnswer")
	}
	err = equalField(result, testAnswer[0], compareAnswers)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllAnswers(t *testing.T) {
	fmt.Println("==RetrieveAllAnswers")
	_, err := testAnswer[1].CreateAnswer(testDb)
	if err != nil {
		t.Errorf("cannot CreateAnswer " + err.Error())
	} else {
		fmt.Println("  Done: CreateAnswer")
	}
	result, err := RetrieveAllAnswers(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllAnswers " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllAnswers")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseAnswers(testAnswer[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareAnswers)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateAnswer(t *testing.T) {
	fmt.Println("==UpdateAnswer")
	result, err := updateAnswer.UpdateAnswer(testDb)
	if err != nil {
		t.Errorf("cannot UpdateAnswer " + err.Error())
	} else {
		fmt.Println("  Done: UpdateAnswer")
	}
	err = equalField(result, updateAnswer, compareAnswers)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const survey_question_linkstableName = "survey_question_links"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testSurvey_question_link = [2]Survey_question_link{  Id: 1, SurveyUuid: "EDPtGpBG743zroNJ", QuestionUuid: "N4NzeMDE13S2qNFt",  Id: 2, SurveyUuid: "BFAOXA0JcjOiBVyU", QuestionUuid: "IKgRkfmlY0QnbtB8" }

var updateSurvey_question_link = Survey_question_link Id: 1, SurveyUuid: "rs7ySIl8SpP6Yzjs", QuestionUuid: "ggU7K0y8WhPBWC61"

//compare functions
var compareSurvey_question_links = map[string]compareType{
	"Id": defaultCompare,
	"SurveyUuid": defaultCompare,
	"QuestionUuid": defaultCompare,

}

func reverseSurvey_question_links(survey_question_links []Survey_question_link) (result []Survey_question_link) {

	for i := len(survey_question_links) - 1; i >= 0; i-- {
		result = append(result, survey_question_links[i])
	}
	return
}

// ======= tests: Survey_question_link =======

func TestCreateTableSurvey_question_links(t *testing.T) {
	fmt.Println("==CreateTableSurvey_question_links")
	err := CreateTableSurvey_question_links(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableSurvey_question_links " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableSurvey_question_links")
	}
	exists, err := tableExists(testDb, "survey_question_links")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(survey_question_links) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateSurvey_question_link(t *testing.T) {
	fmt.Println("==CreateSurvey_question_link")
	result, err := testSurvey_question_link[0].CreateSurvey_question_link(testDb)
	if err != nil {
		t.Errorf("cannot CreateSurvey_question_link " + err.Error())
	} else {
		fmt.Println("  Done: CreateSurvey_question_link")
	}
	err = equalField(result, testSurvey_question_link[0], compareSurvey_question_links)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveSurvey_question_link(t *testing.T) {
	fmt.Println("==RetrieveSurvey_question_link")
	result, err := testSurvey_question_link[0].RetrieveSurvey_question_link(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveSurvey_question_link " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveSurvey_question_link")
	}
	err = equalField(result, testSurvey_question_link[0], compareSurvey_question_links)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllSurvey_question_links(t *testing.T) {
	fmt.Println("==RetrieveAllSurvey_question_links")
	_, err := testSurvey_question_link[1].CreateSurvey_question_link(testDb)
	if err != nil {
		t.Errorf("cannot CreateSurvey_question_link " + err.Error())
	} else {
		fmt.Println("  Done: CreateSurvey_question_link")
	}
	result, err := RetrieveAllSurvey_question_links(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllSurvey_question_links " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllSurvey_question_links")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseSurvey_question_links(testSurvey_question_link[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareSurvey_question_links)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateSurvey_question_link(t *testing.T) {
	fmt.Println("==UpdateSurvey_question_link")
	result, err := updateSurvey_question_link.UpdateSurvey_question_link(testDb)
	if err != nil {
		t.Errorf("cannot UpdateSurvey_question_link " + err.Error())
	} else {
		fmt.Println("  Done: UpdateSurvey_question_link")
	}
	err = equalField(result, updateSurvey_question_link, compareSurvey_question_links)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


const survey_question_answer_linkstableName = "survey_question_answer_links"

//test data - note: double brackets in test data need space between otherwise are interpreted as template action
var testSurvey_question_answer_link = [2]Survey_question_answer_link{  Id: 1, SurveyUuid: "imkTE50fcvJn2Va3", QuestionUuid: "MqifxRVeY8P6qUAN", AnswerUuid: "GCXA2nZQvT5XEMk9", UserId: sql.NullInt32{1, true},  Id: 2, SurveyUuid: "ZEr6DSQSbwaSTYPJ", QuestionUuid: "wxdAkZQ2fDv1vKls", AnswerUuid: "AB3dM3gTClWqOa2i", UserId: sql.NullInt32{2, true} }

var updateSurvey_question_answer_link = Survey_question_answer_link Id: 1, SurveyUuid: "h5PPDwd11j7AqfAe", QuestionUuid: "MfymcwsAEeXbgZGw", AnswerUuid: "wLMF6vPOJrQNFMlF", UserId: sql.NullInt32{1, true}

//compare functions
var compareSurvey_question_answer_links = map[string]compareType{
	"Id": defaultCompare,
	"SurveyUuid": defaultCompare,
	"QuestionUuid": defaultCompare,
	"AnswerUuid": defaultCompare,
	"UserId": defaultCompare,

}

func reverseSurvey_question_answer_links(survey_question_answer_links []Survey_question_answer_link) (result []Survey_question_answer_link) {

	for i := len(survey_question_answer_links) - 1; i >= 0; i-- {
		result = append(result, survey_question_answer_links[i])
	}
	return
}

// ======= tests: Survey_question_answer_link =======

func TestCreateTableSurvey_question_answer_links(t *testing.T) {
	fmt.Println("==CreateTableSurvey_question_answer_links")
	err := CreateTableSurvey_question_answer_links(testDb)
	if err != nil {
		t.Errorf("cannot CreateTableSurvey_question_answer_links " + err.Error())
	} else {
		fmt.Println("  Done: CreateTableSurvey_question_answer_links")
	}
	exists, err := tableExists(testDb, "survey_question_answer_links")
	if err != nil {
		t.Errorf("cannot tableExists " + err.Error())
	}
	if !exists {
		t.Errorf("tableExists(survey_question_answer_links) returned wrong status code: got %v want %v", exists, true)
	} else {
		fmt.Println("  Done: tableExists")
	}
}

func TestCreateSurvey_question_answer_link(t *testing.T) {
	fmt.Println("==CreateSurvey_question_answer_link")
	result, err := testSurvey_question_answer_link[0].CreateSurvey_question_answer_link(testDb)
	if err != nil {
		t.Errorf("cannot CreateSurvey_question_answer_link " + err.Error())
	} else {
		fmt.Println("  Done: CreateSurvey_question_answer_link")
	}
	err = equalField(result, testSurvey_question_answer_link[0], compareSurvey_question_answer_links)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveSurvey_question_answer_link(t *testing.T) {
	fmt.Println("==RetrieveSurvey_question_answer_link")
	result, err := testSurvey_question_answer_link[0].RetrieveSurvey_question_answer_link(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveSurvey_question_answer_link " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveSurvey_question_answer_link")
	}
	err = equalField(result, testSurvey_question_answer_link[0], compareSurvey_question_answer_links)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}

func TestRetrieveAllSurvey_question_answer_links(t *testing.T) {
	fmt.Println("==RetrieveAllSurvey_question_answer_links")
	_, err := testSurvey_question_answer_link[1].CreateSurvey_question_answer_link(testDb)
	if err != nil {
		t.Errorf("cannot CreateSurvey_question_answer_link " + err.Error())
	} else {
		fmt.Println("  Done: CreateSurvey_question_answer_link")
	}
	result, err := RetrieveAllSurvey_question_answer_links(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllSurvey_question_answer_links " + err.Error())
	} else {
		fmt.Println("  Done: RetrieveAllSurvey_question_answer_links")
	}
	//reverse because api is DESC, [:] is slice of all array elements
	expect := reverseSurvey_question_answer_links(testSurvey_question_answer_link[:])
	for i, _ := range expect {
		err = equalField(result[i], expect[i], compareSurvey_question_answer_links)
		if err != nil {
			t.Errorf("api returned unexpected result. " + err.Error())
		}
	}
}


func TestUpdateSurvey_question_answer_link(t *testing.T) {
	fmt.Println("==UpdateSurvey_question_answer_link")
	result, err := updateSurvey_question_answer_link.UpdateSurvey_question_answer_link(testDb)
	if err != nil {
		t.Errorf("cannot UpdateSurvey_question_answer_link " + err.Error())
	} else {
		fmt.Println("  Done: UpdateSurvey_question_answer_link")
	}
	err = equalField(result, updateSurvey_question_answer_link, compareSurvey_question_answer_links)
	if err != nil {
		t.Errorf("api returned unexpected result. " + err.Error())
	}
}


//delete all data in reverse order to accommodate foreign keys


func TestDeleteSurvey_question_answer_link(t *testing.T) {
	fmt.Println("==DeleteSurvey_question_answer_link")
	err := testSurvey_question_answer_link[0].DeleteSurvey_question_answer_link(testDb)
	if err != nil {
		t.Errorf("cannot DeleteSurvey_question_answer_link " + err.Error())
	} else {
		fmt.Println("  Done: DeleteSurvey_question_answer_link")
	}
	_, err = testSurvey_question_answer_link[0].RetrieveSurvey_question_answer_link(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveSurvey_question_answer_link with no result")
		} else {
			t.Errorf("cannot RetrieveSurvey_question_answer_link " + err.Error())
		}
	}
}

func TestDeleteAllSurvey_question_answer_links(t *testing.T) {
	fmt.Println("==DeleteAllSurvey_question_answer_links")
	err := DeleteAllSurvey_question_answer_links(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllSurvey_question_answer_links " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllSurvey_question_answer_links")
	}
	result, err := RetrieveAllSurvey_question_answer_links(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllSurvey_question_answer_links " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllSurvey_question_answer_links with no result")
	}
}


func TestDeleteSurvey_question_link(t *testing.T) {
	fmt.Println("==DeleteSurvey_question_link")
	err := testSurvey_question_link[0].DeleteSurvey_question_link(testDb)
	if err != nil {
		t.Errorf("cannot DeleteSurvey_question_link " + err.Error())
	} else {
		fmt.Println("  Done: DeleteSurvey_question_link")
	}
	_, err = testSurvey_question_link[0].RetrieveSurvey_question_link(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveSurvey_question_link with no result")
		} else {
			t.Errorf("cannot RetrieveSurvey_question_link " + err.Error())
		}
	}
}

func TestDeleteAllSurvey_question_links(t *testing.T) {
	fmt.Println("==DeleteAllSurvey_question_links")
	err := DeleteAllSurvey_question_links(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllSurvey_question_links " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllSurvey_question_links")
	}
	result, err := RetrieveAllSurvey_question_links(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllSurvey_question_links " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllSurvey_question_links with no result")
	}
}


func TestDeleteAnswer(t *testing.T) {
	fmt.Println("==DeleteAnswer")
	err := testAnswer[0].DeleteAnswer(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAnswer " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAnswer")
	}
	_, err = testAnswer[0].RetrieveAnswer(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveAnswer with no result")
		} else {
			t.Errorf("cannot RetrieveAnswer " + err.Error())
		}
	}
}

func TestDeleteAllAnswers(t *testing.T) {
	fmt.Println("==DeleteAllAnswers")
	err := DeleteAllAnswers(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllAnswers " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllAnswers")
	}
	result, err := RetrieveAllAnswers(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllAnswers " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllAnswers with no result")
	}
}


func TestDeleteQuestion(t *testing.T) {
	fmt.Println("==DeleteQuestion")
	err := testQuestion[0].DeleteQuestion(testDb)
	if err != nil {
		t.Errorf("cannot DeleteQuestion " + err.Error())
	} else {
		fmt.Println("  Done: DeleteQuestion")
	}
	_, err = testQuestion[0].RetrieveQuestion(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveQuestion with no result")
		} else {
			t.Errorf("cannot RetrieveQuestion " + err.Error())
		}
	}
}

func TestDeleteAllQuestions(t *testing.T) {
	fmt.Println("==DeleteAllQuestions")
	err := DeleteAllQuestions(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllQuestions " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllQuestions")
	}
	result, err := RetrieveAllQuestions(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllQuestions " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllQuestions with no result")
	}
}


func TestDeleteSurvey(t *testing.T) {
	fmt.Println("==DeleteSurvey")
	err := testSurvey[0].DeleteSurvey(testDb)
	if err != nil {
		t.Errorf("cannot DeleteSurvey " + err.Error())
	} else {
		fmt.Println("  Done: DeleteSurvey")
	}
	_, err = testSurvey[0].RetrieveSurvey(testDb)
	if err == nil {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		if err == sql.ErrNoRows {
			fmt.Println("  Done: RetrieveSurvey with no result")
		} else {
			t.Errorf("cannot RetrieveSurvey " + err.Error())
		}
	}
}

func TestDeleteAllSurveys(t *testing.T) {
	fmt.Println("==DeleteAllSurveys")
	err := DeleteAllSurveys(testDb)
	if err != nil {
		t.Errorf("cannot DeleteAllSurveys " + err.Error())
	} else {
		fmt.Println("  Done: DeleteAllSurveys")
	}
	result, err := RetrieveAllSurveys(testDb)
	if err != nil {
		t.Errorf("cannot RetrieveAllSurveys " + err.Error())
	}
	if len(result) > 0 {
		t.Errorf("api returned unexpected result: got Row want NoRow")
	} else {
		fmt.Println("  Done: RetrieveAllSurveys with no result")
	}
}
