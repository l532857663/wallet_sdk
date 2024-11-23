package main

import (
	"fmt"
	"strconv"
	"strings"
	"wallet_sdk"
	"wallet_sdk/client"
	"wallet_sdk/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shopspring/decimal"
)

var (
	a = app.New()
)

func main() {
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
		container.NewTabItem("Multi to multi transaction", MultiToMultiTransfer()),
		container.NewTabItem("Test TMP", E_G_Box()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	// 设置窗口的内容
	w.SetContent(tabs)
}

func GenerateWallet() *fyne.Container {
	var (
		// 助记词
		mnemonic = ""
		// 助记词数量
		mnemonicLen = []string{"12", "24"}
		// TODO: 助记词显示语言
		// langs = []string{"EN", "CN_S", "CN_T"}
		// 网络可选列表
		networks = []string{"BTC", "BTCTest", "BTCRegt", "ETH", "TRON"}
	)
	/* ------------------------------- CONTENT ------------------------------- */
	// 展示生成的助记词
	content := container.New(layout.NewGridLayout(3))
	/* ------------------------------- CONTENT ------------------------------- */
	/* ------------------------------- LEFT ------------------------------- */
	// 设置助记词类型
	length := widget.NewRadioGroup(mnemonicLen, func(value string) {})
	length.SetSelected("12")
	// 选择网络
	params := widget.NewSelect(networks, func(selected string) {})
	params.SetSelected("BTCTest")
	// TODO: 中文暂时显示不出来
	// language := widget.NewSelect(langs, func(value string) {
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
		params,
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
		var addressIndex uint32 = 0
		res1 := wallet_sdk.GenerateAccountByMnemonic(mnemonic, params.Selected, &addressIndex)
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
	leftLabel2 := widget.NewLabel("")
	// 地址输入框
	addressInput := widget.NewEntry()
	// UTXO列表
	resultContainer := container.NewVBox()
	var useUTXOList []*client.UnspendUTXOList
	// 请求按钮
	query := widget.NewButton("QUERY", func() {
		addr := addressInput.Text
		fmt.Printf("wch---- addr: %+v\n", addr)
		res2 := wallet_sdk.GetUTXOListByAddress(chainName, addr)
		sum := int64(0)
		checkSum := int64(0)
		utxoList := res2.Data.([]*client.UnspendUTXOList)
		resultContainer.Objects = nil
		// 排序未花费的UTXO
		client.DescSortUnspendUTXO(utxoList)
		for _, unspentUTXO := range utxoList {
			// UTXO展示内容
			val := utils.EncodeStringByUtxoInfo(unspentUTXO.TxHash, unspentUTXO.Vout, unspentUTXO.Amount)
			sum += unspentUTXO.Amount
			// 多选框处理
			checkbox := widget.NewCheck(val, func(c bool) {
				_, _, amount := utils.DecodeUtxoInfoByString(val)
				if c {
					checkSum += amount
				} else {
					checkSum -= amount
				}
				checkSumRes := fmt.Sprintf("Total of selected UTXOs:\n %+v BTC", checkSum)
				leftLabel2.SetText(checkSumRes)
			})
			resultContainer.Add(checkbox)
		}
		resultContainer.Refresh()
		useUTXOList = utxoList
		// 侧边栏统计内容
		leftContent := fmt.Sprintf("Number: %v\n Sum: %v", len(utxoList), sum)
		leftLabel.SetText(leftContent)
	})
	send := widget.NewButton("Transaction", func() {
		ChooseUTXOToTransfer(resultContainer, useUTXOList)
	})
	button := container.New(layout.NewGridLayout(2), query, send)
	// 顶部提示
	top := container.NewVBox(tip, addressInput)
	// 侧边统计
	left := container.NewVBox(leftLabel, leftLabel2)
	return container.NewBorder(top, button, left, nil, resultContainer)
}

func ChooseUTXOToTransfer(utxoList *fyne.Container, useUTXOList []*client.UnspendUTXOList) {
	fromInputs := container.NewVBox()
	var fromEntry []*widget.Label
	var useUTXOIndex []int
	for i, obj := range utxoList.Objects {
		checkbox, ok := obj.(*widget.Check)
		if ok && checkbox.Checked {
			entry := widget.NewLabel(checkbox.Text)
			fromInputs.Add(entry)
			fromEntry = append(fromEntry, entry)
			useUTXOIndex = append(useUTXOIndex, i)
		}
	}
	from := container.NewVBox(fromInputs)

	toInputs := container.NewVBox()
	var toEntry []*widget.Entry
	outButton := widget.NewButton("Add output", func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("toAddr:amount")
		toInputs.Add(entry)
		toInputs.Refresh() // 刷新容器以显示新的输入框
		toEntry = append(toEntry, entry)
	})
	to := container.NewVBox(outButton, toInputs)
	/* ------------------------------- TOP ------------------------------- */
	// 输入找零地址
	inputC := widget.NewEntry()
	inputC.SetPlaceHolder("Enter change address...")
	// 输入私钥
	inputM := widget.NewEntry()
	inputM.SetPlaceHolder("Enter privateKey [e.g. cUAxLxQT6W...]")
	// 输入手续费率
	inputG := widget.NewEntry()
	inputG.SetPlaceHolder("Enter gas price [e.g. 0.001] BTC/KB")
	alert := binding.NewString()
	alertBox := container.NewVBox(
		widget.NewLabel("************* Alert **************************"),
		widget.NewLabelWithData(alert),
		widget.NewLabel("**********************************************"),
	)
	top := container.NewVBox(inputC, inputM, inputG, alertBox)
	// 创建一个限制大小的容器
	limitedSizeContainer := container.NewMax(
		top,
	)
	limitedSizeContainer.Resize(fyne.NewSize(300, 200)) // 设置容器的大小
	/* ------------------------------- TOP ------------------------------- */
	signData := ""
	BuildBtn := widget.NewButton("1.BuildTransaction", func() {
		// from
		var vins []*client.UnspendUTXOList
		var inAmount int64
		for _, in := range useUTXOIndex {
			v := useUTXOList[in]
			vins = append(vins, v)
			inAmount += v.Amount
		}
		// to
		var vouts, amounts []string
		outAmount := decimal.Zero
		for _, out := range toEntry {
			toInfo := strings.Split(out.Text, ":")
			if len(toInfo) < 2 {
				continue
			}
			vouts = append(vouts, toInfo[0])
			amounts = append(amounts, toInfo[1])
			a, _ := utils.StringToDecimal(toInfo[1])
			outAmount = outAmount.Add(a)
		}
		// 没填找零地址报错
		if inputC.Text == "" {
			alert.Set("Please enter the change address!")
			return
		}
		// 查询节点gas price
		gasPriceData := wallet_sdk.GetGasPrice(chainName)
		gasPrice := gasPriceData.Data.Average
		if inputG.Text != "" {
			gasPrice = inputG.Text
		}
		// // 构建交易数据
		res1 := wallet_sdk.MultiToMultiTransfer(chainName, vins, vouts, amounts, gasPrice, inputC.Text)
		signData = res1.Data
		size := len(res1.Data) / 2
		// 提示交易数据
		transferInfo := fmt.Sprintf("Get BTC transferInfo\n[In amount] %v\n[Out amount]%s BTC\n[Gas price] %s BTC/vKB, size: %v\n[Flinally fee] %v/1000*%v", utils.Int64ToSatoshi(inAmount), outAmount.String(), gasPrice, size, gasPrice, size)
		alert.Set(transferInfo)
	})
	SignBtn := widget.NewButton("2.SignTransaction", func() {
		priKey := inputM.Text
		// 没填私钥报错
		if inputM.Text == "" {
			alert.Set("Please enter the address private key!")
			return
		}
		fmt.Printf("wch---- sign: %+v\n", signData)
		res1 := wallet_sdk.SignTransferInfo(chainName, priKey, signData)
		fmt.Printf("wch------ res1 data: %+v\n", res1.Data)
		// 提示签名数据
		if res1.Status.Code == wallet_sdk.RES_CODE_FAILED {
			alert.Set(res1.Status.Message)
			return
		} else {
			signData := ""
			d := []byte(res1.Data)
			for i := 0; i < len(d); i++ {
				if i%64 == 0 {
					signData += "\n"
				}
				signData += string(d[i])
			}
			alert.Set("Signature Successful\n" + signData)
		}
		signData = res1.Data
	})
	BroadcastBtn := widget.NewButton("3.BroadcastTransaction", func() {
		res1 := wallet_sdk.BroadcastTransaction(chainName, signData)
		fmt.Printf("wch------ res1 data: %+v\n", res1.Data)
		// 提示交易HASH
		if res1.Status.Code == wallet_sdk.RES_CODE_FAILED {
			alert.Set(res1.Status.Message)
			return
		} else {
			alert.Set(res1.Data)
		}
	})
	button := container.New(layout.NewGridLayout(3), BuildBtn, SignBtn, BroadcastBtn)
	split := container.NewHSplit(from, to)
	split.SetOffset(0.5)
	content := container.NewBorder(limitedSizeContainer, button, nil, nil, split)
	// 在新页面显示选定的数据
	w := a.NewWindow("Selected Data")
	w.SetContent(content)
	w.Resize(fyne.NewSize(1000, 850))
	w.Show()
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
		// // 查询主币余额
		// res2 := wallet_sdk.GetBalanceByAddress(chainName, fromAddr.Text)
		// 查询节点gas price
		gasPriceData := wallet_sdk.GetGasPrice(chainName)
		gasPrice := gasPriceData.Data.Average
		// 构建交易
		res5 := wallet_sdk.BuildTransferInfoByBTC(chainName, fromAddr.Text, toAddr.Text, amount.Text, gasPrice)
		if res5.Status.Code == 0 {
			signData = res5.Data
			// 提示内容
			// alertStr := fmt.Sprintf("Balance: %+v\n gasPrice: %+v\n", res2.Data, gasPrice)
			alertStr := fmt.Sprintf("Balance: %+v\n gasPrice: %+v\n", 0, gasPrice)
			str.Set(alertStr)
		} else {
			str.Set(res5.Status.Message)
		}

	})
	btn2 := widget.NewButton("Sign&Broadcast", func() {
		if priKey.Text == "" || signData == "" {
			str.Set("Please check what you entered!")
		}
		res7 := wallet_sdk.SignTransferInfo(chainName, priKey.Text, signData)
		str.Set(res7.Data)
	})
	button := container.New(layout.NewGridLayout(2), btn1, btn2)

	from := container.NewVBox(chainCombo, fromAddr, amount, toAddr, priKey, button, alert)
	return from
}

