package elastic

import "github.com/shopspring/decimal"

type BaseResp struct {
	Took    int64       `json:"took"`
	TimeOut bool        `json:"time_out"`
	Shards  Shards      `json:"_shards"`
	Hits    Hits        `json:"Hits"`
	Error   interface{} `json:"error"`
	Type    string      `json:"type"`
	Reason  string      `json:"reason"`
}

type Hits struct {
	Total Total      `json:"total"` // 匹配到的文档总数
	Hits  []HitsInfo `json:"hits"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type HitsInfo struct {
	Index       string          `json:"_index"`   // 索引
	Type        string          `json:"_type"`    // type
	Id          string          `json:"_id"`      // 主键
	Score       decimal.Decimal `json:"_score"`   // 分数
	Source      interface{}     `json:"_source"`  // 内容
	Version     int             `json:"_version"` // 操作版本
	SeqNo       int64           `json:"_seq_no"`  // 插入序号
	PrimaryTerm int             `json:"_primary_term"`
	Found       bool            `json:"found"`
	Result      string          `json:"result"`
	Shards      Shards          `json:"_shards"`
	Error       *Error          `json:"error"error,omitempty`
}

type Shards struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Skipped    int64 `json:"skipped"`
	Failed     int64 `json:"failed"`
}

type Error struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type SearchInfo struct {
	From   int      `json:"from,omitempty"`
	Size   int      `json:"size,omitempty"`
	Source []string `json:"_source,omitempty"`
	Query  *Query   `json:"query"`
}

type Query struct {
	Match map[string]interface{} `json:"match,omitempty"`
	Term  map[string]interface{} `json:"term,omitempty"`
}

type UpdateInfo struct {
	Doc map[string]interface{} `json:"doc,omitempty"`
}
