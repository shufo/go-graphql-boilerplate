package models

import (
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/volatiletech/null"
)

// MarshalNullableString marshals NullableString Type
func MarshalNullableString(ns null.String) graphql.Marshaler {
	if !ns.Valid {
		// this is also important, so we can detect if this scalar is used in a not null context and return an appropriate error
		return graphql.Null
	}
	return graphql.MarshalString(ns.String)
}

// UnmarshalNullableString unmarshals NullableString Type
func UnmarshalNullableString(v interface{}) (null.String, error) {
	if v == nil {
		return null.String{Valid: false}, nil
	}
	// again you can delegate to the default implementation to save yourself some work.
	s, err := graphql.UnmarshalString(v)
	return null.String{String: s}, err
}

// MarshalNullableString marshals NullableString Type
func MarshalNullableTime(nt null.Time) graphql.Marshaler {
	if !nt.Valid {
		// this is also important, so we can detect if this scalar is used in a not null context and return an appropriate error
		return graphql.Null
	}
	return graphql.MarshalString(nt.Time.String())
}

// UnmarshalNullableString unmarshals NullableString Type
func UnmarshalNullableTime(v interface{}) (null.Time, error) {
	if v == nil {
		return null.Time{Valid: false}, nil
	}
	// again you can delegate to the default implementation to save yourself some work.
	s, err := graphql.UnmarshalString(v)
	t, _ := time.Parse("2019-12-01 00:00:00", s)
	return null.Time{Time: t}, err
}

type Ownable interface {
	OwnerID() *int
}

func (u User) OwnerID() *int {
	return &u.ID
}
