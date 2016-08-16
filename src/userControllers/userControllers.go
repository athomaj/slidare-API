package userControllers

import (
    "models"
    "encoding/json"
    "net/http"
    "log"
    "core/authentication"
    "github.com/gorilla/mux"
    "github.com/lonelycode/go-uuid/uuid"
    "api/parameters"
    "mydb"
    "gopkg.in/mgo.v2/bson"
    "github.com/codegangsta/negroni"
)

type FBTokenValidationResponse struct {
  Name string `json: "name"`
  ID string `json: "id"`
}

type SignInResponse struct {
  Token string `json:"token"`
  ID string `json:"id" bson:"_id,omitempty"`
}

func initAndGenerateToken() (int, []byte, string) {
  authBackend := authentication.InitJWTAuthenticationBackend()
  token, err := authBackend.GenerateToken(uuid.New())
  if err != nil {
      return http.StatusInternalServerError, []byte(""), string("")
  } else {
      response, _ := json.Marshal(parameters.TokenAuthentication{token})
      return http.StatusOK, response, token
  }
  return http.StatusUnauthorized, []byte(""), string("")
}

func checkUserInformationsForCreation(user models.UserModel) (bool, string){
  log.Println(user.Firstname, ":", user.LastName, ":", user.Email, ":", user.FBToken, ":", user.FBUserID, ":", user.Password)
  if len(user.Firstname) == 0 || len(user.LastName) == 0 || len(user.Email) == 0 || (len(user.FBToken) == 0 && len(user.Password) == 0)|| (len(user.FBUserID) == 0  && len(user.Password) == 0){
    return false, "User informations not filled"
  } else if database.DoesEmailExistInDB(user.Email) == true {
    return false, "User email already exist"
  }
  if (len(user.Password) == 0) {
    resp, err := http.Get("https://graph.facebook.com/me?access_token=" + user.FBToken)
    if (err != nil) {
      log.Println("err:", err)
    } else {
      var jsonResp FBTokenValidationResponse
      json.NewDecoder(resp.Body).Decode(&jsonResp)
      if (resp.StatusCode != 200) {
        return false, "Token not valid"
      } else if (jsonResp.ID != user.FBUserID) {
        return false, "ID not valid"
      }
    }
  }
  return true, ""
}

func checkUserInformationsForLogin(user models.UserModel) (bool, string){
  if len(user.Email) == 0 || (len(user.FBToken) == 0 && len(user.Password) == 0)|| (len(user.FBUserID) == 0  && len(user.Password) == 0){
    return false, "User informations not filled"
  } else if database.DoesEmailExistInDB(user.Email) == false {
    return false, "User does not exist"
  }
  if (len(user.Password) != 0) {
    if (database.ValidateUserPassword(&user.Email, &user.Password) == false) {
      return false, "Incorect password"
    }
  } else {
    resp, err := http.Get("https://graph.facebook.com/me?access_token=" + user.FBToken)
    if (err != nil) {
      log.Println("err:", err)
    } else {
      var jsonResp FBTokenValidationResponse
      json.NewDecoder(resp.Body).Decode(&jsonResp)
      if (resp.StatusCode != 200) {
        return false, "Token not valid"
      } else if (jsonResp.ID != user.FBUserID) {
        return false, "ID not valid"
      }
    }
  }
  return true, ""
}

/*
** POST Request on /updateUser
** Arguments: UserModel
** Header: Authorization: Bearer token
** Response: Updated UserModel
*/
func UpdateUser(token *string) negroni.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    decoder := json.NewDecoder(r.Body)
    var userToUpdate models.UserModel
    err := decoder.Decode(&userToUpdate)
    if err != nil {
      log.Println("Error while decoding Request");
    }

    user := database.GetUserFromToken(token)
    database.UpdateUser(&userToUpdate, user)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
    } else {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      respJson, err := json.Marshal(userToUpdate)
      if err != nil {
          return
      }
      w.Write([]byte(respJson))
    }
  }
}

/*
** POST Request on /updateUserName
** Arguments: username
** Header: Authorization: Bearer token
*/
func UpdateUserName(token *string) negroni.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
      return ;
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      if (params["username"] == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("mandatory field not provided"))
        return ;
      }
      userName := params["username"].(string)

      if (len(userName) == 0) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("Username is empty"))
        return
      }

      database.UpdateUserName(&userName, user)
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
    }
  }
}

/*
** POST Request on /updateUserEmail
** Arguments: email
** Header: Authorization: Bearer token
*/
func UpdateUserEmail(token *string) negroni.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      if (params["email"] == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("mandatory field not provided"))
        return ;
      }
      userEmail := params["email"].(string)

      if (len(userEmail) == 0) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("Email is empty"))
        return
      }

      err := database.UpdateUserEmail(&userEmail, user)
      if (err != nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
      } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("Email already exist"))
      }
    }
  }
}

