package main

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"sort"
	"strings"
)

var db *leveldb.DB

func openDB() {
	var err error
	db, err = leveldb.OpenFile("database.db", nil)
	check(err)
}

func closeDB() {
	db.Close()
}

func getValueOrBlank(family string) FontRatings {
	fontRatings, err := getValueAsStruct(family)
	if err == nil {
		return fontRatings
	} else if strings.HasSuffix(err.Error(), ": not found") {
		fontRatings.Family = family
		return fontRatings
	} else {
		panic(err)
		return FontRatings{}
	}
}

func getValueAsStruct(key string) (FontRatings, error) {
	data, err := db.Get([]byte(key), nil)

	var finalVal FontRatings
	if err != nil {
		return FontRatings{}, err
	} else {
		json.Unmarshal(data, &finalVal)
		return finalVal, nil
	}
}

func getEntireDatabaseAsJSON() string {
	finalJSON, _ := json.MarshalIndent(getEntireDatabaseAsSlice(), "", "\t")
	return string(finalJSON)
}

func getEntireDatabaseAsSlice() []FontRatings {
	var slice []FontRatings
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()
		var structValue FontRatings
		json.Unmarshal(value, &structValue)
		slice = append(slice, structValue)
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].AveragePoints > slice[j].AveragePoints
	})
	return slice
}
