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
    "github.com/antigloss/go/logger"
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
    logger.Info("UpdateUser called")
    decoder := json.NewDecoder(r.Body)
    var userToUpdate models.UserModel
    err := decoder.Decode(&userToUpdate)
    if err != nil {
      logger.Info("Error while decoding Request")
      log.Println("Error while decoding Request");
    }

    user := database.GetUserFromToken(token)
    database.UpdateUser(&userToUpdate, user)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
      logger.Info("cant find user from token")
    } else {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      respJson, err := json.Marshal(userToUpdate)
      if err != nil {
          return
      }
      w.Write([]byte(respJson))
      logger.Info("User Updated")
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
    logger.Info("UpdateUserName called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
      logger.Info("cant find user from token")
      return ;
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      if (params["username"] == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("mandatory field not provided"))
        logger.Info("mandatory field <username> not provided")
        return ;
      }
      userName := params["username"].(string)

      if (len(userName) == 0) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("Username is empty"))
        logger.Info("Username is empty")
        return
      }

      database.UpdateUserName(&userName, user)
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      logger.Info("Username updated")
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
    logger.Info("UpdateUserEmail called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
      logger.Info("cant find user from token")
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      if (params["email"] == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("mandatory field not provided"))
        logger.Info("mandatory field <email> not provided")
        return ;
      }
      userEmail := params["email"].(string)

      if (len(userEmail) == 0) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("Email is empty"))
        logger.Info("Email is empty")
        return
      }

      err := database.UpdateUserEmail(&userEmail, user)
      if (err != nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        logger.Info("User Email updated")
      } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("Email already exist"))
        logger.Info("Email already exist")
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
    logger.Info("UpdateUserPicture called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
      logger.Info("cant find user from token")
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      if (params["profile_picture_url"] == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("mandatory field not provided"))
        logger.Info("mandatory field profile_picture_url not provided")
      }
      userPicture := params["profile_picture_url"].(string)

      if (len(userPicture) == 0) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("picture url is empty"))
        logger.Info("picture url is empty")
        return
      }
      database.UpdateUserPicture(&userPicture, user)
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      w.Write([]byte("Picture Updated"))
      logger.Info("picture Updated")
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
    logger.Info("UpdateUserPassword called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
      logger.Info("cant find user from token")
      return ;
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)
      if (params["old_password"] == nil || params["new_password"] == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("mandatory fields not provided"))
        logger.Info("mandatory fields old_password not provided")
        return ;
      }
      oldPassword := params["old_password"].(string)
      newPassword := params["new_password"].(string)

      if (len(newPassword) == 0) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("new password is empty"))
        logger.Info("new password is empty")
        return
      }

      err := database.UpdateUserPassword(&oldPassword, &newPassword, user)
      if (err != nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte(err.Error()))
        logger.Info("Error updating password: %s", err.Error())
      } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        logger.Info("User Password updated")
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
    logger.Info("FetchUser called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("cant find user from token"))
      logger.Info("cant find user from token")
    } else {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(200)
      respJson, err := json.Marshal(user)
      if err != nil {
        logger.Info("Error when fetching user:%s", err)
        return
      }
      logger.Info("User Fetched")
      w.Write([]byte(respJson))
    }
  }
}


/*
** PUT Request on /acceptContactInvite/{contact_identifier}
** Header: Authorization: Bearer token
*/
func AcceptContactInvite(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  logger.Info("AcceptContactInvite called")
  // vars := mux.Vars(r)
  //
  // contact_identifier := vars["contact_identifier"]
}

/*
** PUT Request on /refuseContactInvite/{contact_identifier}
** Header: Authorization: Bearer token
*/
func RefuseContactInvite(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  logger.Info("RefuseContactInvite called")
  // vars := mux.Vars(r)
  //
  // contact_identifier := vars["contact_identifier"]
}

