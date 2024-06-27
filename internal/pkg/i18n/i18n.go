package i18n

import (
	"errors"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

var TranslationsI18n = map[string]*Translation{}

func NewTranslation(localize *goi18n.Localizer) *Translation {
	translation := &Translation{}
	translation.loadTranslations(localize)
	return translation
}

type Translation struct {
	ErrGeneric              error
	ErrInvalidId            error
	ErrInvalidData          error
	ErrManyRequest          error
	ErrorNonexistentRoute   error
	ErrUndefinedColumn      error
	ErrExpiredToken         error
	ErrDisabledUser         error
	ErrPassUnmatch          error
	ErrUserHasPass          error
	ErrUserHasNoPhoto       error
	ErrFileNotFound         error
	ErrIncorrectCredentials error

	ErrProfileUsed       error
	ErrProfileNotFound   error
	ErrProfileRegistered error

	ErrUserUsed       error
	ErrUserNotFound   error
	ErrUserRegistered error

	ErrDepartmentUsed       error
	ErrDepartmentNotFound   error
	ErrDepartmentRegistered error
}

func (s *Translation) loadTranslations(local *goi18n.Localizer) {
	s.ErrGeneric = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrGeneric"}, PluralCount: 1}))
	s.ErrInvalidId = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrInvalidId"}, PluralCount: 1}))
	s.ErrInvalidData = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrInvalidDatas"}, PluralCount: 1}))
	s.ErrManyRequest = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrManyRequest"}, PluralCount: 1}))
	s.ErrorNonexistentRoute = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrorNonexistentRoute"}, PluralCount: 1}))
	s.ErrUndefinedColumn = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrUndefinedColumn"}, PluralCount: 1}))
	s.ErrExpiredToken = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrExpiredToken"}, PluralCount: 1}))
	s.ErrDisabledUser = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrDisabledUser"}, PluralCount: 1}))
	s.ErrPassUnmatch = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrPassUnmatch"}, PluralCount: 1}))
	s.ErrUserHasPass = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrUserHasPass"}, PluralCount: 1}))
	s.ErrUserHasNoPhoto = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrUserHasNoPhoto"}, PluralCount: 1}))
	s.ErrFileNotFound = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrFileNotFound"}, PluralCount: 1}))
	s.ErrIncorrectCredentials = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrIncorrectCredentials"}, PluralCount: 1}))

	s.ErrProfileUsed = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrProfileUsed"}, PluralCount: 1}))
	s.ErrProfileNotFound = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrProfileNotFound"}, PluralCount: 1}))
	s.ErrProfileRegistered = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrProfileRegistered"}, PluralCount: 1}))

	s.ErrUserUsed = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrUserUsed"}, PluralCount: 1}))
	s.ErrUserNotFound = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrUserNotFound"}, PluralCount: 1}))
	s.ErrUserRegistered = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrUserRegistered"}, PluralCount: 1}))

	s.ErrDepartmentUsed = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrDepartmentUsed"}, PluralCount: 1}))
	s.ErrDepartmentNotFound = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrDepartmentNotFound"}, PluralCount: 1}))
	s.ErrDepartmentRegistered = errors.New(local.MustLocalize(&goi18n.LocalizeConfig{DefaultMessage: &goi18n.Message{ID: "ErrDepartmentRegistered"}, PluralCount: 1}))
}
