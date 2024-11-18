package keeper

import (
	"context"

	"blog/x/blog/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Posts(goCtx context.Context, req *types.QueryPostsRequest) (*types.QueryPostsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	//Define a variable that will store a list of posts
	var posts []*types.Post

	//Get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(goCtx)

	//Get the key-value module store using the store key
	store := ctx.KVStore(k.storeKey)

	//Get the part of the store that keeps posts(using post key, which is "Post-value-")
	postStore := prefix.NewStore(store, []byte(types.PostKey))

	pageRes, err := query.Paginate(postStore, req.Pagination, func(key []byte, value []byte) error {
		var post types.Post
		if err := k.cdc.Unmarshal(value, &post); err != nil {
			return err
		}

		posts = append(posts, &post)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPostsResponse{Post: posts, Pagination: pageRes}, nil
}
