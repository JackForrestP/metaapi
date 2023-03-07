//Auto generated with MetaApi https://github.com/exyzzy/metaapi
package main

import (
	"database/sql"
	"time"
)

// ======= Survey =======

//Create Table
func CreateTableSurveys(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS surveys CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table surveys ( id integer generated always as identity primary key , uuid text not null UNIQUE , version text not null , title text not null , description text not null , open boolean not null , creation_date timestamptz not null ) ; `)
	return
}

//Drop Table
func DropTableSurveys(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS surveys CASCADE")
	return
}

//Struct
type Survey struct {
	Id           int32     `xml:"Id" json:"id"`
	Uuid         string    `xml:"Uuid" json:"uuid"`
	Version      string    `xml:"Version" json:"version"`
	Title        string    `xml:"Title" json:"title"`
	Description  string    `xml:"Description" json:"description"`
	Open         bool      `xml:"Open" json:"open"`
	CreationDate time.Time `xml:"CreationDate" json:"creationdate"`
}

//Create
func (survey *Survey) CreateSurvey(db *sql.DB) (result Survey, err error) {
	stmt, err := db.Prepare("INSERT INTO surveys ( uuid, version, title, description, open, creation_date) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, uuid, version, title, description, open, creation_date")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(survey.Uuid, survey.Version, survey.Title, survey.Description, survey.Open, survey.CreationDate).Scan(&result.Id, &result.Uuid, &result.Version, &result.Title, &result.Description, &result.Open, &result.CreationDate)
	return
}

//Retrieve
func (survey *Survey) RetrieveSurvey(db *sql.DB) (result Survey, err error) {
	result = Survey{}
	err = db.QueryRow("SELECT id, uuid, version, title, description, open, creation_date FROM surveys WHERE (id = $1)", survey.Id).Scan(&result.Id, &result.Uuid, &result.Version, &result.Title, &result.Description, &result.Open, &result.CreationDate)
	return
}

//RetrieveAll
func RetrieveAllSurveys(db *sql.DB) (surveys []Survey, err error) {
	rows, err := db.Query("SELECT id, uuid, version, title, description, open, creation_date FROM surveys ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Survey{}
		if err = rows.Scan(&result.Id, &result.Uuid, &result.Version, &result.Title, &result.Description, &result.Open, &result.CreationDate); err != nil {
			return
		}
		surveys = append(surveys, result)
	}
	rows.Close()
	return
}

//Update
func (survey *Survey) UpdateSurvey(db *sql.DB) (result Survey, err error) {
	stmt, err := db.Prepare("UPDATE surveys SET uuid = $2, version = $3, title = $4, description = $5, open = $6, creation_date = $7 WHERE (id = $1) RETURNING id, uuid, version, title, description, open, creation_date")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(survey.Id, survey.Uuid, survey.Version, survey.Title, survey.Description, survey.Open, survey.CreationDate).Scan(&result.Id, &result.Uuid, &result.Version, &result.Title, &result.Description, &result.Open, &result.CreationDate)
	return
}

//Delete
func (survey *Survey) DeleteSurvey(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM surveys WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(survey.Id)
	return
}

//DeleteAll
func DeleteAllSurveys(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM surveys")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}

// ======= Question =======

//Create Table
func CreateTableQuestions(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS questions CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table questions ( id integer generated always as identity primary key , uuid text not null UNIQUE , question text not null , focus text not null , unit text not null , back_image text not null , min_value integer not null , max_value integer not null , default_value integer not null , warp_function text not null , creation_date timestamptz not null ) ; `)
	return
}

//Drop Table
func DropTableQuestions(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS questions CASCADE")
	return
}

//Struct
type Question struct {
	Id           int32     `xml:"Id" json:"id"`
	Uuid         string    `xml:"Uuid" json:"uuid"`
	Question     string    `xml:"Question" json:"question"`
	Focus        string    `xml:"Focus" json:"focus"`
	Unit         string    `xml:"Unit" json:"unit"`
	BackImage    string    `xml:"BackImage" json:"backimage"`
	MinValue     int32     `xml:"MinValue" json:"minvalue"`
	MaxValue     int32     `xml:"MaxValue" json:"maxvalue"`
	DefaultValue int32     `xml:"DefaultValue" json:"defaultvalue"`
	WarpFunction string    `xml:"WarpFunction" json:"warpfunction"`
	CreationDate time.Time `xml:"CreationDate" json:"creationdate"`
}

