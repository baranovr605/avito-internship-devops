package main

// Imports for server and work with redis
import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/go-redis/redis"
    "net/http"
)

// Vars for setup connect with redis and listen port server
var redisAddr string = "redis-go:6379"
var listenAddrServ string = ":8000"

// Function setup database
func setupDB() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: "",
		DB: 0,
	})

  return client
}

// Function add key-value in redis by API
func set_key(w http.ResponseWriter, r *http.Request) {

  w.Header().Set("Content-Type", "application/json")

  var requestBody map[string]string
  err := json.NewDecoder(r.Body).Decode(&requestBody)

  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }

  db := setupDB()
  for k, v := range requestBody {
    err := db.Set(k, v, 0).Err()
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Fprintln(w, "Key-val correctly write in redis!")
    }
  }

  db.Close()
   
}

// Function get key in redis by API
func get_key(w http.ResponseWriter, r *http.Request) {

  w.Header().Set("Content-Type", "application/json")

  id := r.URL.Query()["key"][0]
  
  db := setupDB()
  value, _ := db.Get(id).Result()
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

  w.Header().Set("Content-Type", "application/json")

  var requestBody map[string]string
  err := json.NewDecoder(r.Body).Decode(&requestBody)

  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  }

  db := setupDB()
  for k, _ := range requestBody {
      _, err := db.Del(k).Result()
      if err != nil { 
        fmt.Println(err)
      } else {
        fmt.Fprintln(w, "Key correctly deleted!")
      }
  }

  db.Close()

}

// Main function with endpoint's
func main() {

  router := mux.NewRouter()

  router.HandleFunc("/set_key", set_key).Methods("POST")

  router.HandleFunc("/get_key", get_key).Methods("GET")

  router.HandleFunc("/del_key", del_key).Methods("DELETE")

  http.ListenAndServe(listenAddrServ, router)
}
