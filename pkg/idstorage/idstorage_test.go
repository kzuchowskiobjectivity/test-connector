package idstorage_test

import (
	"connectors/pkg/idstorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Fail(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
	}{
		{
			name:     "no file",
			filePath: "invalid/file/path",
		},

		{
			name:     "invalid json file",
			filePath: "testdata/invalid.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			allIds, err := idstorage.Load(idstorage.FromFile(tc.filePath))
			assert.Nil(t, allIds)
			assert.Error(t, err)
		})
	}

}

func TestLoad_ValidIds(t *testing.T) {
	testCases := []struct {
		name        string
		filePath    string
		expectedIds []idstorage.Id
	}{
		{
			name:     "single id",
			filePath: "testdata/simple.json",
			expectedIds: []idstorage.Id{
				{
					Owner:  "OwnerInfo",
					ApiKey: "TopSecretApiKey",
				},
			},
		},

		{
			name:     "multiple id",
			filePath: "testdata/complex.json",
			expectedIds: []idstorage.Id{
				{
					Owner:  "OwnerInfo_1",
					ApiKey: "TopSecretApiKey_1",
				},
				{
					Owner:  "OwnerInfo_2",
					ApiKey: "TopSecretApiKey_2",
				},
				{
					Owner:  "OwnerInfo_3",
					ApiKey: "TopSecretApiKey_3",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			allIds, err := idstorage.Load(idstorage.FromFile(tc.filePath))
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedIds, allIds)
		})
	}

}
