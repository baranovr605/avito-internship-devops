package main

// Imports for server and work with redis
import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/redis/go-redis/v9"
    "net/http"
    "os"
    "log"
    "io/ioutil"
    "context"
    "crypto/tls"
)

// DB sctruct
type Database struct {
  Client *redis.Client
}

// Context work with for DB
var (
  ctx = context.TODO()
)

// Var for setup Content-Type, listend port for app, redis data
var (
  contentType string      = "application/json"
  app_port    string      = os.Getenv("APP_PORT")
  redis_host  string      = os.Getenv("REDIS_HOST")
  redis_user  string      = os.Getenv("REDIS_USER")
  redis_pass_file string  = os.Getenv("REDIS_PASS_FILE")
  redis_cert  string      = os.Getenv("REDIS_CRT")
  redis_key   string      = os.Getenv("REDIS_KEY")
)

// Function setup database
func setupDB(host string, username string, pass string, 
            cert string, key string) (*Database, error) {

  // Get certs and key for TLS
  tls_cert, err := tls.LoadX509KeyPair(cert, key)
  if err != nil {
    log.Fatal(err)
  }

  // Setup TLS config
  tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{tls_cert},
    MinVersion:   tls.VersionTLS12,
    InsecureSkipVerify: true,
  }

  // Setup DB client
  client := redis.NewClient(&redis.Options{
    Addr:     host,
    Username: username, 
    Password: pass,
    TLSConfig: tlsConfig,
    })
  
  // Ping to Redis for check connect
  if err := client.Ping(ctx).Err(); err != nil {
    return nil, err
  }

  return &Database{
    Client: client,
  }, nil
}

// Functions for return pass from secret file
func getPassFile(fileName string) string {

  redis_pass_file, err := ioutil.ReadFile(fileName)
  if err != nil {
      panic(err)
  }

  return string(redis_pass_file)
}

// Function for return error 405, if not correct request
func returnErr405(w http.ResponseWriter) {

  w.WriteHeader(http.StatusMethodNotAllowed)
  fmt.Fprintln(w, "405 Method not allowed")
  
  return
}

// Function add key-value in redis by API
func set_key(w http.ResponseWriter, r *http.Request, db *Database) {

  w.Header().Set("Content-Type", contentType)
  var requestBody map[string]string
  err := json.NewDecoder(r.Body).Decode(&requestBody)

  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }

  if len(requestBody) > 1 {
    returnErr405(w)
    return
  }

  for key, value := range requestBody {
    err := db.Client.Set(ctx, key, value, 0).Err()
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Fprintln(w, "Key-val correctly write in redis!")
    }
  }
   
}

// Function get key in redis by API
func get_key(w http.ResponseWriter, r *http.Request, db *Database) {

  w.Header().Set("Content-Type", contentType)

  id := r.URL.Query()["key"][0]
  
  value, _ := db.Client.Get(ctx, id).Result()

  if value != "" {
    fmt.Fprintln(w, value)
  } else {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintln(w, "404 page not found")
    return
  }
  
}

// Function del key in redis by API
func del_key(w http.ResponseWriter, r *http.Request, db *Database) {

  w.Header().Set("Content-Type", contentType)

  var requestBody map[string]string
  err := json.NewDecoder(r.Body).Decode(&requestBody)

  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }

  if len(requestBody) > 1 {
    returnErr405(w)
    return 
  }

  _, errbd := db.Client.Del(ctx, requestBody["key"]).Result()


  if errbd != nil { 
    fmt.Println(errbd)
  } else {
    fmt.Fprintln(w, "Key correctly deleted!")
  }

}

// Main function with endpoint's
func main() {
  
  // Setup DB client for make APi requests
  redis_pass := getPassFile(redis_pass_file)
  db, err := setupDB(redis_host, redis_user, redis_pass, redis_cert, redis_key)

  if err != nil {
    log.Fatalf("Failed to connect to redis: %s", err.Error())
  }

  // Setup routers with API
  router := mux.NewRouter()

  router.HandleFunc("/set_key", func(w http.ResponseWriter, r *http.Request) {
    set_key(w, r, db)
    }).Methods("POST")

  router.HandleFunc("/get_key", func(w http.ResponseWriter, r *http.Request) {
    get_key(w, r, db)
    }).Methods("GET")

  router.HandleFunc("/del_key", func(w http.ResponseWriter, r *http.Request) {
    del_key(w, r, db)
    }).Methods("DELETE")

  http.ListenAndServe(app_port, router)

  db.Client.Close()

}
