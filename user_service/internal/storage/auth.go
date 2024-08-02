package storage

import (
	"auth/genproto"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	coll *mongo.Collection
}

type Profile struct {
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
	Address   string `bson:"address"`
}

// UserReq struct
type UserReq struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	Profile      *Profile           `bson:"profile,omitempty"`
}

func NewAuthRepo(coll *mongo.Collection) *AuthRepo {
	return &AuthRepo{coll: coll}
}

func (a *AuthRepo) CreateUser(ctx context.Context, req *genproto.UserReq) (*genproto.UserResp, error) {
	result, err := a.coll.InsertOne(ctx, bson.M{
		"username":      req.Username,
		"email":         req.Email,
		"password_hash": req.PasswordHash,
		"profile": bson.M{
			"first_name": req.Profile.FirstName,
			"last_name":  req.Profile.LastName,
			"address":    req.Profile.Address,
		},
	})
	if err != nil {
		return nil, err
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	return &genproto.UserResp{
		UserId: insertedID.Hex(),
		Status: "User Created",
		Email:  req.Email,
	}, nil
}

func (a *AuthRepo) UpdateUser(ctx context.Context, req *genproto.UserReq) (*genproto.UserResp, error) {
	id, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"username":      req.Username,
			"email":         req.Email,
			"password_hash": req.PasswordHash,
			"profile": bson.M{
				"first_name": req.Profile.FirstName,
				"last_name":  req.Profile.LastName,
				"address":    req.Profile.Address,
			},
		},
	}

	if _, err = a.coll.UpdateByID(ctx, id, update); err != nil {
		return nil, err
	}

	return &genproto.UserResp{
		UserId: req.UserId,
		Status: "User Updated",
		Email:  req.Email,
	}, nil
}

func (a *AuthRepo) DeleteUser(ctx context.Context, req *genproto.UserReq) (*genproto.UserResp, error) {
	id, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	if _, err = a.coll.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return nil, err
	}

	return &genproto.UserResp{
		UserId: req.UserId,
		Status: "User Deleted",
	}, nil
}

func (a *AuthRepo) GetUserById(ctx context.Context, req *genproto.UserReq) (*genproto.UserReq, error) {
	id, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return nil, err
	}

	var user UserReq
	if err = a.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &genproto.UserReq{
		UserId:       user.ID.Hex(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Profile:      convertProfileToGenProto(user.Profile),
	}, nil
}

func (a *AuthRepo) GetUserByFilter(ctx context.Context, req *genproto.UserReq) (*genproto.Users, error) {
	filter := constructFilter(req)

	cursor, err := a.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // ensure cursor is properly closed

	var users genproto.Users

	for cursor.Next(ctx) {
		var user UserReq
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users.Users = append(users.Users, &genproto.UserReq{
			UserId:       user.ID.Hex(),
			Username:     user.Username,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			Profile:      convertProfileToGenProto(user.Profile),
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func constructFilter(req *genproto.UserReq) bson.M {
	filter := bson.M{}
	if req.Username != "" {
		filter["username"] = req.Username
	}
	if req.Email != "" {
		filter["email"] = req.Email
	}
	if req.Profile != nil {
		if req.Profile.FirstName != "" {
			filter["profile.first_name"] = req.Profile.FirstName
		}
		if req.Profile.LastName != "" {
			filter["profile.last_name"] = req.Profile.LastName
		}
		if req.Profile.Address != "" {
			filter["profile.address"] = req.Profile.Address
		}
	}
	return filter
}

func convertProfileToGenProto(profile *Profile) *genproto.Profile {
	if profile == nil {
		return nil
	}
	return &genproto.Profile{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Address:   profile.Address,
	}
}
