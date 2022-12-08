package query_test

import (
	"reflect"
	"testing"

	"github.com/maniizu3110/attendance/rpc/project/internal/repository/query"
)

func TestParseOrder(t *testing.T) {
	type parentStruct struct {
		Name string
	}
	type childStruct struct {
		Name string
		Age  int
		N1   parentStruct
		N2   *parentStruct
		N3   []parentStruct
	}
	var testStruct childStruct
	testcases := []struct {
		name          string
		member        string // メンバー
		expectedQuery string // 期待するクエリの結果
		exectedError  bool   // エラーを期待するならtrue
	}{
		{
			name:          "Name(ASC)",
			member:        "Name",
			expectedQuery: "name asc",
			exectedError:  false,
		},
		{
			name:          "Name(desc)",
			member:        "- Name",
			expectedQuery: "name desc",
			exectedError:  false,
		},
		{
			name:          "Age(asc)",
			member:        "Age",
			expectedQuery: "age asc",
			exectedError:  false,
		},
		{
			name:          "Age(desc)",
			member:        "- Age",
			expectedQuery: "age desc",
			exectedError:  false,
		},
		{
			name:          "Dummy(asc)",
			member:        "Dummy",
			expectedQuery: "",
			exectedError:  true,
		},
		{
			name:          "Dummy(desc)",
			member:        "-Dummy",
			expectedQuery: "",
			exectedError:  true,
		},
		{
			name:          "N1(asc)",
			member:        "N1",
			expectedQuery: "",
			exectedError:  true,
		},
		{
			name:          "N1(desc)",
			member:        "-N1",
			expectedQuery: "",
			exectedError:  true,
		},
		{
			name:          "N2(asc)",
			member:        "N2",
			expectedQuery: "",
			exectedError:  true,
		},
		{
			name:          "N2(desc)",
			member:        "-N2",
			expectedQuery: "",
			exectedError:  true,
		},
		{
			name:          "N3(asc)",
			member:        "N3",
			expectedQuery: "",
			exectedError:  true,
		},
		{
			name:          "N3(desc)",
			member:        "-N3",
			expectedQuery: "",
			exectedError:  true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			res, err := query.CreateOrder(reflect.TypeOf(testStruct), testcase.member)
			if err != nil {
				if !testcase.exectedError { // 期待されたエラーでないなら不正
					t.Errorf("You got unexpected error. Details of Error:%v\n", err)
				}
			} else {
				if testcase.exectedError { // エラーを期待されたのにエラーが出現していないなら不正
					t.Error("You got unexpexced result. You should recieve an error, but recieved no error.")
				} else if res != testcase.expectedQuery { // 結果が期待するものと違うなら不正
					t.Errorf("You got %s, but expected query is %s. Details of Error:%v\n", res, testcase.expectedQuery, err)
				}
			}
		})
	}
}

func TestParseOrderQuery(t *testing.T) {
	type parentStruct struct {
		Name string
	}
	type childStruct struct {
		Name string
		Age  int
		N1   parentStruct
		N2   *parentStruct
		N3   []parentStruct
	}
	testcases := []struct {
		name            string
		members         []string
		expectedMembers []string
		expextedError   bool
	}{
		{
			name:            "Name Only",
			members:         []string{"Name", "-Name"},
			expectedMembers: []string{"name asc", "name desc"},
			expextedError:   false,
		},
		{
			name:            "Age Only",
			members:         []string{"Age", "-Age"},
			expectedMembers: []string{"age asc", "age desc"},
			expextedError:   false,
		},
		{
			name:            "Name And Age",
			members:         []string{"Name", "-Age", "-Name", "Age"},
			expectedMembers: []string{"name asc", "age desc", "name desc", "age asc"},
			expextedError:   false,
		},
		{
			name:            "N1 Only",
			members:         []string{"N1", "-N1"},
			expectedMembers: []string{},
			expextedError:   true,
		},
		{
			name:            "N2 Only",
			members:         []string{"N2", "-N2"},
			expectedMembers: []string{},
			expextedError:   true,
		},
		{
			name:            "N3 Only",
			members:         []string{"N3", "-N3"},
			expectedMembers: []string{},
			expextedError:   true,
		},
		{
			name:            "Name, Age And N1",
			members:         []string{"Name", "-Name", "Age", "-Age", "N1", "-N1"},
			expectedMembers: []string{},
			expextedError:   true,
		},
		{
			name:            "Name, Age And N2",
			members:         []string{"Name", "-Name", "Age", "-Age", "N2", "-N2"},
			expectedMembers: []string{},
			expextedError:   true,
		},
		{
			name:            "Name, Age And N3",
			members:         []string{"Name", "-Name", "Age", "-Age", "N3", "-N3"},
			expectedMembers: []string{},
			expextedError:   true,
		},
	}
	var testStruct childStruct
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			res, err := query.CreateOrderArray(testStruct, testcase.members)
			if err != nil {
				if !testcase.expextedError { // 期待されたエラーでないなら不正
					t.Errorf("You got unexpected error. Details of Error:%v\n", err)
				}
			} else {
				if testcase.expextedError { // エラーを期待されたのにエラーが出現していないなら不正
					t.Error("You got unexpexced result. You should recieve an error, but recieved no error.")
				} else {
					for i := range res {
						if res[i] != testcase.expectedMembers[i] { // 結果が期待するものと違うなら不正
							t.Errorf("You got %s, but expected query is %s. Details of Error:%v\n", res[i], testcase.expectedMembers[i], err)
						}
					}
				}
			}
		})
	}
}
