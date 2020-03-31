package crawler

import (
	"github.com/araddon/dateparse"
	"regexp"
	"strings"
	"time"
)

func DateParser(s string) string {

	// 1. Make strings to Lower case
	ls := strings.ToLower(s)

	//2. updated 부분은 날리고 published 날짜만 저장한다
	updR, _ := regexp.Compile("\\bupd")

	findIdx := updR.FindStringIndex(ls)
	if len(findIdx) != 0 {
		if findIdx[0] < 2 {
			ls = ls[findIdx[1]:]
		} else {
			ls = ls[:findIdx[0]]
		}
	}

	//3. find month & day 대부분 월, 일은 붙어 있으니 함께 찾도록 한다
	monR, _ := regexp.Compile("\\bjan|\\bfeb|\\bmar|\\bapr|\\bmay|\\bjun|\\bjul|\\baug|\\bsep|\\boct|\\bnov|\\bdec")
	monIdx := monR.FindStringIndex(ls)
	month := ls[monIdx[0]:monIdx[1]]

	dayR, _ := regexp.Compile("[0-9]{1,2}")
	day_idx := dayR.FindStringIndex(ls[monIdx[1]:])
	day := ls[monIdx[1]+day_idx[0] : monIdx[1]+day_idx[1]]

	ls = ls[:monIdx[0]] + ls[monIdx[1]+day_idx[1]:]

	//4. find year
	yearR, _ := regexp.Compile("[0-9]{4}")
	yearIdx := yearR.FindStringIndex(ls)
	year := ls[yearIdx[0]:yearIdx[1]]

	ls = ls[:yearIdx[0]] + ls[yearIdx[1]:]

	//5. find time 00:00 의 형태를 찾고 am,pm 을 찾는다
	timeR, _ := regexp.Compile("[0-9]{1,2}:[0-9]{1,2}")
	amR, _ := regexp.Compile("a.?m.?")
	pmR, _ := regexp.Compile("p.?m.?")

	t := timeR.FindString(ls)
	isAm := amR.MatchString(ls)
	isPm := pmR.MatchString(ls)

	//6. 동일한 형태로 저장한다. (월 일 년도 시간(am/pm))
	tot := []string{month, day, year, t}
	date := strings.Join(tot, " ")
	if isAm {
		date = date + " am"
	} else if isPm {
		date = date + " pm"
	}

	//7. dateparse 모듈을 돌려 보기 좋은 형태로 최종 변환한다. (github.com/araddon/dateparse)
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	timeVar, err := dateparse.ParseIn(date, loc)
	if err != nil {
		panic(err)
	}
	return timeVar.String()
}
