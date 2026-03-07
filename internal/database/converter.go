// Package database contains interactions for data storage
package database

import dto "github.com/bitztec/chirpy/internal/dataTransfer"

func (c *Chirp) ToDTO() dto.DTOChirp {
	return dto.DTOChirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
}

func (u *User) ToDTO() dto.DTOUser {
	return dto.DTOUser{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}