/*
** POST Request on /updateUserPicture
** Arguments: profile_picture_url
** Header: Authorization: Bearer token
*/
func UpdateUserPicture(token *string) negroni.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      if (params["profile_picture_url"] == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("mandatory field not provided"))
      }
      userPicture := params["profile_picture_url"].(string)

      if (len(userPicture) == 0) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("picture url is empty"))
        return
      }
      database.UpdateUserPicture(&userPicture, user)
    }
  }
}

/*
** POST Request on /updateUserPassword
** Arguments: old_password, new_password
** Header: Authorization: Bearer token
*/
func UpdateUserPassword(token *string) negroni.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
      return ;
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      if (params["old_password"] == nil || params["new_password"] == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("mandatory fields not provided"))
        return ;
      }
      oldPassword := params["old_password"].(string)
      newPassword := params["new_password"].(string)

      if (len(newPassword) == 0) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("new password is empty"))
        return
      }

      err := database.UpdateUserPassword(&oldPassword, &newPassword, user)
      if (err != nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte(err.Error()))
      } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
      }
    }
  }
}
/*
** GET Request on /fetchUser
** Header: Authorization: Bearer token
** Response: UserModel
*/
func FetchUser(token *string) negroni.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
    } else {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      respJson, err := json.Marshal(user)
      if err != nil {
          return
      }
      w.Write([]byte(respJson))
    }
  }
}


/*
** PUT Request on /acceptContactInvite/{contact_identifier}
** Header: Authorization: Bearer token
*/
func AcceptContactInvite(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // contact_identifier := vars["contact_identifier"]
}

/*
** PUT Request on /refuseContactInvite/{contact_identifier}
** Header: Authorization: Bearer token
*/
func RefuseContactInvite(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // contact_identifier := vars["contact_identifier"]
}

/*
** PUT Request on /refuseGroupInvite/{group_identifier}
** Header: Authorization: Bearer token
*/
func RefuseGroupInvite(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** PUT Request on /acceptGroupInvite/{group_identifier}
** Header: Authorization: Bearer token
*/
func AcceptGroupInvite(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** PUT Request on /renameGroup/{group_identifier}
** Arguments: name
** Header: Authorization: Bearer token
*/
func RenameGroup(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** PUT Request on /removeFromGroup/{group_identifier}
** Arguments: contact_identifier
** Header: Authorization: Bearer token
*/
func RemoveFromGroup(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** PUT Request on /leaveGroup/{group_identifier}
** Header: Authorization: Bearer token
*/
func LeaveGroup(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** PUT Request on /addToGroup/{group_identifier}
** Arguments: contact_identifier
** Header: Authorization: Bearer token
*/
func AddToGroup(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** DELETE Request on /removeGroup/{group_identifier}
** Header: Authorization: Bearer token
*/
func RemoveGroup(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** POST Request on /createGroup
** Arguments: name
** Header: Authorization: Bearer token
** Response: group identifier
*/
func CreateGroup(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      groupName := params["name"].(string)

      if (database.IsExistingGroup(&groupName, &user.ID) == true) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("You already have a group with this name"))
      } else {
        var newGroup models.GroupModel
        newGroup.Name = groupName
        newGroup.Owner = user.ID
        newGroup.ID = bson.ObjectId.Hex(bson.NewObjectId())
        database.CreateGroup(&newGroup)
      }
    }
  }
}

func isContactInArray(contactId *string, contactList *[]string) (bool, int){
  for idx, tmpContact := range *contactList {
    if tmpContact == *contactId {
        return true, idx
    }
  }
  return false, 0
}

/*
** DELETE Request on /removeContact/{contact_identifier}
** Header: Authorization: Bearer token
*/
func RemoveContact(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    vars := mux.Vars(r)
    contact_identifier := vars["contact_identifier"]

    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
    } else {
      isInside, index := isContactInArray(&contact_identifier, &user.Contacts)
      if (isInside) {
        user.Contacts = append(user.Contacts[:index], user.Contacts[index+1:]...)
        database.UpdateUserContacts(user)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
      } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("Contact not in user's contact list"))
      }
    }
  }
}

