package lobby

import (
	context "context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/tanerius/dungeonforge/pkg/config"
	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/logging"
	"github.com/tanerius/dungeonforge/pkg/tokenmanager"
	"github.com/tanerius/dungeonforge/src/entities"
)

func buildLoginResponse(tokens *tokenmanager.Tokens, user *entities.User) *LoginResponse {

	return &LoginResponse{
		Response: &BaseResponse{
			ResponseCode: 200,
			ApiVersion:   1,
		},
		Tokens: &TokenResponse{
			AccessToken:           tokens.Access,
			AccessTokenExpiresIn:  tokens.AccessExires,
			RefreshToken:          &tokens.Refresh,
			RefreshTokenExpiresIn: &tokens.RefreshExpires,
		},
		User: &UserEntry{
			Id:          user.GetId(),
			Displayname: user.DisplayName,
			Email:       &user.Email,
			IsOnline:    &user.IsOnline,
			IsValidated: &user.Validated,
		},
	}
}

type LobbyServerNode struct {
	UnimplementedLobbyServer
	conf   config.IConfig
	id     string
	db     *database.MongoDB
	logger logging.ILogger
}

func NewMockedLobby(logger logging.ILogger) *LobbyServerNode {
	l := &LobbyServerNode{}
	l.logger = logger
	logger.LogInfo("Setting up IConfig")
	l.conf = config.NewIConfig(true)
	logger.LogInfo("Setting up ID")
	l.id = uuid.NewString()
	root_user, _ := l.conf.ReadKeyString("root_user")
	root_pass, _ := l.conf.ReadKeyString("root_password")
	hostname, _ := l.conf.ReadKeyString("host")
	hostname = fmt.Sprintf("mongodb://%s:%s@%s:27017/", root_user, root_pass, hostname)
	logger.LogInfo("Setting up DB " + hostname)
	if db, err := database.NewMongoDBWrapper(context.Background(), hostname, 100); err == nil {
		l.db = db
		logger.LogInfo("Lobby ready!")
		return l
	}

	return nil
}

func NewLobby(logger logging.ILogger) *LobbyServerNode {
	l := &LobbyServerNode{}
	l.conf = config.NewIConfig(false)
	l.id = uuid.NewString()
	root_user, _ := l.conf.ReadKeyString("root_user")
	root_pass, _ := l.conf.ReadKeyString("root_password")
	hostname, _ := l.conf.ReadKeyString("host")
	hostname = fmt.Sprintf("mongodb://%s:%s@%s:27017/", root_user, root_pass, hostname)
	logger.LogInfo("Setting up DB " + hostname)
	if db, err := database.NewMongoDBWrapper(context.Background(), hostname, 100); err == nil {
		l.db = db
		return l
	}

	return nil
}

func (r *LobbyServerNode) GetUsers(ctx context.Context, token *TokenRequest) (*UserListResponse, error) {
	br := &BaseResponse{
		ResponseCode: 0,
		ApiVersion:   1,
	}

	var users []*UserEntry = make([]*UserEntry, 0)

	return &UserListResponse{Response: br, Users: users}, nil
}

func (r *LobbyServerNode) JoinLobby(context.Context, *TokenRequest) (*JoinResponse, error) {
	return nil, nil
}

func (r *LobbyServerNode) LeaveLobby(context.Context, *TokenRequest) (*LeaveResponse, error) {
	return nil, nil
}

func (r *LobbyServerNode) Login(ctx context.Context, request *UserLoginRequest) (*LoginResponse, error) {
	res := &LoginResponse{
		Response: &BaseResponse{
			ResponseCode: 404,
			ApiVersion:   1,
		},
	}

	if usr, err := entities.GetUser(ctx, r.db, request.GetUsername(), request.GetPassword()); err == nil {
		res.Response.ResponseCode = 200
		tokens, errToken := tokenmanager.GenerateTokens(usr.Username)

		if errToken != nil {
			s := errToken.Error()
			return &LoginResponse{
				Response: &BaseResponse{
					ResponseCode: 409,
					ApiVersion:   1,
					ErrorMsg:     &s,
				},
			}, errToken
		}

		return buildLoginResponse(tokens, usr), nil
	}

	return res, nil
}

func (r *LobbyServerNode) RefreshToken(context.Context, *TokenRequest) (*LoginResponse, error) {
	return nil, nil
}

func (r *LobbyServerNode) Register(_ context.Context, newUser *NewUserMessage) (*LoginResponse, error) {
	ret := &LoginResponse{
		Response: &BaseResponse{
			ResponseCode: 400,
			ApiVersion:   1,
		},
	}

	usr, errUser := entities.RegisterUser(context.Background(), r.db, newUser.GetEmail(), newUser.GetUsername(), newUser.GetPassword())

	if errUser != nil {
		s := errUser.Error()
		ret.Response.ResponseCode = 409
		ret.Response.ErrorMsg = &s
		return ret, nil
	}

	tokens, errToken := tokenmanager.GenerateTokens(usr.Username)

	if errToken != nil {
		s := errToken.Error()
		ret.Response.ErrorMsg = &s
		return ret, errToken
	}

	return buildLoginResponse(tokens, usr), nil
}

func (r *LobbyServerNode) Roll(ctx context.Context, request *RollRequest) (*RollResponse, error) {
	br := &BaseResponse{
		ResponseCode: 200,
		ApiVersion:   1,
	}

	a := make([]int64, request.GetDice())

	for j := 0; j < len(a); j++ {
		a[j] = rand.Int63n(request.GetSides()) + 1
	}

	return &RollResponse{Response: br, Value: a}, nil
}
