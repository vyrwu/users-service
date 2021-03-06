package users

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/iafoosball/users-service/models"
	"github.com/iafoosball/users-service/restapi/operations"
	"time"
)

func MakeFriendRequest() func(params operations.PostFriendsUserIDFriendIDParams) middleware.Responder {
	return func(params operations.PostFriendsUserIDFriendIDParams) middleware.Responder {
		friend := params.Body
		friend.Accepted = false
		friend.DatetimeRequest = time.Now().Format(time.RFC3339)
		friend.Key = params.UserID + params.FriendID
		if _, err := friendsCol.CreateDocument(nil, friend); err != nil {
			panic(err)
		}
		return operations.NewPostFriendsUserIDFriendIDOK()
	}
}

func AcceptFriendRequest() func(params operations.PatchFriendsUserIDFriendIDParams) middleware.Responder {
	return func(params operations.PatchFriendsUserIDFriendIDParams) middleware.Responder {
		query := "Update {_key: \"" + params.UserID + params.FriendID + "\"} WITH { accepted: true, datetime_accepted: \"" +
			time.Now().Format(time.RFC3339) + "\" } IN friends "
		if _, err := database.Query(nil, query, nil); err != nil {
			panic(err)
		}
		return operations.NewPatchFriendsUserIDFriendIDOK()
	}
}

func DeleteFriend() func(params operations.DeleteFriendsFriendshipIDParams) middleware.Responder {
	return func(params operations.DeleteFriendsFriendshipIDParams) middleware.Responder {
		if _, err := friendsCol.RemoveDocument(nil, params.FriendshipID); err != nil {
			panic(err)
		}
		return operations.NewDeleteFriendsFriendshipIDOK()
	}
}

func ErrorHandling(err error) {
	go GetFriends()
}

func GetFriends() func(params operations.GetFriendsUserIDParams) middleware.Responder {
	return func(params operations.GetFriendsUserIDParams) middleware.Responder {
		query := "FOR users, edge, edgesArray IN 1 ANY 'users/" + params.UserID + "' GRAPH 'friends' FILTER edgesArray.edges[*].accepted ALL == true Return {users}"
		var friends []*models.User
		if cursor, err := database.Query(nil, query, nil); err != nil {
			panic(err)
		} else {
			for cursor.HasMore() {
				var friend *models.User
				cursor.ReadDocument(nil, friend)
				fmt.Println(friend)
				friends = append(friends, friend)
			}
		}
		return operations.NewGetFriendsUserIDOK().WithPayload(friends)
	}
}

//func MakeFriendRequest() func(params operations.GetUsersUserIDParams) middleware.Responder {
//	var f friend
//	vars := mux.Vars(r)
//	uid := vars["uid"]
//	friendid := vars["friendid"]
//	f.ACCEPTED = false
//	f.From = "users/" + uid
//	f.To = "users/" + friendid
//	_, _ = colFriends.CreateDocument(nil, f)
//}
//
//func AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	uid := vars["uid"]
//	friendid := vars["friendid"]
//	query := "FOR f IN friends FILTER f._from == 'users/" + friendid + "' and f._to == 'users/" + uid + "' RETURN f"
//	cursor, _ := dbUsers.Query(context.Background(), query, nil)
//	defer cursor.Close()
//	var f friend
//	meta, _ := cursor.ReadDocument(nil, &f)
//	f.ACCEPTED = true
//	f.From = "users/" + friendid
//	f.To = "users/" + uid
//	_, _ = colFriends.UpdateDocument(nil, meta.Key, f)
//	w.WriteHeader(http.StatusOK)
//}
//
//func ListOpenFriendRequests(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	uid := vars["uid"]
//	openRequests := []user{}
//	query := "FOR v, e, p IN 1 ANY 'users/" + uid + "' GRAPH 'friends' FILTER p.edges[*].accepted ALL == false Return {username: v.username, uid: v.uid}"
//	cursor, _ := dbUsers.Query(context.Background(), query, nil)
//	for {
//		var u user
//		_, err := cursor.ReadDocument(nil, &u)
//		if driver.IsNoMoreDocuments(err) {
//			break
//		}
//		openRequests = append(openRequests, u)
//	}
//	w.WriteHeader(http.StatusOK)
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(openRequests)
//}
