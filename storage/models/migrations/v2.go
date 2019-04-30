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
	data	   JSONB				    NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

	CONSTRAINT pipelines_pkey PRIMARY KEY (id),
	CONSTRAINT pipelines_owner_unique UNIQUE (owner_id, lower_name)
);
`
	return tx.Exec(sql).Error
}