func MultiToMultiTransfer() *fyne.Container {
	/******************************** Input UTXO *************************************/
	// 创建一个容器，用于存放动态创建的输入框
	fromInputs := container.NewVBox()
	var fromEntry []*widget.Entry
	// 创建一个按钮，当按钮被点击时，添加一个新的输入框到容器中
	inButton := widget.NewButton("Add in UTXO", func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("txHash:vout amount")
		fromInputs.Add(entry)
		fromInputs.Refresh() // 刷新容器以显示新的输入框
		fromEntry = append(fromEntry, entry)
	})
	from := container.NewVBox(inButton, fromInputs)
	/******************************** Input UTXO *************************************/
	/******************************** Output *************************************/
	toInputs := container.NewVBox()
	var toEntry []*widget.Entry
	outButton := widget.NewButton("Add output", func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("toAddr:amount")
		toInputs.Add(entry)
		toInputs.Refresh() // 刷新容器以显示新的输入框
		toEntry = append(toEntry, entry)
	})
	to := container.NewVBox(outButton, toInputs)
	/******************************** Output *************************************/
	// signData := ""
	BuildBtu := widget.NewButton("BuildTransaction", func() {
		// from
		var vins []wallet_sdk.ChooseUTXO
		var inputs []int64
		for _, in := range fromEntry {
			txHash, vout, amount := utils.DecodeUtxoInfoByString(in.Text)
			vins = append(vins, wallet_sdk.ChooseUTXO{
				TxHash: txHash,
				Vout:   vout,
			})
			inputs = append(inputs, amount)
		}
		// to
		var vouts, amounts []string
		for _, out := range toEntry {
			toInfo := strings.Split(out.Text, ":")
			if len(toInfo) < 2 {
				continue
			}
			vouts = append(vouts, toInfo[0])
			amounts = append(amounts, toInfo[1])
		}
		// 查询节点gas price
		gasPriceData := wallet_sdk.GetGasPrice(chainName)
		gasPrice := gasPriceData.Data.Average
		fmt.Printf("wch----- gasPrice: %+v\n", gasPrice)
		// 构建交易数据
		fmt.Printf("wch----- vins: %+v\n", vins)
		fmt.Printf("wch----- inputs: %+v\n", inputs)
		fmt.Printf("wch----- vouts: %+v\n", vouts)
		fmt.Printf("wch----- amounts: %+v\n", amounts)
		// res1 := wallet_sdk.MultiToMultiTransfer(chainName, vins, inputs, vouts, amounts, gasPrice, "tb1pfzl0rw44mkgevdauhrtzy5kdztjezyq0rnfqfppzxtnrwzdj553qvz6lux")
		// fmt.Printf("wch------ res1 data: %+v\n", res1.Data)
		// signData = res1.Data
	})
	split := container.NewHSplit(from, to)
	split.SetOffset(0.5)
	content := container.NewBorder(nil, BuildBtu, nil, nil, split)
	return content
}

