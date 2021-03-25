package centerlog

import (
	"encoding/json"
	"fmt"
	"github.com/MashiroC/begonia/config"
	"github.com/olivere/elastic/v7"
)

func NewMsg() *Msg {
	return &Msg{
	}
}

// InitEs 初始化Es
func init() {
	var err error
	clientEs, err = elastic.NewClient(elastic.SetURL(config.C.Log.Es.Addr))
	if err != nil {
		panic(err)
	}
}

//PutMsg 增加Msg
func (m *Msg) PutMsg() error {
	_, err := clientEs.Index().
		Index("server").
		BodyJson(m).
		Do(ctx)
	fmt.Println(m)
	if err != nil {
		return err
	}
	return nil
}

// GetAllMsg 拿该服务的所有msg
func (m *Msg) GetAllMsg() ([]Msg, error) {
	q := elastic.NewTermQuery("server_name", m.ServerName)

	result, err := clientEs.
		Search("server").
		Query(q).
		Pretty(true).
		Do(ctx)

	if err != nil {
		//putErrorMsg(m.ServerName, err)
		return nil, err
	}

	return praiseResToMsg(result)
}

// 根据对应的Filed查
func (m *Msg) QueryField(field map[string]string) ([]Msg, error) {
	q := elastic.NewBoolQuery()
	for k, v := range field {
		q.Must(elastic.NewMatchQuery("fields."+k,v))
	}
	result, err := clientEs.
		Search("server").
		Query(q).
		Pretty(true).
		Do(ctx)

	if err != nil {
		//putErrorMsg(m.ServerName, err)
		return nil, err
	}
	return praiseResToMsg(result)
}

func praiseResToMsg(result *elastic.SearchResult) ([]Msg, error) {
	var messages []Msg
	if result.TotalHits() > 0 {
		for _, v := range result.Hits.Hits {
			var msg Msg
			if err := json.Unmarshal(v.Source, &msg); err != nil {
				return messages, err
			}
			messages = append(messages, msg)
		}
	}
	return messages, nil
}
