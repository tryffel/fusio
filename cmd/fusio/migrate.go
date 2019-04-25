package main

import (
	"github.com/akamensky/argparse"
	"github.com/jinzhu/gorm"
	"github.com/tryffel/fusio/storage/models/migrations"
)

type Migrator struct {
	command *argparse.Command
	n       *int
	t       *string
	all     *bool
	skip    *bool
}

func NewMigrator(parser *argparse.Parser) *Migrator {
	m := &Migrator{
		command: parser.NewCommand("migrate", "run update migrations. This should be run whenever "+
			"server is updated. To migrate to latest version, run 'fusio migrate -a'"),
	}

	m.n = parser.Int("n", "N", &argparse.Options{Required: false, Help: "How many steps to run", Default: -1})
	m.t = parser.Selector("t", "Type", []string{"up", "down"}, &argparse.Options{Required: false,
		Default: "up", Help: "which operation to run. Up/Down will upgrade/downgrade. Default up"})
	m.all = parser.Flag("a", "All", &argparse.Options{Required: false, Default: true, Help: "Update to latest version. Default true"})
	m.skip = parser.Flag("s", "Skip", &argparse.Options{Required: false, Default: false, Help: "Skip automatic migrations. Default false"})

	return m
}

func (m *Migrator) RunMigrations(db *gorm.DB) error {
	return migrations.RunMigrations(db, -1)
}
