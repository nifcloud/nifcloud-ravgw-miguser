package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/nifcloud/nifcloud-sdk-go/service/computing"
	"github.com/nifcloud/nifcloud-sdk-go/service/computing/types"
)

type user struct{
    name string
    description string
}

func newOutput(users []user) *computing.DescribeRemoteAccessVpnGatewaysOutput {
    remoteUserSet := make([]types.RemoteUserSet, 0)
    for _, u := range users {
        //  to avoid appending the address of the same variable multiple times
        u := u 
        remoteUserSet = append(remoteUserSet, types.RemoteUserSet{
            Description:&u.description,
            UserName: &u.name,
        })
    }
    return &computing.DescribeRemoteAccessVpnGatewaysOutput{
        RemoteAccessVpnGatewaySet: []types.RemoteAccessVpnGatewaySetOfDescribeRemoteAccessVpnGateways {
            {
                RemoteUserSet: remoteUserSet, 
            },
        },
    }
}

func TestOutputToCsv(t *testing.T) {
    numMultiUser := 100
    tests := []struct {
        name string
        users [] user
    } {
        {
            name: "single user",
            users: []user{
                {name: "test_user", description: "test_description"},
            },
        },
        {
            name: fmt.Sprintf("%d users", numMultiUser),
            users: func() []user{
                users := make([]user, numMultiUser)
                for i := range users {
                    users[i] = user{
                        name: fmt.Sprintf("test_user_%d", i), 
                        description: fmt.Sprintf("test_description_%d", i),
                    }
                }
                return users
            }(),
        },
    }

    for _, tc := range tests {
        op := newOutput(tc.users)
        got := outputToCsv(op)

        want := make([][]string, len(tc.users)+1)
        want[0] = CSVHEADER
        for i := range tc.users {
            want[i+1] = []string{tc.users[i].name, "", tc.users[i].description}
        }

        if !reflect.DeepEqual(want, got) {
            t.Errorf("%s: expected: %v, got: %v", tc.name, want, got)
        }
    }
}

func TestCsvToInput(t *testing.T) {

    tests := []struct {
        name string
        csvData [][]string
        want error    
    } {
        {
            name: "no error",
            csvData: [][]string{
                CSVHEADER,
                {"test_user1", "password", "test_description"},
                {"test_user2", "password", ""},
                {"test_user3", "password", "test_description"},
            },
            want: nil,
        },
        {
            name: "error (empty name)",
            csvData: [][]string{
                CSVHEADER,
                {"", "password", "test_description"},
            },
            want: emptyUserNameError{"", "password", "test_description"},
        },
        {
            name: "error (empty password)",
            csvData: [][]string{
                CSVHEADER,
                {"test_user", "", "test_description"},
            },
            want: emptyPasswordError{"test_user", "", "test_description"},
        },
    }

    dummyRavgwId := "dummyRemoteAccessVpnGateway"

    for _, tc := range tests {
        _, err := csvToInput(tc.csvData, dummyRavgwId)

        if tc.want == nil && err == nil{
            return
        }
        if tc.want == nil {
            t.Fatalf("%s: expected: %t, got: nil", tc.name, tc.want)
        }
        if err == nil {
            t.Fatalf("%s: expected: nil, got: %t", tc.name, err)
        }
        if errors.Is(err, tc.want) {
            t.Fatalf("%s: expected: %t, got: %t", tc.name, tc.want, err)
        }
    }
}
