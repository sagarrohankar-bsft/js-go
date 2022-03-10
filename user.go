package jsgo

import "encoding/json"

var userjson = `{
	"name": "John Doe",
	"email": "jhon@company.com",
	"dob": "1/2/75",
	"address": "new york"
}`

func createUser() map[string]interface{} {
	user := map[string]interface{}{}
	err := json.Unmarshal([]byte(userjson), &user)
	if err != nil {
		panic(err)
	}
	return user
}
