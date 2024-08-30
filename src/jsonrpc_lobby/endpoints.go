package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tanerius/dungeonforge/pkg/tokenmanager"
	"github.com/tanerius/dungeonforge/src/entities"
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
func (r *JsonRpcService) Login(ctx context.Context, params json.RawMessage) (*tokenmanager.Tokens, error) {
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

func (r *JsonRpcService) RefreshToken(ctx context.Context, params json.RawMessage) (*tokenmanager.Tokens, error) {
	return nil, nil
}
