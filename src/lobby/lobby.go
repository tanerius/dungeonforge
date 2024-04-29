package lobby

import (
	context "context"
	"math/rand"
)

type LobbyServerNode struct {
	UnimplementedLobbyServer
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

func (r *LobbyServerNode) Login(context.Context, *UserLoginRequest) (*LoginResponse, error) {
	return nil, nil
}

func (r *LobbyServerNode) Matchmaking(context.Context, *TokenRequest) (*MatchmakingResponse, error) {
	return nil, nil
}

func (r *LobbyServerNode) Roll(ctx context.Context, request *RollRequest) (*RollResponse, error) {
	br := &BaseResponse{
		ResponseCode: 0,
		ApiVersion:   1,
	}

	a := make([]int64, request.GetDice())

	for j := 0; j < len(a); j++ {
		a[j] = rand.Int63n(request.GetSides()) + 1
	}

	return &RollResponse{Response: br, Value: a}, nil
}
