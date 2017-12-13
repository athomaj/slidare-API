package models

type UserModel struct {
  Firstname string `json:"first_name,omitempty" bson:"first_name,omitempty"`
  LastName string `json:"last_name,omitempty" bson:"last_name,omitempty"`
  UserName string `json:"username,omitempty" bson:"username,omitempty"`
  Password string `json:"password,omitempty" bson:"password,omitempty"`
  Email string `json:"email,omitempty" bson:"email,omitempty"`
  ProfilePictureURL string `json:"profile_picture_url,omitempty" bson:"profile_picture_url,omitempty"`
  PhoneNumber string `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
  Token string `json:"token,omitempty" bson:"token,omitempty"`
  ID string `json:"id" bson:"_id,omitempty" `
  FBToken string `json:"fb_token,omitempty" bson:"fb_token,omitempty"`
  FBUserID string `json:"fb_user_id,omitempty" bson:"fb_user_id,omitempty"`
  Contacts []string `json:"contacts,omitempty" bson:"contacts,omitempty"`
  FileUrls []string `json:"file_urls,omitempty" bson:"file_urls,omitempty"`
  Senders []string `json:"senders,omitempty" bson:"senders,omitempty"`
}
