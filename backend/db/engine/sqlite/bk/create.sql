
CREATE TABLE product (
	jan_code        integer,
	product_id      integer  DEFAULT '0'  NOT NULL,
	is_valid        bool     DEFAULT '0'  NOT NULL,
	stock           integer  DEFAULT '0'  NOT NULL,
	place           text     DEFAULT ''   NOT NULL,
	calorie         integer  DEFAULT '0'  NOT NULL,
	cost_price      integer  DEFAULT '0'  NOT NULL,
	retail_price    integer  DEFAULT '0'  NOT NULL,
	operation_cost  integer  DEFAULT '0'  NOT NULL,
	energy_cost     integer  DEFAULT '0'  NOT NULL,
	recipe          text     DEFAULT ''   NOT NULL,
	tag             text     DEFAULT ''   NOT NULL,
	allergy         text     DEFAULT ''   NOT NULL,
	name_ja         text     DEFAULT ''   NOT NULL,
	name_ja_short   text     DEFAULT ''   NOT NULL,
	information_ja  text     DEFAULT ''   NOT NULL,
	updated_at      text     DEFAULT (DATETIME('now', 'localtime')) NOT NULL,
	PRIMARY KEY (jan_code)
);

CREATE TABLE tag_primary (
	id            integer,
	name          text     DEFAULT ''  NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE tag_secondary (
	id              integer,
	primary_tag_id  integer,
	name            text     DEFAULT ''  NOT NULL,
	PRIMARY KEY (id)
);

