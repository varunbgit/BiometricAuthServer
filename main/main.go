package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-webauthn/webauthn/User"
	"github.com/go-webauthn/webauthn/main/helper/db"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	webAuthn *webauthn.WebAuthn
	err      error
	user     webauthn.User
	session  *webauthn.SessionData
	userID   string
	userName string
	wconfig  webauthn.Config
)

// Your initialization function
func main() {
	wconfig = webauthn.Config{
		RPDisplayName: "localhost:5173", // Display Name for your site
		RPID:          "localhost",      // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:5173"},
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			AuthenticatorAttachment: "platform",
		},
		AttestationPreference: protocol.PreferDirectAttestation, // The origin URLs allowed for WebAuthn requests
	}
	if webAuthn, err = webauthn.New(&wconfig); err != nil {
		fmt.Println(err)
	}

	userID = "0001"
	userName = "bansal"

	user = User.NewUser(userName, userID) //

	db.Redis = map[string]webauthn.SessionData{}
	db.CdpDB = map[string]webauthn.Credential{}

	http.Handle("/register", corsMiddleware(http.HandlerFunc(RegisterHandler)))
	//http.Handle("/verification", corsMiddleware(http.HandlerFunc(VerificationHandler)))
	http.Handle("/save", corsMiddleware(http.HandlerFunc(SaveHandler)))
	http.Handle("/create_verify_options", corsMiddleware(http.HandlerFunc(BeginLogin)))
	http.Handle("/complete_verification", corsMiddleware(http.HandlerFunc(FinishLogin)))

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

	//	==================================================================================	=======================================================================
	// Now lets start with the Verification of credentials received.

	parsedCredentialData, err1 := protocol.ParseCredentialCreationResponse(r)
	if err1 != nil {
		e := fmt.Errorf("Error while parsing client Response: %s", err)
		fmt.Println(e)
		http.Error(w, "Failed to unmarshal clientDataJSON", http.StatusBadRequest)
		return
	}

	//todo get the session details from Redis
	sessionFromDB := db.Redis[userID]

	//todo this needs to be removed
	//parsedCredentialData.Response.CollectedClientData.Challenge = sessionFromDB.Challenge
	//webAuthn.Config.RPID = "http://localhost:5173"
	//todo above line is just for testing

	credential, err2 := webAuthn.CreateCredential(user, sessionFromDB, parsedCredentialData)
	if err2 != nil {
		fmt.Println("There is some error verifying Credential after registration: ", err2.Error())
	}

	fmt.Println(userID)
	fmt.Println(*credential)

	db.CdpDB[userID] = *credential
	w.WriteHeader(http.StatusCreated)

	fmt.Println("Hurray! verification of registration credentials successful")
	write, err := w.Write([]byte("Successful"))
	if err != nil {
		return
	}
	fmt.Println(write)
	return
}

func BeginLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	options, session1, err1 := webAuthn.BeginLogin(user)
	if err1 != nil {
		// Handle Error and return.
		fmt.Println("Some error occurred in starting the login")
		return
	}

	fmt.Println("options generated are: ", options.Response)

	jsonData, err2 := json.Marshal(options)
	if err2 != nil {
		fmt.Println("There is some error in unmarshalling")
	}

	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, "error while writing to json", 500)
		return
	}
	w.WriteHeader(http.StatusCreated)

	//storing session data
	key := userID
	db.Redis[key] = *session1

	w.WriteHeader(http.StatusOK)

}

func FinishLogin(w http.ResponseWriter, r *http.Request) {
	//	//user := datastore.GetUser() // Get the user
	//	user := User.NewUser("Varun,", "8955")
	//	//
	//	//// Get the session data stored from the function above
	session := db.Redis[userID]
	//
	//

	credential, err := webAuthn.FinishLogin(user, session, r)

	if err != nil {
		//Handle Error and return.
		http.Error(w, "error while getting session from Redis", 500)
	}

	fmt.Println("Credential is: ", credential)
	//	//	return
	//	//}
	//	//
	//	//// Handle credential.Authenticator.CloneWarning
	//	//
	//	//// If login was successful, update the credential object
	//	//// Pseudocode to update the user credential.
	//	//user.UpdateCredential(credential)
	//	//datastore.SaveUser(user)
	//	//
	w.WriteHeader(http.StatusOK)
	write, err := w.Write([]byte("Login Successful"))
	if err != nil {
		return
	}
	fmt.Println(write)
	return
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var createOptions *protocol.CredentialCreation
	createOptions, session, err = BeginRegistration()
	fmt.Println("Printing the Create Options:  ", createOptions)
	fmt.Println("Printing the Create Options challenge :  ", createOptions.Response.Challenge)
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
	db.Redis[key] = *session
}

func BeginRegistration() (*protocol.CredentialCreation, *webauthn.SessionData, error) {

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
