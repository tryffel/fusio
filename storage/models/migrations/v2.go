package migrations

import "github.com/jinzhu/gorm"

func createPipelines(tx *gorm.DB) error {

	sql := `
CREATE TABLE pipelines (
    id         TEXT                     NOT NULL,
    owner_id   INTEGER                  NOT NULL,
    enabled    BOOLEAN                  NOT NULL,
    name       TEXT                     NOT NULL,
    lower_name TEXT                     NOT NULL,
    info       TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

	CONSTRAINT pipelines_pkey PRIMARY KEY (id),
	CONSTRAINT pipelines_owner_unique UNIQUE (owner_id, lower_name)
);

CREATE TABLE pipeline_blocks (
    id            TEXT                     NOT NULL,
    owner_id      INTEGER				   NOT NULL,
    block_type    TEXT                     NOT NULL,
    pipeline_id   TEXT					   NOT NULL,
    name          TEXT                     NOT NULL,
    block_model   TEXT                     NOT NULL,
    data          JSONB                    NOT NULL,
    next_block_id TEXT,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

	CONSTRAINT pipeline_blocks_pkey PRIMARY KEY (id),
	-- same block can be multiple times in same pipeline, but must have different next block
	CONSTRAINT pipeline_block_unique UNIQUE (id, pipeline_id, next_block_id)
);
`
	return tx.Exec(sql).Error
}
