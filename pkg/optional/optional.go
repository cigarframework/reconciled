package optional

import "time"

func String(s string) *string {
	p := s
	return &p
}

func Time(t time.Time) *time.Time {
	return &t
}

func Now() *time.Time {
	return Time(time.Now())
}

func Bool(b bool) *bool {
	return &b
}

func True() *bool {
	return Bool(true)
}

func False() *bool {
	return Bool(false)
}

func Error(err error) *string {
	return String(err.Error())
}

func UseBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func UseString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
