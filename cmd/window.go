package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"wallet_sdk"
	"wallet_sdk/client"
	"wallet_sdk/global"
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
	a              = app.New()
	srv            *wallet_sdk.GetUtxoInfo
	mainPageInfo   *fyne.Container
	getBalancePage *fyne.Container
	mainSize       = fyne.NewSize(1100, 700)
)

var (
	// 次级页面
	getAddressBalaceTab *container.TabItem
)

func main() {
	wallet_sdk.MustLoad("config.yml")
	// 设置自定义字体
	myTheme, err := utils.NewMyTheme()
	if err != nil {
		log.Fatalf("Failed to load font: %v", err)
	}
	a.Settings().SetTheme(myTheme)
	w := a.NewWindow("Wallet 钱包")

	MainContent(w)

	w.Resize(mainSize)

	w.Show()
	a.Run()
	exit()
}

func MainContent(w fyne.Window) {
	getAddressBalaceTab := container.NewTabItem("Get address balance", GetAddressBalance(w))
	tabs := container.NewAppTabs(
		container.NewTabItem("Generate wallet", GenerateWallet()),
		getAddressBalaceTab,
		container.NewTabItem("Get address UnUTXO list", GetAddressUTXO()),
		container.NewTabItem("Sync address UTXO list", SyncAddressUTXO()),
		container.NewTabItem("Transaction info", TransactionInfo()),
		container.NewTabItem("Multi to multi transaction", MultiToMultiTransfer()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	// 设置窗口的内容
	mainPageInfo = container.NewVBox(tabs)
	mainPageInfo.Resize(mainSize)
	w.SetContent(mainPageInfo)
}

func GenerateWallet() *fyne.Container {
	var (
		// 助记词
		mnemonic = ""
		// 助记词数量
		mnemonicLen = []string{"12", "24"}
		// TODO: 助记词显示语言
		langs = []string{"EN", "CN_S", "CN_T"}
		// 网络可选列表
		networks = []string{"BTC", "BTCTest", "BTCRegT", "ETH", "TRON"}
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
	language := widget.NewSelect(langs, func(value string) {
		fmt.Println("Select set to", value)
	})
	language.SetSelected("EN")
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
		widget.NewLabel("Default language is EN"),
		language,
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
		if err := priKey.Set(accountInfo.PrivateKey); err != nil {
			fmt.Println(err)
		}
		if err := address.Set(accountInfo.Address); err != nil {
			fmt.Println(err)
		}
	})
	button := container.New(layout.NewGridLayout(2), btn1, btn2)
	/* ------------------------------- BUTTON ------------------------------- */
	return container.NewBorder(nil, button, left, right, content)
}

func getCenter(data string) *fyne.Container {
	return container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), widget.NewLabel(data))
}

