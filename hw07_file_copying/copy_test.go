package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	srcPath     = "testdata/input.txt"
	testdataDir = "testdata"
)

var tmpDir string

type CopyTestSuite struct {
	suite.Suite
}

func (suite *CopyTestSuite) SetupSuite() {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal(err)
	}

	tmpDir = dir
}

func (suite *CopyTestSuite) TearDownSuite() {
	if err := os.RemoveAll(tmpDir); err != nil {
		log.Fatal(err)
	}
}

func (suite *CopyTestSuite) TestCopyWithValidParameters() {
	testCases := []struct {
		fileName      string
		limit, offset int64
	}{
		{
			fileName: "/out_offset0_limit0.txt",
		},
		{
			fileName: "/out_offset0_limit10.txt",
			limit:    10,
		},
		{
			fileName: "/out_offset0_limit1000.txt",
			limit:    1000,
		},
		{
			fileName: "/out_offset0_limit10000.txt",
			limit:    10000,
		},
		{
			fileName: "/out_offset100_limit1000.txt",
			limit:    1000,
			offset:   100,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(fmt.Sprintf("limit: %d offset: %d", tc.limit, tc.offset), func(t *testing.T) {
			dstPath := tmpDir + tc.fileName

			err := Copy(srcPath, dstPath, tc.offset, tc.limit)
			require.NoError(suite.T(), err)

			expected, _ := ioutil.ReadFile("testdata" + tc.fileName)
			actual, _ := ioutil.ReadFile(dstPath)

			require.Equal(t, expected, actual)
		})
	}
}

func (suite *CopyTestSuite) TestCopyWithInvalidParams() {
	testCases := []struct {
		errMessage, srcPath string
		limit, offset       int64
	}{
		{
			errMessage: ErrOffsetExceedsFileSize.Error(),
			srcPath:    testdataDir + "/out_offset6000_limit1000.txt",
			limit:      1000,
			offset:     6000,
		},
		{
			errMessage: ErrUnknownFileSize.Error(),
			srcPath:    "/dev/urandom",
			limit:      10,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.errMessage, func(t *testing.T) {
			err := Copy(tc.srcPath, "", tc.offset, tc.limit)
			require.EqualError(suite.T(), err, tc.errMessage)
		})
	}
}

func TestCopyTestSuite(t *testing.T) {
	suite.Run(t, new(CopyTestSuite))
}
