package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andlabs/ui"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	err := ui.Main(func() {
		vBox := ui.NewVerticalBox()
		hBox := ui.NewHorizontalBox()

		urlEntry := ui.NewEntry()
		jsonLabel := ui.NewLabel("Input json:")
		jsonEntry := ui.NewEntry()
		goButton := ui.NewButton("Go")
		resLabel := ui.NewLabel("")
		selectCombobox := ui.NewCombobox()

		hBox.Append(selectCombobox, false)
		hBox.Append(urlEntry, true)
		hBox.Append(goButton, false)
		hBox.SetPadded(true)

		vBox.Append(hBox, false)
		vBox.Append(jsonLabel, false)
		vBox.Append(jsonEntry, false)
		vBox.Append(resLabel, false)
		vBox.SetPadded(true)

		jsonLabel.Hide()
		jsonEntry.Hide()

		selectCombobox.Append("Get")
		selectCombobox.Append("Post")
		selectCombobox.SetSelected(0)

		window := ui.NewWindow("doreq - do request", 800, 600, false)
		window.SetChild(vBox)
		window.SetMargined(true)

		selectCombobox.OnSelected(func(*ui.Combobox) {
			selectedIndex := selectCombobox.Selected()
			fmt.Println("selectedIndex:", selectedIndex)
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
				handleResp(resp, resLabel)
			case 1:
				urlStr := urlEntry.Text()
				jsonStr := jsonEntry.Text()

				var jsonBytes = []byte(jsonStr)
				req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonBytes))
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				handleResp(resp, resLabel)
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

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func handleResp(resp *http.Response, label *ui.Label) {
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	if resp.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		bodyString := string(bodyBytes)
		if !isJSON(bodyString) {
			bodyString = "Response is not json."
		}
		label.SetText(bodyString)
	} else {
		label.SetText("Error with code " + strconv.Itoa(resp.StatusCode))
	}
}
