package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"wallet_sdk"
	"wallet_sdk/client"
	"wallet_sdk/utils"

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

	w.Resize(fyne.NewSize(1100, 700))

	w.Show()
	a.Run()
	exit()
}

func MainContent(w fyne.Window) {
	tabs := container.NewAppTabs(
		container.NewTabItem("Generate wallet", GenerateWallet()),
		container.NewTabItem("Get address Unutxo list", GetAddressUTXO()),
		container.NewTabItem("Transaction info", TransactionInfo()),
		container.NewTabItem("Test TMP", E_G_Box()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(tabs)
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

func GetAddressUTXO() *fyne.Container {
	tip := widget.NewLabel("Enter address to query UTXO")
	// 结果信息
	leftLabel := widget.NewLabel("")
	// 地址输入框
	addressInput := widget.NewEntry()
	// UTXO列表
	data := binding.BindStringList(
		&[]string{},
	)
	list := utils.NewDataList(data)
	// 请求按钮
	query := widget.NewButton("QUERY", func() {
		fmt.Printf("wch------ data: %+v\n", data.Length())
		if data.Length() > 0 {
			data.Set([]string{})
			err := data.Reload()
			if err != nil {
				return
			}
		}
		fmt.Printf("wch------ data1: %+v\n", data.Length())
		addr := addressInput.Text
		fmt.Printf("wch---- addr: %+v\n", addr)
		res2 := wallet_sdk.GetUTXOListByAddress(chainName, addr)
		sum := int64(0)
		utxoList := res2.Data.([]*client.UnspendUTXOList)
		for _, unspentUTXO := range utxoList {
			val := fmt.Sprintf("%s:%d %d", unspentUTXO.TxHash, unspentUTXO.Vout, unspentUTXO.Amount)
			sum += unspentUTXO.Amount
			data.Append(val)
		}
		// 侧边栏统计内容
		leftContent := fmt.Sprintf("Number: %v\n Sum: %v", len(utxoList), sum)
		leftLabel.SetText(leftContent)
	})
	send := widget.NewButton("Transaction", func() {
		for i := 0; i < data.Length(); i++ {
			d, _ := data.GetItem(i)
			fmt.Printf("d: %+v\n", d)
		}
	})
	button := container.NewHBox(query, send)
	// 顶部提示
	top := container.NewVBox(tip, addressInput)
	// 侧边统计
	left := container.NewVBox(leftLabel)
	return container.NewBorder(top, button, left, nil, list)
}

func TransactionInfo() *fyne.Container {
	// 选择网络
	chainName := wallet_sdk.BTC_Testnet
	chainCombo := widget.NewSelect(wallet_sdk.ChainCombo, func(value string) {
		chainName = value
	})
	// 输入地址
	fromAddr := widget.NewEntry()
	fromAddr.SetPlaceHolder("Enter from address")
	fromAddr.SetText("n1HE1YJ1zF5U5aiX2DNu5WhjE9KFrkSKkx")
	// 转账金额
	amount := widget.NewEntry()
	amount.SetPlaceHolder("Enter from amount")
	amount.SetText("0.00004")
	// 转出地址
	toAddr := widget.NewEntry()
	toAddr.SetPlaceHolder("Enter to address")
	toAddr.SetText("2NBeoUKGLyk5ZfSDtAvsfWteYQaAKdUAniF")
	// 输入地址私钥
	priKey := widget.NewEntry()
	priKey.SetPlaceHolder("Enter from private key")

	// 结果提示
	str := binding.NewString()
	text := widget.NewLabelWithData(str)
	text.Wrapping = fyne.TextWrapWord // 设置为单词换行
	alert := container.NewVBox()
	alert.Resize(fyne.NewSize(300, 0))
	alert.Add(text)

	// 交易内容
	var signData string

	// 操作按钮
	btn1 := widget.NewButton("Builder", func() {
		// 查询主币余额
		res2 := wallet_sdk.GetBalanceByAddress(chainName, fromAddr.Text)
		// 查询节点gas price
		gasPriceData := wallet_sdk.GetGasPrice(chainName)
		gasPrice := gasPriceData.Data.Average
		// 构建交易
		res5 := wallet_sdk.BuildTransferInfoByBTC(chainName, fromAddr.Text, toAddr.Text, amount.Text, gasPrice)
		if res5.Status.Code == 0 {
			signData = res5.Data
			// 提示内容
			alertStr := fmt.Sprintf("Balance: %+v\n gasPrice: %+v\n", res2.Data, gasPrice)
			str.Set(alertStr)
		} else {
			str.Set(res5.Status.Message)
		}

	})
	btn2 := widget.NewButton("Sign&Broadcast", func() {
		if priKey.Text == "" || signData == "" {
			str.Set("Please check what you entered!")
		}
		res7 := wallet_sdk.SignAndSendTransferInfo(chainName, priKey.Text, signData, fromAddr.Text)
		str.Set(res7.Data)
	})
	button := container.New(layout.NewGridLayout(2), btn1, btn2)

	from := container.NewVBox(chainCombo, fromAddr, amount, toAddr, priKey, button, alert)
	return from
}

func E_G_Box() *fyne.Container {
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
