package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/types/types.yaml

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplicationSetWalk(t *testing.T) {
	var (
		value = make(ApplicationSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*Application) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*Application) error { return fmt.Errorf("walk error") }))
}

func TestApplicationSetFilter(t *testing.T) {
	var (
		value = make(ApplicationSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*Application) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*Application) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*Application) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestApplicationSetIDs(t *testing.T) {
	var (
		value = make(ApplicationSet, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new(Application)
	value[1] = new(Application)
	value[2] = new(Application)
	// set ids
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	// Find existing
	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	// List IDs from set
	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestAttachmentSetWalk(t *testing.T) {
	var (
		value = make(AttachmentSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*Attachment) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*Attachment) error { return fmt.Errorf("walk error") }))
}

func TestAttachmentSetFilter(t *testing.T) {
	var (
		value = make(AttachmentSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*Attachment) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*Attachment) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*Attachment) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestAttachmentSetIDs(t *testing.T) {
	var (
		value = make(AttachmentSet, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new(Attachment)
	value[1] = new(Attachment)
	value[2] = new(Attachment)
	// set ids
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	// Find existing
	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	// List IDs from set
	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestCredentialsSetWalk(t *testing.T) {
	var (
		value = make(CredentialsSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*Credentials) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*Credentials) error { return fmt.Errorf("walk error") }))
}

func TestCredentialsSetFilter(t *testing.T) {
	var (
		value = make(CredentialsSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*Credentials) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*Credentials) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*Credentials) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestCredentialsSetIDs(t *testing.T) {
	var (
		value = make(CredentialsSet, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new(Credentials)
	value[1] = new(Credentials)
	value[2] = new(Credentials)
	// set ids
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	// Find existing
	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	// List IDs from set
	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestReminderSetWalk(t *testing.T) {
	var (
		value = make(ReminderSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*Reminder) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*Reminder) error { return fmt.Errorf("walk error") }))
}

func TestReminderSetFilter(t *testing.T) {
	var (
		value = make(ReminderSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*Reminder) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*Reminder) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*Reminder) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestReminderSetIDs(t *testing.T) {
	var (
		value = make(ReminderSet, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new(Reminder)
	value[1] = new(Reminder)
	value[2] = new(Reminder)
	// set ids
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	// Find existing
	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	// List IDs from set
	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestRoleSetWalk(t *testing.T) {
	var (
		value = make(RoleSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*Role) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*Role) error { return fmt.Errorf("walk error") }))
}

func TestRoleSetFilter(t *testing.T) {
	var (
		value = make(RoleSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*Role) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*Role) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*Role) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestRoleSetIDs(t *testing.T) {
	var (
		value = make(RoleSet, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new(Role)
	value[1] = new(Role)
	value[2] = new(Role)
	// set ids
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	// Find existing
	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	// List IDs from set
	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}

func TestRoleMemberSetWalk(t *testing.T) {
	var (
		value = make(RoleMemberSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*RoleMember) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*RoleMember) error { return fmt.Errorf("walk error") }))
}

func TestRoleMemberSetFilter(t *testing.T) {
	var (
		value = make(RoleMemberSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*RoleMember) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*RoleMember) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*RoleMember) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestSettingValueSetWalk(t *testing.T) {
	var (
		value = make(SettingValueSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*SettingValue) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*SettingValue) error { return fmt.Errorf("walk error") }))
}

func TestSettingValueSetFilter(t *testing.T) {
	var (
		value = make(SettingValueSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*SettingValue) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*SettingValue) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*SettingValue) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestUserSetWalk(t *testing.T) {
	var (
		value = make(UserSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*User) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*User) error { return fmt.Errorf("walk error") }))
}

func TestUserSetFilter(t *testing.T) {
	var (
		value = make(UserSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*User) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*User) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*User) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestUserSetIDs(t *testing.T) {
	var (
		value = make(UserSet, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new(User)
	value[1] = new(User)
	value[2] = new(User)
	// set ids
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	// Find existing
	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	// List IDs from set
	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}