//Create
func (question *Question) CreateQuestion(db *sql.DB) (result Question, err error) {
	stmt, err := db.Prepare("INSERT INTO questions ( uuid, question, focus, unit, back_image, min_value, max_value, default_value, warp_function, creation_date) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id, uuid, question, focus, unit, back_image, min_value, max_value, default_value, warp_function, creation_date")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(question.Uuid, question.Question, question.Focus, question.Unit, question.BackImage, question.MinValue, question.MaxValue, question.DefaultValue, question.WarpFunction, question.CreationDate).Scan(&result.Id, &result.Uuid, &result.Question, &result.Focus, &result.Unit, &result.BackImage, &result.MinValue, &result.MaxValue, &result.DefaultValue, &result.WarpFunction, &result.CreationDate)
	return
}

//Retrieve
func (question *Question) RetrieveQuestion(db *sql.DB) (result Question, err error) {
	result = Question{}
	err = db.QueryRow("SELECT id, uuid, question, focus, unit, back_image, min_value, max_value, default_value, warp_function, creation_date FROM questions WHERE (id = $1)", question.Id).Scan(&result.Id, &result.Uuid, &result.Question, &result.Focus, &result.Unit, &result.BackImage, &result.MinValue, &result.MaxValue, &result.DefaultValue, &result.WarpFunction, &result.CreationDate)
	return
}

//RetrieveAll
func RetrieveAllQuestions(db *sql.DB) (questions []Question, err error) {
	rows, err := db.Query("SELECT id, uuid, question, focus, unit, back_image, min_value, max_value, default_value, warp_function, creation_date FROM questions ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Question{}
		if err = rows.Scan(&result.Id, &result.Uuid, &result.Question, &result.Focus, &result.Unit, &result.BackImage, &result.MinValue, &result.MaxValue, &result.DefaultValue, &result.WarpFunction, &result.CreationDate); err != nil {
			return
		}
		questions = append(questions, result)
	}
	rows.Close()
	return
}

//Update
func (question *Question) UpdateQuestion(db *sql.DB) (result Question, err error) {
	stmt, err := db.Prepare("UPDATE questions SET uuid = $2, question = $3, focus = $4, unit = $5, back_image = $6, min_value = $7, max_value = $8, default_value = $9, warp_function = $10, creation_date = $11 WHERE (id = $1) RETURNING id, uuid, question, focus, unit, back_image, min_value, max_value, default_value, warp_function, creation_date")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(question.Id, question.Uuid, question.Question, question.Focus, question.Unit, question.BackImage, question.MinValue, question.MaxValue, question.DefaultValue, question.WarpFunction, question.CreationDate).Scan(&result.Id, &result.Uuid, &result.Question, &result.Focus, &result.Unit, &result.BackImage, &result.MinValue, &result.MaxValue, &result.DefaultValue, &result.WarpFunction, &result.CreationDate)
	return
}

//Delete
func (question *Question) DeleteQuestion(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM questions WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(question.Id)
	return
}

