package files

import (
	"AdviserBot/lib/e"
	"AdviserBot/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const defaultPerm = 0777

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = e.WrapIfErr("can't save page", err)
	}()

	filePath := filepath.Join(s.basePath, page.UserName)
	if err = os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	if err = gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (p *storage.Page, err error) {
	defer func() {
		err = e.WrapIfErr("can't pick random page", err)
	}()

	path := filepath.Join(s.basePath, userName)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return nil, storage.ErrorNotSaved
	}
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrorNotSaved
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rand.Intn(len(files))
	file := files[index]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fileName, err := fileName(page)
	if err != nil {
		return e.Wrap("can't remove page", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)
	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file by path: %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fileName, err := fileName(page)
	if err != nil {
		return false, e.Wrap("can't check file if file exists", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check file if file %s exists", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't open page", err)
	}

	defer func() {
		_ = file.Close()
	}()

	var page storage.Page

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &page, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
