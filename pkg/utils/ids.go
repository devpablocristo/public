package pkgutils

import "errors"

func ValidateID(id int64) error {
	if id <= 0 {
		return errors.New("id must be greater than 0")
	}
	return nil
}
