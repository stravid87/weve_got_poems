package controller

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"

	"proxy_slave/config"
	"proxy_slave/models"

	"github.com/go-playground/validator/v10"

	pb "proxy_slave/service"
)

// PingHandler handles ping requests
func Ping(w http.ResponseWriter, r *http.Request) {
	// Send response to client
	_, err := w.Write([]byte("ping successful"))
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("SERVICE_PROVIDER_PORT")
	// Make request to the content server
	resp, err := http.Get("http://localhost:" + port + "/") // + "/image" + "?id=" + req.Choice)
	if err != nil {
		log.Printf("failed to get '/': %v", err)
		return
	}
	defer resp.Body.Close()

	// Convert the response body to a string
	RespBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return
	}

	// Send the response back to the WASM module
	_, err = w.Write(RespBodyBytes)
	if err != nil {
		log.Printf("failed to send response: %v", err)
		return
	}
}

// SimpleProxy just sends the request to the backend
func PingServiceProvider(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("SERVICE_PROVIDER_PORT")
	// Make request to the content server
	fmt.Println("Pinged?")
	resp, err := http.Get("http://localhost:" + port + "/api/v1/ping-service-provider") // + "/image" + "?id=" + req.Choice)
	if err != nil {
		log.Printf("failed to get /route2: %v", err)
		return
	}
	defer resp.Body.Close()

	// Convert the response body to a string
	RespBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return
	}

	// Send the response back to the WASM module
	_, err = w.Write(RespBodyBytes)
	if err != nil {
		log.Printf("failed to send response: %v", err)
		return
	}
}

