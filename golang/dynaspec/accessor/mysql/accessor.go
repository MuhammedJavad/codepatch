package dynaspec

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MuhammedJavad/codepatch/dynaspec/tree"
)

type Accessor struct {
	db *sql.DB
}

type treeModel struct {
	ID        uint
	CreatedAt time.Time
	Active    bool
	Name      string
	StartTime *time.Time
	EndTime   *time.Time
	Result    json.RawMessage
	Structure json.RawMessage
}

func New(db *sql.DB) Accessor {
	return Accessor{db: db}
}

func (a *Accessor) Get(ctx context.Context, id uint) (*tree.Tree, error) {
	var dao treeModel
	const query = `SELECT id, active, name, start_time, end_time, result, structure FROM trees WHERE id = ? AND active`
	err := a.db.QueryRowContext(ctx, query, id).Scan(
		&dao.ID,
		&dao.Active,
		&dao.Name,
		&dao.StartTime,
		&dao.EndTime,
		&dao.Result,
		&dao.Structure,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tree %d: %w", id, err)
	}

	var root tree.Node
	if err = json.Unmarshal(dao.Structure, &root); err != nil {
		return nil, fmt.Errorf("failed to create tree node: %w", err)
	}

	tree := tree.NewTree(dao.ID, dao.Name, dao.StartTime, dao.EndTime, dao.Active, dao.Result, root)
	return &tree, nil
}

func (a *Accessor) List(ctx context.Context) ([]tree.Tree, error) {
	const query = `SELECT id, active, name, start_time, end_time, result, structure FROM trees WHERE active ORDER BY created_at DESC`
	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trees: %w", err)
	}
	defer rows.Close()

	var trees []tree.Tree
	for rows.Next() {
		var m treeModel
		err := rows.Scan(
			&m.ID,
			&m.Active,
			&m.Name,
			&m.StartTime,
			&m.EndTime,
			&m.Result,
			&m.Structure,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tree: %w", err)
		}

		var root tree.Node
		if err = json.Unmarshal(m.Structure, &root); err != nil {
			return nil, fmt.Errorf("failed to create tree node: %w", err)
		}

		tree := tree.NewTree(m.ID, m.Name, m.StartTime, m.EndTime, m.Active, m.Result, root)
		trees = append(trees, tree)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trees: %w", err)
	}

	return trees, nil
}

func (a *Accessor) Create(ctx context.Context, tree *tree.Tree) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = create(ctx, tx, tree)
	if err != nil {
		return fmt.Errorf("failed to delete tree: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *Accessor) Update(ctx context.Context, tree *tree.Tree) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = delete(ctx, tx, tree)
	if err != nil {
		return fmt.Errorf("failed to delete tree: %w", err)
	}

	err = create(ctx, tx, tree)
	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *Accessor) Delete(ctx context.Context, tree *tree.Tree) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = delete(ctx, tx, tree)
	if err != nil {
		return fmt.Errorf("failed to delete tree: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func delete(ctx context.Context, tx *sql.Tx, tree *tree.Tree) error {
	const query = `DELETE FROM trees WHERE id = ?`
	result, err := tx.ExecContext(ctx, query, tree.ID)
	if err != nil {
		return fmt.Errorf("failed to delete tree %d: %w", tree.ID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tree with id %d not found", tree.ID)
	}

	return nil
}

func create(ctx context.Context, tx *sql.Tx, tree *tree.Tree) error {
	structureJSON, err := json.Marshal(tree.Root)
	if err != nil {
		return fmt.Errorf("failed to marshal tree structure: %w", err)
	}

	const query = `INSERT INTO trees (name, start_time, end_time, result, structure, active, created_at) VALUES (?,?,?,?,?,?,?)`
	result, err := tx.ExecContext(ctx, query, tree.Name, tree.Start, tree.End, tree.Result, structureJSON, tree.Active, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create tree: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	tree.ID = uint(id)
	return nil
}
