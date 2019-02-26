package cron

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/stringutil"
)

// ParseString parses a cron formatted string into a schedule.
// The string must be at least 5 components, whitespace separated.
// If the string has 5 components a 0 will be prepended for the seconds component, and a * appended for the year component.
// If the string has 6 components a * appended for the year component.
// (seconds) (minutes) (hours) (day of month) (month) (day of week) (year)
/*
	Field name     Mandatory?   Allowed values    Allowed special characters
	----------     ----------   --------------    --------------------------
	Seconds        No           0-59              * / , -
	Minutes        Yes          0-59              * / , -
	Hours          Yes          0-23              * / , -
	Day of month   Yes          1-31              * / , - L W
	Month          Yes          1-12 or JAN-DEC   * / , -
	Day of week    Yes          0-6 or SUN-SAT    * / , - L #
	Year           No           1970–2099         * / , -
*/
/*
You can also use shorthands for the cron string:
	@yearly is equivalent to "0 0 0 1 1 * *"
	@monthly is equivalent to "0 0 0 1 * * *"
	@weekly is equivalent to "0 0 0 * * 0 *"
	@daily is equivalent to "0 0 0 * * * *"
	@hourly is equivalent to "0 0 * * * * *"
	@every xyz will parse the `xyz` value as a duration and return an every schedule for that
*/
func ParseString(cronString string) (Schedule, error) {
	// escape shorthands.
	if shorthand, ok := StringScheduleShorthands[strings.TrimSpace(cronString)]; ok {
		cronString = shorthand
	}

	if strings.HasPrefix(cronString, "@every") {
		duration, err := time.ParseDuration(strings.TrimSpace(strings.TrimPrefix(cronString, "@every")))
		if err != nil {
			return nil, exception.New(ErrStringScheduleInvalid).WithInner(err)
		}
		return Every(duration), nil
	}

	parts := stringutil.SplitSpace(cronString)
	if len(parts) < 5 || len(parts) > 7 {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(ErrStringScheduleComponents).WithMessagef("provided string; %s", cronString)
	}
	// fill in optional components
	if len(parts) == 5 {
		parts = append([]string{"0"}, parts...)
		parts = append(parts, "*")
	} else if len(parts) == 6 {
		parts = append(parts, "*")
	}

	seconds, err := parsePart(parts[0], parseInt, below(60))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("seconds invalid")
	}

	minutes, err := parsePart(parts[1], parseInt, below(60))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("minutes invalid")
	}

	hours, err := parsePart(parts[2], parseInt, below(24))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("hours invalid")
	}

	days, err := parsePart(parts[3], parseInt, between(1, 32))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("days invalid")
	}

	months, err := parsePart(parts[4], parseMonth, between(1, 13))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("months invalid")
	}

	daysOfWeek, err := parsePart(parts[5], parseDayOfWeek, between(0, 7))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("days of week invalid")
	}

	years, err := parsePart(parts[6], parseInt, between(1970, 2100))
	if err != nil {
		return nil, exception.New(ErrStringScheduleInvalid).WithInner(err).WithMessage("years invalid")
	}

	schedule := &StringSchedule{
		Original:    cronString,
		Seconds:     seconds,
		Minutes:     minutes,
		Hours:       hours,
		DaysOfMonth: days,
		Months:      months,
		DaysOfWeek:  daysOfWeek,
		Years:       years,
	}
	return schedule, nil
}

// Error Constants
const (
	ErrStringScheduleInvalid         exception.Class = "cron: schedule string invalid"
	ErrStringScheduleComponents      exception.Class = "cron: must have at least (5) components space delimited; ex: '0 0 * * * * *'"
	ErrStringScheduleValueOutOfRange exception.Class = "cron: string schedule part out of range"
	ErrStringScheduleInvalidRange    exception.Class = "cron: range (from-to) invalid"
)

