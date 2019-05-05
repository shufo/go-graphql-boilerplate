package resolver_test

import (
	"context"
	"database/sql"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/go-testfixtures/testfixtures"
	"github.com/machinebox/graphql"
	"github.com/volatiletech/sqlboiler/queries/qm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shufo/go-graphql-boilerplate/models"
	"github.com/shufo/go-graphql-boilerplate/testutils"
	"github.com/stretchr/testify/suite"
)

type UserResolverSuite struct {
	suite.Suite
	db       *sql.DB
	ts       *httptest.Server
	client   *graphql.Client
	fixtures *testfixtures.Context
}

func (suite *UserResolverSuite) SetupSuite() {
	suite.db = testutils.PrepareDB()
	m := testutils.PrepareRouter(suite.db)
	suite.ts = httptest.NewServer(m)
	suite.client = graphql.NewClient(suite.ts.URL + "/query")

	fixtures, err := testfixtures.NewFolder(suite.db, &testfixtures.MySQL{}, "../fixtures")
	if err != nil {
		log.Fatal(err)
	}
	suite.fixtures = fixtures
}

func (suite *UserResolverSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *UserResolverSuite) SetupTest() {
	if err := suite.fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func (suite *UserResolverSuite) TearDownTest() {
	// testutils.PopulateRecords(suite.db)
}

// Basic create user test
func (suite *UserResolverSuite) TestCreateUser() {

	// test cases
	cases := []struct {
		email       string
		password    string
		lastName    string
		firstName   string
		phoneNumber string
		valid       bool
		expected    string
	}{
		{
			email:       "test@example.com",
			password:    "123456",
			firstName:   "shuhei",
			lastName:    "hayashibara",
			phoneNumber: "03-1234-5678",
			valid:       true,
		},
		{
			email:       "test@example.com",
			password:    "1234",
			firstName:   "shuhei",
			lastName:    "hayashibara",
			phoneNumber: "03-1234-5678",
			valid:       false,
			expected:    "Password",
		},
		{
			email:       "test@example",
			password:    "123456",
			firstName:   "shuhei",
			lastName:    "hayashibara",
			phoneNumber: "03-1234-5678",
			valid:       false,
			expected:    "email format",
		},
	}

	for _, c := range cases {
		// valid request
		req := graphql.NewRequest(`
			mutation createUser {
				createUser(input: {
					email: "` + c.email + `", 
					password: "` + c.password + `"
					firstName: "` + c.firstName + `"
					lastName: "` + c.lastName + `"
					phoneNumber: "` + c.phoneNumber + `"
				}) {
					id
					token
				}
			}
		`)

		ctx := context.Background()
		var userResponse map[string]map[string]interface{}

		if c.valid {
			err := suite.client.Run(ctx, req, &userResponse)
			suite.NoError(err)
			var t float64
			suite.IsType(t, userResponse["createUser"]["id"])
			var s string
			suite.IsType(s, userResponse["createUser"]["token"])
		} else {
			err := suite.client.Run(ctx, req, &userResponse)
			suite.Error(err)
			suite.Contains(err.Error(), c.expected)
		}

	}
}

// Accept-Language Header test
func (suite *UserResolverSuite) TestLanguageHeader() {

	// test cases
	cases := []struct {
		email       string
		password    string
		lastName    string
		firstName   string
		phoneNumber string
		valid       bool
		expected    string
	}{
		{
			email:       "test@example.com",
			password:    "1234",
			firstName:   "shuhei",
			lastName:    "hayashibara",
			phoneNumber: "03-1234-5678",
			valid:       false,
			expected:    "パスワード",
		},
		{
			email:       "test@example",
			password:    "123456",
			firstName:   "shuhei",
			lastName:    "hayashibara",
			phoneNumber: "03-1234-5678",
			valid:       false,
			expected:    "メールアドレス",
		},
	}

	for _, c := range cases {
		// valid request
		req := graphql.NewRequest(`
			mutation createUser {
				createUser(input: {
					email: "` + c.email + `", 
					password: "` + c.password + `"
					firstName: "` + c.firstName + `"
					lastName: "` + c.lastName + `"
					phoneNumber: "` + c.phoneNumber + `"

				}) {
					id
					token
				}
			}
		`)
		req.Header.Add("Accept-Language", "ja")

		ctx := context.Background()
		var userResponse map[string]map[string]interface{}

		err := suite.client.Run(ctx, req, &userResponse)

		suite.Error(err)
		suite.Contains(err.Error(), c.expected)
	}
}

// Email Authentication test
func (suite *UserResolverSuite) TestEmailAuthentication() {
	// test cases
	cases := []struct {
		email    string
		password string
		valid    bool
		expected float64
	}{
		// use fixture user
		{email: "success@simulator.amazonses.com", password: "123456", valid: true, expected: 1},
		{email: "success@simulator.amazonses.com", password: "12345678", valid: false, expected: 0},
	}

	for _, c := range cases {
		// valid request
		req := graphql.NewRequest(`
			mutation createUser {
				authUser(input: {email: "` + c.email + `", password: "` + c.password + `"}) {
					id
					token
				}
			}
		`)

		ctx := context.Background()
		var userResponse map[string]map[string]interface{}

		if c.valid {
			res := suite.client.Run(ctx, req, &userResponse)
			suite.NoError(res)
			// assert the response is expected type/value
			var t float64
			suite.IsType(t, userResponse["authUser"]["id"])
			suite.Equal(c.expected, userResponse["authUser"]["id"])
			var s string
			suite.IsType(s, userResponse["authUser"]["token"])

			// JWT token test
			secret, found := os.LookupEnv("JWT_SECRET")

			if !found {
				log.Fatal("There is no JWT_SECRET variable in environment variables")
			}

			tokenAuth := jwtauth.New("HS256", []byte(secret), nil)
			token := userResponse["authUser"]["token"].(string)
			_, err := tokenAuth.Decode(token)
			suite.NoError(err)

			// reauthenticate with given token
			req.Header.Add("Authorization", "Bearer "+token)
			res = suite.client.Run(ctx, req, &userResponse)
			suite.NoError(res)

			exists, _ := models.AuthTokens(
				qm.Where("user_id = ?", 1),
				qm.Where("token = ?", token),
			).Exists(ctx, suite.db)

			// assert the previously used token has deleted
			suite.False(exists)

			// assert the token auth failed if secret is different
			tokenAuth = jwtauth.New("HS256", []byte("INVALID_SECRET"), nil)
			_, err = tokenAuth.Decode(userResponse["authUser"]["token"].(string))

			suite.Error(err)
			suite.EqualError(err, "signature is invalid")

		} else {
			// assert invalid request returns error
			err := suite.client.Run(ctx, req, &userResponse)
			suite.EqualError(err, "graphql: Email or Password is incorrect")
		}
	}
}

func TestUserResolverSuite(t *testing.T) {
	suite.Run(t, new(UserResolverSuite))
}