/*
** PUT Request on /refuseGroupInvite/{group_identifier}
** Header: Authorization: Bearer token
*/
func RefuseGroupInvite(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  logger.Info("RefuseGroupInvite called")
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** PUT Request on /acceptGroupInvite/{group_identifier}
** Header: Authorization: Bearer token
*/
func AcceptGroupInvite(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  logger.Info("AcceptGroupInvite called")
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

func FetchGroups(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("FetchGroups called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      data := database.FetchGroupsFromUser(&user.ID, &user.Email)
      respJson, err := json.Marshal(bson.M{"groups": data})
      if err != nil {
          return;
      }
      w.WriteHeader(200)
      w.Write([]byte(respJson))
    }
  }
}


/*
** PUT Request on /renameGroup/{group_identifier}
** Arguments: name
** Header: Authorization: Bearer token
*/
func RenameGroup(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("CreateGroup called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)

      log.Println(params["name"])

      if (params["name"] == nil){
        w.WriteHeader(400)
        w.Write([]byte("No Group Name specified"))
        logger.Info("No Group Name specified")
        return ;
      }
      if (params["new_name"] == nil){
        w.WriteHeader(400)
        w.Write([]byte("No New Group Name specified"))
        logger.Info("No New Group Name specified")
        return ;
      }
      groupName := params["name"].(string)
      groupId := params["id"].(string)
      newGroupName := params["new_name"].(string)

      if (groupName) {
        if (database.IsExistingGroup(&groupName, &user.ID) == false) {
          w.WriteHeader(400)
          w.Write([]byte("You do not have a group with this name"))
          logger.Info("You do not have a group with this name: %s", groupName)
        } else if (database.IsExistingGroup(&newGroupName, &user.ID) == true) {
          w.WriteHeader(400)
          w.Write([]byte("You already have a group with this name"))
          logger.Info("You already have a group with this name: %s", groupName)
        } else {
          database.UpdateGroupName(&groupName, &user.ID, &newGroupName)
          respJson, err := json.Marshal(bson.M{"group_name": newGroupName})
           if err != nil {
               return
           }
          logger.Info("Group Renamed: name:%s", newGroupName)
          w.Header().Set("Content-Type", "application/json")
          w.WriteHeader(200)
          w.Write([]byte(respJson))
        }
      } else if (groupId) {
        if (database.IsExistingGroupId(&groupId, &user.ID) == false) {
          w.WriteHeader(400)
          w.Write([]byte("You do not have a group with this name"))
          logger.Info("You do not have a group with this name: %s", groupName)
        } else if (database.IsExistingGroup(&newGroupName, &user.ID) == true) {
          w.WriteHeader(400)
          w.Write([]byte("You already have a group with this name"))
          logger.Info("You already have a group with this name: %s", groupName)
        } else {
          database.UpdateGroupNameById(&groupId, &user.ID, &newGroupName)
          respJson, err := json.Marshal(bson.M{"group_name": newGroupName})
           if err != nil {
               return
           }
          logger.Info("Group Renamed: name:%s", newGroupName)
          w.Header().Set("Content-Type", "application/json")
          w.WriteHeader(200)
          w.Write([]byte(respJson))
        }
      }
    }
  }
}


/*
** PUT Request on /leaveGroup/{group_identifier}
** Header: Authorization: Bearer token
*/
func LeaveGroup(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  logger.Info("LeaveGroup called")
  // vars := mux.Vars(r)
  //
  // group_identifier := vars["group_identifier"]
}

/*
** PUT Request on /removeFromGroup/{group_identifier}
** Arguments: contact_identifier
** Header: Authorization: Bearer token
*/
func RemoveFromGroup(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("RemoveFromGroup called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      vars := mux.Vars(r)
      groupName := vars["group_identifier"]

      if (database.IsExistingGroup(&groupName, &user.ID) == true) {
        var params map[string]interface{}
        decoder := json.NewDecoder(r.Body)
        decoder.Decode(&params)

        if (params["contact_identifier"] == nil) {
          w.WriteHeader(400)
          w.Write([]byte("No contact indentifier specified"))
          logger.Info("No contact indentifier specified")
          return ;
        }
        contactId := params["contact_identifier"].(string)
        strErr := database.RemoveFromGroup(&groupName, &user.ID, &contactId)
        if (strErr == "") {
          w.WriteHeader(200)
          w.Write([]byte("User removed from group"))
        } else {
          w.WriteHeader(400)
          w.Write([]byte(strErr))
        }
      } else {
        w.WriteHeader(400)
        w.Write([]byte("No group with that name"))
        logger.Info("You already have a group with this name: %s", groupName)
        log.Println("No group with that name")
      }
    }
  }
}


