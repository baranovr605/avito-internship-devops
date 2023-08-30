package main

// Imports for server and work with redis
import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/redis/go-redis/v9"
    "net/http"
    "os"
    "io/ioutil"
    "context"
)

// Var for setup Content-Type, listend port for app, redis data
var (
  contentType string      = "application/json"
  app_port    string      = os.Getenv("APP_PORT")
  redis_host  string      = os.Getenv("REDIS_HOST")
  redis_user  string      = os.Getenv("REDIS_USER")
  redis_pass_file string  = os.Getenv("REDIS_PASS_FILE")
)

// Function setup database
func setupDB() *redis.Client {

  // Get Pass for redis from file (before mount for docker secret)
	redis_pass_file, err := ioutil.ReadFile(redis_pass_file)
    if err != nil {
        panic(err)
    }
	redis_pass := string(redis_pass_file)

  // Create redis client
  client := redis.NewClient(&redis.Options{
    Addr:     redis_host,
    Username: redis_user, 
    Password: redis_pass, 
})
  
  return client
}

// Function add key-value in redis by API
func set_key(w http.ResponseWriter, r *http.Request) {

  ctx := context.Background()
  w.Header().Set("Content-Type", contentType)
  var requestBody map[string]string
  err := json.NewDecoder(r.Body).Decode(&requestBody)

  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }

  db := setupDB()
  if len(requestBody) > 1: {
    w.WriteHeader(http.StatusMethodNotAllowed)
    fmt.Fprintln(w, "405 Method not allowed")
  } else {
    for key, value := range requestBody {
      err := db.Set(ctx, key, value, 0).Err()
      if err != nil {
        fmt.Println(err)
      } else {
        fmt.Fprintln(w, "Key-val correctly write in redis!")
      }
    }
  }

  db.Close()
   
}

// Function get key in redis by API
func get_key(w http.ResponseWriter, r *http.Request) {

  ctx := context.Background()
  w.Header().Set("Content-Type", contentType)

  id := r.URL.Query()["key"][0]
  
  db := setupDB()
  value, _ := db.Get(ctx, id).Result()
  if value != "" {
    fmt.Fprintln(w, value)
  } else {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintln(w, "404 page not found")
  }

  db.Close()
  
}

// Function del key in redis by API
func del_key(w http.ResponseWriter, r *http.Request) {

  ctx := context.Background()
  w.Header().Set("Content-Type", contentType)

  var requestBody map[string]string
  err := json.NewDecoder(r.Body).Decode(&requestBody)

  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }

  db := setupDB()

  if len(requestBody) > 1: {
    w.WriteHeader(http.StatusMethodNotAllowed)
    fmt.Fprintln(w, "405 Method not allowed")
  } else {
    _, errbd := db.Del(ctx, requestBody["key"]).Result()
  } 
  
  if errbd != nil { 
    fmt.Println(errbd)
  } else {
    fmt.Fprintln(w, "Key correctly deleted!")
  }

  db.Close()

}

// Main function with endpoint's
func main() {

  router := mux.NewRouter()

  router.HandleFunc("/set_key", set_key).Methods("POST")

  router.HandleFunc("/get_key", get_key).Methods("GET")

  router.HandleFunc("/del_key", del_key).Methods("DELETE")

  http.ListenAndServe(app_port, router)
}
