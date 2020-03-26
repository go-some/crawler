package crawler

import(
  "time"
  "strings"
  "regexp"
  "github.com/araddon/dateparse"
)

func dateParser(s string) string{

  // 1. Make strings to Lower case
  ls := strings.ToLower(s)

  //2. updated 부분은 날리고 published 날짜만 저장한다
  upd_r, _ := regexp.Compile("\\bupd")

  find_idx := upd_r.FindStringIndex(ls)
  if len(find_idx) != 0 {
    ls = ls[:find_idx[0]]
  }

  //3. find month & day 대부분 월, 일은 붙어 있으니 함께 찾도록 한다
  mon_r, _ := regexp.Compile("\\bjan|\\bfeb|\\bmar|\\bapr|\\bmay|\\bjun|\\bjul|\\baug|\\bsep|\\boct|\\bnov|\\bdec")
  mon_idx := mon_r.FindStringIndex(ls)
  month := ls[mon_idx[0]:mon_idx[1]]

  day_r, _ := regexp.Compile("[0-9]{1,2}")
  day_idx := day_r.FindStringIndex(ls[mon_idx[1]:])
  day := ls[mon_idx[1]+day_idx[0]:mon_idx[1]+day_idx[1]]

  ls = ls[:mon_idx[0]] + ls[mon_idx[1]+day_idx[1]:]

  //4. find year
  year_r, _ := regexp.Compile("[0-9]{4}")
  year_idx := year_r.FindStringIndex(ls)
  year := ls[year_idx[0]:year_idx[1]]

  ls = ls[:year_idx[0]] + ls[year_idx[1]:]

  //5. find time 00:00 의 형태를 찾고 am,pm 을 찾는다
  time_r, _ := regexp.Compile("[0-9]{1,2}:[0-9]{1,2}")
  am_r,_ := regexp.Compile("a.?m.?")
  pm_r,_ := regexp.Compile("p.?m.?")

  t := time_r.FindString(ls)
  is_am := am_r.MatchString(ls)
  is_pm := pm_r.MatchString(ls)

  //6. 동일한 형태로 저장한다. (월 일 년도 시간(am/pm))
  tot := []string{month,day,year,t}
  date := strings.Join(tot," ")
  if is_am {
    date = date + " am"
  }
  if is_pm {
    date = date + " pm"
  }

  //7. dateparse 모듈을 돌려 보기 좋은 형태로 최종 변환한다. (github.com/araddon/dateparse)
  loc, err := time.LoadLocation("America/New_York")
  if err != nil {
    panic(err)
  }


  timeVar, err := dateparse.ParseIn(date,loc)
	if err != nil {
		panic(err)
	}
	return timeVar.String()
}
