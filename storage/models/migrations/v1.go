package migrations

import (
	"github.com/jinzhu/gorm"
)

func initialSchema(tx *gorm.DB) error {

	sql := `
CREATE TABLE users
(
  id         SERIAL                   NOT NULL,
  name       TEXT                     UNIQUE,
  lower_name TEXT                     UNIQUE,
  email      TEXT,
  password   TEXT,
  last_seen  TIMESTAMP WITH TIME ZONE,
  is_active  BOOLEAN                  NOT NULL,
  is_admin   BOOLEAN                  NOT NULL,
  can_delete BOOLEAN                  NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT user_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_users_deleted_at
  ON users (deleted_at);

CREATE table groups
(
  id         TEXT                     NOT NULL UNIQUE ,
  name       TEXT                     NOT NULL,
  lower_name TEXT                     NOT NULL,
  info       TEXT,
  owner_id   INTEGER                  NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

  CONSTRAINT groups_pkey PRIMARY KEY (id),
  CONSTRAINT owner_groups_unique UNIQUE (lower_name, owner_id)
);


CREATE TABLE devices
(
  id          TEXT                     NOT NULL UNIQUE ,
  name        TEXT                     NOT NULL,
  lower_name  TEXT                     NOT NULL,
  info        TEXT,
  owner_id    INTEGER                  NOT NULL,
  device_type TEXT                     NOT NULL,
  created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,

  CONSTRAINT devices_pkey PRIMARY KEY (id),
  CONSTRAINT owner_devices_unique UNIQUE (lower_name, owner_id)
);


CREATE TABLE groups_devices
(
  device_id TEXT  NOT NULL,
  group_id  TEXT  NOT NULL,

  CONSTRAINT groups_devices_pkey
    PRIMARY KEY (device_id, group_id)
);


CREATE TABLE alarms
(
  id           TEXT                     NOT NULL,
  name         TEXT                     NOT NULL,
  lower_name   TEXT                     NOT NULL,
  info         TEXT,
  message      TEXT                     NOT NULL,
  owner_id     INTEGER                  NOT NULL,
  "group"      TEXT                     NOT NULL,
  fired        BOOLEAN                  NOT NULL,
  enabled      BOOLEAN,
  query        TEXT,
  run_interval BIGINT                   NOT NULL,
  last_run     TIMESTAMP WITH TIME ZONE,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL,

  CONSTRAINT alarms_pkey
    PRIMARY KEY (id),
  CONSTRAINT alarms_groups_unique UNIQUE ("group", lower_name)
);


CREATE TABLE alarm_histories
(
  id         SERIAL          NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE,
  alarm_id   TEXT            NOT NULL,
  value      TEXT            NOT NULL,
  cleared    BOOLEAN         NOT NULL,
  CONSTRAINT alarm_histories_pkey
    PRIMARY KEY (id)
);

CREATE INDEX idx_alarm_histories_deleted_at
  ON alarm_histories (deleted_at);


CREATE TABLE api_keys
(
  id         SERIAL          NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at timestamp WITH TIME ZONE,
  name       TEXT,
  device_id  TEXT            NOT NULL,
  key        TEXT            NOT NULL,
  last_seen  TIMESTAMP WITH TIME ZONE,
  CONSTRAINT api_keys_pkey
    PRIMARY key (id)
);

CREATE INDEX idx_api_keys_deleted_at
  ON api_keys (deleted_at);

CREATE TABLE output_channels
(
  id          TEXT                     NOT NULL,
  owner_id    INTEGER                  NOT NULL,
  name        TEXT                     NOT NULL,
  lower_name  TEXT                     NOT NULL,
  output_type TEXT                     NOT NULL,
  data        TEXT,
  created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,

  CONSTRAINT output_channels_pkey
    PRIMARY KEY (id),
  CONSTRAINT owner_channel_unique UNIQUE (owner_id, lower_name)
);


CREATE TABLE outputs
(
  id                TEXT                     NOT NULL,
  owner_id          INTEGER                  NOT NULL,
  alarm_id          TEXT                     NOT NULL,
  name              TEXT,
  lower_name        TEXT                     NOT NULL,
  output_channel_id TEXT,
  fire_template     TEXT,
  clear_template    TEXT,
  error_template    TEXT,
  repeat            BIGINT,
  last_pushed       TIMESTAMP WITH TIME ZONE,
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at        TIMESTAMP WITH TIME ZONE NOT NULL,
  enabled           BOOLEAN DEFAULT TRUE,
  on_fire           BOOLEAN DEFAULT TRUE,
  on_clear          BOOLEAN DEFAULT TRUE,
  on_error          BOOLEAN DEFAULT FALSE,

  CONSTRAINT outputs_pkey
    PRIMARY KEY (id),
  CONSTRAINT name_channel_unique UNIQUE (owner_id, lower_name, output_channel_id)
);


CREATE TABLE output_histories
(
  id         SERIAL  NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE,
  output_id  TEXT    NOT NULL,
  success    BOOLEAN NOT NULL,
  message    TEXT,
  CONSTRAINT output_histories_pkey
    PRIMARY KEY (id)
);

CREATE INDEX idx_output_histories_deleted_at
  ON output_histories (deleted_at);

`
	return tx.Exec(sql).Error
}
