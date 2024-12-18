package elastic

// 查询数据
func GetDataByFilter(_index string, _query interface{}) (*BaseResp, error) {
	filter := UrlFilter{
		Index: _index,
		Type:  "_search",
	}
	res := &BaseResp{}
	err := AskHttpJson(HttpPost, filter, _query, res)
	return res, err
}

// 使用ID查询数据
func GetDataById(_index, _id string) (*HitsInfo, error) {
	filter := UrlFilter{
		Index: _index,
		Id:    _id,
	}
	res := &HitsInfo{}
	err := AskHttp(filter, res)
	return res, err
}
