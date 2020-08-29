package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/credentials.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

var _ = errors.Is

const (
	TriggerBeforeCredentialsCreate triggerKey = "credentialsBeforeCreate"
	TriggerBeforeCredentialsUpdate triggerKey = "credentialsBeforeUpdate"
	TriggerBeforeCredentialsUpsert triggerKey = "credentialsBeforeUpsert"
	TriggerBeforeCredentialsDelete triggerKey = "credentialsBeforeDelete"
)

// SearchCredentials returns all matching rows
//
// This function calls convertCredentialsFilter with the given
// types.CredentialsFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchCredentials(ctx context.Context, f types.CredentialsFilter) (types.CredentialsSet, types.CredentialsFilter, error) {
	var scap uint
	q, err := s.convertCredentialsFilter(f)
	if err != nil {
		return nil, f, err
	}

	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*types.Credentials, 0, scap)
		// Paging is disabled in definition yaml file
		// {search: {enablePaging:false}} and this allows
		// a much simpler row fetching logic
		fetch = func() error {
			var (
				res       *types.Credentials
				rows, err = s.Query(ctx, q)
			)

			if err != nil {
				return err
			}

			for rows.Next() {
				if err = rows.Err(); err == nil {
					res, err = s.internalCredentialsRowScanner(rows)
				}

				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return err
				}

				// If check function is set, call it and act accordingly
				set = append(set, res)
			}

			return rows.Close()
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
}

// LookupCredentialsByID searches for credentials by ID
//
// It returns credentials even if deleted
func (s Store) LookupCredentialsByID(ctx context.Context, id uint64) (*types.Credentials, error) {
	return s.execLookupCredentials(ctx, squirrel.Eq{
		s.preprocessColumn("crd.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateCredentials creates one or more rows in credentials table
func (s Store) CreateCredentials(ctx context.Context, rr ...*types.Credentials) (err error) {
	for _, res := range rr {
		err = s.checkCredentialsConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateCredentials(ctx, s.internalCredentialsEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateCredentials updates one or more existing rows in credentials
func (s Store) UpdateCredentials(ctx context.Context, rr ...*types.Credentials) error {
	return s.config.ErrorHandler(s.PartialCredentialsUpdate(ctx, nil, rr...))
}

// PartialCredentialsUpdate updates one or more existing rows in credentials
func (s Store) PartialCredentialsUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Credentials) (err error) {
	for _, res := range rr {
		err = s.checkCredentialsConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateCredentials(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("crd.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalCredentialsEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertCredentials updates one or more existing rows in credentials
func (s Store) UpsertCredentials(ctx context.Context, rr ...*types.Credentials) (err error) {
	for _, res := range rr {
		err = s.checkCredentialsConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.config.ErrorHandler(s.execUpsertCredentials(ctx, s.internalCredentialsEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteCredentials Deletes one or more rows from credentials table
func (s Store) DeleteCredentials(ctx context.Context, rr ...*types.Credentials) (err error) {
	for _, res := range rr {

		err = s.execDeleteCredentials(ctx, squirrel.Eq{
			s.preprocessColumn("crd.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteCredentialsByID Deletes row from the credentials table
func (s Store) DeleteCredentialsByID(ctx context.Context, ID uint64) error {
	return s.execDeleteCredentials(ctx, squirrel.Eq{
		s.preprocessColumn("crd.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateCredentials Deletes all rows from the credentials table
func (s Store) TruncateCredentials(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.credentialsTable()))
}

// execLookupCredentials prepares Credentials query and executes it,
// returning types.Credentials (or error)
func (s Store) execLookupCredentials(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Credentials, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.credentialsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalCredentialsRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateCredentials updates all matched (by cnd) rows in credentials with given data
func (s Store) execCreateCredentials(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.credentialsTable()).SetMap(payload)))
}

// execUpdateCredentials updates all matched (by cnd) rows in credentials with given data
func (s Store) execUpdateCredentials(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.credentialsTable("crd")).Where(cnd).SetMap(set)))
}

// execUpsertCredentials inserts new or updates matching (by-primary-key) rows in credentials with given data
func (s Store) execUpsertCredentials(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.credentialsTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteCredentials Deletes all matched (by cnd) rows in credentials with given data
func (s Store) execDeleteCredentials(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.credentialsTable("crd")).Where(cnd)))
}

func (s Store) internalCredentialsRowScanner(row rowScanner) (res *types.Credentials, err error) {
	res = &types.Credentials{}

	if _, has := s.config.RowScanners["credentials"]; has {
		scanner := s.config.RowScanners["credentials"].(func(_ rowScanner, _ *types.Credentials) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.OwnerID,
			&res.Kind,
			&res.Label,
			&res.Credentials,
			&res.Meta,
			&res.LastUsedAt,
			&res.ExpiresAt,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for Credentials: %w", err)
	} else {
		return res, nil
	}
}

// QueryCredentials returns squirrel.SelectBuilder with set table and all columns
func (s Store) credentialsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.credentialsTable("crd"), s.credentialsColumns("crd")...)
}

// credentialsTable name of the db table
func (Store) credentialsTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "credentials" + alias
}

// CredentialsColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) credentialsColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_owner",
		alias + "kind",
		alias + "label",
		alias + "credentials",
		alias + "meta",
		alias + "last_used_at",
		alias + "expires_at",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false false false}

// internalCredentialsEncoder encodes fields from types.Credentials to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeCredentials
// func when rdbms.customEncoder=true
func (s Store) internalCredentialsEncoder(res *types.Credentials) store.Payload {
	return store.Payload{
		"id":           res.ID,
		"rel_owner":    res.OwnerID,
		"kind":         res.Kind,
		"label":        res.Label,
		"credentials":  res.Credentials,
		"meta":         res.Meta,
		"last_used_at": res.LastUsedAt,
		"expires_at":   res.ExpiresAt,
		"created_at":   res.CreatedAt,
		"updated_at":   res.UpdatedAt,
		"deleted_at":   res.DeletedAt,
	}
}

func (s *Store) checkCredentialsConstraints(ctx context.Context, res *types.Credentials) error {

	return nil
}
