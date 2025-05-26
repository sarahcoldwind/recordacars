-- +goose Up
-- +goose StatementBegin
CREATE TABLE acars_messages (
	id bigserial primary key,
	timestamp timestamptz not null,
	station_id text,
	channel int not null,
	freq float not null,
	level float not null,
	error int not null,
	mode text not null,
	label text not null,
	block_id text,
	tail text,
	flight text,
	msgno text,
	text text,
	"end" bool not null,
	depa text,
	dsta text,
	eta text,
	gtout text,
	gtin text,
	wloff text,
	won text,
	sublabel text,
	mfi text,
	assstat text,
	libacars jsonb,
	app_name text not null,
	app_ver text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE acars_messages;
-- +goose StatementEnd
