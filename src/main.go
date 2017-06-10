package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "encoding/json"
    "controllers"
    "userControllers"
    "github.com/codegangsta/negroni"
    "core/authentication"
    "mydb"
    "github.com/antigloss/go/logger"
  //  "github.com/davecgh/go-spew/spew"
)

type AuthModel struct
{
  Username, Password string
}

func main() {
    router := mux.NewRouter().StrictSlash(true)

    logger.Init("./log", // specify the directory to save the logfiles
                400, // maximum logfiles allowed under the specified log directory
                20, // number of logfiles to delete when number of logfiles exceeds the configured limit
                100, // maximum size of a logfile in MB
                false) // whether logs with Trace level are written down

    defer func() { //catch or finally
      if err := recover(); err != nil { //catch
        logger.Info("initDatabaseError: %s", err)
        log.Println("initDatabaseError:", err)
        return ;
      }
    }()
    database.Init()

    logger.Info("initDatabase success")
    log.Println("initDatabase success")


//    logger.Info("Failed to find player! uid=%d plid=%d cmd=%s xxx=%d", 1234, 678942, "getplayer", 102020101)

    var token string
    router.Handle("/getUserContacts", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(getUserContacts),
        )).Methods("POST")
    router.Handle("/userContacts", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUserContacts(&token)),
        )).Methods("GET")
    router.Handle("/userContact/{contact_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUserContact(&token)),
        )).Methods("GET")
    router.Handle("/addContact", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.AddContact(&token)),
        )).Methods("POST")
    router.Handle("/removeContact/{contact_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.RemoveContact(&token)),
        )).Methods("DELETE")
    router.Handle("/removeContactByEmail/{contact_email}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.RemoveContactByEmail(&token)),
        )).Methods("DELETE")
    router.Handle("/fetchGroups", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchGroups(&token)),
        )).Methods("GET")
    router.Handle("/createGroup", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.CreateGroup(&token)),
        )).Methods("POST")
    router.Handle("/removeGroup/{group_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.RemoveGroup(&token)),
        )).Methods("DELETE")
    router.Handle("/addToGroup/{group_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.AddToGroup(&token)),
        )).Methods("PUT")
    router.Handle("/leaveGroup/{group_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUserContact(&token)),
        )).Methods("PUT")
    router.Handle("/removeFromGroup/{group_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.RemoveFromGroup(&token)),
        )).Methods("PUT")
    router.Handle("/renameGroup/{group_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUserContact(&token)),
        )).Methods("PUT")
    router.Handle("/acceptGroupInvite/{group_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUserContact(&token)),
        )).Methods("PUT")
    router.Handle("/refuseGroupInvite/{group_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUserContact(&token)),
        )).Methods("PUT")
    router.Handle("/refuseContactInvite/{contact_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUserContact(&token)),
        )).Methods("PUT")
    router.Handle("/acceptContactInvite/{contact_identifier}", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUserContact(&token)),
        )).Methods("PUT")
    router.Handle("/fetchUser", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.FetchUser(&token)),
        )).Methods("GET")
    // router.Handle("/updateUser", negroni.New(
    //     negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
    //     negroni.HandlerFunc(userControllers.UpdateUser(&token)),
    //     )).Methods("POST")
    router.Handle("/updateUserName", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.UpdateUserName(&token)),
        )).Methods("POST")
    router.Handle("/updateUserEmail", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.UpdateUserEmail(&token)),
        )).Methods("POST")
    router.Handle("/updateUserPicture", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.UpdateUserPicture(&token)),
        )).Methods("POST")
    router.Handle("/updateUserPassword", negroni.New(
        negroni.HandlerFunc(authentication.RequireTokenAuthentication(&token)),
        negroni.HandlerFunc(userControllers.UpdateUserPassword(&token)),
        )).Methods("POST")

    router.HandleFunc("/createUser", userControllers.CreateUser).Methods("POST")
    router.HandleFunc("/loginUser", userControllers.LoginUser).Methods("POST")
    router.HandleFunc("/token-auth", controllers.Login).Methods("POST")

    logger.Info(http.ListenAndServe(":50000", router).Error())
}

func getUserContacts(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
}

func Index(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var logs AuthModel
  err := decoder.Decode(&logs)
  if err != nil {
    log.Println("Error while decoding Request")
  }
  if len(logs.Username) != 0 && len(logs.Password) != 0 {
    log.Printf("%s %s\n", logs.Username, logs.Password);
  } else {
    fmt.Fprintf(w, "Wrong parameters on authentification request");
    log.Printf("Wrong parameters");
  }
//  fmt.Fprintf(w, "Hellolol, %q", html.EscapeString(r.URL.Path))
}
