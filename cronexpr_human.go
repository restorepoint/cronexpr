package cronexpr

import (
	"fmt"
	"sort"
	"strings"
)

var (
	weekdayName = map[int]string{
		0: "Sunday",
		1: "Monday",
		2: "Tuesday",
		3: "Wednesday",
		4: "Thursday",
		5: "Friday",
		6: "Saturday",
	}
	monthName = map[int]string{
		1:  "January",
		2:  "February",
		3:  "March",
		4:  "April",
		5:  "May",
		6:  "June",
		7:  "July",
		8:  "August",
		9:  "September",
		10: "October",
		11: "November",
		12: "December",
	}
)

// English returns a string representation of expr
func (expr *Expression) English() (output string) {
	if expr == nil {
		return
	}
	output = "Every "
	if len(expr.hourList) > 0 && len(expr.minuteList) > 0 && len(expr.hourList) <= 2 && len(expr.minuteList) <= 2 {
		output = "At "
		var hm []string
		for _, h := range expr.hourList {
			for _, m := range expr.minuteList {
				hm = append(hm, fmt.Sprintf("%02d:%02d", h, m))
			}
		}
		output += joinList(hm)
		if len(expr.daysOfWeek) == 7 && len(expr.daysOfMonth) == 31 {
			output += " every day"
		}

	} else {
		if len(expr.hourList) != 24 {
			if len(expr.minuteList) != 60 {
				output += numberList(expr.minuteList) + " minute past the " + numberList(expr.hourList) + " hour"
			} else {
				output += "minute of " + numberList(expr.hourList) + " hour"
			}
		} else if len(expr.minuteList) != 60 {
			if len(expr.minuteList) == 1 && expr.minuteList[0] == 0 {
				output += "hour, on the hour"
			} else {
				output += numberList(expr.minuteList) + " minute past every hour"
			}
		} else if len(expr.secondList) != 60 {
			if len(expr.secondList) == 1 && expr.secondList[0] == 0 {
				output += "minute, on the minute"
			} else {
				output += numberList(expr.secondList) + " second past every minute"
			}
		} else {
			output += "second"
		}
	}

	if len(expr.daysOfMonth) < 31 {
		var dom []int
		for r, t := range expr.daysOfMonth {
			if t {
				dom = append(dom, r)
			}
		}
		sort.Ints(dom)
		output += " on the " + numberList(dom)
		if len(expr.monthList) == 12 {
			output += " of every month"
		}
	}

	l := len(expr.daysOfWeek)
	if l > 0 && l != 7 {
		if len(expr.actualDaysOfMonthList) > 0 {
			output += " and every "
		} else {
			output += " on "
		}
		var d []int
		for n, v := range expr.daysOfWeek {
			if v {
				d = append(d, n)
			}
		}
		output += dowList(d) + " every week"
	}

	if len(expr.monthList) != 12 {
		output += " in " + monthList(expr.monthList)
	}

	return
}

func numberList(n []int) (s string) {
	return makeList(n, ordinal)
}

func monthList(n []int) (s string) {
	return makeList(n, func(n int) string { return monthName[n] })
}

func dowList(n []int) (s string) {
	return makeList(n, func(n int) string { return weekdayName[n] })
}

func makeList(n []int, f func(int) string) (s string) {
	var o = make([]string, len(n))
	for i, v := range n {
		o[i] = f(v)
	}
	return joinList(o)
}

func joinList(n []string) string {
	if len(n) == 1 {
		return n[0]
	}
	return strings.Join(n[:len(n)-1], ", ") + " and " + n[len(n)-1]
}

func ordinal(n int) string {
	ext := "th"
	if n != 11 && n%10 == 1 {
		ext = "st"
	} else if n != 12 && n%10 == 2 {
		ext = "nd"
	} else if n != 13 && n%10 == 3 {
		ext = "rd"
	}
	return fmt.Sprintf("%d%s", n, ext)
}
