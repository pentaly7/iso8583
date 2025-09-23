package iso8583

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

var isoHeader = []byte("ISO")

type IsoPackager struct {
	HasHeader         bool                 `json:"hasHeader"`
	HeaderLength      int                  `json:"headerLength"`
	MessageKey        []int                `json:"messageKey"`
	PackagerConfig    map[string]BitConfig `json:"packagerConfig"` // from json
	MandatoryBit      []int                `json:"mandatoryBit"`
	IsoPackagerConfig [129]BitConfig
	PrefixLengths     [129]int // Pre-computed prefix lengths
	MaxLengths        [129]int // Pre-computed max lengths
}

type BitConfig struct {
	IsMandatory bool      `json:"isMandatory"`
	Type        BitType   `json:"type"`
	Length      BitLength `json:"length"`
}

func NewPackager(r io.Reader) (*IsoPackager, error) {
	buffer, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Join(err, ErrCreatingNewPackager)
	}

	var packager IsoPackager
	err = json.Unmarshal(buffer, &packager)
	if err != nil {
		return nil, errors.Join(err, ErrCreatingNewPackager)
	}

	packager.MandatoryBit = make([]int, 0)
	for k, v := range packager.PackagerConfig {
		key, err := strconv.Atoi(k)
		if err != nil {
			return nil, errors.Join(err, ErrCreatingNewPackager)
		}
		packager.IsoPackagerConfig[key] = v
		if v.IsMandatory {
			packager.MandatoryBit = append(packager.MandatoryBit, key)
		}
		// Pre-compute values for faster access
		packager.PrefixLengths[key] = v.Length.Type.GetPrefixLen()
		packager.MaxLengths[key] = v.Length.Max
	}

	// clear packager config that read from reader
	packager.PackagerConfig = nil

	return &packager, nil
}

func (p *IsoPackager) GetMandatoryBitsFromConfig() []int {
	mandatoryBit := make([]int, 0)
	for k, v := range p.IsoPackagerConfig {
		if v.IsMandatory {
			mandatoryBit = append(mandatoryBit, k)
		}
	}
	return mandatoryBit
}
