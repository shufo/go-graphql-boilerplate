package resolver

import (
	"context"
	"fmt"

	"github.com/shufo/go-graphql-boilerplate/graph/generated"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shufo/go-graphql-boilerplate/models"
	"github.com/go-chi/jwtauth"
)

func NewDirectives() generated.DirectiveRoot {
	return generated.DirectiveRoot{
		HasRole:         HasRole,
		HasMinimumRole:  HasMinimumRole,
		IsResourceOwner: IsResourceOwner,
		Length:          Length,
	}
}

type Roles struct {
	ADMIN int
	USER  int
}

type RoleMap map[models.RoleType]int

var roles = RoleMap{
	models.RoleTypeSuperAdmin:         50,
	models.RoleTypeOrganizationAdmin:  40,
	models.RoleTypeOrganizationMember: 30,
	models.RoleTypeResourceOwner:      20,
	models.RoleTypeUser:               10,
}

func HasMinimumRole(ctx context.Context, obj interface{}, next graphql.Resolver, role models.RoleType) (interface{}, error) {
	_, claims, _ := jwtauth.FromContext(ctx)

	// Check if user has permission with required role
	if roles[models.RoleType(claims["role"].(string))] < roles[role] {
		return nil, fmt.Errorf("You are not granted to access this resource with your role")
	}

	return next(ctx)
}

func isAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	if _, _, err := jwtauth.FromContext(ctx); err != nil {
		return nil, err
	}

	return next(ctx)
}

func HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role models.RoleType) (interface{}, error) {
	switch role {

	case models.RoleTypeOrganizationAdmin:
		_, claims, err := jwtauth.FromContext(ctx)
		if err != nil {
			return nil, err
		}

		if claims["role"] != nil && claims["role"] != models.RoleTypeOrganizationAdmin.String() {
			return nil, fmt.Errorf("You are not granted access this resource")
		}

	case models.RoleTypeResourceOwner:
		ownable, isOwnable := obj.(models.Ownable)

		if !isOwnable {
			return nil, fmt.Errorf("This object can't be owned")
		}

		_, claims, err := jwtauth.FromContext(ctx)

		if err != nil {
			return nil, err
		}

		if claims["user_id"] == nil {
			return nil, fmt.Errorf("Invalid token")
		}

		if *ownable.OwnerID() != int(claims["user_id"].(float64)) {
			return nil, fmt.Errorf("You are not own this resource")
		}

	}

	// or let it pass through
	return next(ctx)
}

func Length(ctx context.Context, input interface{}, next graphql.Resolver, min *int, max *int) (interface{}, error) {

	i, err := input.(int64)

	if !err {
		return next(ctx)
	}

	if i < int64(*min) || i > int64(*max) {
		return nil, fmt.Errorf("Arguments length must satisfied with min length: %d and max length: %d", *min, *max)
	}
	return next(ctx)
}

func IsResourceOwner(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// super admin
	_, claims, err := jwtauth.FromContext(ctx)

	if err != nil {
		return nil, err
	}

	if claims["roles"] == nil {
		return nil, fmt.Errorf("no roles found")
	}

	for _, v := range claims["roles"].([]interface{}) {
		if v.(string) == "SUPER_ADMIN" {
			return next(ctx)
		}
	}

	// space owner
	// check by casbin
	for _, v := range claims["roles"].([]interface{}) {
		if v.(string) == "ORGANIZATION_ADMIN" || v.(string) == "ORGANIZATION_MEMBER" {
			// TODO implement casbin authz
			return next(ctx)
		}
	}

	// resource subject

	ownable, isOwnable := obj.(models.Ownable)

	if !isOwnable {
		return nil, fmt.Errorf("This object can't be owned")
	}

	if err != nil {
		return nil, err
	}

	if claims["user_id"] == nil {
		return nil, fmt.Errorf("Invalid token")
	}

	if *ownable.OwnerID() != int(claims["user_id"].(float64)) {
		return nil, fmt.Errorf("You are not own this resource")
	}

	return next(ctx)
}