//DeleteAll
func DeleteAllQuestions(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM questions")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}

// ======= Answer =======

//Create Table
func CreateTableAnswers(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS answers CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table answers ( id integer generated always as identity primary key , uuid text not null UNIQUE , answer_column text not null , int_answer int , string_answer text , creation_date timestamptz not null ) ; `)
	return
}

//Drop Table
func DropTableAnswers(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS answers CASCADE")
	return
}

//Struct
type Answer struct {
	Id           int32          `xml:"Id" json:"id"`
	Uuid         string         `xml:"Uuid" json:"uuid"`
	AnswerColumn string         `xml:"AnswerColumn" json:"answercolumn"`
	IntAnswer    sql.NullInt32  `xml:"IntAnswer" json:"intanswer"`
	StringAnswer sql.NullString `xml:"StringAnswer" json:"stringanswer"`
	CreationDate time.Time      `xml:"CreationDate" json:"creationdate"`
}

//Create
func (answer *Answer) CreateAnswer(db *sql.DB) (result Answer, err error) {
	stmt, err := db.Prepare("INSERT INTO answers ( uuid, answer_column, int_answer, string_answer, creation_date) VALUES ($1,$2,$3,$4,$5) RETURNING id, uuid, answer_column, int_answer, string_answer, creation_date")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(answer.Uuid, answer.AnswerColumn, answer.IntAnswer, answer.StringAnswer, answer.CreationDate).Scan(&result.Id, &result.Uuid, &result.AnswerColumn, &result.IntAnswer, &result.StringAnswer, &result.CreationDate)
	return
}

//Retrieve
func (answer *Answer) RetrieveAnswer(db *sql.DB) (result Answer, err error) {
	result = Answer{}
	err = db.QueryRow("SELECT id, uuid, answer_column, int_answer, string_answer, creation_date FROM answers WHERE (id = $1)", answer.Id).Scan(&result.Id, &result.Uuid, &result.AnswerColumn, &result.IntAnswer, &result.StringAnswer, &result.CreationDate)
	return
}

//RetrieveAll
func RetrieveAllAnswers(db *sql.DB) (answers []Answer, err error) {
	rows, err := db.Query("SELECT id, uuid, answer_column, int_answer, string_answer, creation_date FROM answers ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Answer{}
		if err = rows.Scan(&result.Id, &result.Uuid, &result.AnswerColumn, &result.IntAnswer, &result.StringAnswer, &result.CreationDate); err != nil {
			return
		}
		answers = append(answers, result)
	}
	rows.Close()
	return
}

//Update
func (answer *Answer) UpdateAnswer(db *sql.DB) (result Answer, err error) {
	stmt, err := db.Prepare("UPDATE answers SET uuid = $2, answer_column = $3, int_answer = $4, string_answer = $5, creation_date = $6 WHERE (id = $1) RETURNING id, uuid, answer_column, int_answer, string_answer, creation_date")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(answer.Id, answer.Uuid, answer.AnswerColumn, answer.IntAnswer, answer.StringAnswer, answer.CreationDate).Scan(&result.Id, &result.Uuid, &result.AnswerColumn, &result.IntAnswer, &result.StringAnswer, &result.CreationDate)
	return
}

//Delete
func (answer *Answer) DeleteAnswer(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM answers WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(answer.Id)
	return
}

//DeleteAll
func DeleteAllAnswers(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM answers")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}

// ======= Survey_question_link =======

//Create Table
func CreateTableSurvey_question_links(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS survey_question_links CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table survey_question_links ( id integer generated always as identity primary key , survey_uuid text not null REFERENCES surveys(uuid) , question_uuid text not null REFERENCES questions(uuid) ) ; `)
	return
}

//Drop Table
func DropTableSurvey_question_links(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS survey_question_links CASCADE")
	return
}

//Struct
type Survey_question_link struct {
	Id           int32  `xml:"Id" json:"id"`
	SurveyUuid   string `xml:"SurveyUuid" json:"surveyuuid"`
	QuestionUuid string `xml:"QuestionUuid" json:"questionuuid"`
}

//Create
func (survey_question_link *Survey_question_link) CreateSurvey_question_link(db *sql.DB) (result Survey_question_link, err error) {
	stmt, err := db.Prepare("INSERT INTO survey_question_links ( survey_uuid, question_uuid) VALUES ($1,$2) RETURNING id, survey_uuid, question_uuid")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(survey_question_link.SurveyUuid, survey_question_link.QuestionUuid).Scan(&result.Id, &result.SurveyUuid, &result.QuestionUuid)
	return
}

//Retrieve
func (survey_question_link *Survey_question_link) RetrieveSurvey_question_link(db *sql.DB) (result Survey_question_link, err error) {
	result = Survey_question_link{}
	err = db.QueryRow("SELECT id, survey_uuid, question_uuid FROM survey_question_links WHERE (id = $1)", survey_question_link.Id).Scan(&result.Id, &result.SurveyUuid, &result.QuestionUuid)
	return
}

//RetrieveAll
func RetrieveAllSurvey_question_links(db *sql.DB) (survey_question_links []Survey_question_link, err error) {
	rows, err := db.Query("SELECT id, survey_uuid, question_uuid FROM survey_question_links ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Survey_question_link{}
		if err = rows.Scan(&result.Id, &result.SurveyUuid, &result.QuestionUuid); err != nil {
			return
		}
		survey_question_links = append(survey_question_links, result)
	}
	rows.Close()
	return
}

//Update
func (survey_question_link *Survey_question_link) UpdateSurvey_question_link(db *sql.DB) (result Survey_question_link, err error) {
	stmt, err := db.Prepare("UPDATE survey_question_links SET survey_uuid = $2, question_uuid = $3 WHERE (id = $1) RETURNING id, survey_uuid, question_uuid")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(survey_question_link.Id, survey_question_link.SurveyUuid, survey_question_link.QuestionUuid).Scan(&result.Id, &result.SurveyUuid, &result.QuestionUuid)
	return
}

//Delete
func (survey_question_link *Survey_question_link) DeleteSurvey_question_link(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM survey_question_links WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(survey_question_link.Id)
	return
}

