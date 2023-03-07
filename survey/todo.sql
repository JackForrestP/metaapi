
create table surveys (
  id           integer generated always as identity primary key,
  uuid text  not null UNIQUE,
  version text not null,
  title         text not null,
  description        text not null,
  open        boolean not null,
  creation_date timestamptz  not null
);

create table questions (
  id           integer generated always as identity primary key,
  uuid text not null UNIQUE,
  question text not null,
  focus text not null,
  unit text not null,
  back_image text not null,
  min_value integer not null,
  max_value integer not null,
  default_value integer not null,
  warp_function text not null,
  creation_date timestamptz  not null
);

create table answers (

  id integer generated always as identity primary key,
  uuid text not null UNIQUE,
  answer_column text not null,
  int_answer int,
  string_answer text,
  creation_date timestamptz not null
);

create table survey_question_links (
  id       integer generated always as identity primary key,
  survey_uuid text not null REFERENCES surveys(uuid),
  question_uuid text not null REFERENCES questions(uuid)
);

create table survey_question_answer_links  (
  id integer generated always as identity primary key,
  survey_uuid text not null REFERENCES surveys(uuid),
  question_uuid text not null REFERENCES questions(uuid),
  answer_uuid text not null REFERENCES answers(uuid),
  user_id int REFERENCES users(id)
);

