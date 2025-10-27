package constants

import "errors"

var (
	ErrBookNotFound   = errors.New("книга не найдена")
	ErrBookOutOfStock = errors.New("книга закончилась")
	ErrReaderNotFound = errors.New("читатель не найден")
)