// RegisterUserHandler handles user registration requests
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// Unmarshal request
	var req models.RegisterUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.(*validator.InvalidValidationError).Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Make connection to database
	db := config.SetupDatabaseConnection()
	// Close connection database
	defer config.CloseDatabaseConnection(db)
	// Save user to database
	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Salt:     req.Salt,
	}
	if err := db.Create(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// LoginPrecheckHandler handles login precheck requests and get the salt of the user from the database using the username from the request URL
func LoginPrecheckHandler(w http.ResponseWriter, r *http.Request) {
	// Get username from body
	var req models.LoginPrecheckDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.(*validator.InvalidValidationError).Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Make connection to database
	db := config.SetupDatabaseConnection()
	// Close connection database
	defer config.CloseDatabaseConnection(db)
	// Using the username, find the user in the database
	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	resp := models.LoginPrecheckResponseDTO{
		Username: user.Username,
		Salt:     user.Salt,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Unmarshal request
	var req models.LoginUserDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.(*validator.InvalidValidationError).Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Make connection to database
	db := config.SetupDatabaseConnection()
	// Close connection database
	defer config.CloseDatabaseConnection(db)
	// Using the username, find the user in the database
	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Compare the password with the password in the database
	if user.Password != req.SaltedHashedPassword {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Invalid credentials"))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	err := db.Model(&user).Update("public_key", req.PubKey).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
	masterPort := os.Getenv("LAYER8_MASTER_PORT")
	// Get JWT_SECRET from the Layer8 Master gRPC server
	conn, err := grpc.Dial("localhost:"+masterPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Create gRPC client
	client := pb.NewLayer8MasterServiceClient(conn)
	// Make gRPC request
	jwtSecretResp, err := client.GetJwtSecret(context.Background(), &pb.Empty{})
	if err != nil {
		log.Printf("failed to get jwt secret: %v", err)
		return
	}
	JWT_SECRET_STR := jwtSecretResp.JwtSecret
	// Save the JWT_SECRET as a byte array
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &models.Claims{
		UserName: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "GlobeAndCitizen",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWT_SECRET_BYTE)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
	resp := models.LoginUserResponseDTO{
		Token: tokenString,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
}

func GetContentHandler(w http.ResponseWriter, r *http.Request) {
	// Unmarshal request
	var req models.ContentReqDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.(*validator.InvalidValidationError).Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	masterPort := os.Getenv("LAYER8_MASTER_PORT")
	// Get JWT_SECRET from the Layer8 Master gRPC server
	conn, err := grpc.Dial("localhost:"+masterPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Create gRPC client
	client := pb.NewLayer8MasterServiceClient(conn)
	// Make gRPC request
	jwtSecretResp, err := client.GetJwtSecret(context.Background(), &pb.Empty{})
	if err != nil {
		log.Printf("failed to get jwt secret: %v", err)
		return
	}
	JWT_SECRET_STR := jwtSecretResp.JwtSecret
	// Save the JWT_SECRET as a byte array
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)
	// Separate the ECDSA signature from the rest of the token
	JwtSignedToken := strings.Split(req.Token, ".")[0] + "." + strings.Split(req.Token, ".")[1] + "." + strings.Split(req.Token, ".")[2]
	fmt.Println("JwtSignedToken: ", JwtSignedToken)
	// Parse the token
	token, err := jwt.ParseWithClaims(JwtSignedToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET_BYTE, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Check if the token is valid
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Invalid token"))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Get the user id from the token
	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error getting user id"))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	fmt.Println("User id: ", claims.UserID)
	fmt.Println("Token expires at: ", claims.ExpiresAt)

	// Make connection to database
	db := config.SetupDatabaseConnection()
	// Close connection database
	defer config.CloseDatabaseConnection(db)
	// Get the public key of the user from the database saved during the login precheck
	var user models.User
	if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}

	PUBLIC_KEY := user.PublicKey
	fmt.Printf("PUBLIC_KEY: %v\n", PUBLIC_KEY)
	publicKeyBytes, _ := hex.DecodeString(PUBLIC_KEY)

	// Separate the R and S components from the SignedToken using the dot separator
	encodedR := strings.Split(req.Token, ".")[3]
	encodedS := strings.Split(req.Token, ".")[4]

	rBytes, _ := base64.RawURLEncoding.DecodeString(encodedR)
	sBytes, _ := base64.RawURLEncoding.DecodeString(encodedS)

	// Create a new ECDSA public key
	pubKey := new(ecdsa.PublicKey)
	pubKey.Curve = elliptic.P256()
	pubKey.X, pubKey.Y = elliptic.Unmarshal(elliptic.P256(), publicKeyBytes)

	// Deserialize the original claims
	originalClaims := map[string]interface{}{
		"WasmSignature": "Signed by Go-WASM Client, Globe&Citizen",
	}

	// Serialize the original claims
	encodedOriginalClaims := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"WasmSignature":"%s"}`, originalClaims["WasmSignature"])))

	// Hash the original claims data
	hash := sha256.Sum256([]byte(encodedOriginalClaims))

	// Verify the ECDSA signature
	isValid := ecdsa.Verify(pubKey, hash[:], new(big.Int).SetBytes(rBytes), new(big.Int).SetBytes(sBytes))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("ECDSA signature is invalid"))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	fmt.Println("ECDSA signature is valid: ", encodedR, encodedS)

	port := os.Getenv("CONTENT_SERVER_PORT")
	// Make request to the content server
	resp, err := http.Get("http://localhost:" + port + "/image" + "?id=" + req.Choice)
	if err != nil {
		log.Printf("failed to get picture: %v", err)
		return
	}
	defer resp.Body.Close()

	// Convert the response body to a string
	RespBodyByteImg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return
	}

	// Convert RespBodyByte to string
	RespBodyString := string(RespBodyByteImg)

	// Send the response back to the WASM module
	_, err = w.Write([]byte(RespBodyString))
	if err != nil {
		log.Printf("failed to send response: %v", err)
		return
	}
}

func InitializeECDHTunnelHandler(w http.ResponseWriter, r *http.Request) {

	// Unmarshal request
	var req models.ECDHReqDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// validate request
	if err := validator.New().Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.(*validator.InvalidValidationError).Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	masterPort := os.Getenv("LAYER8_MASTER_PORT")
	// Get JWT_SECRET from the Layer8 Master gRPC server
	conn, err := grpc.Dial("localhost:"+masterPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Create gRPC client
	client := pb.NewLayer8MasterServiceClient(conn)
	// Make gRPC request
	jwtSecretResp, err := client.GetJwtSecret(context.Background(), &pb.Empty{})
	if err != nil {
		log.Printf("failed to get jwt secret: %v", err)
		return
	}
	JWT_SECRET_STR := jwtSecretResp.JwtSecret
	// Save the JWT_SECRET as a byte array
	JWT_SECRET_BYTE := []byte(JWT_SECRET_STR)
	// Separate the ECDSA signature from the rest of the token
	JwtSignedToken := strings.Split(req.Token, ".")[0] + "." + strings.Split(req.Token, ".")[1] + "." + strings.Split(req.Token, ".")[2]
	fmt.Println("JwtSignedToken: ", JwtSignedToken)
	// Parse the token
	token, err := jwt.ParseWithClaims(JwtSignedToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET_BYTE, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Check if the token is valid
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Invalid token"))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	// Get the user id from the token
	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error getting user id"))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	fmt.Println("User id: ", claims.UserID)
	fmt.Println("Token expires at: ", claims.ExpiresAt)

	// Make connection to database
	db := config.SetupDatabaseConnection()
	// Close connection database
	defer config.CloseDatabaseConnection(db)
	// Get the public key of the user from the database saved during the login precheck
	var user models.User
	if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}

	PUBLIC_KEY := user.PublicKey
	fmt.Printf("PUBLIC_KEY: %v\n", PUBLIC_KEY)
	publicKeyBytes, _ := hex.DecodeString(PUBLIC_KEY)

	// Separate the R and S components from the SignedToken using the dot separator
	encodedR := strings.Split(req.Token, ".")[3]
	encodedS := strings.Split(req.Token, ".")[4]

	rBytes, _ := base64.RawURLEncoding.DecodeString(encodedR)
	sBytes, _ := base64.RawURLEncoding.DecodeString(encodedS)

	// Create a new ECDSA public key
	pubKey := new(ecdsa.PublicKey)
	pubKey.Curve = elliptic.P256()
	pubKey.X, pubKey.Y = elliptic.Unmarshal(elliptic.P256(), publicKeyBytes)

	// Deserialize the original claims
	originalClaims := map[string]interface{}{
		"WasmSignature": "Signed by Go-WASM Client, Globe&Citizen",
	}

	// Serialize the original claims
	encodedOriginalClaims := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"WasmSignature":"%s"}`, originalClaims["WasmSignature"])))

	// Hash the original claims data
	hash := sha256.Sum256([]byte(encodedOriginalClaims))

	// Verify the ECDSA signature
	isValid := ecdsa.Verify(pubKey, hash[:], new(big.Int).SetBytes(rBytes), new(big.Int).SetBytes(sBytes))
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("ECDSA signature is invalid"))
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
		return
	}
	fmt.Println("ECDSA signature is valid: ", encodedR, encodedS)

	payload := struct {
		PubKeyWasmX *big.Int `json:"pub_key_wasm_x"`
		PubKeyWasmY *big.Int `json:"pub_key_wasm_y"`
	}{
		PubKeyWasmX: req.PubKeyWasmX,
		PubKeyWasmY: req.PubKeyWasmY,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload: %v", err)
		return
	}
	port := os.Getenv("CONTENT_SERVER_PORT")
	// Make request to the content server
	resp, err := http.Post("http://localhost:"+port+"/initialize-ecdh-key-exchange", "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Printf("failed to initialize ecdh key exchange: %v", err)
		return
	}
	defer resp.Body.Close()

	// Convert the response body to a string
	RespBodyByteKey, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return
	}
	var ecdhKeyExchangeOutput models.ECDHKeyExchangeOutput
	if err := json.Unmarshal(RespBodyByteKey, &ecdhKeyExchangeOutput); err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
		return
	}
	// Send the response back to the WASM module
	err = json.NewEncoder(w).Encode(ecdhKeyExchangeOutput)
	if err != nil {
		http.Error(w, "Error sending server's public key", http.StatusInternalServerError)
		return
	}
}
