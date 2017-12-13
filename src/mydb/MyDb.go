package database

import (
    "sync"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "log"
    "models"
    "fmt"
    "errors"
    "github.com/antigloss/go/logger"
    "math/rand"
    "strconv"
    "net/smtp"
)

var instance mongodb
var once sync.Once

type mongodb struct {
  Session *mgo.Session
}

func Init() {
  instance = mongodb{}
  session, err := mgo.Dial("mongodb://localhost:27018")
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
    logger.Info("Email already Exist")
    return ;
  }
  c := instance.Session.DB("slidare").C("users")
  err := c.Insert(&user)
  if err != nil {
    log.Fatal(err)
    logger.Info("CreateUserError: %s", err)
    return ;
  }
  logger.Info("User Created")
}

func UpdateUserContacts(user *models.UserModel) {
  c := instance.Session.DB("slidare").C("users")
  c.UpdateId(user.ID, bson.M{"$set": bson.M{"contacts": user.Contacts}})
  logger.Info("User Contacts Updated")
}

func UpdateUserName(userName *string, user *models.UserModel) {
  c := instance.Session.DB("slidare").C("users")
  c.UpdateId(user.ID, bson.M{"$set": bson.M{"username": *userName}})
  logger.Info("UserName Updated")
}

func UpdateUserPicture(userPicture *string, user *models.UserModel) {
  c := instance.Session.DB("slidare").C("users")
  c.UpdateId(user.ID, bson.M{"$set": bson.M{"profile_picture_url": *userPicture}})
  logger.Info("UserName Updated")
}

func UpdateUserEmail(userEmail *string, user *models.UserModel) error {
  c := instance.Session.DB("slidare").C("users")
  var result models.UserModel
  err := c.Find(bson.M{"email": userEmail}).One(&result)
  if err != nil {
    c.UpdateId(user.ID, bson.M{"$set": bson.M{"email": *userEmail}})
    logger.Info("UserEmail Updated")
  } else {
    logger.Info("UpdateUserEmail failed: %s", err)
  }
  return err
}

func ResetUserPassword(userEmail *string) error {
  c := instance.Session.DB("slidare").C("users")
  var result models.UserModel
  err := c.Find(bson.M{"email": userEmail}).One(&result)
  if err == nil {
    newPwd := strconv.Itoa(rand.Int())
    c.UpdateId(result.ID, bson.M{"$set": bson.M{"password": newPwd}})
    auth := smtp.PlainAuth(
  		"",
  		"julien.slidare@gmail.com",
  		"slidaredev",
  		"smtp.gmail.com",
  	)
  	// Connect to the server, authenticate, set the sender and recipient,
  	// and send the email all in one step.
    msg := []byte("To: " + *userEmail + "\r\n" +
		"Subject: new slidare password!\r\n" +
		"\r\n" +
		"Here is your new password:" + newPwd + "\r\n")
  	err := smtp.SendMail(
  		"smtp.gmail.com:587",
  		auth,
  		"julien.slidare@gmail.com",
  		[]string{*userEmail},
  		[]byte(msg),
  	)
  	if err != nil {
      logger.Info("Send Email Error: %s", err)
  	}
    logger.Info("UserPassword Reseted")
  } else {
    logger.Info("UpdateUserEmail failed: %s", err)
  }
  return err
}

