package elastic

import "fmt"

// 插入数据
func CreateData(_index, _id string, _source interface{}) error {
	filter := UrlFilter{
		Index: _index,
		// Type:  "_create",
		Id: _id,
	}
	res := &HitsInfo{}
	// err := AskHttpJson(HttpPut, filter, _source, res)
	err := AskHttpJson(HttpPost, filter, _source, res)
	if res.Error != nil {
		return fmt.Errorf("CreateData error: [%s]%s", res.Error.Type, res.Error.Reason)
	}
	return err
}

// 插入数据(已存在则失败)
func CreateDataByPut(_index, _id string, _source interface{}) error {
	filter := UrlFilter{
		Index: _index,
		Type:  "_create",
		Id:    _id,
	}
	res := &HitsInfo{}
	err := AskHttpJson(HttpPut, filter, _source, res)
	if res.Error != nil {
		return fmt.Errorf("CreateDataByPut error: [%s]%s", res.Error.Type, res.Error.Reason)
	}
	return err
}

// 修改数据
func UpdateData(_index, _id string, _source interface{}) error {
	filter := UrlFilter{
		Index: _index,
		Type:  "_update",
		Id:    _id,
	}
	res := &HitsInfo{}
	err := AskHttpJson(HttpPost, filter, _source, res)
	if res.Error != nil {
		return fmt.Errorf("UpdateData error: [%s]%s", res.Error.Type, res.Error.Reason)
	}
	return err
}

// 删除数据
func DeleteData(_index, _id string, _source interface{}) error {
	filter := UrlFilter{
		Index: _index,
		Id:    _id,
	}
	res := &BaseResp{}
	err := AskHttpJson(HttpDelete, filter, _source, res)
	if res.Error != nil {
		return fmt.Errorf("DeleteData error: [%s]%s", res.Type, res.Reason)
	}
	return err
}