func GetAddressBalance(w fyne.Window) *fyne.Container {
	// 创建一个按钮和一个标签
	myLabel := widget.NewLabel("Please choose Chain")

	// 创建一个下拉菜单，内容为 BTC 和 ETH
	options := []string{"BTC", "ETH"}
	selectOption := widget.NewSelect(options, func(selected string) {
		// 创建一个新的标签页容器，并添加标签页
		tabContainer := GetAddressUTXO()
		// 设置窗口的内容为新的标签页容器
		w.SetContent(tabContainer)
	})
	// 将按钮和标签放入一个垂直布局中
	getBalancePage = container.NewVBox(myLabel, selectOption)

	return getBalancePage
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
		res2 := wallet_sdk.GetUTXOListByAddress(global.ChainName, addr)
		sum := int64(0)
		checkSum := int64(0)
		utxoList := res2.Data.([]*client.UnspendUTXOList)
		resultContainer.Objects = nil
		// 排序未花费的UTXO
		client.DescSortUnspendUTXO(utxoList)
		for _, unspentUTXO := range utxoList {
			// UTXO展示内容
			a := unspentUTXO.Amount.CoefficientInt64()
			val := utils.EncodeStringByUtxoInfo(unspentUTXO.TxHash, unspentUTXO.Vout, a)
			sum += a
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
	chainName := global.ChainName
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
	limitedSizeContainer := container.NewStack(
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
			inAmount += v.Amount.CoefficientInt64()
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
			if err := alert.Set("Please enter the change address!"); err != nil {
				fmt.Println(err)
			}
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
		if res1.Status.Code == wallet_sdk.RES_CODE_FAILED {
			if err := alert.Set(res1.Status.Message); err != nil {
				fmt.Println(err)
			}
			return
		}
		signData = res1.Data
		size := len(res1.Data) / 2
		// 提示交易数据
		transferInfo := fmt.Sprintf("Get BTC transferInfo\n[In amount] %v\n[Out amount]%s BTC\n[Gas price] %s BTC/vKB, size: %v\n[Flinally fee] %v/1000*%v", utils.Int64ToSatoshi(inAmount), outAmount.String(), gasPrice, size, gasPrice, size)
		if err := alert.Set(transferInfo); err != nil {
			fmt.Println(err)
		}
	})
	SignBtn := widget.NewButton("2.SignTransaction", func() {
		priKey := inputM.Text
		// 没填私钥报错
		if inputM.Text == "" {
			if err := alert.Set("Please enter the address private key!"); err != nil {
				fmt.Println(err)
			}
			return
		}
		fmt.Printf("wch---- sign: %+v\n", signData)
		res1 := wallet_sdk.SignTransferInfo(chainName, priKey, signData)
		fmt.Printf("wch------ res1 data: %+v\n", res1.Data)
		// 提示签名数据
		if res1.Status.Code == wallet_sdk.RES_CODE_FAILED {
			if err := alert.Set(res1.Status.Message); err != nil {
				fmt.Println(err)
			}
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
			if err := alert.Set("Signature Successful\n" + signData); err != nil {
				fmt.Println(err)
			}
		}
		signData = res1.Data
	})
	BroadcastBtn := widget.NewButton("3.BroadcastTransaction", func() {
		res1 := wallet_sdk.BroadcastTransaction(chainName, signData)
		fmt.Printf("wch------ res1 data: %+v\n", res1.Data)
		// 提示交易HASH
		if res1.Status.Code == wallet_sdk.RES_CODE_FAILED {
			if err := alert.Set(res1.Status.Message); err != nil {
				fmt.Println(err)
			}
			return
		} else {
			if err := alert.Set(res1.Data); err != nil {
				fmt.Println(err)
			}
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
			if err := str.Set(alertStr); err != nil {
				fmt.Println(err)
			}
		} else {
			if err := str.Set(res5.Status.Message); err != nil {
				fmt.Println(err)
			}
		}

	})
	btn2 := widget.NewButton("Sign&Broadcast", func() {
		if priKey.Text == "" || signData == "" {
			if err := str.Set("Please check what you entered!"); err != nil {
				fmt.Println(err)
			}
		}
		res7 := wallet_sdk.SignTransferInfo(chainName, priKey.Text, signData)
		if err := str.Set(res7.Data); err != nil {
			fmt.Println(err)
		}
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
		gasPriceData := wallet_sdk.GetGasPrice(global.ChainName)
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

func SyncAddressUTXO() *fyne.Container {
	tip := widget.NewLabel("Enter address to query UTXO")
	// 结果信息
	leftLabel := widget.NewLabel("")
	// 地址输入框
	addressInput := widget.NewEntry()
	addressInput.SetPlaceHolder("Enter address")
	// 获取最新块高
	startHeightInput := widget.NewEntry()
	startHeightInput.SetPlaceHolder("Enter start height")
	// 请求按钮
	checkInfo := widget.NewButton("Sync UTXO for address", func() {
		addr := addressInput.Text
		fmt.Printf("wch---- addr: %+v\n", addr)
		srv = wallet_sdk.NewGetUtxoInfo(addr)
		// 检查是否存在历史块高
		startHeight := srv.GetUserHeightByAddress()
		// 获取最新块高
		res1 := wallet_sdk.GetBlockHeight(global.ChainName)
		newHigh, err := strconv.ParseInt(res1.Data, 0, 64)
		if err != nil {
			fmt.Println(err)
			leftLabel.SetText(err.Error())
			return
		}
		inputH := int64(0)
		// 优先输入框
		if startHeightInput.Text != "" {
			inputH, _ = strconv.ParseInt(startHeightInput.Text, 0, 64)
			startHeight = inputH
		}
		// 其次历史数据，最后从最新块开始同步
		if startHeight == 0 && inputH == 0 {
			startHeight = newHigh
		}
		fmt.Printf("startHeight: %+v\n", startHeight)
		go srv.GetTransferByBlockHeight(startHeight, newHigh)
	})
	button := container.New(layout.NewGridLayout(2), checkInfo)
	// 顶部提示
	top := container.NewVBox(tip, addressInput, startHeightInput)
	// 侧边统计
	left := container.NewVBox(leftLabel)
	return container.NewBorder(top, button, left, nil, nil)
}

// 退出应用后调用 Run()方法不会执行后续的代码
func exit() {
	fmt.Println("Wait Stop server...")
	if srv != nil {
		srv.Stop = true
		srv.Wg.Wait()
		for {
			if srv.StopOK {
				break
			}
		}
	}
	a.Quit()
	fmt.Println("Exited")
}
