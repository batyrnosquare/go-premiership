package mongodb

import (
	"batyrnosquare/go-premiership/pkg/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserModel struct {
	DB  *mongo.Collection
	ctx context.Context
}

func (u *UserModel) Insert(ctx context.Context, users *models.User) error {

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(users.HashedPassword), 12)
	if err != nil {
		return err
	}

	user := bson.M{
		"name":            users.Name,
		"email":           users.Email,
		"hashed_password": string(hashedPw),
		"role":            users.Role,
	}

	_, err = u.DB.InsertOne(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "users_uc_email") {
			return models.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (u *UserModel) Authenticate(ctx context.Context, users *models.User) (int, error) {
	var user bson.M
	filter := bson.M{"email": users.Email}
	err := u.DB.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}
	hashedPw := []byte(user["hashed_password"].(string))
	err = bcrypt.CompareHashAndPassword(hashedPw, []byte(users.HashedPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return user["_id"].(int), nil
}

func (u *UserModel) Get(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	filter := bson.M{"_id": id}
	err := u.DB.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Users(ctx context.Context) ([]*models.User, error) {
	cursor, err := m.DB.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
