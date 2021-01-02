package keychain

import (
	"fmt"
	"os/user"
	"runtime"

	"dfl/lib/cher"

	"github.com/keybase/go-keychain"
)

const (
	service     = "mn.dfl.auth"
	accessGroup = "mn.dfl.auth"
	prefix      = "DFL"
)

func Supported() bool {
	return runtime.GOOS != "windows"
}

func NewItem(name string, data []byte) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetLabel(fmt.Sprintf("%s %s", prefix, name))
	item.SetAccessGroup(accessGroup)
	item.SetAccount(user.Username)
	item.SetData(data)
	item.SetSynchronizable(keychain.SynchronizableAny)
	item.SetAccessible(keychain.AccessibleWhenUnlockedThisDeviceOnly)

	if err := keychain.AddItem(item); err != nil {
		return err
	}

	return nil
}

func UpsertItem(name string, data []byte) error {
	item, err := GetItem(name)
	if err != nil {
		// if unknown error, or not_found
		if v, ok := err.(cher.E); !ok || v.Code != cher.NotFound {
			return err
		}
	}

	if item != nil {
		if err := DeleteItem(name); err != nil {
			return err
		}
	}

	return NewItem(name, data)
}

func GetItem(name string) (data []byte, err error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetAccount(user.Username)
	query.SetLabel(fmt.Sprintf("%s %s", prefix, name))
	query.SetAccessGroup(accessGroup)
	query.SetReturnData(true)
	query.SetReturnAttributes(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		fmt.Println("test")
		if err == keychain.ErrorItemNotFound {
			return nil, cher.New(cher.NotFound, nil)
		}

		return nil, err
	}

	if len(results) != 1 {
		return nil, cher.New(cher.NotFound, nil)
	}

	return results[0].Data, nil
}

func DeleteItem(name string) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(user.Username)
	item.SetLabel(fmt.Sprintf("%s %s", prefix, name))

	return keychain.DeleteItem(item)
}
