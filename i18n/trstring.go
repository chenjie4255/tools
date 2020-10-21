package i18n

import (
	"errors"
	"reflect"
	"strings"

	"golang.org/x/text/language"
)

var supportLangs = []language.Tag{
	language.AmericanEnglish,
	language.SimplifiedChinese,
	language.TraditionalChinese,
	language.Spanish,
	language.Japanese,
	language.Korean,
	language.Russian,
}

var matcher language.Matcher

func init() {
	matcher = language.NewMatcher(supportLangs)
}

// TrString 多语言动态适配
type TrString struct {
	En     string `json:"en" bson:"en"`
	ZhHans string `json:"zh-Hans" bson:"zh-hans"`
	ZhHant string `json:"zh-Hant" bson:"zh-hant"`
	Ja     string `json:"ja" bson:"ja"`
	Es     string `json:"es" bson:"es"`
	Ko     string `json:"ko" bson:"ko"`
	Ru     string `json:"ru" bson:"ru"`
	Tr     string `json:"_tr" bson:"-"`
}

func (s *TrString) SetLang(langs ...string) {
	t, _ := language.MatchStrings(matcher, langs...)
	// GOLANG'S BUG: https://github.com/golang/go/issues/24211
	_, i, _ := matcher.Match(t)
	switch strings.ToLower(supportLangs[i].String()) {
	case "en-us":
		if s.En != "" {
			*s = TrString{Tr: s.En}
		}
	case "zh-hans":
		if s.ZhHans != "" {
			*s = TrString{Tr: s.ZhHans}
		}
	case "zh-hant":
		if s.ZhHant != "" {
			*s = TrString{Tr: s.ZhHant}
		}
	case "es":
		if s.Es != "" {
			*s = TrString{Tr: s.Es}
		}
	case "ja":
		if s.Ja != "" {
			*s = TrString{Tr: s.Ja}
		}
	case "ko":
		if s.Ko != "" {
			*s = TrString{Tr: s.Ko}
		}
	case "ru":
		if s.Ru != "" {
			*s = TrString{Tr: s.Ru}
		}
	default:
		if s.En != "" {
			*s = TrString{Tr: s.En}
		}
	}
}

type trObject interface {
	SetLang(langs ...string)
}

func SetLangForObject(obj interface{}, lang ...string) error {
	if len(lang) == 0 {
		return errors.New("lang empty")
	}

	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Slice {
		return errors.New("object must be a struct pointer or slice")
	}

	return setLangForReflectVal(v, lang...)
}

func setLangForReflectVal(val reflect.Value, lang ...string) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			setLangForReflectVal(item, lang...)
		}
		return nil
	}

	if val.Kind() != reflect.Struct {
		return errors.New("object must be a struct pointer or slice")
	}

	trObjectType := reflect.TypeOf(new(trObject)).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldT := field.Type()
		if fieldT.Implements(trObjectType) || reflect.PtrTo(fieldT).Implements(trObjectType) {
			if !field.CanSet() || !field.CanAddr() {
				// fmt.Printf("field(%s) cannot set or addr,%s, %v ,%v \n", val.Type().Field(i).Name, fieldT.Name(), field.CanSet(), field.CanAddr())
				continue
			}

			trObj, ok := field.Addr().Interface().(trObject)
			if !ok {
				panic("should implement trobject")
			}

			// fmt.Printf("field(%s) set lang...%s\n", val.Type().Field(i).Name, fieldT.Name())

			trObj.SetLang(lang...)
		} else if field.Kind() == reflect.Struct || field.Kind() == reflect.Slice || field.Kind() == reflect.Ptr {
			setLangForReflectVal(field, lang...)
		}
	}
	return nil
}

func BuildTrString(str string) TrString {
	ret := TrString{}
	ret.En = str
	ret.ZhHans = str
	ret.ZhHant = str
	ret.Ja = str
	ret.Es = str
	ret.Ko = str
	ret.Ru = str
	return ret
}

func BuildTrStringWithSuffix(str string) TrString {
	ret := TrString{}
	ret.En = str + "_en"
	ret.ZhHans = str + "_zh-hans"
	ret.ZhHant = str + "_zh-hant"
	ret.Ja = str + "ja"
	ret.Es = str + "es"
	ret.Ko = str + "ko"
	ret.Ru = str + "ru"
	return ret
}
