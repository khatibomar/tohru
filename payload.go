package tohru

import (
	"encoding/json"
	"fmt"
)

type JsonPayload map[string]interface{}

func (p *JsonPayload) WithOrder(o order) error {
	if err := o.valid(); err != nil {
		return err
	}
	(*p)["_order_by"] = string(o)
	return nil
}

func (p *JsonPayload) WithOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("negative Offset")
	}
	(*p)["_offset"] = offset
	return nil
}

func (p *JsonPayload) WithLimit(limit int) error {
	if limit <= 0 {
		return fmt.Errorf("negative limit or zero")
	}
	(*p)["_limit"] = limit
	return nil
}

func (p *JsonPayload) WithListType(l listType) error {
	if err := l.valid(); err != nil {
		return err
	}
	(*p)["list_type"] = string(l)
	return nil
}
func (p *JsonPayload) WithJustInfo(info string) error {
	if info == "No" || info == "Yes" {
		(*p)["just_info"] = info
		return nil
	}
	return fmt.Errorf("value must be Yes or No")
}

func (p *JsonPayload) WithName(name string) {
	(*p)["anime_name"] = name
}

func (p *JsonPayload) WithSeason(s season) error {
	if err := s.valid(); err != nil {
		return err
	}
	(*p)["anime_season"] = string(s)
	return nil
}

func (p *JsonPayload) WithReleaseYear(year int) error {
	if year <= 0 {
		return fmt.Errorf("Year must be positive")
	}
	(*p)["anime_release_years"] = year
	return nil
}

func (p *JsonPayload) ToJson() (string, error) {
	json, err := json.Marshal(p)
	return string(json), err
}
