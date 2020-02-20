package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	input := `lab001;[unknown];__libc_start_main;main;write_log;__write 5
lab001;[unknown];__libc_start_main;main;write_log;debugdump;fsync 6
lab001;[unknown];__libc_start_main;main;write_log;debugdump;__write 23
lab001;[unknown];__libc_start_main;main;write_log;fsync 136
`
	p, err := convert(bytes.NewReader([]byte(input)))
	require.NoError(t, err)
	require.NotNil(t, p)

	require.Len(t, p.SampleType, 1)
	require.Equal(t, "samples", p.SampleType[0].Type)
	require.Equal(t, "count", p.SampleType[0].Unit)
	require.Len(t, p.Function, 8)
	require.Len(t, p.Location, 8)
	require.Len(t, p.Mapping, 1)
	require.Len(t, p.Sample, 4)

	require.Equal(t, p.Sample[0].Location[0], p.Location[0])
	require.Equal(t, "__write", p.Sample[0].Location[0].Line[0].Function.Name)
}
