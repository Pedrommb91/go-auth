package reflection

import (
	"reflect"
	"testing"
)

type Credentials struct {
	password string `name:"pw"`
}

type User struct {
	Name        string      `name:"name"`
	ID          int         `name:"id"`
	Credentials Credentials `name:"credentials" ref:"cred"`
}

func TestGetType(t *testing.T) {
	type args struct {
		myvar interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get type of User",
			args: args{
				myvar: User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
			},
			want: "User",
		},
		{
			name: "Get type of User as a pointer",
			args: args{
				myvar: &User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
			},
			want: "*User",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetType(tt.args.myvar); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllFields(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Get the name of all fields of the struct",
			args: args{
				s: User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
			},
			want: []string{"Name", "ID", "Credentials"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllFields(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllTags(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Get all tags",
			args: args{
				s: User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
			},
			want: []string{"name:\"name\"", "name:\"id\"", "name:\"credentials\" ref:\"cred\""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllTags(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllTagsWithName(t *testing.T) {
	type args struct {
		s    interface{}
		name string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Get all tags with name 'name'",
			args: args{
				s: User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
				name: "name",
			},
			want: []string{"name", "id", "credentials"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllTagsWithName(tt.args.s, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTagsWithName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllStructsWithTagName(t *testing.T) {
	type args struct {
		s    interface{}
		name string
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "Get struct with tag name 'ref'",
			args: args{
				s: User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
				name: "ref",
			},
			want: []any{
				Credentials{
					password: "strong-pw",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllStructsWithTagName(tt.args.s, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllStructsWithTagName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTagByTypeName(t *testing.T) {
	type args struct {
		s     interface{}
		field string
		tag   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Get the tag by type name",
			args: args{
				s: User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
				field: "Credentials",
				tag:   "name",
			},
			want: "credentials",
		},
		{
			name: "Return empty string if not found",
			args: args{
				s: User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
				field: "UnexistingField",
				tag:   "name",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTagByTypeName(tt.args.s, tt.args.field, tt.args.tag); got != tt.want {
				t.Errorf("GetTagByTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllValues(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "Get all user values",
			args: args{
				s: User{
					Name: "dummy-name",
					ID:   1,
					Credentials: Credentials{
						password: "strong-pw",
					},
				},
			},
			want: []any{
				"dummy-name",
				int64(1),
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllValues(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllValuesAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