//DeleteAll
func DeleteAllSurvey_question_links(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM survey_question_links")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}

// ======= Survey_question_answer_link =======

//Create Table
func CreateTableSurvey_question_answer_links(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS survey_question_answer_links CASCADE")
	if err != nil {
		return
	}
	_, err = db.Exec(`create table survey_question_answer_links ( id integer generated always as identity primary key , survey_uuid text not null REFERENCES surveys(uuid) , question_uuid text not null REFERENCES questions(uuid) , answer_uuid text not null REFERENCES answers(uuid) , user_id int REFERENCES users(id) ) ; `)
	return
}

//Drop Table
func DropTableSurvey_question_answer_links(db *sql.DB) (err error) {
	_, err = db.Exec("DROP TABLE IF EXISTS survey_question_answer_links CASCADE")
	return
}

//Struct
type Survey_question_answer_link struct {
	Id           int32         `xml:"Id" json:"id"`
	SurveyUuid   string        `xml:"SurveyUuid" json:"surveyuuid"`
	QuestionUuid string        `xml:"QuestionUuid" json:"questionuuid"`
	AnswerUuid   string        `xml:"AnswerUuid" json:"answeruuid"`
	UserId       sql.NullInt32 `xml:"UserId" json:"userid"`
}

//Create
func (survey_question_answer_link *Survey_question_answer_link) CreateSurvey_question_answer_link(db *sql.DB) (result Survey_question_answer_link, err error) {
	stmt, err := db.Prepare("INSERT INTO survey_question_answer_links ( survey_uuid, question_uuid, answer_uuid, user_id) VALUES ($1,$2,$3,$4) RETURNING id, survey_uuid, question_uuid, answer_uuid, user_id")
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(survey_question_answer_link.SurveyUuid, survey_question_answer_link.QuestionUuid, survey_question_answer_link.AnswerUuid, survey_question_answer_link.UserId).Scan(&result.Id, &result.SurveyUuid, &result.QuestionUuid, &result.AnswerUuid, &result.UserId)
	return
}

//Retrieve
func (survey_question_answer_link *Survey_question_answer_link) RetrieveSurvey_question_answer_link(db *sql.DB) (result Survey_question_answer_link, err error) {
	result = Survey_question_answer_link{}
	err = db.QueryRow("SELECT id, survey_uuid, question_uuid, answer_uuid, user_id FROM survey_question_answer_links WHERE (id = $1)", survey_question_answer_link.Id).Scan(&result.Id, &result.SurveyUuid, &result.QuestionUuid, &result.AnswerUuid, &result.UserId)
	return
}

//RetrieveAll
func RetrieveAllSurvey_question_answer_links(db *sql.DB) (survey_question_answer_links []Survey_question_answer_link, err error) {
	rows, err := db.Query("SELECT id, survey_uuid, question_uuid, answer_uuid, user_id FROM survey_question_answer_links ORDER BY id DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		result := Survey_question_answer_link{}
		if err = rows.Scan(&result.Id, &result.SurveyUuid, &result.QuestionUuid, &result.AnswerUuid, &result.UserId); err != nil {
			return
		}
		survey_question_answer_links = append(survey_question_answer_links, result)
	}
	rows.Close()
	return
}

//Update
func (survey_question_answer_link *Survey_question_answer_link) UpdateSurvey_question_answer_link(db *sql.DB) (result Survey_question_answer_link, err error) {
	stmt, err := db.Prepare("UPDATE survey_question_answer_links SET survey_uuid = $2, question_uuid = $3, answer_uuid = $4, user_id = $5 WHERE (id = $1) RETURNING id, survey_uuid, question_uuid, answer_uuid, user_id")
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(survey_question_answer_link.Id, survey_question_answer_link.SurveyUuid, survey_question_answer_link.QuestionUuid, survey_question_answer_link.AnswerUuid, survey_question_answer_link.UserId).Scan(&result.Id, &result.SurveyUuid, &result.QuestionUuid, &result.AnswerUuid, &result.UserId)
	return
}

//Delete
func (survey_question_answer_link *Survey_question_answer_link) DeleteSurvey_question_answer_link(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM survey_question_answer_links WHERE (id = $1)")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(survey_question_answer_link.Id)
	return
}

//DeleteAll
func DeleteAllSurvey_question_answer_links(db *sql.DB) (err error) {
	stmt, err := db.Prepare("DELETE FROM survey_question_answer_links")
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return
}
