package main

import (
	//"encoding/json"
	"errors"
	//"reflect"
	//"io"
	"log"
)

// func search(query string, maxNumResults int, chunkSize int, items []searchItem, cache *Cache) ([]searchItem, int, error) {
// 	if chunkSize < 0 {
// 		return nil, 0, fmt.Errorf("Num results cannot be < 0")
// 	}

// 	if chunkSize > 0 && chunkSize < 100 {
// 		return nil, 0, fmt.Errorf("Requested num results must be > 100")
// 	}

// 	if chunkSize > 5000 {
// 		return nil, -1, fmt.Errorf("ChunkSize number of results requested exceeded")
// 	}

// 	if items == nil {
// 		tmp := new([]searchItem)
// 		items = *tmp
// 	}

// 	cursor := ""

// 	if chunkSize != 0 {
// 		query = query + "&count=" + strconv.Itoa(chunkSize)
// 	}

// 	count := 0
// 	var totalResults int

// 	for {
// 		if totalResults >= maxNumResults {
// 			break
// 		}
// 		log.Println("-------------------New search---------------------")
// 		if count >= maxNumResults {
// 			log.Println("A")
// 			break
// 		}

// 		var tmpItems scrapeItems
// 		url := IA_ScrapeBaseURL + query

// 		if cursor != "" {
// 			url = url + "&cursor=" + cursor
// 			log.Println("-----------Cursor:", cursor)
// 		}

// 		log.Println("search", url)

// 		err := getUrlJSON(url, nil, true, &tmpItems, cursor, cache)
// 		if err != nil {
// 			return nil, 0, err
// 		}

// 		err = cleanSearchItems(tmpItems.Items)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		count += tmpItems.Count
// 		totalResults = tmpItems.Total

// 		log.Println(len(tmpItems.Items))

// 		items = append(items, tmpItems.Items...)

// 		if tmpItems.Cursor == "" {
// 			log.Println("-----------EMPTY Cursor********************************************")
// 			break
// 		}
// 		cursor = tmpItems.Cursor

// 	}
// 	return items, totalResults, nil
// }

type StringFields struct {
	s    *[]string
	sRaw interface{}
}

func cleanSearchItems(items []searchItem) error {
	log.Println(" cleanSearchItems")
	for i, _ := range items {
		item := &(items[i])
		//ints := []*[]int{&item.AvgRating}
		//ints := []*[]int{&item.AvgRating}
		//intsRaw := []*interface{}{&item.AvgRating_Raw}
		//intsRaw = []*interface{}{&item.AvgRating_Raw}

		sf := []StringFields{
			{&item.Title, item.Title_Raw},
			{&item.BackupLocation, item.BackupLocation_Raw},
			{&item.CurateNote, item.CurateNote_Raw},
			{&item.Curation, item.Curation_Raw},
			{&item.Format, item.Format_Raw},
			{&item.Date, item.Date_Raw},
		}

		//cleanInts(ints, intsRaw)
		err := cleanStrings(sf)
		if err != nil {
			log.Printf("------ %#v\n", item)
			log.Println(err)
			log.Fatal(item.Identifier)
		}

		if false {
			// ---------------
			if item.AvgRating_Raw != nil {
				avgi, ok := item.AvgRating_Raw.(int)
				if ok {
					item.AvgRating = []int{avgi}
				} else {
					avgai, ok := item.AvgRating_Raw.([]int)
					if ok {
						item.AvgRating = avgai
					} else {
						return errors.New("AvgRating_Raw not int or []int")
					}
				}
			}

			if item.Title_Raw != nil {
				titlei, ok := item.Title_Raw.(string)
				if ok {
					//log.Println("FFFFFFFFFFFFFFFFFFFFFFFFF", titlei)
					item.Title = []string{titlei}
					//log.Println(item.Title)
				} else {
					titleai, ok := item.Title_Raw.([]string)
					if ok {
						//log.Println("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZ", titleai)
						item.Title = titleai
					} else {

						return errors.New("Title_Raw not string or []string")
					}
				}
			}
			//items[i] = item
		}
	}
	return nil
}

func cleanInts(ints [][]int, intsRaw []interface{}) error {
	for i := 0; i < len(ints); i++ {
		// ---------------
		intv, ok := intsRaw[i].(int)
		if ok {
			ints[i] = []int{intv}
		} else {
			inta, ok := intsRaw[i].([]int)
			if ok {
				ints[i] = inta
			} else {
				return errors.New("Non int or []int")
			}
		}
	}
	return nil
}

func cleanStrings(sf []StringFields) error {
	for i := 0; i < len(sf); i++ {
		if sf[i].sRaw != nil {
			if v, ok := sf[i].sRaw.(string); ok {
				*sf[i].s = []string{v}
			} else {
				if inter, ok := sf[i].sRaw.([]interface{}); ok {
					*sf[i].s = make([]string, len(inter))

					for j := 0; j < len(inter); j++ {
						if v2, ok := inter[j].(string); ok {
							(*sf[i].s)[j] = v2
						}
					}
				}
			}
		}
	}
	return nil
}
