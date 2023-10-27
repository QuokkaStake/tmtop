package types

type ChainValidator struct {
	Moniker string
	Address string
}

type ChainValidators []ChainValidator

func (c ChainValidators) ToMap() map[string]ChainValidator {
	valsMap := make(map[string]ChainValidator, len(c))

	for _, validator := range c {
		valsMap[validator.Address] = validator
	}

	return valsMap
}
