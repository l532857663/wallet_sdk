package utils

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func EncodeStringByUtxoInfo(txHash string, vout uint32, amount int64) string {
	return fmt.Sprintf("%s:%d %d", txHash, vout, amount)
}

func DecodeUtxoInfoByString(val string) (string, uint32, int64) {
	var txHash string
	var vout uint32
	var amount int64
	utxoInfo := strings.Split(val, " ")
	if len(utxoInfo) == 0 {
		return txHash, vout, amount
	}
	utxo := strings.Split(utxoInfo[0], ":")
	if len(utxo) != 2 {
		return txHash, vout, amount
	}
	txHash = utxo[0]
	v, _ := strconv.ParseUint(utxo[1], 0, 0)
	amount, _ = strconv.ParseInt(utxoInfo[1], 0, 0)
	return txHash, uint32(v), amount
}

// filterData 用于根据查询字符串过滤数据
func FilterData(data []string, query string) []string {
	var result []string
	for _, item := range data {
		if strings.Contains(strings.ToLower(item), strings.ToLower(query)) {
			result = append(result, item)
		}
	}
	return result
}

// 自定义主题以使用加载的字体
type MyTheme struct {
	FontResource fyne.Resource
}

func NewMyTheme() (*MyTheme, error) {
	// 加载支持中文的字体文件
	chineseFont, err := fyne.LoadResourceFromPath("/Users/mac/Library/Fonts/Noto_Sans_SC/NotoSansSC-VariableFont_wght.ttf")
	if err != nil {
		return nil, err
	}
	if chineseFont == nil {
		return nil, errors.New("Font is nil")
	}
	return &MyTheme{FontResource: chineseFont}, nil
}

func (m *MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	// 强制返回所有文本使用当前字体
	// if s.Monospace {
	// 	return theme.DefaultTheme().Font(s)
	// }
	// if s.Bold {
	// 	if s.Italic {
	// 		return theme.DefaultTheme().Font(fyne.TextStyle{Bold: true, Italic: true})
	// 	}
	// 	return theme.DefaultTheme().Font(fyne.TextStyle{Bold: true})
	// }
	// if s.Italic {
	// 	return theme.DefaultTheme().Font(fyne.TextStyle{Italic: true})
	// }
	return m.FontResource
}

func (m *MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (m *MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (m *MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
