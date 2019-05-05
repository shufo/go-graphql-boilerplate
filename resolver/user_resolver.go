package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/shufo/go-graphql-boilerplate/graph/generated"
	"github.com/shufo/go-graphql-boilerplate/utils"

	"github.com/shufo/go-graphql-boilerplate/translations"

	"github.com/shufo/go-graphql-boilerplate/configs"

	"golang.org/x/crypto/bcrypt"

	"github.com/volatiletech/null"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/99designs/gqlgen/graphql"

	"github.com/go-chi/jwtauth"

	"github.com/dgrijalva/jwt-go"
	"github.com/shufo/go-graphql-boilerplate/models"
	"github.com/volatiletech/sqlboiler/boil"
)

var tokenAuth *jwtauth.JWTAuth

type userResolver struct{ *Resolver }

func (r *Resolver) User() generated.UserResolver {
	return &userResolver{r}
}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.AuthenticatedUser, error) {
	db := ctx.Value("db").(*sql.DB)

	if err := input.Validate(ctx); err != nil {
		for k, v := range err {
			graphql.AddErrorf(ctx, "%s: %s", k, v)
		}
		return nil, nil
	}

	// check if user is already exists
	if exists, _ := models.AuthenticationProviders(
		qm.Where("provider_type = ?", "email"),
		qm.Where("provider_username = ?", input.Email),
	).Exists(ctx, db); exists {
		return nil, fmt.Errorf(translations.T(ctx, "email_already_exists"))
	}

	// create User record
	u := &models.User{Username: null.String{String: input.Email, Valid: true}}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	// create auth provider record
	ap := models.AuthenticationProvider{
		ProviderType:     "email",
		ProviderUsername: input.Email,
		ProviderPassword: string(hashed),
		Email:            null.StringFrom(input.Email),
		FirstName:        null.StringFrom(input.FirstName),
		LastName:         null.StringFrom(input.LastName),
	}

	if err := ap.SetUser(ctx, db, false, u); err != nil {
		return nil, err
	}

	if err := ap.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	// create a profile for user
	pr := models.Profile{
		FirstName:   null.StringFrom(input.FirstName),
		LastName:    null.StringFrom(input.LastName),
		PhoneNumber: null.StringFrom(input.PhoneNumber),
	}

	if err := pr.SetUser(ctx, db, false, u); err != nil {
		return nil, err
	}

	if err := pr.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	// set roles
	role, _ := models.Roles(qm.Where("type = ?", "USER")).One(ctx, db)

	ur := &models.UserRole{}
	ur.SetUser(ctx, db, false, u)
	ur.SetRole(ctx, db, false, role)

	if err := ur.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	// create token
	token, err := createToken(u)

	if err != nil {
		return nil, err
	}

	// add token to auth token table for revocation
	at := &models.AuthToken{
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(configs.TokenLifetime),
	}

	if err := at.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}

	res := &models.AuthenticatedUser{
		ID:    u.ID,
		Token: token,
	}

	return res, nil
}

func (r *mutationResolver) AuthUser(ctx context.Context, input models.AuthUserInput) (*models.AuthenticatedUser, error) {
	db := ctx.Value("db").(*sql.DB)

	// validate inputs
	if err := input.Validate(ctx); err != nil {
		for k, v := range err {
			graphql.AddErrorf(ctx, "%s: %s", k, v)
		}
		return nil, nil
	}

	// search user if it exists
	ap, err := models.AuthenticationProviders(
		qm.Where("provider_type = ?", "email"),
		qm.Where("provider_username = ?", input.Email),
		qm.Load("User"),
	).One(ctx, db)

	if err != nil {
		return nil, fmt.Errorf(translations.T(ctx, "email_or_password_is_incorrect"))
	}

	// compare hashed password with inputed password
	if err := bcrypt.CompareHashAndPassword([]byte(ap.ProviderPassword), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf(translations.T(ctx, "email_or_password_is_incorrect"))
	}

	// swipe already exists auth token if token exists in header
	if ht, ok := ctx.Value(jwtauth.TokenCtxKey).(jwt.Token); ok {
		if _, err := models.AuthTokens(
			qm.Where("user_id = ?", ap.UserID),
			qm.Where("token = ?", ht.Raw),
		).DeleteAll(ctx, db); err != nil {
			return nil, err
		}
	}

	// create new token
	token, err := createToken(ap.R.User)

	if err != nil {
		return nil, err
	}

	res := &models.AuthenticatedUser{
		ID:    ap.UserID,
		Token: token,
	}

	return res, nil
}

func createToken(u *models.User) (string, error) {
	// initialize jwt
	secret := os.Getenv("JWT_SECRET")
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)

	// set user claims

	claims["user_id"] = &u.ID
	claims["roles"] = []models.RoleType{models.RoleTypeUser}
	claims["uuid"] = utils.RandomUUID()

	jaclaims := jwtauth.Claims(claims)
	jaclaims.SetIssuedNow()

	// expires after 24 hours
	jaclaims.SetExpiryIn(24 * time.Hour)

	_, tokenString, err := tokenAuth.Encode(jaclaims)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (r *queryResolver) User(ctx context.Context, userID *int) (*models.User, error) {

	if userID == nil {
		_, claims, _ := jwtauth.FromContext(ctx)

		if claims["user_id"] == nil {
			return nil, fmt.Errorf("Invalid token")
		}

		res := &models.User{
			ID: int(claims["user_id"].(float64)),
		}

		return res, nil
	}

	db := ctx.Value("db").(*sql.DB)
	u, err := models.FindUser(ctx, db, *userID)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userResolver) AuthenticationProviders(ctx context.Context, u *models.User) ([]models.AuthenticationProvider, error) {
	db := ctx.Value("db").(*sql.DB)

	aps, err := u.AuthenticationProviders().All(ctx, db)

	if err != nil {
		return nil, err
	}

	res := make([]models.AuthenticationProvider, len(aps))
	for i, v := range aps {
		res[i] = *v
	}

	return res, nil
}