/*
** PUT Request on /addToGroup/{group_identifier}
** Arguments: contact_identifier
** Header: Authorization: Bearer token
*/
func AddToGroup(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("AddToGroup called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      vars := mux.Vars(r)
      groupName := vars["group_identifier"]

      if (database.IsExistingGroup(&groupName, &user.ID) == true) {
        var params map[string]interface{}
        decoder := json.NewDecoder(r.Body)
        decoder.Decode(&params)

        if (params["contact_identifier"] == nil) {
          w.WriteHeader(400)
          w.Write([]byte("No contact indentifier specified"))
          logger.Info("No contact indentifier specified")
          return ;
        }
        contactId := params["contact_identifier"].(string)
        strErr := database.AddToGroup(&groupName, &user.ID, &contactId)
        if (strErr == "") {
          w.WriteHeader(200)
          w.Write([]byte("User added to group"))
        } else {
          w.WriteHeader(400)
          w.Write([]byte(strErr))
        }
      } else {
        w.WriteHeader(400)
        w.Write([]byte("No group with that name"))
        logger.Info("You already have a group with this name: %s", groupName)
        log.Println("No group with that name")
      }
    }
  }
}

/*
** DELETE Request on /removeGroup/{group_identifier}
** Header: Authorization: Bearer token
*/
func RemoveGroup(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("RemoveGroup called")
    log.Println("RemoveGroup called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      vars := mux.Vars(r)
      groupName := vars["group_identifier"]

      if (database.IsExistingGroup(&groupName, &user.ID) == true) {
        database.DeleteGroup(&groupName, &user.ID)
        w.WriteHeader(200)
        w.Write([]byte("Group Deleted"))
      } else {
        w.WriteHeader(400)
        w.Write([]byte("No group with that name"))
        logger.Info("You already have a group with this name: %s", groupName)
        log.Println("No group with that name")
      }
    }
  }
}

