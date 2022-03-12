package main

import (
	"fmt"

	gosearchengine "github.com/masamichhhhi/go-search-engine"
)

func main() {
	db, _ := gosearchengine.NewDBClient(
		&gosearchengine.DBConfig{
			User:     "root",
			Password: "password",
			Addr:     "127.0.0.1",
			Port:     "13306",
			DB:       "stalefish",
		},
	)
	storage := gosearchengine.NewStorageRdbImpl(db)
	analyzer := gosearchengine.NewAnalyzer(
		[]gosearchengine.CharFilter{},
		gosearchengine.NewStandardTokenizer(),
		[]gosearchengine.TokenFilter{
			gosearchengine.LowerCaseFilter{},
		},
	)

	// 転置インデックスに登録
	indexer := gosearchengine.NewIndexer(storage, analyzer, 1)
	for _, body := range []string{
		"Ruby PHP JS",
		"Go Ruby",
		"Ruby Go PHP",
		"Go PHP",
	} {
		indexer.AddDocument(
			gosearchengine.NewDocument(body),
		)
	}

	// or検索を実行
	sorter := gosearchengine.NewTfIdSorter(storage)
	mq := gosearchengine.NewMatchQuery(
		"GO Ruby",
		gosearchengine.OR,
		analyzer,
		sorter,
	)

	msearcher := mq.Searcher(storage)
	result, _ := msearcher.Search()
	fmt.Println(result)

	// フレーズ検索を実行
	pq := gosearchengine.NewPhraseQuery(
		"go RUBY", analyzer, nil)
	psearcher := pq.Searcher(storage)
	result, _ = psearcher.Search()
	fmt.Println(result)
}
