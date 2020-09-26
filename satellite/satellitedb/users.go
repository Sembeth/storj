// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package satellitedb

import (
	"context"
	"strings"

	"github.com/zeebo/errs"

	"storj.io/common/uuid"
	"storj.io/storj/satellite/console"
	"storj.io/storj/satellite/satellitedb/dbx"
)

// ensures that users implements console.Users.
var _ console.Users = (*users)(nil)

// implementation of Users interface repository using spacemonkeygo/dbx orm.
type users struct {
	db dbx.Methods
}

// Get is a method for querying user from the database by id.
func (users *users) Get(ctx context.Context, id uuid.UUID) (_ *console.User, err error) {
	defer mon.Task()(&ctx)(&err)
	user, err := users.db.Get_User_By_Id(ctx, dbx.User_Id(id[:]))
	if err != nil {
		return nil, err
	}

	return userFromDBX(ctx, user)
}

// GetByEmail is a method for querying user by email from the database.
func (users *users) GetByEmail(ctx context.Context, email string) (_ *console.User, err error) {
	defer mon.Task()(&ctx)(&err)
	user, err := users.db.Get_User_By_NormalizedEmail_And_Status_Not_Number(ctx, dbx.User_NormalizedEmail(normalizeEmail(email)))

	if err != nil {
		return nil, err
	}

	return userFromDBX(ctx, user)
}

// Insert is a method for inserting user into the database.
func (users *users) Insert(ctx context.Context, user *console.User) (_ *console.User, err error) {
	defer mon.Task()(&ctx)(&err)

	if user.ID.IsZero() {
		return nil, errs.New("user id is not set")
	}

	optional := dbx.User_Create_Fields{
		ShortName: dbx.User_ShortName(user.ShortName),
	}
	if !user.PartnerID.IsZero() {
		optional.PartnerId = dbx.User_PartnerId(user.PartnerID[:])
	}
	if user.ProjectLimit != 0 {
		optional.ProjectLimit = dbx.User_ProjectLimit(user.ProjectLimit)
	}

	createdUser, err := users.db.Create_User(ctx,
		dbx.User_Id(user.ID[:]),
		dbx.User_Email(user.Email),
		dbx.User_NormalizedEmail(normalizeEmail(user.Email)),
		dbx.User_FullName(user.FullName),
		dbx.User_PasswordHash(user.PasswordHash),
		optional,
	)

	if err != nil {
		return nil, err
	}

	return userFromDBX(ctx, createdUser)
}

// Delete is a method for deleting user by Id from the database.
func (users *users) Delete(ctx context.Context, id uuid.UUID) (err error) {
	defer mon.Task()(&ctx)(&err)
	_, err = users.db.Delete_User_By_Id(ctx, dbx.User_Id(id[:]))

	return err
}

// Update is a method for updating user entity.
func (users *users) Update(ctx context.Context, user *console.User) (err error) {
	defer mon.Task()(&ctx)(&err)

	_, err = users.db.Update_User_By_Id(
		ctx,
		dbx.User_Id(user.ID[:]),
		toUpdateUser(user),
	)

	return err
}

// GetProjectLimit is a method to get the users project limit.
func (users *users) GetProjectLimit(ctx context.Context, id uuid.UUID) (limit int, err error) {
	defer mon.Task()(&ctx)(&err)

	row, err := users.db.Get_User_ProjectLimit_By_Id(ctx, dbx.User_Id(id[:]))
	if err != nil {
		return 0, err
	}
	return row.ProjectLimit, nil
}

// toUpdateUser creates dbx.User_Update_Fields with only non-empty fields as updatable.
func toUpdateUser(user *console.User) dbx.User_Update_Fields {
	update := dbx.User_Update_Fields{
		FullName:        dbx.User_FullName(user.FullName),
		ShortName:       dbx.User_ShortName(user.ShortName),
		Email:           dbx.User_Email(user.Email),
		NormalizedEmail: dbx.User_NormalizedEmail(normalizeEmail(user.Email)),
		Status:          dbx.User_Status(int(user.Status)),
		ProjectLimit:    dbx.User_ProjectLimit(user.ProjectLimit),
	}

	// extra password check to update only calculated hash from service
	if len(user.PasswordHash) != 0 {
		update.PasswordHash = dbx.User_PasswordHash(user.PasswordHash)
	}

	return update
}

// userFromDBX is used for creating User entity from autogenerated dbx.User struct.
func userFromDBX(ctx context.Context, user *dbx.User) (_ *console.User, err error) {
	defer mon.Task()(&ctx)(&err)
	if user == nil {
		return nil, errs.New("user parameter is nil")
	}

	id, err := uuid.FromBytes(user.Id)
	if err != nil {
		return nil, err
	}

	result := console.User{
		ID:           id,
		FullName:     user.FullName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Status:       console.UserStatus(user.Status),
		CreatedAt:    user.CreatedAt,
		ProjectLimit: user.ProjectLimit,
	}

	if user.PartnerId != nil {
		result.PartnerID, err = uuid.FromBytes(user.PartnerId)
		if err != nil {
			return nil, err
		}
	}

	if user.ShortName != nil {
		result.ShortName = *user.ShortName
	}

	return &result, nil
}

func normalizeEmail(email string) string {
	return strings.ToUpper(email)
}
