
package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "main/data"
    "main/tools"
    "mime/multipart"
    "net/http"
    "os"
    "strings"
    "time"
)

func waitServer(url string, duration time.Duration) bool {
    deadline := time.Now().Add(duration)
    for time.Now().Before(deadline) {
        resp, err := http.Get(url)

        if err == nil && resp.StatusCode == http.StatusOK {
            resp.Body.Close()
            return true
        }
    }
    return false
}

func main() {

    if !waitServer("http://localhost:9876", 5*time.Second) {
        fmt.Println("Server tidal ada")
        return
    }

    var choice int
    for {
        fmt.Println("Main Menu")
        fmt.Println("1. Get message")
        fmt.Println("2. Send file")
        fmt.Println("3. Quit")
        fmt.Print(">> ")
        fmt.Scanf("%d\n", &choice)

        if choice == 1 {
            go getMessage()  
        } else if choice == 2 {
            go sendFile()  
        } else if choice == 3 {
            break
        } else {
            fmt.Println("Invalid choice")
        }
        
        time.Sleep(1 * time.Second)
    }
}

func getMessage() {
    resp, err := http.Get("http://localhost:9876")
    tools.ErrorHandler(err)
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    tools.ErrorHandler(err)

    fmt.Println("Server: ", string(data))
}

func sendFile() {
    var name string
    var age int

    scanner := bufio.NewReader(os.Stdin)

    fmt.Print("input name: ")
    name, _ = scanner.ReadString('\n')
    name = strings.TrimSpace(name)

    fmt.Print("Input age: ")
    fmt.Scanf("%d\n", &age)

    person := data.Person{Name: name, Age: age}
    jsonData, err := json.Marshal(person)
    tools.ErrorHandler(err)

    temp := new(bytes.Buffer)
    w := multipart.NewWriter(temp)

    personField, err := w.CreateFormField("Person")
    tools.ErrorHandler(err)
    _, err = personField.Write(jsonData)
    tools.ErrorHandler(err)

    file, err := os.Open("./file.txt")
    tools.ErrorHandler(err)
    defer file.Close()

    fileField, err := w.CreateFormFile("File", file.Name())
    tools.ErrorHandler(err)
    _, err = io.Copy(fileField, file)
    tools.ErrorHandler(err)

    err = w.Close()
    tools.ErrorHandler(err)

    req, err := http.NewRequest("POST", "http://localhost:9876/sendFile", temp)
    tools.ErrorHandler(err)
    req.Header.Set("Content-Type", w.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    tools.ErrorHandler(err)
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    tools.ErrorHandler(err)

    fmt.Println("Server: ", string(data))
}