// String schedule shorthands labels
const (
	StringScheduleShorthandAnnually = "@annually"
	StringScheduleShorthandYearly   = "@yearly"
	StringScheduleShorthandMonthly  = "@monthly"
	StringScheduleShorthandWeekly   = "@weekly"
	StringScheduleShorthandDaily    = "@daily"
	StringScheduleShorthandHourly   = "@hourly"
)

// String schedule shorthand values
var (
	StringScheduleShorthands = map[string]string{
		StringScheduleShorthandAnnually: "0 0 0 1 1 * *",
		StringScheduleShorthandYearly:   "0 0 0 1 1 * *",
		StringScheduleShorthandMonthly:  "0 0 0 1 * * *",
		StringScheduleShorthandDaily:    "0 0 0 * * * *",
		StringScheduleShorthandHourly:   "0 0 * * * * *",
	}
)

// Interface assertions.
var (
	_ Schedule     = (*StringSchedule)(nil)
	_ fmt.Stringer = (*StringSchedule)(nil)
)

// StringSchedule is a schedule generated from a cron string.
type StringSchedule struct {
	Original string

	Seconds     []int
	Minutes     []int
	Hours       []int
	DaysOfMonth []int
	Months      []int
	DaysOfWeek  []int
	Years       []int
}

// String returns the original string schedule.
func (ss *StringSchedule) String() string {
	return ss.Original
}

// FullString returns a fully formed string representation of the schedule's components.
// It shows fields as expanded.
func (ss *StringSchedule) FullString() string {
	fields := []string{
		csvOfInts(ss.Seconds, "*"),
		csvOfInts(ss.Minutes, "*"),
		csvOfInts(ss.Hours, "*"),
		csvOfInts(ss.DaysOfMonth, "*"),
		csvOfInts(ss.Months, "*"),
		csvOfInts(ss.DaysOfWeek, "*"),
		csvOfInts(ss.Years, "*"),
	}
	return strings.Join(fields, " ")
}

// Next implements cron.Schedule.
func (ss *StringSchedule) Next(after time.Time) time.Time {
	working := after
	if after.IsZero() {
		working = Now()
	}
	original := working

	if len(ss.Years) > 0 {
		for _, year := range ss.Years {
			if year >= working.Year() {
				working = advanceYearTo(working, year)
				break
			}
		}
	}

	if len(ss.Months) > 0 {
		var didSet bool
		for _, month := range ss.Months {
			if time.Month(month) == working.Month() && working.After(original) {
				didSet = true
				break
			}
			if time.Month(month) > working.Month() {
				working = advanceMonthTo(working, time.Month(month))
				didSet = true
				break
			}
		}
		if !didSet {
			working = advanceYear(working)
			for _, month := range ss.Months {
				if time.Month(month) >= working.Month() {
					working = advanceMonthTo(working, time.Month(month))
					break
				}
			}
		}
	}

	if len(ss.DaysOfMonth) > 0 {
		var didSet bool
		for _, day := range ss.DaysOfMonth {
			if day == working.Day() && working.After(original) {
				didSet = true
				break
			}
			if day > working.Day() {
				working = advanceDayTo(working, day)
				didSet = true
				break
			}
		}
		if !didSet {
			working = advanceMonth(working)
			for _, day := range ss.DaysOfMonth {
				if day >= working.Day() {
					working = advanceDayTo(working, day)
					break
				}
			}
		}
	}

	if len(ss.DaysOfWeek) > 0 {
		var didSet bool
		for _, dow := range ss.DaysOfWeek {
			if dow == int(working.Weekday()) && working.After(original) {
				didSet = true
				break
			}
			if dow > int(working.Weekday()) {
				working = advanceDayBy(working, (dow - int(working.Weekday())))
				didSet = true
				break
			}
		}
		if !didSet {
			working = advanceToNextSunday(working)
			for _, dow := range ss.DaysOfWeek {
				if dow >= int(working.Weekday()) {
					working = advanceDayBy(working, (dow - int(working.Weekday())))
					break
				}
			}
		}
	}

	if len(ss.Hours) > 0 {
		var didSet bool
		for _, hour := range ss.Hours {
			if hour == working.Hour() && working.After(original) {
				didSet = true
				break
			}
			if hour > working.Hour() {
				working = advanceHourTo(working, hour)
				didSet = true
				break
			}
		}
		if !didSet {
			working = advanceDay(working)
			for _, hour := range ss.Hours {
				if hour >= working.Hour() {
					working = advanceHourTo(working, hour)
					break
				}
			}
		}
	}

	if len(ss.Minutes) > 0 {
		var didSet bool
		for _, minute := range ss.Minutes {
			if minute == working.Minute() && working.After(original) {
				didSet = true
				break
			}
			if minute > working.Minute() {
				working = advanceMinuteTo(working, minute)
				didSet = true
				break
			}
		}
		if !didSet {
			working = advanceHour(working)
			for _, minute := range ss.Minutes {
				if minute >= working.Minute() {
					working = advanceMinuteTo(working, minute)
					break
				}
			}
		}
	}

	if len(ss.Seconds) > 0 {
		var didSet bool
		for _, second := range ss.Seconds {
			if second == working.Second() && working.After(original) {
				didSet = true
				break
			}
			if second > working.Second() {
				working = advanceSecondTo(working, second)
				didSet = true
				break
			}
		}
		if !didSet {
			working = advanceMinute(working)
			for _, second := range ss.Hours {
				if second >= working.Second() {
					working = advanceSecondTo(working, second)
					break
				}
			}
		}
	}

	return working
}

