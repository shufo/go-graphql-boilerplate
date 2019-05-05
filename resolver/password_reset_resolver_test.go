package resolver_test

import (
	"context"
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/shufo/go-graphql-boilerplate/testutils"
	"github.com/go-testfixtures/testfixtures"
	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/suite"
)

type PasswordResetResolverSuite struct {
	suite.Suite
	db       *sql.DB
	ts       *httptest.Server
	client   *graphql.Client
	fixtures *testfixtures.Context
}

func (suite *PasswordResetResolverSuite) SetupSuite() {
	suite.db = testutils.PrepareDB()
	m := testutils.PrepareRouter(suite.db)
	suite.ts = httptest.NewServer(m)
	suite.client = graphql.NewClient(suite.ts.URL + "/query")

	fixtures, err := testfixtures.NewFolder(suite.db, &testfixtures.MySQL{}, "../fixtures")
	if err != nil {
		log.Fatal(err)
	}
	suite.fixtures = fixtures

	os.Setenv("APP_ENV", "test")
}

func (suite *PasswordResetResolverSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *PasswordResetResolverSuite) SetupTest() {
	if err := suite.fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func (suite *PasswordResetResolverSuite) TearDownTest() {
	// testutils.PopulateRecords(suite.db)
}

// Basic create user test
func (s *PasswordResetResolverSuite) Test_mutationResolver_RequestPasswordReset() {
	// test cases
	cases := []struct {
		name    string
		email   string
		valid   bool
		want    string
		wantErr string
	}{
		{
			name:  "case: valid request",
			email: "success@simulator.amazonses.com",
			valid: true,
			want:  "sent",
		},
		{
			name:    "case: email address not found",
			email:   "test@example.com",
			valid:   false,
			wantErr: "not found",
		},
	}

	for _, c := range cases {
		req := graphql.NewRequest(`
			mutation requestPasswordReset($email: String!) {
				requestPasswordReset(input: {email: $email}) {
					status
				}
			}
		`)
		req.Var("email", c.email)

		var response map[string]map[string]interface{}

		ctx := context.Background()

		if c.valid {
			err := s.client.Run(ctx, req, &response)
			s.NoError(err, c.name)
			s.Equal(c.want, response["requestPasswordReset"]["status"], c.name)
		} else {
			err := s.client.Run(ctx, req, &response)
			s.Error(err, c.name)
			s.Contains(err.Error(), c.wantErr, c.name)
		}

	}
}

func (s *PasswordResetResolverSuite) Test_mutationResolver_ValidatePasswordReset() {
	// test cases
	cases := []struct {
		name    string
		token   string
		valid   bool
		want    string
		wantErr string
	}{
		{
			name:  "case: valid request",
			token: "valid_password_reset_token",
			valid: true,
			want:  "verified",
		},
		{
			name:    "case: token is invalid",
			token:   "invalid_token",
			valid:   false,
			wantErr: "reset token is invalid",
		},
	}

	for _, c := range cases {
		req := graphql.NewRequest(`
			mutation validatePasswordReset($token: String!) {
				validatePasswordReset(input: {token: $token}) {
					status
				}
			}
		`)
		req.Var("token", c.token)

		var res map[string]map[string]interface{}

		ctx := context.Background()

		if c.valid {
			err := s.client.Run(ctx, req, &res)
			s.NoError(err, c.name)
			s.Equal(c.want, res["validatePasswordReset"]["status"], c.name)
		} else {
			err := s.client.Run(ctx, req, &res)
			s.Error(err, c.name)
			s.Contains(err.Error(), c.wantErr, c.name)
		}

	}
}

func (s *PasswordResetResolverSuite) Test_mutationResolver_CompletePasswordReset() {
	// test cases
	cases := []struct {
		name     string
		token    string
		password string
		valid    bool
		want     string
		wantErr  string
	}{
		{
			name:     "case: valid request",
			token:    "valid_password_reset_token",
			password: "1234abcd",
			valid:    true,
			want:     "success@simulator.amazonses.com",
		},
		{
			name:     "case: password is too short",
			token:    "valid_password_reset_token",
			password: "1234",
			valid:    false,
			wantErr:  "characters",
		},
		{
			name:     "case: verified token is not found",
			token:    "not_verified_reset_token",
			password: "1234abcd",
			valid:    false,
			wantErr:  "no verified password reset token",
		},
	}

	for _, c := range cases {
		req := graphql.NewRequest(`
			mutation completePasswordReset($token: String!, $newPassword: String!) {
				completePasswordReset(input: {token: $token, newPassword: $newPassword}) {
					email
				}
			}
		`)
		req.Var("token", c.token)
		req.Var("newPassword", c.password)

		var res map[string]map[string]interface{}

		ctx := context.Background()

		if c.valid {
			err := s.client.Run(ctx, req, &res)
			s.NoError(err, c.name)
			s.Equal(c.want, res["completePasswordReset"]["email"], c.name)
		} else {
			err := s.client.Run(ctx, req, &res)
			s.Error(err, c.name)
			s.Contains(err.Error(), c.wantErr, c.name)
		}

	}
}

func TestPasswordResetResolverSuite(t *testing.T) {
	suite.Run(t, new(PasswordResetResolverSuite))
}
