// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"

	"github.com/morevault/vaultum/sdk/v2/physical"
)

type LogicalStorage struct {
	underlying physical.Backend
}

var _ Storage = &LogicalStorage{}

func (s *LogicalStorage) Get(ctx context.Context, key string) (*StorageEntry, error) {
	entry, err := s.underlying.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &StorageEntry{
		Key:      entry.Key,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	}, nil
}

func (s *LogicalStorage) Put(ctx context.Context, entry *StorageEntry) error {
	return s.underlying.Put(ctx, &physical.Entry{
		Key:      entry.Key,
		Value:    entry.Value,
		SealWrap: entry.SealWrap,
	})
}

func (s *LogicalStorage) Delete(ctx context.Context, key string) error {
	return s.underlying.Delete(ctx, key)
}

func (s *LogicalStorage) List(ctx context.Context, prefix string) ([]string, error) {
	return s.underlying.List(ctx, prefix)
}

func (s *LogicalStorage) ListPage(ctx context.Context, prefix string, after string, limit int) ([]string, error) {
	return s.underlying.ListPage(ctx, prefix, after, limit)
}

func (s *LogicalStorage) Underlying() physical.Backend {
	return s.underlying
}

type TransactionalLogicalStorage struct {
	LogicalStorage
}

var _ TransactionalStorage = &TransactionalLogicalStorage{}

type LogicalTransaction struct {
	LogicalStorage
}

var _ Transaction = &LogicalTransaction{}

func (s *TransactionalLogicalStorage) BeginReadOnlyTx(ctx context.Context) (Transaction, error) {
	tx, err := s.Underlying().(physical.TransactionalBackend).BeginReadOnlyTx(ctx)
	if err != nil {
		return nil, err
	}

	return &LogicalTransaction{
		LogicalStorage{
			underlying: tx,
		},
	}, nil
}

func (s *TransactionalLogicalStorage) BeginTx(ctx context.Context) (Transaction, error) {
	tx, err := s.Underlying().(physical.TransactionalBackend).BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	return &LogicalTransaction{
		LogicalStorage{
			underlying: tx,
		},
	}, nil
}

func (s *LogicalTransaction) Commit(ctx context.Context) error {
	return s.Underlying().(physical.Transaction).Commit(ctx)
}

func (s *LogicalTransaction) Rollback(ctx context.Context) error {
	return s.Underlying().(physical.Transaction).Rollback(ctx)
}

func NewLogicalStorage(underlying physical.Backend) Storage {
	ls := &LogicalStorage{
		underlying: underlying,
	}

	if _, ok := underlying.(physical.TransactionalBackend); ok {
		return &TransactionalLogicalStorage{
			*ls,
		}
	}

	return ls
}
