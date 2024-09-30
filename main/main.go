package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-webauthn/webauthn/User"
	"github.com/go-webauthn/webauthn/main/helper"
	"github.com/go-webauthn/webauthn/protocol"
	webauthnv1 "github.com/go-webauthn/webauthn/rpc/main/v1/webauthn"
	"github.com/go-webauthn/webauthn/webauthn"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	webAuthn *webauthn.WebAuthn
	err      error
	redis    map[string]interface{}
)

// Your initialization function
func main() {
	wconfig := &webauthn.Config{
		RPDisplayName: "localhost", // Display Name for your site
		RPID:          "localhost", // Generally the FQDN for your site
		RPOrigins:     []string{"localhost"},
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			AuthenticatorAttachment: "platform",
		}, // The origin URLs allowed for WebAuthn requests
	}
	if webAuthn, err = webauthn.New(wconfig); err != nil {
		fmt.Println(err)
	}

	redis = map[string]interface{}{}

	http.Handle("/register", corsMiddleware(http.HandlerFunc(RegisterHandler)))
	//http.Handle("/verification", corsMiddleware(http.HandlerFunc(VerificationHandler)))
	http.Handle("/save", corsMiddleware(http.HandlerFunc(SaveHandler)))
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Printing the request body: ", r.Body)
	// Parse the request body
	//var saveRequest PublicKeyCredential
	//err := json.NewDecoder(body).Decode(&saveRequest)

	var credential webauthnv1.PublicKeyCredential
	err = protojson.Unmarshal(body, &credential)

	if err != nil {
		fmt.Println("the error is ", err)
		http.Error(w, "err unmarshalling into proto ,Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received data: ", credential)
	fmt.Println("Printing the request body ID: ", credential.Id)
	fmt.Println("Printing the request body raw Id: ", credential.RawId)
	fmt.Println("Printing the request body raw Id direct typecasting: ", string(credential.RawId))
	//fmt.Println("Printing the request body raw Id hex format: ", hex.EncodeToString(credential.RawId))
	//fmt.Println("Printing the request body raw Id base64.std Endcoding format: ", base64.StdEncoding.EncodeToString(credential.RawId))

	fmt.Println("Printing the request body Type: ", credential.Type)
	fmt.Println("Printing the request body response.AttestationObject: ", credential.Response.AttestationObject)
	//fmt.Println("Printing the request body AuthenticatorData base64.std Endcoding format: ", base64.StdEncoding.EncodeToString(credential.Response.AuthenticatorData))
	fmt.Println("Printing the request body response.ClientDataJSON: ", credential.Response.ClientDataJson)
	fmt.Println("Printing the request body Client Data Json base64.std Endcoding format: ", base64.StdEncoding.EncodeToString(credential.Response.ClientDataJson))
	//fmt.Println("Printing the request body response.Signature: ", credential.Response.Signature)
	//fmt.Println("Printing the request body Signature base64.std Endcoding format: ", base64.StdEncoding.EncodeToString(credential.Response.Signature))
	//fmt.Println("Printing the request body response.UserHandle: ", credential.Response.UserHandle)
	//fmt.Println("Printing the request body User handler base64.std Endcoding format: ", base64.StdEncoding.EncodeToString(credential.Response.UserHandle))

	//// Step 1: Decode the Base64 attestationObject
	//decoded, err := decodeBase64AttestationObject(credential.Response.AuthenticatorData)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to decode attestationObject: %s", err), http.StatusBadRequest)
	//	return
	//}

	// Step 2: Parse the CBOR attestation object
	attestation, err := helper.ParseAttestationObject(credential.Response.AttestationObject)
	if err != nil {
		fmt.Println("The error while parsing attestionObject", err)
		http.Error(w, fmt.Sprintf("Failed to parse attestationObject: %s", err), http.StatusBadRequest)
		return
	}

	// Step 3: Process and validate the authData
	err = helper.ProcessAuthData(attestation.AuthData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process authData: %s", err), http.StatusBadRequest)
		return
	}

	credId, err := helper.ExtractCredentialID(attestation.AuthData)

	fmt.Println("the credential ID is: ", credId)

	//================client data json========
	var clientDataJsonObject helper.ClientDataJSON

	err = json.Unmarshal(credential.Response.ClientDataJson, &clientDataJsonObject)
	if err != nil {
		http.Error(w, "Failed to unmarshal clientDataJSON", http.StatusBadRequest)
		return
	}

	fmt.Println("The sent challenge was: ", clientDataJsonObject.Challenge)
	fmt.Println("The sent origin was: ", clientDataJsonObject.Origin)
	fmt.Println("The sent type was: ", clientDataJsonObject.Type)
	// Respond with a success message

}

func VerificationHandler(w http.ResponseWriter, r *http.Request) {

}

func BeginLogin(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Find the user
	//
	//options, session, err := webAuthn.BeginLogin(user)
	//if err != nil {
	//	// Handle Error and return.
	//
	//	return
	//}

	// todo store the session values
	//datastore.SaveSession(session)
	//
	//JSONResponse(w, options, http.StatusOK) // todo return the options generated
	// options.publicKey contain our registration options
}

func FinishLogin(w http.ResponseWriter, r *http.Request) {
	//user := datastore.GetUser() // Get the user
	//user := User.NewUser("Varun,", "8955")
	//
	//// Get the session data stored from the function above
	//session := redis[user.ID]
	//
	//credential, err := webAuthn.FinishLogin(user, session.(webauthn.SessionData), r)
	//if err != nil {
	//	// Handle Error and return.
	//
	//	return
	//}
	//
	//// Handle credential.Authenticator.CloneWarning
	//
	//// If login was successful, update the credential object
	//// Pseudocode to update the user credential.
	//user.UpdateCredential(credential)
	//datastore.SaveUser(user)
	//
	//JSONResponse(w, "Login Success", http.StatusOK)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	createOptions, session, err := BeginRegistration()
	fmt.Println("Printing the Create Options:  ", createOptions)
	jsonData, err := json.Marshal(createOptions.Response)
	if err != nil {
		http.Error(w, "Failed to serialize createOptions", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "error while writing to json", 500)
		return
	}
	w.WriteHeader(http.StatusCreated)

	//storing session data
	key := string(session.UserID)
	redis[key] = session
}

func BeginRegistration() (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	user := User.NewUser("Varun,", "8955") // Find or create the new user
	options, session, err := webAuthn.BeginRegistration(user)
	// handle errors if present

	// todo store the sessionData values
	//JSONResponse(w, options, http.StatusOK) // return the options generated
	// options.publicKey contain our registration options

	fmt.Println("options are: ", options)
	fmt.Println("session is:", session)
	fmt.Println("err is: ", err)
	fmt.Println(options.Response.Challenge)

	str := fmt.Sprintf("%s", options.Response.Challenge)
	fmt.Println("challenge in string format is: ", str)
	return options, session, err
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                            // Allow all origins, adjust as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")          // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allowed headers

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r) // Call the next handler
	})
}
