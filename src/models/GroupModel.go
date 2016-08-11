package models

type GroupModel struct {
  Name string `json:"name,omitempty" bson:"name,omitempty"`
  Owner string `json:"owner,omitempty" bson:"owner,omitempty"`
  ID string `json:"id" bson:"_id,omitempty" `
  Users []string `json:"users,omitempty" bson:"users,omitempty"`
}
