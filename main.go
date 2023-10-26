package main

import (
	"encoding/json"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/home", verifyJWT(handlePage))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port:8080", err)

	}

}

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

//`json:"status"`,`json:"info"` , parts are called struct tags .
//they tell the golang  that when the struct is encoded or decoded from json format ,
//it should use the field names "status" and "info" instead of "Status" and "Info".

func handlePage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	//this line sets the content type of the response to `application/json`.
	//this means that the data that the server sends back to the client will be in JSON format
	var message Message
	err := json.NewDecoder(request.Body).Decode(&message)
	//this line reads the JSON data from the request and decodes it into the 'message' variable.
	//the json.NewDecoder function creates a new decoder that reads from the request.Body(Which is where the json data is stored),
	//and the Decode method reads teh json data and puts it into the 'message ' variable.
	//The & symbol is used to pass a pointer to the  'message' variable , which allows the Decode method to modify the variable directly
	if err != nil {
		return
	}
	message.Info += "hello i am the man in the middle "
	//This line checks whether there was an error decoding the JSON data. If there was an
	//error, the function stops running and returns nothing. If there wasn't an error, the
	//function continues running.
	err = json.NewEncoder(writer).Encode(message)
	// thsi lines takes the message variable (which now contains the json data from the request ),
	//and encodes it back into json format using the json.NewEncoder function  and the  Encode method .
	//This encoded json data is then written to the writer object, which sends the data to the client
	if err != nil {
		return
	}
	//This line checks whether there was an error encoding the JSON data. If there was an error, the function stops
	// running and returns nothing. If there wasn't an error, the function finishes running and the JSON data is sent back to the client.

}

func generateJwt()(string, error){
	var sampleSecretKey=[]byte("SecretYouShouldHide")
token := jwt.New(jwt.SigningMethodEdDSA)
claims:=token.Claims.(jwt.MapClaims)
//token.Claims is used to modify the jwt .
//we will able to retrieve the claims when attempting to verify the jwt .
claims["exp"] = time.Now().Add(10*time.Minute)
claims["authorized"] = true
claims["user"] = "username"
tokenString,err:=token.SignedString(sampleSecretKey)
if err != nil {
	return "",err
}
return tokenString , nil 
}

//If there’s an error generating the JWT, the function returns an empty string and the error. 
//If there are no errors, the function returns the JWT string and the nil type.




func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {


	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Token"] != nil {
			token, err := jwt.Parse(request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
				   writer.WriteHeader(http.StatusUnauthorized)
				   _, err := writer.Write([]byte("You're Unauthorized!"))
				   if err != nil {
					  return nil, err
	
				   }

				   
				}
				return "", nil
	
			 })
			 if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
			   if err2 != nil {
					   return
				 }
 }
 if token.Valid {
	endpointHandler(writer, request)
	  } else {
			  writer.WriteHeader(http.StatusUnauthorized)
			  _, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
			  if err != nil {
					  return
			  }
}
		}
	})
}

//its a middleware that takes in the handler function for the 
//request you want to verify . 
//The handler function uses the token parameter from the request 
//header to verify the request and respond based on the status

/*The verifyJWT function returns the handler function passed 
in as a parameter if the request is authorized.
*/