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
		if err == keychain.ErrorDuplicateItem {
			return nil
		}

		return err
	}

	return nil
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

	results, err := keychain.QueryItem(query)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, cher.New(cher.NotFound, nil)
	}

	return results[0].Data, nil
}
