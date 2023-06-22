package slices

import (
	"reflect"
	"strings"
	"testing"
)

type User struct {
	name string
	id   int
}

func TestFilter(t *testing.T) {
	type args struct {
		ss   []User
		fill func(User) bool
	}
	tests := []struct {
		name    string
		args    args
		wantRet []User
	}{
		{
			name: "Filter by names starting with 'A'",
			args: args{
				ss: []User{
					{
						name: "John",
						id:   1,
					},
					{
						name: "Anna",
						id:   2,
					},
				},
				fill: func(u User) bool {
					return strings.HasPrefix(u.name, "A")
				},
			},
			wantRet: []User{
				{
					name: "Anna",
					id:   2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := Filter(tt.args.ss, tt.args.fill); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("Filter() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type args struct {
		ss     []User
		remove func(User) bool
	}
	tests := []struct {
		name    string
		args    args
		wantRet []User
	}{
		{
			name: "Remove names starting with 'A'",
			args: args{
				ss: []User{
					{
						name: "John",
						id:   1,
					},
					{
						name: "Anna",
						id:   2,
					},
				},
				remove: func(u User) bool {
					return strings.HasPrefix(u.name, "A")
				},
			},
			wantRet: []User{
				{
					name: "John",
					id:   1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := Remove(tt.args.ss, tt.args.remove); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("Remove() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		ss     []User
		exists func(User) bool
	}
	tests := []struct {
		name   string
		args   args
		index  int
		exists bool
	}{
		{
			name: "Contains the name John",
			args: args{
				ss: []User{
					{
						name: "John",
						id:   1,
					},
					{
						name: "Anna",
						id:   2,
					},
				},
				exists: func(u User) bool {
					return u.name == "John"
				},
			},
			index:  0,
			exists: true,
		},
		{
			name: "Contains the name Carlos",
			args: args{
				ss: []User{
					{
						name: "John",
						id:   1,
					},
					{
						name: "Anna",
						id:   2,
					},
				},
				exists: func(u User) bool {
					return u.name == "Carlos"
				},
			},
			index:  -1,
			exists: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Contains(tt.args.ss, tt.args.exists)
			if got != tt.index {
				t.Errorf("Contains() got = %v, want %v", got, tt.index)
			}
			if got1 != tt.exists {
				t.Errorf("Contains() got1 = %v, want %v", got1, tt.exists)
			}
		})
	}
}