func parsePart(values string, parser func(string) (int, error), validator func(int) bool) ([]int, error) {
	if values == string(cronSpecialStar) {
		return nil, nil
	}

	// check if we need to expand an "every" pattern
	if strings.HasPrefix(values, cronSpecialEvery) {
		return parseEvery(values, parseInt, validator)
	}

	components := strings.Split(values, string(cronSpecialComma))

	output := map[int]bool{}
	var component string
	for x := 0; x < len(components); x++ {
		component = components[x]
		if strings.Contains(component, string(cronSpecialDash)) {
			rangeValues, err := parseRange(values, parser, validator)
			if err != nil {
				return nil, err
			}

			for _, value := range rangeValues {
				output[value] = true
			}
			continue
		}

		part, err := parser(component)
		if err != nil {
			return nil, exception.New(err)
		}
		if validator != nil && !validator(part) {
			return nil, exception.New(err)
		}
		output[part] = true
	}
	return mapKeysToArray(output), nil
}

func parseEvery(values string, parser func(string) (int, error), validator func(int) bool) ([]int, error) {
	every, err := parser(strings.TrimPrefix(values, "*/"))
	if err != nil {
		return nil, exception.New(err)
	}
	if validator != nil && !validator(every) {
		return nil, exception.New(ErrStringScheduleValueOutOfRange)
	}

	var output []int
	for x := 0; x < 60; x = x + every {
		output = append(output, x)
	}
	return output, nil
}

