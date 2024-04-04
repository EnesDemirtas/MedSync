-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
	user_id       UUID        NOT NULL,
	name          TEXT        NOT NULL,
	email         TEXT UNIQUE NOT NULL,
	roles         TEXT[]      NOT NULL,
	password_hash TEXT        NOT NULL,
    enabled       BOOLEAN     NOT NULL,
	date_created  TIMESTAMP   NOT NULL,
	date_updated  TIMESTAMP   NOT NULL,

	PRIMARY KEY (user_id)
);

-- Version: 1.02
-- Description: Create table tags
CREATE TABLE tags (
	tag_id UUID NOT NULL,
	name   TEXT NOT NULL,

	PRIMARY KEY (tag_id)
);

-- Version: 1.03
-- Description: Create table medicines
CREATE TABLE medicines (
	medicine_id  UUID      NOT NULL,
	name         TEXT      NOT NULL,
    description  TEXT      NULL,
    tags         UUID[]    NULL,
	date_created TIMESTAMP NOT NULL,
	date_updated TIMESTAMP NOT NULL,

	PRIMARY KEY (medicine_id)
);

-- Version: 1.04
-- Description: Create table inventories
CREATE TABLE inventories (
    inventory_id        UUID  NOT NULL,
    name                TEXT  NOT NULL,
    description         TEXT  NULL,
    medicine_quantities JSONB NULL,
    date_created        TIMESTAMP NOT NULL,
    date_updated        TIMESTAMP NOT NULL,

    PRIMARY KEY (inventory_id)
);

