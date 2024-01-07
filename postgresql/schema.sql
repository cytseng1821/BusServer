CREATE TABLE user (
    "id"            BIGSERIAL NOT NULL PRIMARY KEY,
    "user_id"       VARCHAR(50) NOT NULL UNIQUE,
    "account"       VARCHAR(254) NOT NULL,
    "password"      VARCHAR(64) NOT NULL,
    "created_at"    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE citybus_route (
    "id"                    BIGSERIAL NOT NULL PRIMARY KEY,
    "route_id"              VARCHAR(50) NOT NULL,
    "route_name"            VARCHAR(50) NOT NULL,
    "subroute_id"           VARCHAR(50) NOT NULL,
    "subroute_name"         VARCHAR(50) NOT NULL,
    "direction"             SMALLINT NOT NULL,
    "city"                  VARCHAR(50) NOT NULL,
    "departure_stop_name"   VARCHAR(50),
    "destination_stop_name" VARCHAR(50),
    UNIQUE("route_id", "subroute_id", "direction")
);

CREATE TABLE citybus_stop (
    "id"            BIGSERIAL NOT NULL PRIMARY KEY,
    "route_id"      INTEGER NOT NULL REFERENCES "citybus_route" ("id"),
    "stop_id"       VARCHAR(50) NOT NULL UNIQUE,
    "stop_name"     VARCHAR(50) NOT NULL,
    "stop_sequence" SMALLINT NOT NULL
);

CREATE TABLE follow_citybus_stop (
    "id"        BIGSERIAL NOT NULL PRIMARY KEY,
    "user_id"   INTEGER NOT NULL REFERENCES "user" ("id"),
    "stop_id"   INTEGER NOT NULL REFERENCES "citybus_stop" ("id"),
    UNIQUE("stop_id", "user_id")
);