/*
** POST Request on /createGroup
** Arguments: name
** Header: Authorization: Bearer token
** Response: group identifier
*/
func CreateGroup(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("CreateGroup called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      var params map[string]interface{}
      decoder := json.NewDecoder(r.Body)
      decoder.Decode(&params)

      log.Println(params["name"])

      if (params["name"] == nil){
        w.WriteHeader(400)
        w.Write([]byte("No Group Name specified"))
        logger.Info("No Group Name specified")
        return ;
      }
      groupName := params["name"].(string)

      if (database.IsExistingGroup(&groupName, &user.ID) == true) {
        w.WriteHeader(400)
        w.Write([]byte("You already have a group with this name"))
        logger.Info("You already have a group with this name: %s", groupName)
      } else {
        var newGroup models.GroupModel
        newGroup.Name = groupName
        newGroup.Owner = user.ID
        newGroup.ID = bson.ObjectId.Hex(bson.NewObjectId())
        database.CreateGroup(&newGroup)
        respJson, err := json.Marshal(bson.M{"group_id": newGroup.ID})
         if err != nil {
             return
         }
        logger.Info("Group Created: name:%s, owner:%s, id:%s", newGroup.Name, newGroup.Owner, newGroup.ID)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write([]byte(respJson))
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
    logger.Info("RemoveContact called")
    vars := mux.Vars(r)
    contact_identifier := vars["contact_identifier"]

    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
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
** DELETE Request on /removeContactByEmail/{contact_identifier}
** Header: Authorization: Bearer token
*/
func RemoveContactByEmail(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("RemoveContactByEmail called")
    vars := mux.Vars(r)
    contact_email := vars["contact_email"]
    contact := database.GetUsersByEmail(contact_email)
    contact_identifier := contact.ID

    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      isInside, index := isContactInArray(&contact_identifier, &user.Contacts)
      if (isInside) {
        user.Contacts = append(user.Contacts[:index], user.Contacts[index+1:]...)
        database.UpdateUserContacts(user)
	       respJson, err := json.Marshal(bson.M{"response": "Contact deleted succesfully"})
        if err != nil {
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
    	   w.Write([]byte(respJson))
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
    logger.Info("AddContact called")
    var params map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    decoder.Decode(&params)

    email := params["email"]
    if (email == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("No email provided"))
      logger.Info("AddContact: No email provided")
    } else {
      contact := database.GetUsersByEmail(email.(string))
      if (contact == nil) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(400)
        w.Write([]byte("No user with this Email"))
        logger.Info("AddContact: No user with this Email")
      } else {
        user := database.GetUserFromToken(token)
        if (email == user.Email) {
          w.Header().Set("Content-Type", "application/json")
          w.WriteHeader(400)
          w.Write([]byte("Same Email as the logged user"))
          logger.Info("AddContact: Same Email as the logged user")
          return ;
        } else {
          for _,tmpContact := range user.Contacts {
            if (contact.ID == tmpContact) {
              w.Header().Set("Content-Type", "application/json")
              w.WriteHeader(400)
              w.Write([]byte("user already in your contacts"))
              logger.Info("AddContact: user already in your contacts")
              return ;
            }
            // element is the element from someSlice for where we are
          }
        }
        database.AddContactToUser(&user.Email, &contact.ID)
        respJson, err := json.Marshal(bson.M{"contact": contact})
      	if err != nil {
             		return;
            	}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write([]byte(respJson))
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
  logger.Info("FetchUserContact called")
  vars := mux.Vars(r)

  contact_identifier := vars["contact_identifier"]
  user := database.GetUserFromToken(token)
  if (user == nil) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write([]byte("User does not exist"))
    logger.Info("User does not exist")
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
    logger.Info("FetchUserContacts called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
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
** POST Request on /addFileToList
** Arguments: Filename
** Header: Authorization: Bearer token
** Response: {}
*/

func AddFileToList(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("AddFileToList called")
    user := database.GetUserFromToken(token)
    var params map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    decoder.Decode(&params)

    file_url := params["file_url"]
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      database.AddFileToUser(&user.Email, file_url.(string))
      w.WriteHeader(200)
      respJson, err := json.Marshal(bson.M{"success": "file added"})
      if err != nil {
          return
      }
      log.Println("resp", respJson)
      w.Write([]byte(respJson))
    }
  }
}

/*
** POST Request on /removeFileFromList
** Arguments: Filename
** Header: Authorization: Bearer token
** Response: {}
*/

func RemoveFileFromList(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("RemoveFileFromList called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
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
** POST Request on /getUserFiles
** Arguments: none
** Header: Authorization: Bearer token
** Response: {[FileUrls]}
*/

func GetUserFiles(token *string) negroni.HandlerFunc {
  return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logger.Info("AddFileToList called")
    user := database.GetUserFromToken(token)
    if (user == nil) {
      w.Header().Set("Content-Type", "application/json")
      w.WriteHeader(400)
      w.Write([]byte("User does not exist"))
      logger.Info("User does not exist")
    } else {
      w.WriteHeader(200)
      respJson, err := json.Marshal(bson.M{"file_urls": user.FileUrls})
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
  logger.Info("LoginUser called")
  decoder := json.NewDecoder(r.Body)
  var user models.UserModel
  err := decoder.Decode(&user)
  if err != nil {
    log.Println("Error while decoding Request");
    logger.Info("Error while decoding Request")
  } else {
    log.Printf("%s %s %s %s\n", user.Firstname, user.LastName, user.Email, user.PhoneNumber);
    logger.Info("Received user login with Firstname:%s, Lastname:%s, Email:%s, PhoneNumber:%s, Password:%s, FBToken:%s, FBUserID:%s", user.Firstname, user.LastName, user.Email, user.PhoneNumber, user.Password, user.FBToken, user.FBUserID)
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
      logger.Info("User Logged in Failed: %s", err)
        return
    }
    logger.Info("User Logged in")
    w.Write([]byte(respJson))
  } else {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write([]byte(errMsg))
    logger.Info("User Logged in Failed: %s", errMsg)
  }
}

/*
** POST Request on /createUser
** Arguments: firstname, lastname, password(exist when user doesn't login with facebook), email, fbtoken(exist when user login with facebook), fbid(exist when user login with facebook)
** Response: {"token": "xxxxxxx", "identifier": "xxxx"}
*/
func CreateUser(w http.ResponseWriter, r *http.Request)  {
  logger.Info("CreateUser called")
  decoder := json.NewDecoder(r.Body)
  var user models.UserModel
  err := decoder.Decode(&user)
  if err != nil {
    log.Println("Error while decoding Request in create user");
    log.Println(err)
    logger.Info("Error while decoding Request in create user: %s", err)
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
    logger.Info("Create User:%s, LastName:%s, Email:%s", user.Firstname, user.LastName, user.Email)
    w.Write([]byte(respJson))
  } else {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write([]byte(errMsg))
    logger.Info("Couldn't create Facebook User: %s", errMsg)
  }
}

func createFacebookUser(w http.ResponseWriter, r *http.Request, user models.UserModel) {
  logger.Info("createFacebookUser called")
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
    logger.Info("Creating Facebook User With Firstname:%s, LastName:%s, Email:%s", user.Firstname, user.LastName, user.Email)
  } else {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(400)
    w.Write([]byte(errMsg))
    logger.Info("Couldn't create Facebook User: %s", errMsg)
  }
}
