package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/compose_attachments.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/store"
)

var _ = errors.Is

const (
	TriggerBeforeComposeAttachmentCreate triggerKey = "composeAttachmentBeforeCreate"
	TriggerBeforeComposeAttachmentUpdate triggerKey = "composeAttachmentBeforeUpdate"
	TriggerBeforeComposeAttachmentUpsert triggerKey = "composeAttachmentBeforeUpsert"
	TriggerBeforeComposeAttachmentDelete triggerKey = "composeAttachmentBeforeDelete"
)

// SearchComposeAttachments returns all matching rows
//
// This function calls convertComposeAttachmentFilter with the given
// types.AttachmentFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchComposeAttachments(ctx context.Context, f types.AttachmentFilter) (types.AttachmentSet, types.AttachmentFilter, error) {
	var scap uint
	q, err := s.convertComposeAttachmentFilter(f)
	if err != nil {
		return nil, f, err
	}

	if scap == 0 {
		scap = DefaultSliceCapacity
	}

	var (
		set = make([]*types.Attachment, 0, scap)
		// Paging is disabled in definition yaml file
		// {search: {enablePaging:false}} and this allows
		// a much simpler row fetching logic
		fetch = func() error {
			var (
				res       *types.Attachment
				rows, err = s.Query(ctx, q)
			)

			if err != nil {
				return err
			}

			for rows.Next() {
				if rows.Err() == nil {
					res, err = s.internalComposeAttachmentRowScanner(rows)
				}

				if err != nil {
					if cerr := rows.Close(); cerr != nil {
						err = fmt.Errorf("could not close rows (%v) after scan error: %w", cerr, err)
					}

					return err
				}

				// If check function is set, call it and act accordingly

				if f.Check != nil {
					if chk, err := f.Check(res); err != nil {
						if cerr := rows.Close(); cerr != nil {
							err = fmt.Errorf("could not close rows (%v) after check error: %w", cerr, err)
						}

						return err
					} else if !chk {
						// did not pass the check
						// go with the next row
						continue
					}
				}
				set = append(set, res)
			}

			return rows.Close()
		}
	)

	return set, f, s.config.ErrorHandler(fetch())
}

// LookupComposeAttachmentByID searches for attachment by its ID
//
// It returns attachment even if deleted
func (s Store) LookupComposeAttachmentByID(ctx context.Context, id uint64) (*types.Attachment, error) {
	return s.execLookupComposeAttachment(ctx, squirrel.Eq{
		s.preprocessColumn("att.id", ""): s.preprocessValue(id, ""),
	})
}

// CreateComposeAttachment creates one or more rows in compose_attachment table
func (s Store) CreateComposeAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkComposeAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.composeAttachmentHook(ctx, TriggerBeforeComposeAttachmentCreate, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execCreateComposeAttachments(ctx, s.internalComposeAttachmentEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateComposeAttachment updates one or more existing rows in compose_attachment
func (s Store) UpdateComposeAttachment(ctx context.Context, rr ...*types.Attachment) error {
	return s.config.ErrorHandler(s.PartialComposeAttachmentUpdate(ctx, nil, rr...))
}

// PartialComposeAttachmentUpdate updates one or more existing rows in compose_attachment
func (s Store) PartialComposeAttachmentUpdate(ctx context.Context, onlyColumns []string, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkComposeAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.composeAttachmentHook(ctx, TriggerBeforeComposeAttachmentUpdate, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execUpdateComposeAttachments(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("att.id", ""): s.preprocessValue(res.ID, ""),
			},
			s.internalComposeAttachmentEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return
}

// UpsertComposeAttachment updates one or more existing rows in compose_attachment
func (s Store) UpsertComposeAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		err = s.checkComposeAttachmentConstraints(ctx, res)
		if err != nil {
			return err
		}

		// err = s.composeAttachmentHook(ctx, TriggerBeforeComposeAttachmentUpsert, res)
		// if err != nil {
		// 	return err
		// }

		err = s.config.ErrorHandler(s.execUpsertComposeAttachments(ctx, s.internalComposeAttachmentEncoder(res)))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteComposeAttachment Deletes one or more rows from compose_attachment table
func (s Store) DeleteComposeAttachment(ctx context.Context, rr ...*types.Attachment) (err error) {
	for _, res := range rr {
		// err = s.composeAttachmentHook(ctx, TriggerBeforeComposeAttachmentDelete, res)
		// if err != nil {
		// 	return err
		// }

		err = s.execDeleteComposeAttachments(ctx, squirrel.Eq{
			s.preprocessColumn("att.id", ""): s.preprocessValue(res.ID, ""),
		})
		if err != nil {
			return s.config.ErrorHandler(err)
		}
	}

	return nil
}

// DeleteComposeAttachmentByID Deletes row from the compose_attachment table
func (s Store) DeleteComposeAttachmentByID(ctx context.Context, ID uint64) error {
	return s.execDeleteComposeAttachments(ctx, squirrel.Eq{
		s.preprocessColumn("att.id", ""): s.preprocessValue(ID, ""),
	})
}

// TruncateComposeAttachments Deletes all rows from the compose_attachment table
func (s Store) TruncateComposeAttachments(ctx context.Context) error {
	return s.config.ErrorHandler(s.Truncate(ctx, s.composeAttachmentTable()))
}

// execLookupComposeAttachment prepares ComposeAttachment query and executes it,
// returning types.Attachment (or error)
func (s Store) execLookupComposeAttachment(ctx context.Context, cnd squirrel.Sqlizer) (res *types.Attachment, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.composeAttachmentsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalComposeAttachmentRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateComposeAttachments updates all matched (by cnd) rows in compose_attachment with given data
func (s Store) execCreateComposeAttachments(ctx context.Context, payload store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.InsertBuilder(s.composeAttachmentTable()).SetMap(payload)))
}

// execUpdateComposeAttachments updates all matched (by cnd) rows in compose_attachment with given data
func (s Store) execUpdateComposeAttachments(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.UpdateBuilder(s.composeAttachmentTable("att")).Where(cnd).SetMap(set)))
}