func parseRange(values string, parser func(string) (int, error), validator func(int) bool) ([]int, error) {
	parts := strings.Split(values, string(cronSpecialDash))

	if len(parts) != 2 {
		return nil, exception.New(ErrStringScheduleInvalidRange).WithMessagef("invalid range: %s")
	}

	from, err := parser(parts[0])
	if err != nil {
		return nil, exception.New(err)
	}
	to, err := parser(parts[1])
	if err != nil {
		return nil, exception.New(err)
	}

	if validator != nil && !validator(from) {
		return nil, exception.New(ErrStringScheduleValueOutOfRange).WithMessage("invalid range from")
	}
	if validator != nil && !validator(to) {
		return nil, exception.New(ErrStringScheduleValueOutOfRange).WithMessage("invalid range to")
	}

	if from >= to {
		return nil, exception.New(ErrStringScheduleInvalidRange).WithMessage("invalid range; from greater than to")
	}

	var output []int
	for x := from; x <= to; x++ {
		output = append(output, x)
	}
	return output, nil
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func parseMonth(s string) (int, error) {
	if value, ok := validMonths[s]; ok {
		return value, nil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, exception.New(err).WithMessage("month not a valid integer")
	}
	if value < 1 || value > 12 {
		return 0, exception.New(ErrStringScheduleValueOutOfRange).WithMessagef("month out of range (1-12): %s", s)
	}
	return value, nil
}

func parseDayOfWeek(s string) (int, error) {
	if value, ok := validDaysOfWeek[s]; ok {
		return value, nil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, exception.New(err).WithMessage("day of week not a valid integer")
	}
	if value < 0 || value > 6 {
		return 0, exception.New(ErrStringScheduleValueOutOfRange).WithMessagef("day of week out of range (0-6): %s", s)
	}
	return value, nil
}

func below(max int) func(int) bool {
	return between(0, max)
}

func between(min, max int) func(int) bool {
	return func(value int) bool {
		return value >= min && value < max
	}
}

func mapKeysToArray(values map[int]bool) []int {
	output := make([]int, len(values))
	var index int
	for key := range values {
		output[index] = key
		index++
	}
	sort.Ints(output)
	return output
}

//
// time helpers
//

func advanceYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location()).AddDate(1, 0, 0)
}

func advanceYearTo(t time.Time, year int) time.Time {
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

func advanceMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, 0)
}

func advanceMonthTo(t time.Time, month time.Month) time.Time {
	return time.Date(t.Year(), month, 1, 0, 0, 0, 0, t.Location())
}

func advanceDayTo(t time.Time, day int) time.Time {
	return time.Date(t.Year(), t.Month(), day, 0, 0, 0, 0, t.Location())
}

func advanceToNextSunday(t time.Time) time.Time {
	daysUntilSunday := 7 - int(t.Weekday())
	return t.AddDate(0, 0, daysUntilSunday)
}

func advanceDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
}

func advanceDayBy(t time.Time, days int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, days)
}

func advanceHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location()).Add(time.Hour)
}

func advanceHourTo(t time.Time, hour int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())
}

func advanceMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location()).Add(time.Minute)
}

func advanceMinuteTo(t time.Time, minute int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, 0, 0, t.Location())
}

func advanceSecondTo(t time.Time, second int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), second, 0, t.Location())
}

func csvOfInts(values []int, placeholder string) string {
	if len(values) == 0 {
		return placeholder
	}
	valueStrings := make([]string, len(values))
	for x := 0; x < len(values); x++ {
		valueStrings[x] = strconv.Itoa(values[x])
	}
	return strings.Join(valueStrings, ",")
}

// these are special characters
const (
	cronSpecialComma    = ',' //
	cronSpecialDash     = '-'
	cronSpecialStar     = '*'
	cronSpecialSlash    = '/'
	cronSpecialQuestion = '?' // sometimes used as the startup time, sometimes as a *

	cronSpecialLast       = 'L'
	cronSpecialWeekday    = 'W' // nearest weekday to the given day of the month
	cronSpecialDayOfMonth = '#' //

	cronSpecialEvery = "*/"
)

var (
	validMonths = map[string]int{
		"JAN": 1,
		"FEB": 2,
		"MAR": 3,
		"APR": 4,
		"MAY": 5,
		"JUN": 6,
		"JUL": 7,
		"AUG": 8,
		"SEP": 9,
		"OCT": 10,
		"NOV": 11,
		"DEC": 12,
	}

	validDaysOfWeek = map[string]int{
		"SUN": 0,
		"MON": 1,
		"TUE": 2,
		"WED": 3,
		"THU": 4,
		"FRI": 5,
		"SAT": 6,
	}
)
