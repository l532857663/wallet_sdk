package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"wallet_sdk"
	"wallet_sdk/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	// BTC
	chainName = wallet_sdk.BTC_Testnet
)

func main() {
	a := app.New()
	w := a.NewWindow("Wallet 钱包")

	MainContent(w)

	w.Resize(fyne.NewSize(900, 700))

	w.Show()
	a.Run()
	exit()
}

func MainContent(w fyne.Window) {
	tabs := container.NewAppTabs(
		container.NewTabItem("Get address Unutxo list", GetAddressUTXO()),
		container.NewTabItem("Generate wallet", GenerateWallet()),
		container.NewTabItem("Test TMP", testTmp()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(tabs)
}

func GetAddressUTXO() *fyne.Container {
	tip := widget.NewLabel("Enter address to query UTXO")
	count := widget.NewLabel("")
	addressInput := widget.NewEntry()
	data := binding.BindStringList(
		&[]string{},
	)
	query := widget.NewButton("QUERY", func() {
		if data.Length() > 0 {
			err := data.Reload()
			if err != nil {
				fmt.Printf("data.Reload error: %+v\n", err)
				return
			}
		}
		addr := addressInput.Text
		fmt.Printf("wch---- addr: %+v\n", addr)
		res2 := wallet_sdk.GetUTXOListByAddress(chainName, addr)
		sum := int64(0)
		utxoList := res2.Data.([]*client.UnspendUTXOList)
		for _, unspentUTXO := range utxoList {
			val := fmt.Sprintf("UTXO %s:%d, Amount: %d", unspentUTXO.TxHash, unspentUTXO.Vout, unspentUTXO.Amount)
			sum += unspentUTXO.Amount
			data.Append(val)
		}
		count.SetText(fmt.Sprintf("Number: %v\n Sum: %v", len(utxoList), sum))
	})
	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})
	top := container.NewVBox(tip, addressInput)
	left := container.NewVBox(count)
	return container.NewBorder(top, query, left, nil, list)
}

func GenerateWallet() *fyne.Container {
	var (
		// 助记词
		mnemonic = ""
	)
	/* ------------------------------- CONTENT ------------------------------- */
	// 展示生成的助记词
	content := container.New(layout.NewGridLayout(3))
	/* ------------------------------- CONTENT ------------------------------- */
	/* ------------------------------- LEFT ------------------------------- */
	// 设置助记词类型
	length := widget.NewRadioGroup([]string{"12", "24"}, func(value string) {
		fmt.Println("Radio set to", value)
	})
	length.SetSelected("12")
	// TODO: 中文暂时显示不出来
	// language := widget.NewSelect([]string{"EN", "CN_S", "CN_T"}, func(value string) {
	// 	fmt.Println("Select set to", value)
	// })
	// language.SetSelected("EN")
	// 导入助记词
	inputM := widget.NewEntry()
	inputM.SetPlaceHolder("Enter mnemonic...")
	importMnemonic := container.NewVBox(inputM, widget.NewButton("Import", func() {
		mnemonic = inputM.Text
		mList := strings.Split(inputM.Text, " ")
		for _, m := range mList {
			content.Add(getCenter(m))
		}
	}))
	left := container.NewVBox(
		widget.NewLabel("Default length is 12, configurable to 24"),
		length,
		// widget.NewLabel("Default language is EN"),
		// language,
		importMnemonic,
	)
	/* ------------------------------- LEFT ------------------------------- */
	/* ------------------------------- RIGHT ------------------------------- */
	priKey := binding.NewString()
	address := binding.NewString()
	right := container.NewVBox(
		widget.NewLabel("PrivateKey:"),
		widget.NewLabelWithData(priKey),
		widget.NewLabel("Address:"),
		widget.NewLabelWithData(address),
	)
	/* ------------------------------- RIGHT ------------------------------- */
	/* ------------------------------- BUTTON ------------------------------- */
	// 设置按钮
	btn1 := widget.NewButton("GenerateMnemonic", func() {
		l, _ := strconv.Atoi(length.Selected)
		// res := wallet_sdk.GenerateMnemonic(l, language.Selected)
		res := wallet_sdk.GenerateMnemonic(l, "")
		mnemonic = res.Data
		fmt.Printf("res.Data: %+v\n", mnemonic)
		mList := strings.Split(mnemonic, " ")
		for _, m := range mList {
			content.Add(getCenter(m))
		}
	})
	btn2 := widget.NewButton("Generate wallet", func() {
		var purpose uint32 = 86
		params := "BTCTest"
		res1 := wallet_sdk.GenerateAccountByMnemonic(mnemonic, params, &purpose)
		fmt.Printf("res: %+v\n", res1)
		fmt.Printf("res Data: %+v\n", res1.Data)
		accountInfo := res1.Data
		priKey.Set(accountInfo.PrivateKey)
		address.Set(accountInfo.Address)
	})
	button := container.New(layout.NewGridLayout(2), btn1, btn2)
	/* ------------------------------- BUTTON ------------------------------- */
	return container.NewBorder(nil, button, left, right, content)
}

func getCenter(data string) *fyne.Container {
	return container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), widget.NewLabel(data))
}

func testTmp() *fyne.Container {
	text1 := canvas.NewText("Hello", color.Black)
	text2 := canvas.NewText("There", color.Black)
	text3 := canvas.NewText("(right)", color.Black)
	content := container.New(layout.NewHBoxLayout(), text1, text2, layout.NewSpacer(), text3)

	text4 := canvas.NewText("centered", color.Black)
	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), text4, layout.NewSpacer())
	return container.New(layout.NewVBoxLayout(), content, centered)
}

// 退出应用后调用 Run()方法不会执行后续的代码
func exit() {
	fmt.Println("Exited")
}
