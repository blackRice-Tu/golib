package xcommon

import (
	"strings"

	"github.com/fatih/structs"
)

func GetStructTags(s any, tagName string) []string {
	tagList := make([]string, 0)

	fields := structs.Fields(s)
	for _, field := range fields {
		tagList = append(tagList, field.Tag(tagName))
	}
	return tagList
}

func GetStructJsonFields(s any) []string {
	rawTagList := GetStructTags(s, "json")
	tagList := make([]string, 0)
	for i, _ := range rawTagList {
		tag := strings.TrimSpace(rawTagList[i])
		if tag == "-" {
			continue
		}
		tagList = append(tagList, strings.Split(tag, ",")[0])
	}
	return tagList
}
