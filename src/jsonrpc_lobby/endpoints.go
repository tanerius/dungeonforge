package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tanerius/dungeonforge/pkg/tokenmanager"
	"github.com/tanerius/dungeonforge/src/entities"
	"golang.org/x/exp/rand"
)

/*
service Lobby {
    rpc Login(UserLoginRequest) returns (LoginResponse);
    rpc RefreshToken(TokenRequest) returns (LoginResponse);
    rpc JoinLobby(TokenRequest) returns (JoinResponse);
    rpc LeaveLobby(TokenRequest) returns (LeaveResponse);
    rpc GetUsers(TokenRequest) returns (UserListResponse);
    rpc Register(NewUserMessage) returns (LoginResponse);
    rpc Roll(RollRequest) returns (RollResponse);
}
*/

// Example function that the RPC server will expose.
func (r *JsonRpcService) login(ctx context.Context, params json.RawMessage) (*tokenmanager.Tokens, error) {
	var credentials []string
	if err := json.Unmarshal(params, &credentials); err != nil {
		return nil, fmt.Errorf("invalid params")
	}
	if len(credentials) != 2 {
		return nil, fmt.Errorf("exactly two strings are required")
	}

	if usr, err := entities.GetUser(ctx, r.db, credentials[0], credentials[1]); err == nil {
		tokens, errToken := tokenmanager.GenerateTokens(usr.Username)

		if errToken != nil {
			return nil, errToken
		}

		return tokens, nil
	} else {
		return nil, err
	}
}

func (r *JsonRpcService) getUsers(ctx context.Context, params json.RawMessage) (interface{}, error) {
	return nil, nil
}

func (r *JsonRpcService) refreshToken(ctx context.Context, params json.RawMessage) (*tokenmanager.Tokens, error) {
	return nil, nil
}

func (r *JsonRpcService) register(_ context.Context, params json.RawMessage) (*tokenmanager.Tokens, error) {
	var credentials []string
	if err := json.Unmarshal(params, &credentials); err != nil {
		return nil, fmt.Errorf("invalid params")
	}
	if len(credentials) != 3 {
		return nil, fmt.Errorf("exactly three strings are required (email, username, password)")
	}

	usr, errUser := entities.RegisterUser(context.Background(), r.db, credentials[0], credentials[1], credentials[2])

	if errUser != nil {
		return nil, errUser
	}

	tokens, errToken := tokenmanager.GenerateTokens(usr.Username)

	if errToken != nil {
		return nil, errToken
	}

	return tokens, nil
}

func (r *JsonRpcService) roll(_ context.Context, params json.RawMessage) ([]int64, error) {
	var diceInfo []int64
	if err := json.Unmarshal(params, &diceInfo); err != nil {
		return nil, fmt.Errorf("invalid params")
	}
	if len(diceInfo) != 3 {
		return nil, fmt.Errorf("exactly two unsigned integers are required (number of dice, sides per dice)")
	}

	a := make([]int64, diceInfo[0])

	for j := 0; j < len(a); j++ {
		a[j] = rand.Int63n(diceInfo[1]) + 1
	}

	return a, nil
}
