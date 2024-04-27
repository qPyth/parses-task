package types

import (
	"fmt"
	"strconv"
)

type Influencer struct {
	Rank      int
	Info      Info
	Category  []string
	Followers string
	Country   string
	EngAuth   string
	EngAvg    string
}

type Info struct {
	IGUsername string
	Name       string
}

func (i Influencer) ToStringSlice() []string {
	var slice []string
	slice = append(slice, strconv.Itoa(i.Rank))
	slice = append(slice, fmt.Sprintf("%s\n%s", i.Info.IGUsername, i.Info.Name))

	var category string
	for _, s := range i.Category {
		category += s + "\n"
	}
	slice = append(slice, category, i.Followers, i.Country, i.EngAuth, i.EngAvg)

	return slice
}
