CREATE TABLE IF NOT EXISTS requests(
	id SERIAL PRIMARY KEY,
	request_type varchar(255) NOT NULL,
	request_version varchar(255) NOT NULL,
	request_hash varchar(255) NOT NULL,
	request_hashes_are_equal bool,
	creation_date timestamp default current_timestamp
)