func E_G_Box() *fyne.Container {
	// 创建标签和输入框，用于第一页
	label1 := widget.NewLabel("Enter query and click search:")
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter query")

	// 模拟的数据库数据
	data := []string{"Apple", "Banana", "Cherry", "Date", "Elderberry", "Fig", "Grape", "Honeydew"}

	// 创建一个容器，用于显示第一页内容
	page1 := container.NewVBox(label1, entry)

	// 创建第二页的标签和结果容器
	label2 := widget.NewLabel("Search Results:")
	resultContainer := container.NewVBox()
	page2 := container.NewVBox(label2, resultContainer)

	// 创建 Tab 容器，用于管理多个页面
	tabs := container.NewAppTabs(
		container.NewTabItem("Page 1", page1),
		container.NewTabItem("Page 2", page2),
	)

	// 隐藏第二页 Tab 以便于跳转控制
	tabs.Items[1].Content.Hide()

	// 创建一个按钮，当按钮被点击时，进行查询并跳转到第二页
	button := widget.NewButton("Search", func() {
		query := entry.Text
		filteredData := utils.FilterData(data, query)
		resultContainer.Objects = nil

		for _, item := range filteredData {
			checkbox := widget.NewCheck(item, func(bool) {})
			resultContainer.Add(checkbox)
		}
		resultContainer.Refresh()

		// 显示第二页内容并切换到第二页
		tabs.Items[1].Content.Show()
		tabs.SelectIndex(1)
	})

	hideButton := widget.NewButton("Hide Tabs", func() {
		if tabs.Visible() {
			tabs.Hide()
		} else {
			tabs.Show()
		}
	})

	// 将按钮添加到第一页容器
	page1.Add(button)
	page1.Add(hideButton)

	return container.NewMax(tabs)
}

// 退出应用后调用 Run()方法不会执行后续的代码
func exit() {
	a.Quit()
	fmt.Println("Exited")
}
