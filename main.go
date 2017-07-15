package main

import (
    "fmt"
    "bytes"
    "net/http"
    "io/ioutil"
    "github.com/andlabs/ui"
)

func main() {
    err := ui.Main(func() {
        box := ui.NewVerticalBox()

        urlEntry := ui.NewEntry()
        jsonLabel := ui.NewLabel("Input json:")
        jsonEntry := ui.NewEntry()
        goButton := ui.NewButton("Go")
        resLabel := ui.NewLabel("")
        selectCombobox := ui.NewCombobox()

        box.Append(ui.NewLabel("Input url:"), false)
        box.Append(urlEntry, false)
        box.Append(jsonLabel, false)
        box.Append(jsonEntry, false)
        box.Append(selectCombobox, false)
        box.Append(goButton, false)
        box.Append(resLabel, false)

        jsonLabel.Hide()
        jsonEntry.Hide()

        selectCombobox.Append("Get")
        selectCombobox.Append("Post")

        selectCombobox.SetSelected(0)

        window := ui.NewWindow("Doreq - Do request", 400, 300, false)
        window.SetChild(box)

        selectCombobox.OnSelected(func(*ui.Combobox) {
            selectedIndex := selectCombobox.Selected()
            fmt.Println("selectedIndex:>", selectedIndex)
            switch selectedIndex {
                case 1:
                    jsonLabel.Show()
                    jsonEntry.Show()
                default:
                    jsonLabel.Hide()
                    jsonEntry.Hide()
            }
        })

        goButton.OnClicked(func(*ui.Button) {
            fmt.Printf("goButton pressed\n")
            selectedIndex := selectCombobox.Selected()
            switch selectedIndex {
                case 0:
                    resp, err := http.Get(urlEntry.Text())
                    if err != nil {
                        resLabel.SetText(err.Error())
                        return
                    }
                    if resp.StatusCode == 200 {
                        bodyBytes, err := ioutil.ReadAll(resp.Body)
                        if err != nil {
                            return
                        }
                        bodyString := string(bodyBytes)
                        fmt.Printf(bodyString)
                        resLabel.SetText(bodyString)
                    }
                case 1:
                    url := urlEntry.Text()
                    fmt.Println("URL:>", url)
                    jsonStr := jsonEntry.Text()

                    var jsonBytes = []byte(jsonStr)
                    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
                    req.Header.Set("Content-Type", "application/json")

                    client := &http.Client{}
                    resp, err := client.Do(req)
                    if err != nil {
                        panic(err)
                    }
                    defer resp.Body.Close()

                    fmt.Println("response Status:", resp.Status)
                    fmt.Println("response Headers:", resp.Header)
                    body, _ := ioutil.ReadAll(resp.Body)
                    resLabel.SetText(string(body))
                    fmt.Println("response Body:", string(body))
                default:
                    fmt.Printf("no selected")
            }

        })
        window.OnClosing(func(*ui.Window) bool {
            ui.Quit()
            return true
        })
        window.Show()
    })
    if err != nil {
        panic(err)
    }
}