// Code generated by ent, DO NOT EDIT.

package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/keepcalmist/chat-service/internal/store/failedjob"
	"github.com/keepcalmist/chat-service/internal/types"
)

// FailedJobCreate is the builder for creating a FailedJob entity.
type FailedJobCreate struct {
	config
	mutation *FailedJobMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetName sets the "name" field.
func (fjc *FailedJobCreate) SetName(s string) *FailedJobCreate {
	fjc.mutation.SetName(s)
	return fjc
}

// SetPayload sets the "payload" field.
func (fjc *FailedJobCreate) SetPayload(s string) *FailedJobCreate {
	fjc.mutation.SetPayload(s)
	return fjc
}

// SetReason sets the "reason" field.
func (fjc *FailedJobCreate) SetReason(s string) *FailedJobCreate {
	fjc.mutation.SetReason(s)
	return fjc
}

// SetCreatedAt sets the "created_at" field.
func (fjc *FailedJobCreate) SetCreatedAt(t time.Time) *FailedJobCreate {
	fjc.mutation.SetCreatedAt(t)
	return fjc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (fjc *FailedJobCreate) SetNillableCreatedAt(t *time.Time) *FailedJobCreate {
	if t != nil {
		fjc.SetCreatedAt(*t)
	}
	return fjc
}

// SetID sets the "id" field.
func (fjc *FailedJobCreate) SetID(tji types.FailedJobID) *FailedJobCreate {
	fjc.mutation.SetID(tji)
	return fjc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (fjc *FailedJobCreate) SetNillableID(tji *types.FailedJobID) *FailedJobCreate {
	if tji != nil {
		fjc.SetID(*tji)
	}
	return fjc
}

// Mutation returns the FailedJobMutation object of the builder.
func (fjc *FailedJobCreate) Mutation() *FailedJobMutation {
	return fjc.mutation
}

// Save creates the FailedJob in the database.
func (fjc *FailedJobCreate) Save(ctx context.Context) (*FailedJob, error) {
	fjc.defaults()
	return withHooks(ctx, fjc.sqlSave, fjc.mutation, fjc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (fjc *FailedJobCreate) SaveX(ctx context.Context) *FailedJob {
	v, err := fjc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fjc *FailedJobCreate) Exec(ctx context.Context) error {
	_, err := fjc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fjc *FailedJobCreate) ExecX(ctx context.Context) {
	if err := fjc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fjc *FailedJobCreate) defaults() {
	if _, ok := fjc.mutation.CreatedAt(); !ok {
		v := failedjob.DefaultCreatedAt()
		fjc.mutation.SetCreatedAt(v)
	}
	if _, ok := fjc.mutation.ID(); !ok {
		v := failedjob.DefaultID()
		fjc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fjc *FailedJobCreate) check() error {
	if _, ok := fjc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`store: missing required field "FailedJob.name"`)}
	}
	if v, ok := fjc.mutation.Name(); ok {
		if err := failedjob.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`store: validator failed for field "FailedJob.name": %w`, err)}
		}
	}
	if _, ok := fjc.mutation.Payload(); !ok {
		return &ValidationError{Name: "payload", err: errors.New(`store: missing required field "FailedJob.payload"`)}
	}
	if v, ok := fjc.mutation.Payload(); ok {
		if err := failedjob.PayloadValidator(v); err != nil {
			return &ValidationError{Name: "payload", err: fmt.Errorf(`store: validator failed for field "FailedJob.payload": %w`, err)}
		}
	}
	if _, ok := fjc.mutation.Reason(); !ok {
		return &ValidationError{Name: "reason", err: errors.New(`store: missing required field "FailedJob.reason"`)}
	}
	if v, ok := fjc.mutation.Reason(); ok {
		if err := failedjob.ReasonValidator(v); err != nil {
			return &ValidationError{Name: "reason", err: fmt.Errorf(`store: validator failed for field "FailedJob.reason": %w`, err)}
		}
	}
	if _, ok := fjc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`store: missing required field "FailedJob.created_at"`)}
	}
	if v, ok := fjc.mutation.ID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`store: validator failed for field "FailedJob.id": %w`, err)}
		}
	}
	return nil
}

func (fjc *FailedJobCreate) sqlSave(ctx context.Context) (*FailedJob, error) {
	if err := fjc.check(); err != nil {
		return nil, err
	}
	_node, _spec := fjc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fjc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*types.FailedJobID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	fjc.mutation.id = &_node.ID
	fjc.mutation.done = true
	return _node, nil
}

