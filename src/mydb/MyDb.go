package database

import (
    "sync"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "log"
    "models"
    "fmt"
    "errors"
)


var instance mongodb
var once sync.Once

type mongodb struct {
  Session *mgo.Session
}

func Init() {
  instance = mongodb{}
  session, err := mgo.Dial("mongodb://localhost:27017")
  instance.Session = session
          if err != nil {
              panic(err)
          }
}

func DoesTokenExist(token string) bool {
  c := instance.Session.DB("slidare").C("users")
  result := models.UserModel{}
  err := c.Find(bson.M{"token": token}).One(&result)
  if err == nil {
    return true
  }
  return false
}

func GetUsersByEmail(email string) *models.UserModel {
  c := instance.Session.DB("slidare").C("users")
  result := models.UserModel{}
  err := c.Find(bson.M{"email": email}).One(&result)
  if err != nil {
    return nil
  } else {
    fmt.Println("Phone:", result.Firstname)
  }
  return &result
}

func CreateNewUser(user models.UserModel) {
  if DoesEmailExistInDB(user.Email) {
    log.Println("Email already Exist")
    return ;
  }
  c := instance.Session.DB("slidare").C("users")
  err := c.Insert(&user)
  if err != nil {
    log.Fatal(err)
  }
}

func UpdateUserContacts(user *models.UserModel) {
  c := instance.Session.DB("slidare").C("users")
  c.UpdateId(user.ID, bson.M{"$set": bson.M{"contacts": user.Contacts}})
}

func UpdateUserName(userName *string, user *models.UserModel) {
  c := instance.Session.DB("slidare").C("users")
  c.UpdateId(user.ID, bson.M{"$set": bson.M{"username": *userName}})
}

func UpdateUserPicture(userPicture *string, user *models.UserModel) {
  c := instance.Session.DB("slidare").C("users")
  c.UpdateId(user.ID, bson.M{"$set": bson.M{"profile_picture_url": *userPicture}})
}

func UpdateUserEmail(userEmail *string, user *models.UserModel) error {
  c := instance.Session.DB("slidare").C("users")
  var result models.UserModel
  err := c.Find(bson.M{"email": userEmail}).One(&result)
  if err != nil {
    c.UpdateId(user.ID, bson.M{"$set": bson.M{"email": *userEmail}})
  }
  return err
}

func UpdateUserPassword(oldPassword *string, newPassword *string, user *models.UserModel) error {
  c := instance.Session.DB("slidare").C("users")
  if (ValidateUserPassword(&user.Email, oldPassword) == true) {
    c.UpdateId(user.ID, bson.M{"$set": bson.M{"password": *newPassword}})
  } else {
    return errors.New("Wrong password provided")
  }
  return nil
}

func UpdateUser(userToUpdate *models.UserModel, mainUser *models.UserModel) {
  c := instance.Session.DB("slidare").C("users")

  userToUpdate.Token = mainUser.Token
  userToUpdate.ID = mainUser.ID
  userToUpdate.FBUserID = mainUser.FBUserID
  if (len(userToUpdate.Firstname) == 0) {
    userToUpdate.Firstname = mainUser.Firstname
  } else {
    userToUpdate.Firstname = userToUpdate.Firstname
  }
  if (len(userToUpdate.LastName) == 0) {
    userToUpdate.LastName = mainUser.LastName
  } else {
    userToUpdate.LastName = userToUpdate.LastName
  }
  if (len(userToUpdate.Email) == 0) {
    userToUpdate.Email = mainUser.Email
  } else {
    userToUpdate.Email = userToUpdate.Email
  }
  if (len(userToUpdate.FBToken) == 0) {
    userToUpdate.FBToken = mainUser.FBToken
  } else {
    userToUpdate.FBToken = userToUpdate.FBToken
  }
  err := c.UpdateId(mainUser.ID, &userToUpdate)
  if err != nil {
    log.Fatal("err", err)
  }
}

func DoesEmailExistInDB(email string) bool {
  c := instance.Session.DB("slidare").C("users")
  result := models.UserModel{}
  err := c.Find(bson.M{"email": email}).One(&result)
  if err == nil {
    return true
  }
  return false
}

func GetUserFromToken(token *string) *models.UserModel {
  c := instance.Session.DB("slidare").C("users")
  result := models.UserModel{}
  err := c.Find(bson.M{"token": *token}).One(&result)
  if err == nil {
    return &result
  }
  return nil
}

func GetUserFromId(id *string) *models.UserModel {
  c := instance.Session.DB("slidare").C("users")
  result := models.UserModel{}
  err := c.Find(bson.M{"_id": id}).One(&result)
  if err == nil {
    return &result
  }
  return nil
}

func AddContactToUser(userEmail *string, contactID *string) {
    c := instance.Session.DB("slidare").C("users")
    result := models.UserModel{}
    c.Find(bson.M{"email": *userEmail}).One(&result)
    result.Contacts = append(result.Contacts, *contactID);
    c.UpdateId(result.ID, &result)
}

func IsExistingGroup(groupName *string, userId *string) bool {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  err := c.Find(bson.M{"name": *groupName, "owner": *userId}).One(&result)

  return err == nil
}

func CreateGroup(groupModel *models.GroupModel) {
  c := instance.Session.DB("slidare").C("groups")
  c.Insert(groupModel)
}

func ValidateUserPassword(userEmail *string, userPassword *string) bool {
  c := instance.Session.DB("slidare").C("users")
  var result models.UserModel
  err := c.Find(bson.M{"email": *userEmail, "password": *userPassword}).One(&result)
  return err == nil
}
