package utils

import (
	"crypto/rand"
	"math/big"
)

// Iterate through the slice of strings, and generate random characters from each, and append them to the final byte slice
func RandomCharsFromSlice(sets []string, final_chars []byte) ([]byte, error) {
	for _, set := range sets {
		rand_index, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		if err != nil {
			return []byte{}, err
		}

		char := set[rand_index.Int64()]

		final_chars = append(final_chars, char)

	}

	return final_chars, nil

}

// Iterate through the given string, and append the random characters to the parent slice
func RandomCharsFromString(set string, parent []byte) ([]byte, error) {
	for range 8 {
		random_index, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		if err != nil {
			return []byte{}, err
		}

		parent = append(parent, set[int(random_index.Int64())])

	}
	return parent, nil

}

// Shuffle through the final string to make the pattern unrecognisable
func ShuffleGeneratedPass(generated []byte) ([]byte, error) {
	for i := len(generated) - 1; i > 0; i-- {
		random_index, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return []byte{}, err
		}

		j := int(random_index.Int64())

		generated[i], generated[j] = generated[j], generated[i]
	}
	return generated, nil
}
