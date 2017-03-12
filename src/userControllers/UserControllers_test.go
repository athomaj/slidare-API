package userControllers

import ("testing"
        "models"
        "mydb")

func TestEmptyUserModelForCreation(t *testing.T) {
  var user models.UserModel

  valid, err := checkUserInformationsForCreation(user)
  if (valid == true) {
    t.Error("Check should return false: ", err)
  }
}

func TestEmptyFieldUserModelForCreation(t *testing.T) {
  var user models.UserModel

  user.Firstname = ""
  user.LastName = ""
  user.Password = ""
  user.Email = ""
  valid, err := checkUserInformationsForCreation(user)
  if (valid == true) {
    t.Error("Check should return false: ", err)
  }
}

func TestValidUserCreationWithPassword(t *testing.T) {
  database.Init()
  var user models.UserModel

  user.Firstname = "Julien"
  user.LastName = "Athomas"
  user.Password = "1234"
  user.Email = "julien.athomas45@epitech.eu"
  valid, err := checkUserInformationsForCreation(user)
  if (valid == false) {
    t.Error("Check should return false: ", err)
  }
}

func TestCreationWithExistingEmail(t *testing.T) {
  database.Init()
  var user models.UserModel

  user.Firstname = "Julien"
  user.LastName = "Athomas"
  user.Password = "1234"
  user.Email = "julien.athomas@epitech.eu"
  valid, err := checkUserInformationsForCreation(user)
  if (valid == true) {
    t.Error("Check should return false: ", err)
  }
}

func TestCreationWithoutEmail(t *testing.T) {
  database.Init()
  var user models.UserModel

  user.Firstname = "Julien"
  user.LastName = "Athomas"
  user.Password = "1234"
  user.Email = ""
  valid, err := checkUserInformationsForCreation(user)
  if (valid == true) {
    t.Error("Check should return false: ", err)
  }
}

func TestCreationWithoutPassword(t *testing.T) {
  database.Init()
  var user models.UserModel

  user.Firstname = "Julien"
  user.LastName = "Athomas"
  user.Password = ""
  user.Email = "julien.athomas@epitech.eu"
  valid, err := checkUserInformationsForCreation(user)
  if (valid == true) {
    t.Error("Check should return false: ", err)
  }
}

func TestCreationWithoutFirstname(t *testing.T) {
  database.Init()
  var user models.UserModel

  user.Firstname = ""
  user.LastName = "Athomas"
  user.Password = "1234"
  user.Email = "julien.athomas@epitech.eu"
  valid, err := checkUserInformationsForCreation(user)
  if (valid == true) {
    t.Error("Check should return false: ", err)
  }
}

func TestCreationWithoutLastname(t *testing.T) {
  database.Init()
  var user models.UserModel

  user.Firstname = "Julien"
  user.LastName = ""
  user.Password = "1234"
  user.Email = "julien.athomas@epitech.eu"
  valid, err := checkUserInformationsForCreation(user)
  if (valid == true) {
    t.Error("Check should return false: ", err)
  }
}