func UpdateUserPassword(oldPassword *string, newPassword *string, user *models.UserModel) error {
  c := instance.Session.DB("slidare").C("users")
  if (ValidateUserPassword(&user.Email, oldPassword) == true) {
    c.UpdateId(user.ID, bson.M{"$set": bson.M{"password": *newPassword}})
    logger.Info("UserPassword Updated")
  } else {
    logger.Info("Wrong password provided")
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

func AddFileToUser(userEmail *string, fileUrl string, sender string) {
    c := instance.Session.DB("slidare").C("users")
    result := models.UserModel{}
    c.Find(bson.M{"email": *userEmail}).One(&result)
    result.FileUrls = append(result.FileUrls, fileUrl);
    result.Senders = append(result.Senders, sender);
    c.UpdateId(result.ID, &result)
}

func IsExistingGroup(groupName *string, userId *string) bool {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  err := c.Find(bson.M{"name": *groupName, "owner": *userId}).One(&result)

  return err == nil
}

func IsExistingGroupId(groupId *string, userId *string) bool {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  err := c.Find(bson.M{"_id": *groupId, "owner": *userId}).One(&result)

  return err == nil
}

func IsExistingGroupById(groupId *string) bool {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  err := c.Find(bson.M{"_id": *groupId}).One(&result)

  return err == nil
}


func FetchGroupsFromUser(userId *string, userEmail *string) []models.GroupModel {
  c := instance.Session.DB("slidare").C("groups")
  result := []models.GroupModel{}
  c.Find(bson.M{
    "owner": *userId,}).All(&result)

  result2 := []models.GroupModel{}
  c.Find(bson.M{
    "users": *userEmail,}).All(&result2)

    result = append(result, result2...)

  return result
  // for _,tmpUser := range result.Users {
  //   if (*userToAdd == tmpUser) {
  //     return "User Already in Group";
  //   }
  // }
  // result.Users = append(result.Users, *userToAdd)
  // c.UpdateId(result.ID, &result);
  // return "";
}

func AddToGroup(groupName *string, userId *string, userToAdd *string) string {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  c.Find(bson.M{"name": *groupName, "owner": *userId}).One(&result)
  for _,tmpUser := range result.Users {
    if (*userToAdd == tmpUser) {
      return "User Already in Group";
    }
  }
  result.Users = append(result.Users, *userToAdd)
  c.UpdateId(result.ID, &result);
  return "";
}

func RemoveFromGroup(groupName *string, userId *string, userToAdd *string) string {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  c.Find(bson.M{"name": *groupName, "owner": *userId}).One(&result)
  user := GetUserFromId(userId);

  if (result.Owner != *userId) {
    return "You are not the group owner";
  }
  if (user.Email == *userToAdd) {
    return "You are the owner of the group, you cant leave the group";
  }

  for idx,tmpUser := range result.Users {
    if (*userToAdd == tmpUser) {
//      result.Users = append(result.Users, *userToAdd)
      result.Users = append(result.Users[:idx], result.Users[idx+1:]...)
      c.UpdateId(result.ID, &result);
      return "";
    }
  }
  return "User not in Group";
}

func LeaveGroup(groupId *string, userId *string, userEmail *string) string {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  c.Find(bson.M{"_id": *groupId}).One(&result)
  user := GetUserFromId(userId);

  if (result.Owner == *userId) {
    return "You are the group owner";
  }
  // if (user.Email == *userToAdd) {
  //   return "You are the owner of the group, you cant leave the group";
  // }

  for idx,tmpUser := range result.Users {
    if (*userId == tmpUser.ID || *userEmail == tmpUser.Email) {
      result.Users = append(result.Users[:idx], result.Users[idx+1:]...)
      c.UpdateId(result.ID, &result);
      return "";
    }
  }
  return "User not in Group";
}

func UpdateGroupName(groupName *string, userId *string, newGroupName *string) string {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  c.Find(bson.M{"name": *groupName, "owner": *userId}).One(&result)
  result.Name = *newGroupName
  c.UpdateId(result.ID, &result);
  return "";
}

func UpdateGroupNameById(groupId *string, userId *string, newGroupName *string) string {
  c := instance.Session.DB("slidare").C("groups")
  result := models.GroupModel{}
  c.Find(bson.M{"_id": *groupId, "owner": *userId}).One(&result)
  result.Name = *newGroupName
  c.UpdateId(result.ID, &result);
  return "";
}


func CreateGroup(groupModel *models.GroupModel) {
  c := instance.Session.DB("slidare").C("groups")
  c.Insert(groupModel)
}

func DeleteGroup(groupName *string, userId *string) {
  c := instance.Session.DB("slidare").C("groups")
  err := c.Remove(bson.M{"name": *groupName, "owner": *userId})
  log.Println(err)
}

func ValidateUserPassword(userEmail *string, userPassword *string) bool {
  c := instance.Session.DB("slidare").C("users")
  var result models.UserModel
  err := c.Find(bson.M{"email": *userEmail, "password": *userPassword}).One(&result)
  return err == nil
}