func (fjc *FailedJobCreate) createSpec() (*FailedJob, *sqlgraph.CreateSpec) {
	var (
		_node = &FailedJob{config: fjc.config}
		_spec = sqlgraph.NewCreateSpec(failedjob.Table, sqlgraph.NewFieldSpec(failedjob.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = fjc.conflict
	if id, ok := fjc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := fjc.mutation.Name(); ok {
		_spec.SetField(failedjob.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := fjc.mutation.Payload(); ok {
		_spec.SetField(failedjob.FieldPayload, field.TypeString, value)
		_node.Payload = value
	}
	if value, ok := fjc.mutation.Reason(); ok {
		_spec.SetField(failedjob.FieldReason, field.TypeString, value)
		_node.Reason = value
	}
	if value, ok := fjc.mutation.CreatedAt(); ok {
		_spec.SetField(failedjob.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.FailedJob.Create().
//		SetName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.FailedJobUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (fjc *FailedJobCreate) OnConflict(opts ...sql.ConflictOption) *FailedJobUpsertOne {
	fjc.conflict = opts
	return &FailedJobUpsertOne{
		create: fjc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.FailedJob.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (fjc *FailedJobCreate) OnConflictColumns(columns ...string) *FailedJobUpsertOne {
	fjc.conflict = append(fjc.conflict, sql.ConflictColumns(columns...))
	return &FailedJobUpsertOne{
		create: fjc,
	}
}

type (
	// FailedJobUpsertOne is the builder for "upsert"-ing
	//  one FailedJob node.
	FailedJobUpsertOne struct {
		create *FailedJobCreate
	}

	// FailedJobUpsert is the "OnConflict" setter.
	FailedJobUpsert struct {
		*sql.UpdateSet
	}
)

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.FailedJob.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(failedjob.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *FailedJobUpsertOne) UpdateNewValues() *FailedJobUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(failedjob.FieldID)
		}
		if _, exists := u.create.mutation.Name(); exists {
			s.SetIgnore(failedjob.FieldName)
		}
		if _, exists := u.create.mutation.Payload(); exists {
			s.SetIgnore(failedjob.FieldPayload)
		}
		if _, exists := u.create.mutation.Reason(); exists {
			s.SetIgnore(failedjob.FieldReason)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(failedjob.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.FailedJob.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *FailedJobUpsertOne) Ignore() *FailedJobUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *FailedJobUpsertOne) DoNothing() *FailedJobUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the FailedJobCreate.OnConflict
// documentation for more info.
func (u *FailedJobUpsertOne) Update(set func(*FailedJobUpsert)) *FailedJobUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&FailedJobUpsert{UpdateSet: update})
	}))
	return u
}

// Exec executes the query.
func (u *FailedJobUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("store: missing options for FailedJobCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *FailedJobUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *FailedJobUpsertOne) ID(ctx context.Context) (id types.FailedJobID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("store: FailedJobUpsertOne.ID is not supported by MySQL driver. Use FailedJobUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *FailedJobUpsertOne) IDX(ctx context.Context) types.FailedJobID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// FailedJobCreateBulk is the builder for creating many FailedJob entities in bulk.
type FailedJobCreateBulk struct {
	config
	err      error
	builders []*FailedJobCreate
	conflict []sql.ConflictOption
}

// Save creates the FailedJob entities in the database.
func (fjcb *FailedJobCreateBulk) Save(ctx context.Context) ([]*FailedJob, error) {
	if fjcb.err != nil {
		return nil, fjcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(fjcb.builders))
	nodes := make([]*FailedJob, len(fjcb.builders))
	mutators := make([]Mutator, len(fjcb.builders))
	for i := range fjcb.builders {
		func(i int, root context.Context) {
			builder := fjcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FailedJobMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, fjcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = fjcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fjcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, fjcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fjcb *FailedJobCreateBulk) SaveX(ctx context.Context) []*FailedJob {
	v, err := fjcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fjcb *FailedJobCreateBulk) Exec(ctx context.Context) error {
	_, err := fjcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fjcb *FailedJobCreateBulk) ExecX(ctx context.Context) {
	if err := fjcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.FailedJob.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.FailedJobUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (fjcb *FailedJobCreateBulk) OnConflict(opts ...sql.ConflictOption) *FailedJobUpsertBulk {
	fjcb.conflict = opts
	return &FailedJobUpsertBulk{
		create: fjcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.FailedJob.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (fjcb *FailedJobCreateBulk) OnConflictColumns(columns ...string) *FailedJobUpsertBulk {
	fjcb.conflict = append(fjcb.conflict, sql.ConflictColumns(columns...))
	return &FailedJobUpsertBulk{
		create: fjcb,
	}
}

// FailedJobUpsertBulk is the builder for "upsert"-ing
// a bulk of FailedJob nodes.
type FailedJobUpsertBulk struct {
	create *FailedJobCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.FailedJob.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(failedjob.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *FailedJobUpsertBulk) UpdateNewValues() *FailedJobUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(failedjob.FieldID)
			}
			if _, exists := b.mutation.Name(); exists {
				s.SetIgnore(failedjob.FieldName)
			}
			if _, exists := b.mutation.Payload(); exists {
				s.SetIgnore(failedjob.FieldPayload)
			}
			if _, exists := b.mutation.Reason(); exists {
				s.SetIgnore(failedjob.FieldReason)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(failedjob.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.FailedJob.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *FailedJobUpsertBulk) Ignore() *FailedJobUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *FailedJobUpsertBulk) DoNothing() *FailedJobUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the FailedJobCreateBulk.OnConflict
// documentation for more info.
func (u *FailedJobUpsertBulk) Update(set func(*FailedJobUpsert)) *FailedJobUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&FailedJobUpsert{UpdateSet: update})
	}))
	return u
}

// Exec executes the query.
func (u *FailedJobUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("store: OnConflict was set for builder %d. Set it on the FailedJobCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("store: missing options for FailedJobCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *FailedJobUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
