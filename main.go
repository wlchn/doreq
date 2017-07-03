package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "github.com/andlabs/ui"
)

func main() {
    err := ui.Main(func() {
        name := ui.NewEntry()
        button := ui.NewButton("Greet")
        greeting := ui.NewLabel("")
        box := ui.NewVerticalBox()

        urlEntry := ui.NewEntry()
        getButton := ui.NewButton("Get")
        resLabel := ui.NewLabel("")

        box.Append(ui.NewLabel("Enter your name:"), false)
        box.Append(name, false)
        box.Append(button, false)
        box.Append(greeting, false)

        box.Append(ui.NewLabel("Input url:"), false)
        box.Append(urlEntry, false)
        box.Append(getButton, false)
        box.Append(resLabel, false)

        window := ui.NewWindow("Hello", 400, 300, false)
        window.SetChild(box)
        button.OnClicked(func(*ui.Button) {
            fmt.Printf("button pressed\n")
            greeting.SetText("Hello, " + name.Text() + "!")
        })
        getButton.OnClicked(func(*ui.Button) {
            fmt.Printf("getButton pressed\n")
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