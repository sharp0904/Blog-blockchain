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

func (k Keeper) Comments(goCtx context.Context, req *types.QueryCommentsRequest) (*types.QueryCommentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	//Define a variable that will store a list of posts
	var comments []*types.Comment

	//get context with the information about the environment
	ctx := sdk.UnwrapSDKContext(goCtx)

	//Get the key-value module store using the store key
	store := ctx.KVStore(k.storeKey)

	//Get the part of the store that keeps posts
	commentStore := prefix.NewStore(store, []byte(types.CommentKey))

	//Get the post by ID
	post, _ := k.GetPost(ctx, req.Id)

	//Get the postId
	postID := post.Id

	//paginate the posts store based on PageRequest
	pageRes, err := query.Paginate(commentStore, req.Pagination, func(key []byte, value []byte) error {
		var comment types.Comment
		if err := k.cdc.Unmarshal(value, &comment); err != nil {
			return err
		}

		if comment.PostID == postID {
			comments = append(comments, &comment)
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCommentsResponse{Post: &post, Comment: comments, Pagination: pageRes}, nil
}
