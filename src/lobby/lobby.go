package lobby

import (
	context "context"
	"math/rand"

	common "github.com/tanerius/dungeonforge/src/common"
)

type LobbyServerNode struct {
	UnimplementedLobbyServiceServer
}

func (r *LobbyServerNode) GetUsers(ctx context.Context, token *common.TokenRequest) (*UserListResponse, error) {
	br := &common.BaseResponse{
		ResponseCode: 0,
		ApiVersion:   1,
	}

	var users []*UserEntry = make([]*UserEntry, 0)

	return &UserListResponse{Response: br, Users: users}, nil
}

func (r *LobbyServerNode) JoinLobby(context.Context, *common.TokenRequest) (*JoinResponse, error) {
	return nil, nil
}

func (r *LobbyServerNode) LeaveLobby(context.Context, *common.TokenRequest) (*LeaveResponse, error) {
	return nil, nil
}

func (r *LobbyServerNode) Roll(ctx context.Context, request *RollRequest) (*RollResponse, error) {
	br := &common.BaseResponse{
		ResponseCode: 200,
		ApiVersion:   1,
	}

	a := make([]int64, request.GetDice())

	for j := 0; j < len(a); j++ {
		a[j] = rand.Int63n(request.GetSides()) + 1
	}

	return &RollResponse{Response: br, Value: a}, nil
}