/*
** POST Request on /addContact
** Arguments: email
** Header: Authorization: Bearer token
*/
func AddContact(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    var params map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    decoder.Decode(&params)

    email := params["email"]
    if (email == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("No email provided"))
    } else {
      contact := database.GetUsersByEmail(email.(string))
      if (contact == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("No user with this Email"))
      } else {
        user := database.GetUserFromToken(token)
        if (email == user.Email) {
          w.Header().Set("Content-Type", "application/json")
          w.WriteHeader(400)
          w.Write([]byte("Same Email as the logged user"))
          return ;
        } else {
          for _,tmpContact := range user.Contacts {
            if (contact.ID == tmpContact) {
              w.Header().Set("Content-Type", "application/json")
              w.WriteHeader(400)
              w.Write([]byte("user already in your contacts"))
              return ;
            }
            // element is the element from someSlice for where we are
          }
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        database.AddContactToUser(&user.Email, &contact.ID)
      }
    }
  }
}
/*
** GET Request on /userContact/{contact_identifier}
** Header: Authorization: Bearer token
** Response: {"contact: ContactModel"}
*/
func FetchUserContact(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  vars := mux.Vars(r)

  contact_identifier := vars["contact_identifier"]
  user := database.GetUserFromToken(token)
  if (user == nil) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write([]byte("User does not exist"))
  } else {
    isInside, _ := isContactInArray(&contact_identifier, &user.Contacts)
    if (isInside) {
      contact := database.GetUserFromId(&contact_identifier)
      respJson, err := json.Marshal(bson.M{"contact": contact})
      if err != nil {
          return
      }
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      w.Write([]byte(respJson))
    } else {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("Contact identifier not in user contact list"))
    }
  }
  }
}

/*
** GET Request on /userContacts
** Arguments: NONE
** Header: Authorization: Bearer token
** Response: {"contacts: [ContactModel]"}
*/
func FetchUserContacts(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    log.Println("here")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
    } else {
      var contactList []models.UserModel
      for _, tmpContact := range user.Contacts {
        contactList = append(contactList, *database.GetUserFromId(&tmpContact))
      }
      w.WriteHeader(200)
      respJson, err := json.Marshal(bson.M{"contacts": contactList})
      if err != nil {
          return
      }
      log.Println("resp", respJson)
      w.Write([]byte(respJson))
    }
  }
}

/*
** POST Request on /loginUser
** Arguments: firstname, lastname, password(exist when user doesn't login with facebook), email, fbtoken(exist when user login with facebook)
** Response: {"token": "xxxxxxx", "identifier": "xxxx"}
*/
func LoginUser(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var user models.UserModel
  err := decoder.Decode(&user)
  if err != nil {
    log.Println("Error while decoding Request");
  } else {
    log.Printf("%s %s %s %s\n", user.Firstname, user.LastName, user.Email, user.PhoneNumber);
  }

  valid, errMsg := checkUserInformationsForLogin(user)
  if (valid == false && errMsg == "User does not exist" && len(user.FBToken) != 0) {
    createFacebookUser(w, r, user)
    return ;
  }
  if (valid){
    user := database.GetUsersByEmail(user.Email)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    var resp SignInResponse
    resp.Token = user.Token
    resp.ID = user.ID
    respJson, err := json.Marshal(resp)
    if err != nil {
        return
    }
    w.Write([]byte(respJson))
  } else {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write([]byte(errMsg))
  }
}

/*
** POST Request on /createUser
** Arguments: firstname, lastname, password(exist when user doesn't login with facebook), email, fbtoken(exist when user login with facebook), fbid(exist when user login with facebook)
** Response: {"token": "xxxxxxx", "identifier": "xxxx"}
*/
func CreateUser(w http.ResponseWriter, r *http.Request)  {
  log.Println(r)
  decoder := json.NewDecoder(r.Body)
  var user models.UserModel
  err := decoder.Decode(&user)
  if err != nil {
    log.Println("Error while decoding Request in create user");
    log.Println(err)
  }
  valid, errMsg := checkUserInformationsForCreation(user)
  if (valid){
    responseStatus, _, tokenString := initAndGenerateToken()
    user.ID = bson.ObjectId.Hex(bson.NewObjectId())
    user.Token = tokenString;

    if (len(user.UserName) == 0) {
      user.UserName = user.Email
    }

    if (len(user.FBUserID) != 0) {
      user.ProfilePictureURL = "https://graph.facebook.com/" + user.FBUserID + "/picture?type=large"
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(responseStatus)

    database.CreateNewUser(user)

    var resp SignInResponse
    resp.Token = user.Token
    resp.ID = user.ID
    respJson, err := json.Marshal(resp)
    if err != nil {
        return
    }
    w.Write([]byte(respJson))
  } else {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write([]byte(errMsg))
  }
}

func createFacebookUser(w http.ResponseWriter, r *http.Request, user models.UserModel) {
  valid, errMsg := checkUserInformationsForCreation(user)
  log.Println(valid)
  if (valid){
    responseStatus, _, tokenString := initAndGenerateToken()
    user.ID = bson.ObjectId.Hex(bson.NewObjectId())
    user.Token = tokenString;

    if (len(user.UserName) == 0) {
      user.UserName = user.Email
    }

    if (len(user.FBUserID) != 0) {
      user.ProfilePictureURL = "https://graph.facebook.com/" + user.FBUserID + "/picture?type=large"
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(responseStatus)

    database.CreateNewUser(user)

    var resp SignInResponse
    resp.Token = user.Token
    resp.ID = user.ID
    respJson, err := json.Marshal(resp)
    if err != nil {
        return
    }
    w.Write([]byte(respJson))
  } else {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write([]byte(errMsg))
  }
}
