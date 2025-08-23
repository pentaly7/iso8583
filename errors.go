package iso8583

import (
	"errors"
)

var (
	ErrCreatingNewPackager = errors.New("error creating new packager")
)

var (
	ErrNotIsoMessage               = errors.New("not iso message")
	ErrInsufficientDataMti         = errors.New("insufficient data for mti")
	ErrInsufficientDataFirstBitmap = errors.New("insufficient data for first bitmap")
	ErrParsingFirstBitmap          = errors.New("cannot parse first bitmap")
	ErrParsingSecondBitmap         = errors.New("cannot parse second bitmap")
	ErrInsufficientDataBitmap      = errors.New("insufficient data for parsing bitmap")
	ErrFailedToParseBitmapData     = errors.New("failed parse bitmap data")
	ErrInvalidBitNumber            = errors.New("invalid bit number")
	ErrNoMtiToPack                 = errors.New("no mti to pack")
	ErrNotDefaultMti               = errors.New("not default mti to pack")
	ErrInvalidPackager             = errors.New("invalid packager value")
)

var (
	ErrInvalidBitType = errors.New("invalid bit type")
	ErrInvalidBitMap  = errors.New("invalid bitmap")
)