// execUpsertComposeAttachments inserts new or updates matching (by-primary-key) rows in compose_attachment with given data
func (s Store) execUpsertComposeAttachments(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.composeAttachmentTable(),
		set,
		"id",
	)

	if err != nil {
		return err
	}

	return s.config.ErrorHandler(s.Exec(ctx, upsert))
}

// execDeleteComposeAttachments Deletes all matched (by cnd) rows in compose_attachment with given data
func (s Store) execDeleteComposeAttachments(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.config.ErrorHandler(s.Exec(ctx, s.DeleteBuilder(s.composeAttachmentTable("att")).Where(cnd)))
}

func (s Store) internalComposeAttachmentRowScanner(row rowScanner) (res *types.Attachment, err error) {
	res = &types.Attachment{}

	if _, has := s.config.RowScanners["composeAttachment"]; has {
		scanner := s.config.RowScanners["composeAttachment"].(func(_ rowScanner, _ *types.Attachment) error)
		err = scanner(row, res)
	} else {
		err = row.Scan(
			&res.ID,
			&res.NamespaceID,
			&res.Kind,
			&res.Url,
			&res.PreviewUrl,
			&res.Name,
			&res.Meta,
			&res.OwnerID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
		)
	}

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("could not scan db row for ComposeAttachment: %w", err)
	} else {
		return res, nil
	}
}

// QueryComposeAttachments returns squirrel.SelectBuilder with set table and all columns
func (s Store) composeAttachmentsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.composeAttachmentTable("att"), s.composeAttachmentColumns("att")...)
}

// composeAttachmentTable name of the db table
func (Store) composeAttachmentTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "compose_attachment" + alias
}

// ComposeAttachmentColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) composeAttachmentColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_namespace",
		alias + "kind",
		alias + "url",
		alias + "preview_url",
		alias + "name",
		alias + "meta",
		alias + "rel_owner",
		alias + "created_at",
		alias + "updated_at",
		alias + "deleted_at",
	}
}

// {true true false false true}

// internalComposeAttachmentEncoder encodes fields from types.Attachment to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeComposeAttachment
// func when rdbms.customEncoder=true
func (s Store) internalComposeAttachmentEncoder(res *types.Attachment) store.Payload {
	return store.Payload{
		"id":            res.ID,
		"rel_namespace": res.NamespaceID,
		"kind":          res.Kind,
		"url":           res.Url,
		"preview_url":   res.PreviewUrl,
		"name":          res.Name,
		"meta":          res.Meta,
		"rel_owner":     res.OwnerID,
		"created_at":    res.CreatedAt,
		"updated_at":    res.UpdatedAt,
		"deleted_at":    res.DeletedAt,
	}
}

func (s *Store) checkComposeAttachmentConstraints(ctx context.Context, res *types.Attachment) error {

	return nil
}

// func (s *Store) composeAttachmentHook(ctx context.Context, key triggerKey, res *types.Attachment) error {
// 	if fn, has := s.config.TriggerHandlers[key]; has {
// 		return fn.(func (ctx context.Context, s *Store, res *types.Attachment) error)(ctx, s, res)
// 	}
//
// 	return nil
// }