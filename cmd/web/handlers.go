package main

import (
	"batyrnosquare/go-premiership/pkg/models"
	"batyrnosquare/go-premiership/pkg/models/mongodb"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.news.Latest(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(s)
}

func (app *application) showNews(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.URL.Query().Get(":id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	s, err := app.news.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(s)

}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {

	var newPost models.News

	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err = json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	createPost, err := app.news.Insert(r.Context(), &newPost)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(createPost)
	if err != nil {
		app.serverError(w, err)
		return
	}

}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.users.Insert(r.Context(), &newUser)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "You signed up successfully. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) (*application, error) {

	var user models.User

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	id, err := app.users.Authenticate(r.Context(), &user)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.clientError(w, http.StatusBadRequest)
			return &application{
				session: nil,
			}, err
		} else {
			app.serverError(w, err)
		}

	}
	app.session.Put(r, "authenticatedUserID", id)
	return &application{
			session: app.session, users: app.users,
		},
		nil
}

func (app *application) commentNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID := app.session.GetInt(r, "authenticatedUserID")
	comment.UserID = strconv.Itoa(userID)
	err = app.comments.Insert(&models.User{ID: comment.UserID}, &models.News{ID: comment.NewsID}, comment.Text)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", comment.NewsID), http.StatusSeeOther)
}

func (app *application) deleteComment(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")
	user, err := app.users.Get(r.Context(), userID)
	commentID, err := strconv.Atoi(r.FormValue("commentID"))
	if err != nil || commentID < 1 {
		app.serverError(w, err)
		return
	}
	newsID, err := app.comments.GetNewsId(&models.Comment{ID: strconv.Itoa(commentID)})
	if err != nil {
		app.serverError(w, err)
		return
	}
	commUserID, err := app.comments.GetUserId(&models.Comment{ID: strconv.Itoa(commentID)})
	if err != nil {
		app.serverError(w, err)
		return
	}
	if user.Role != string(mongodb.Admin) && strconv.Itoa(userID) != commUserID {
		app.session.Put(r, "flash", "You can only delete your own comments!")
		http.Redirect(w, r, fmt.Sprintf("/news?id=%d", newsID), http.StatusSeeOther)
		return
	}
	err = app.comments.Delete(&models.Comment{ID: strconv.Itoa(commentID)})
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", newsID), http.StatusSeeOther)
}

func loginUserHandler(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := app.loginUser(w, r)
		if err != nil {
			return
		}
		_ = result
	}
}

func (app application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
