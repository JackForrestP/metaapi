create table community_results (
  id           integer generated always as identity primary key,
  uuid varchar(36) not null UNIQUE,
  product_uuid  varchar(36) not null UNIQUE REFERENCES products(uuid),
  metric_graphs text not null
);

create table product_reviews (
  id           integer generated always as identity primary key,
  uuid varchar(36) not null UNIQUE,
  title text not null UNIQUE,
  author text not null,
  creation_date timestamptz not null,
  last_updated timestamptz not null,
  hero_image text not null,
  summary text not null,
  rating text not null,  
  score decimal not null,
  body_html text not null,
  user_id integer not null REFERENCES users(id),
  product_uuid varchar(36) not null REFERENCES products(uuid),
  expected_tst integer,
  expected_deep integer,
  expected_rem integer,
  expected_onset integer,
  expected_wakefulness integer,
  expected_efficiency integer,
  expected_accuracy integer,
  reviewed_results text,
  validated boolean not null,
  start_date timestamptz not null,
  end_date timestamptz not null,
  is_preview boolean not null
);

create table product_listings (
  id           integer generated always as identity primary key,
  uuid varchar(36) not null UNIQUE,
  vendor_uuid varchar(36) not null REFERENCES vendors(uuid),
  creation_date timestamptz not null,
  last_updated timestamptz not null,
  purchase_link text not null,
  price decimal not null,
  product_uuid varchar(36) not null REFERENCES products(uuid),
  is_manufacterer boolean not null
);

create table photo_urls (
  id           integer generated always as identity primary key,
  uuid varchar(36) not null UNIQUE,
  url text not null ,
  title text not null ,
  alt_text text not null,
  keywords text not null,
  user_id integer not null REFERENCES users(id),
  creation_date timestamptz not null
);

create table vendors (
  id           integer generated always as identity primary key,
  uuid varchar(36) not null UNIQUE,
  name text not null UNIQUE,
  link text,
  creation_date timestamptz not null,
  last_updated timestamptz not null,
  affiliate_id text
);

create table solutions ( 
		id integer generated always as identity primary key , 
		uuid varchar(36) not null unique ,
		name text not null,
		product_uuid varchar(36)  ,
		technique_uuid varchar(36) not null  ,
		user_id integer not null references users(id) ,
		created_at timestamptz not null , 
		is_public boolean not null , 
		validated boolean not null , 
		last_updated timestamptz not null ,
		archived boolean not null,
    description text
  ) ;

create table products ( 
		id integer generated always as identity primary key , 
		uuid varchar(36) not null unique ,
		name text not null,
		model text not null,
		category integer not null references product_categories(id)  , 
		manufacturer text not null , 
		image_link text not null , 
		purchase_link text not null , 
		affiliated_id text not null ,
		description text not null , 
		user_id integer not null references users(id) ,
		created_at timestamptz not null , 
		is_public boolean not null,
		validated boolean not null , 
		archived boolean not null ,
    manufacturer_uuid varchar(36) refrences vendor(uuid) ,
    manufacturer_product_id text,
    manufacturer_product_link text,
    msrp decimal
  ) ;

  create table articles ( 
		id integer generated always as identity primary key , 
		uuid varchar(36) not null unique ,
		title text not null,
    summary text not null,
    urls text not null,
		author text not null,
		category text not null , 
		image_link text not null , 
		creation_date timestamptz not null , 
		user_id integer not null references users(id) ,
    article_link text not null,
    body_html text not null,
    validated boolean not null
  ) ;

    create table product_trials ( 
      id integer generated always as identity primary key , 
      uuid varchar(36) not null unique ,
      user_id integer not null references users(id), 
      product_uuid varchar(36) not null references products(uuid),
      start_date timestamptz,
      end_date timestamptz,
      creation_date timestamptz not null,
      archived boolean not null
  ) ;

  create table hypnostats (
    id integer generated always as identity primary key , 
    user_id integer not null references users(id) , 
    created_at timestamptz not null , 
    source text not null , 
    use boolean not null , 
    hypno jsonb not null , 
    hypno_model int not null , 
    motion_len int not null , 
    energy_len int not null , 
    heartrate_len int not null , 
    stand_len int not null , 
    sleepsleep_len int not null , 
    sleepinbed_len int not null , 
    begin_bed_rel decimal not null , 
    begin_bed decimal not null , 
    end_bed decimal not null , 
    hr_min decimal not null , 
    hr_max decimal not null , 
    hr_avg decimal not null , 
    tst decimal not null , 
    time_awake decimal not null , 
    time_rem decimal not null , 
    time_light decimal not null , 
    time_deep decimal not null , 
    num_awakes decimal not null , 
    score decimal not null , 
    score_duration decimal not null , 
    score_efficiency decimal not null , 
    score_continuity decimal not null , 
    utc_offset integer not null ,
    sleep_onset decimal not null ,
    sleep_efficiency decimal not null
  ) ;

 create table product_comments (
    id integer generated always as identity primary key , 
    uuid varchar(36) not null unique ,
    user_id integer not null references users(id) , 
    created_at timestamptz not null , 
    last_updated timestamptz not null ,
    comments text ,
    rating int not null ,
    product_uuid text not null refrences products(uuid)
  ) ;

  create table community_product_trial_stats (
    id integer generated always as identity primary key , 
    created_at timestamptz not null , 
    product_uuid varchar(36) not null references products(uuid) ,
    source text not null , 
    days int not null ,
    user_count int not null ,
    tst decimal , 
    time_awake decimal , 
    time_rem decimal , 
    time_light decimal , 
    time_deep decimal , 
    num_awakes decimal , 
    score decimal , 
    sleep_onset decimal  ,
    sleep_efficiency decimal 
  ) ;

   create table habits ( 
      id integer generated always as identity primary key , 
      uuid varchar(36) not null unique ,
      user_id integer not null references users(id), 
      product_uuid varchar(36) not null references products(uuid),
      start_date timestamptz,
      end_date timestamptz,
      creation_date timestamptz not null,
      archived boolean not null
  ) ;

  create table community_product_results (
    id integer generated always as identity primary key , 
    user_count integer not null , 
    count integer not null ,
    product_uuid varchar(36) not null references products(uuid) ,
    tst decimal , 
    time_awake decimal , 
    time_rem decimal , 
    time_light decimal , 
    time_deep decimal , 
    num_awakes decimal , 
    score decimal , 
    sleep_onset decimal  ,
    sleep_efficiency decimal 
  ) ;
