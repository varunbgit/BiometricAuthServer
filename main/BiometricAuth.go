package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type RelyingParty struct {
	Name string
	Id   string
}

type User struct {
	Id          string
	Name        string
	DisplayName string
}

type PublicCredParam struct {
	Alg  int
	Type string
}

type AuthenticatorSelection struct {
	AuthenticatorAttachment string
}

// CreateOptions represents a user in the system
type CreateOptions struct {
	Challenge              string
	Rp                     RelyingParty
	User                   User
	PublicKeyParams        []PublicCredParam
	AuthenticatorSelection AuthenticatorSelection
	Timeout                int
	Attestation            string
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	publicCredParam := PublicCredParam{
		Alg:  -7,
		Type: "public-key",
	}

	var publicCredParams []PublicCredParam
	publicCredParams = append(publicCredParams, publicCredParam)

	createOptions := CreateOptions{
		Challenge: uuid.New().String(),
		Rp: RelyingParty{
			Name: "Rzp",
			Id:   "Rzp.com",
		},
		User: User{
			Id:          "user_rzp_6088",
			Name:        "8955496900",
			DisplayName: "8955496900",
		},
		PublicKeyParams:        publicCredParams,
		AuthenticatorSelection: AuthenticatorSelection{AuthenticatorAttachment: "cross-platform"},
		Timeout:                60000,
		Attestation:            "direct",
	}

	jsonData, err := json.Marshal(createOptions)
	if err != nil {
		http.Error(w, "Failed to serialize createOptions", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
	w.WriteHeader(http.StatusCreated)
}

// VerificationHandler handles user verification
func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User verified successfully"))
}

// SaveHandler handles the /save endpoint
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type AuthenticatorAssertionResponse struct {
		AuthenticatorData string `json:"authenticatorData"`
		ClientDataJSON    string `json:"clientDataJSON"`
		Signature         string `json:"signature"`
		UserHandle        string `json:"userHandle,omitempty"` // May be optional
	}
	// Define a struct to hold the incoming data
	type PublicKeyCredential struct {
		ID       string
		RawId    string
		Type     string
		response AuthenticatorAssertionResponse
	}

	fmt.Println("Printing the request body: ", r.Body)
	// Parse the request body
	var saveRequest PublicKeyCredential
	err := json.NewDecoder(r.Body).Decode(&saveRequest)
	if err != nil {
		fmt.Println("the error is ", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rawIdd, err := base64.StdEncoding.DecodeString(saveRequest.RawId)
	if err != nil {
		http.Error(w, "Failed to decode authenticatorData", http.StatusBadRequest)
		return
	}
	fmt.Println("The raw id after data after decoding is ", rawIdd)

	// Decode base64-encoded fields
	//authData, err := base64.StdEncoding.DecodeString(saveRequest.response.AuthenticatorData)
	//if err != nil {
	//	http.Error(w, "Failed to decode authenticatorData", http.StatusBadRequest)
	//	return
	//}
	//fmt.Println("The auth data after decoding is ", authData)
	//
	//clientData, err := base64.StdEncoding.DecodeString(saveRequest.response.ClientDataJSON)
	//if err != nil {
	//	http.Error(w, "Failed to decode clientDataJSON", http.StatusBadRequest)
	//	return
	//}
	//fmt.Println("The client data after decoding is ", clientData)

	log.Printf("Received data: ", saveRequest)
	fmt.Println("Printing the request body ID: ", saveRequest.ID)
	fmt.Println("Printing the request body raw Id: ", saveRequest.RawId)
	fmt.Println("Printing the request body Type: ", saveRequest.Type)
	fmt.Println("Printing the request body response.AuthenticatorData: ", saveRequest.response.AuthenticatorData)
	fmt.Println("Printing the request body response.ClientDataJSON: ", saveRequest.response.ClientDataJSON)
	fmt.Println("Printing the request body response.Signature: ", saveRequest.response.Signature)
	fmt.Println("Printing the request body response.UserHandle: ", saveRequest.response.UserHandle)

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "id": saveRequest.ID})
}

func main() {

	http.Handle("/register", corsMiddleware(http.HandlerFunc(RegisterHandler)))
	http.Handle("/verification", corsMiddleware(http.HandlerFunc(VerificationHandler)))
	http.Handle("/save", corsMiddleware(http.HandlerFunc(SaveHandler)))
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
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
