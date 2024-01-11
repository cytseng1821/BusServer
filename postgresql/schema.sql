CREATE TABLE citybus_user (
    "id"            BIGSERIAL NOT NULL PRIMARY KEY,
    "user_id"       VARCHAR(50) NOT NULL UNIQUE,
    "account"       VARCHAR(254) NOT NULL,
    "password"      VARCHAR(64) NOT NULL,
    "created_at"    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE citybus_route (
    "id"                    BIGSERIAL NOT NULL PRIMARY KEY,
    "route_id"              VARCHAR(50) NOT NULL UNIQUE,
    "route_name"            VARCHAR(50) NOT NULL,
    "city"                  VARCHAR(50) NOT NULL,
    "departure_stop_name"   VARCHAR(50),
    "destination_stop_name" VARCHAR(50)
);

CREATE TABLE citybus_subroute (
    "id"            BIGSERIAL NOT NULL PRIMARY KEY,
    "route_id"      INTEGER NOT NULL REFERENCES "citybus_route" ("id"),
    "subroute_id"   VARCHAR(50) NOT NULL,
    "subroute_name" VARCHAR(50) NOT NULL,
    "direction"     SMALLINT NOT NULL,
    UNIQUE("subroute_id", "direction")
);

CREATE TABLE citybus_stop (
    "id"            BIGSERIAL NOT NULL PRIMARY KEY,
    "route_id"      INTEGER NOT NULL REFERENCES "citybus_route" ("id"),
    "stop_id"       VARCHAR(50) NOT NULL UNIQUE,
    "stop_name"     VARCHAR(50) NOT NULL
);

CREATE TABLE citybus_subroute_stop_relation (
    "id"            BIGSERIAL NOT NULL PRIMARY KEY,
    "subroute_id"   INTEGER NOT NULL REFERENCES "citybus_subroute" ("id"),
    "stop_id"       INTEGER NOT NULL REFERENCES "citybus_stop" ("id"),
    "stop_sequence" SMALLINT NOT NULL,
    UNIQUE("subroute_id", "stop_sequence")
)

CREATE TABLE citybus_user_stop (
    "id"        BIGSERIAL NOT NULL PRIMARY KEY,
    "user_id"   INTEGER NOT NULL REFERENCES "citybus_user" ("id"),
    "stop_id"   INTEGER NOT NULL REFERENCES "citybus_stop" ("id"),
    UNIQUE("stop_id", "user_id")
);