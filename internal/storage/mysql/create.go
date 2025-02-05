package mysql

import (
	"bytes"
	"embed"
	"errors"
	"text/template"
)

//go:embed create.sql
var sql embed.FS

func (s *Store) createSQL() (string, error) {
	queryContent, _ := sql.ReadFile("create.sql")
	tmpl, err := template.New("create").Parse(string(queryContent))
	if err != nil {
		return "", err
	}
	conf := s.DSNSource.(Config)
	// The connection we need to use for create must be schema-less so
	// that we can create the database, so we need to override that with
	// the default schema here.
	if conf.TargetSchema == "" {
		return "", errors.New("can't create database, empty target schema")
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, conf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *Store) Create() error {
	q, err := s.createSQL()
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(s.Ctx, q)
	return err
}

func (s *Store) Maintain() error {
	return errors.New("mysql: maintain not implemented")
}
