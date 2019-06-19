package services

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"

	"git.ve.home/nicolasc/linotte/libs/helpers"
	"git.ve.home/nicolasc/linotte/services/user/configuration"
	"git.ve.home/nicolasc/linotte/services/user/models"
	proto_user "git.ve.home/nicolasc/linotte/services/user/proto"
)

type keys struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

// UserService give an access to every possible user operations
type UserService struct {
	accessor *models.UserAccessor
	keys     *keys
}

// NewUserService creates a new user service from the given accessor
func NewUserService(accessor *models.UserAccessor) *UserService {
	settings := config.Get()
	return &UserService{
		accessor: accessor,
		keys: &keys{
			private: helpers.ParsePrivateKey(settings.PrivateKeyPath),
			public:  helpers.ParsePublicKey(settings.PublicKeyPath),
		},
	}
}

// Login authenticated the given user and returns true if succeed
func (service *UserService) Login(req *proto_user.LoginRequest) (*proto_user.UserReply, error) {
	if len(req.Login) > 0 && len(req.Password) > 0 {
		user := service.accessor.GetByUsername(req.Login)

		if user.ID > 0 {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
				token, err := helpers.GenerateToken(service.keys.private)
				return &proto_user.UserReply{
					Success: true,
					Token:   token,
				}, err
			}
		}
	} else {
		return nil, errors.New("Invalid user")
	}

	return &proto_user.UserReply{
		Success: false,
		Token:   "",
	}, errors.New("Unable to log in")
}

// Logout helps to logout the given user and return true if succeed
func (service *UserService) Logout(req *proto_user.LogoutRequest) (*proto_user.UserReply, error) {
	fmt.Println("Trying to logout. Not implemented.")
	return nil, errors.New("Logout not implemented")
}

// CheckAuthentication verifies if the current token is valid and returns true if it is
func (service *UserService) CheckAuthentication(req *proto_user.CheckRequest) (*proto_user.CheckReply, error) {
	strs := strings.Split(string(req.Token), "Bearer ")

	if len(strs) >= 2 {
		token, err := jwt.Parse(
			strs[1],
			func(token *jwt.Token) (interface{}, error) {
				return service.keys.public, nil
			},
		)

		return &proto_user.CheckReply{
			Success: err == nil && token.Valid,
		}, err
	}

	return &proto_user.CheckReply{
		Success: false,
	}, errors.New("No token")
}
