package mongodb

import (
	"batyrnosquare/go-premiership/pkg/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentModel struct {
	DB *mongo.Collection
}

func (m *CommentModel) Insert(users *models.User, news *models.News, text string) error {
	_, err := m.DB.InsertOne(context.TODO(), bson.M{
		"user_id": users.ID,
		"news_id": news.ID,
		"text":    text,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) GetComment(comments *models.Comment) ([]*models.Comment, error) {
	cursor, err := m.DB.Find(context.TODO(), bson.M{"news_id": comments.NewsID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	commentList := []*models.Comment{}
	for cursor.Next(context.TODO()) {
		var comment models.Comment
		err := cursor.Decode(&comment)
		if err != nil {
			return nil, err
		}
		commentList = append(commentList, &comment)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	fmt.Println(commentList)
	return commentList, nil
}

func (m *CommentModel) Delete(comments *models.Comment) error {
	_, err := m.DB.DeleteOne(context.TODO(), bson.M{"_id": comments.ID})
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentModel) GetNewsId(comments *models.Comment) (string, error) {
	var comment models.Comment
	err := m.DB.FindOne(context.TODO(), bson.M{"_id": comments.ID}).Decode(&comment)
	if err != nil {
		return "", err
	}
	return comment.NewsID, nil
}

func (m *CommentModel) GetUserId(comments *models.Comment) (string, error) {
	var comment models.Comment
	err := m.DB.FindOne(context.TODO(), bson.M{"_id": comments.ID}).Decode(&comment)
	if err != nil {
		return "", err
	}
	return comment.UserID, nil

}
