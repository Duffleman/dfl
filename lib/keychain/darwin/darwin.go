// +build darwin

package darwin

import (
	"fmt"
	"os/user"

	"github.com/cuvva/cuvva-public-go/lib/cher"
	gokeychain "github.com/keybase/go-keychain"
)

const (
	service     = "mn.dfl.auth"
	accessGroup = "mn.dfl.auth"
	prefix      = "DFL"
)

type Keychain struct{}

func (k Keychain) NewItem(name string, data []byte) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	item := gokeychain.NewItem()
	item.SetSecClass(gokeychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetLabel(fmt.Sprintf("%s %s", prefix, name))
	item.SetAccessGroup(accessGroup)
	item.SetAccount(user.Username)
	item.SetData(data)
	item.SetSynchronizable(gokeychain.SynchronizableAny)
	item.SetAccessible(gokeychain.AccessibleWhenUnlockedThisDeviceOnly)

	if err := gokeychain.AddItem(item); err != nil {
		return err
	}

	return nil
}

func (k Keychain) UpsertItem(name string, data []byte) error {
	item, err := k.GetItem(name)
	if err != nil {
		// if unknown error, or not_found
		if v, ok := err.(cher.E); !ok || v.Code != cher.NotFound {
			return err
		}
	}

	if item != nil {
		if err := k.DeleteItem(name); err != nil {
			return err
		}
	}

	return k.NewItem(name, data)
}

func (k Keychain) GetItem(name string) (data []byte, err error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	query := gokeychain.NewItem()
	query.SetSecClass(gokeychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetAccount(user.Username)
	query.SetLabel(fmt.Sprintf("%s %s", prefix, name))
	query.SetAccessGroup(accessGroup)
	query.SetReturnData(true)
	query.SetReturnAttributes(true)

	results, err := gokeychain.QueryItem(query)
	if err != nil {
		fmt.Println("test")
		if err == gokeychain.ErrorItemNotFound {
			return nil, cher.New(cher.NotFound, nil)
		}

		return nil, err
	}

	if len(results) != 1 {
		return nil, cher.New(cher.NotFound, nil)
	}

	return results[0].Data, nil
}

func (k Keychain) DeleteItem(name string) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	item := gokeychain.NewItem()
	item.SetSecClass(gokeychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(user.Username)
	item.SetLabel(fmt.Sprintf("%s %s", prefix, name))

	return gokeychain.DeleteItem(item)
}
