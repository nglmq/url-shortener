package storage

import "errors"

var (
	UrlNotFound      = errors.New("url not found")
	UrlAlreadyExists = errors.New("url already exists")
)
