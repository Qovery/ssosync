// Copyright (c) 2020, Amazon.com, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package internal ...
package internal

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/awslabs/ssosync/internal/aws"
	"github.com/awslabs/ssosync/internal/config"
	"github.com/awslabs/ssosync/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	admin "google.golang.org/api/admin/directory/v1"
)

// toJSON return a json pretty of the stc
func toJSON(stc interface{}) []byte {
	JSON, err := json.MarshalIndent(stc, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return JSON
}

func Test_getGroupOperations(t *testing.T) {
	type args struct {
		awsGroups    []*aws.Group
		googleGroups []*admin.Group
	}
	tests := []struct {
		name       string
		args       args
		wantAdd    []*aws.Group
		wantDelete []*aws.Group
		wantEquals []*aws.Group
	}{
		{
			name: "equal groups google and aws",
			args: args{
				awsGroups: []*aws.Group{
					aws.NewGroup("Group-1"),
					aws.NewGroup("Group-2"),
				},
				googleGroups: []*admin.Group{
					{Name: "Group-1"},
					{Name: "Group-2"},
				},
			},
			wantAdd:    nil,
			wantDelete: nil,
			wantEquals: []*aws.Group{
				aws.NewGroup("Group-1"),
				aws.NewGroup("Group-2"),
			},
		},
		{
			name: "add two new aws groups",
			args: args{
				awsGroups: nil,
				googleGroups: []*admin.Group{
					{Name: "Group-1"},
					{Name: "Group-2"},
				},
			},
			wantAdd: []*aws.Group{
				aws.NewGroup("Group-1"),
				aws.NewGroup("Group-2"),
			},
			wantDelete: nil,
			wantEquals: nil,
		},
		{
			name: "delete two aws groups",
			args: args{
				awsGroups: []*aws.Group{
					aws.NewGroup("Group-1"),
					aws.NewGroup("Group-2"),
				}, googleGroups: nil,
			},
			wantAdd: nil,
			wantDelete: []*aws.Group{
				aws.NewGroup("Group-1"),
				aws.NewGroup("Group-2"),
			},
			wantEquals: nil,
		},
		{
			name: "add one, delete one and one equal",
			args: args{
				awsGroups: []*aws.Group{
					aws.NewGroup("Group-2"),
					aws.NewGroup("Group-3"),
				},
				googleGroups: []*admin.Group{
					{Name: "Group-1"},
					{Name: "Group-2"},
				},
			},
			wantAdd: []*aws.Group{
				aws.NewGroup("Group-1"),
			},
			wantDelete: []*aws.Group{
				aws.NewGroup("Group-3"),
			},
			wantEquals: []*aws.Group{
				aws.NewGroup("Group-2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdd, gotDelete, gotEquals := getGroupOperations(tt.args.awsGroups, tt.args.googleGroups)
			if !reflect.DeepEqual(gotAdd, tt.wantAdd) {
				t.Errorf("getGroupOperations() gotAdd = %s, want %s", toJSON(gotAdd), toJSON(tt.wantAdd))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("getGroupOperations() gotDelete = %s, want %s", toJSON(gotDelete), toJSON(tt.wantDelete))
			}
			if !reflect.DeepEqual(gotEquals, tt.wantEquals) {
				t.Errorf("getGroupOperations() gotEquals = %s, want %s", toJSON(gotEquals), toJSON(tt.wantEquals))
			}
		})
	}
}

func Test_getUserOperations(t *testing.T) {
	type args struct {
		awsUsers    []*aws.User
		googleUsers []*admin.User
	}
	tests := []struct {
		name       string
		args       args
		wantAdd    []*aws.User
		wantDelete []*aws.User
		wantUpdate []*aws.User
		wantEquals []*aws.User
	}{
		{
			name: "equal user google and aws",
			args: args{
				awsUsers: []*aws.User{
					aws.NewUser("name-1", "lastname-1", "user-1@email.com", true),
					aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
				},
				googleUsers: []*admin.User{
					{Name: &admin.UserName{
						GivenName:  "name-1",
						FamilyName: "lastname-1",
					},
						Suspended:    false,
						PrimaryEmail: "user-1@email.com",
					},
					{Name: &admin.UserName{
						GivenName:  "name-2",
						FamilyName: "lastname-2",
					},
						Suspended:    false,
						PrimaryEmail: "user-2@email.com",
					},
				},
			},
			wantAdd:    nil,
			wantDelete: nil,
			wantUpdate: nil,
			wantEquals: []*aws.User{
				aws.NewUser("name-1", "lastname-1", "user-1@email.com", true),
				aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
			},
		},
		{
			name: "add two new aws users",
			args: args{
				awsUsers: nil,
				googleUsers: []*admin.User{
					{Name: &admin.UserName{
						GivenName:  "name-1",
						FamilyName: "lastname-1",
					},
						Suspended:    false,
						PrimaryEmail: "user-1@email.com",
					},
					{Name: &admin.UserName{
						GivenName:  "name-2",
						FamilyName: "lastname-2",
					},
						Suspended:    false,
						PrimaryEmail: "user-2@email.com",
					},
				},
			},
			wantAdd: []*aws.User{
				aws.NewUser("name-1", "lastname-1", "user-1@email.com", true),
				aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
			},
			wantDelete: nil,
			wantUpdate: nil,
			wantEquals: nil,
		},
		{
			name: "delete two aws users",
			args: args{
				awsUsers: []*aws.User{
					aws.NewUser("name-1", "lastname-1", "user-1@email.com", true),
					aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
				},
				googleUsers: nil,
			},
			wantAdd: nil,
			wantDelete: []*aws.User{
				aws.NewUser("name-1", "lastname-1", "user-1@email.com", true),
				aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
			},
			wantUpdate: nil,
			wantEquals: nil,
		},
		{
			name: "add on, delete one, update one and one equal",
			args: args{
				awsUsers: []*aws.User{
					aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
					aws.NewUser("name-3", "lastname-3", "user-3@email.com", true),
					aws.NewUser("name-4", "lastname-4", "user-4@email.com", true),
				},
				googleUsers: []*admin.User{
					{
						Name: &admin.UserName{
							GivenName:  "name-1",
							FamilyName: "lastname-1",
						},
						Suspended:    false,
						PrimaryEmail: "user-1@email.com",
					},
					{
						Name: &admin.UserName{
							GivenName:  "name-2",
							FamilyName: "lastname-2",
						},
						Suspended:    false,
						PrimaryEmail: "user-2@email.com",
					},
					{
						Name: &admin.UserName{
							GivenName:  "name-4",
							FamilyName: "lastname-4",
						},
						Suspended:    true,
						PrimaryEmail: "user-4@email.com",
					},
				},
			},
			wantAdd: []*aws.User{
				aws.NewUser("name-1", "lastname-1", "user-1@email.com", true),
			},
			wantDelete: []*aws.User{
				aws.NewUser("name-3", "lastname-3", "user-3@email.com", true),
			},
			wantUpdate: []*aws.User{
				aws.NewUser("name-4", "lastname-4", "user-4@email.com", false),
			},
			wantEquals: []*aws.User{
				aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdd, gotDelete, gotUpdate, gotEquals := getUserOperations(tt.args.awsUsers, tt.args.googleUsers)
			if !reflect.DeepEqual(gotAdd, tt.wantAdd) {
				t.Errorf("getUserOperations() gotAdd = %s, want %s", toJSON(gotAdd), toJSON(tt.wantAdd))
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("getUserOperations() gotDelete = %s, want %s", toJSON(gotDelete), toJSON(tt.wantDelete))
			}
			if !reflect.DeepEqual(gotUpdate, tt.wantUpdate) {
				t.Errorf("getUserOperations() gotUpdate = %s, want %s", toJSON(gotUpdate), toJSON(tt.wantUpdate))
			}
			if !reflect.DeepEqual(gotEquals, tt.wantEquals) {
				t.Errorf("getUserOperations() gotEquals = %s, want %s", toJSON(gotEquals), toJSON(tt.wantEquals))
			}
		})
	}
}

func Test_getGroupUsersOperations(t *testing.T) {
	type args struct {
		gGroupsUsers   map[string][]*admin.User
		awsGroupsUsers map[string][]*aws.User
	}
	tests := []struct {
		name       string
		args       args
		wantDelete map[string][]*aws.User
		wantEquals map[string][]*aws.User
	}{
		{
			name: "one add, one delete, one equal",
			args: args{
				gGroupsUsers: map[string][]*admin.User{
					"group-1": {
						{
							Name: &admin.UserName{
								GivenName:  "name-1",
								FamilyName: "lastname-1",
							},
							Suspended:    false,
							PrimaryEmail: "user-1@email.com",
						},
					},
				},
				awsGroupsUsers: map[string][]*aws.User{
					"group-1": {
						aws.NewUser("name-1", "lastname-1", "user-1@email.com", true),
						aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
					},
				},
			},
			wantDelete: map[string][]*aws.User{
				"group-1": {
					aws.NewUser("name-2", "lastname-2", "user-2@email.com", true),
				},
			},
			wantEquals: map[string][]*aws.User{
				"group-1": {
					aws.NewUser("name-1", "lastname-1", "user-1@email.com", true),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDelete, gotEquals := getGroupUsersOperations(tt.args.gGroupsUsers, tt.args.awsGroupsUsers)
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Errorf("getGroupUsersOperations() gotDelete = %s, want %s", toJSON(gotDelete), toJSON(tt.wantDelete))
			}
			if !reflect.DeepEqual(gotEquals, tt.wantEquals) {
				t.Errorf("getGroupUsersOperations() gotEquals = %s, want %s", toJSON(gotEquals), toJSON(tt.wantEquals))
			}
		})
	}
}

func Test_GetGroupsWithoutPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	groupId1 := "group-1-test-id"
	groupId2 := "group-2-test-id"
	groupId3 := "group-3-test-id"
	groupId4 := "group-4-test-id"
	displayName1 := "group-1-test-displayname"
	displayName2 := "group-2-test-displayname"
	displayName3 := "group-3-test-displayname"
	displayName4 := "group-4-test-displayname"

	// >1 group response with no pagination (<100 groups returned)
	sampleResponseNoPagination := &identitystore.ListGroupsOutput{
		Groups: []types.Group{
			{GroupId: &groupId1, DisplayName: &displayName1},
			{GroupId: &groupId2, DisplayName: &displayName2},
			{GroupId: &groupId3, DisplayName: &displayName3},
			{GroupId: &groupId4, DisplayName: &displayName4},
		},
	}

	expectedOutput := []*aws.Group{
		{ID: "group-1-test-id", Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}, DisplayName: "group-1-test-displayname", Members: []string{}},
		{ID: "group-2-test-id", Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}, DisplayName: "group-2-test-displayname", Members: []string{}},
		{ID: "group-3-test-id", Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}, DisplayName: "group-3-test-displayname", Members: []string{}},
		{ID: "group-4-test-id", Schemas: []string{"urn:ietf:params:scim:schemas:core:2.0:Group"}, DisplayName: "group-4-test-displayname", Members: []string{}},
	}

	mockIdentityStoreClient.EXPECT().ListGroups(gomock.Any(), gomock.Any()).MaxTimes(1).Return(sampleResponseNoPagination, nil)

	actualOutput, err := mockClient.GetGroups()

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
	assert.NoError(t, err)
}

func Test_GetGroupsWithPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	// >1 group response with pagination (150 groups returned)
	sampleGroupsResponseA := make([]types.Group, 0, 100)
	sampleGroupsResponseB := make([]types.Group, 0, 50)
	expectedOutput := make([]*aws.Group, 0, 150)

	// Populate responses
	for i := 0; i < 150; i++ {
		id := strconv.Itoa(i)
		grp := types.Group{
			GroupId:     &id,
			DisplayName: &id,
		}
		grpNative := aws.Group{
			ID:          strconv.Itoa(i),
			Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
			DisplayName: strconv.Itoa(i),
			Members:     []string{},
		}

		expectedOutput = append(expectedOutput, &grpNative)

		if i < 100 {
			sampleGroupsResponseA = append(sampleGroupsResponseA, grp)
		} else {
			sampleGroupsResponseB = append(sampleGroupsResponseB, grp)
		}
	}

	nextToken := "sample NextToken"
	sampleResponsePaginationA := &identitystore.ListGroupsOutput{
		Groups:    sampleGroupsResponseA,
		NextToken: &nextToken,
	}

	sampleResponsePaginationB := &identitystore.ListGroupsOutput{
		Groups: sampleGroupsResponseB,
	}

	gomock.InOrder(
		mockIdentityStoreClient.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(sampleResponsePaginationA, nil),
		mockIdentityStoreClient.EXPECT().ListGroups(gomock.Any(), gomock.Any()).Return(sampleResponsePaginationB, nil),
	)

	actualOutput, err := mockClient.GetGroups()

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
	assert.NoError(t, err)
}

func Test_GetGroupsEmptyResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	sampleResponseNoGroups := &identitystore.ListGroupsOutput{Groups: []types.Group{}}

	expectedOutput := []*aws.Group{}

	mockIdentityStoreClient.EXPECT().ListGroups(gomock.Any(), gomock.Any()).MaxTimes(1).Return(sampleResponseNoGroups, nil)

	actualOutput, err := mockClient.GetGroups()

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
	assert.NoError(t, err)
}

func Test_GetGroupsErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	sampleResponseError := errors.New("Sample error")

	expectedOutput := errors.New("Sample error")

	mockIdentityStoreClient.EXPECT().ListGroups(gomock.Any(), gomock.Any()).MaxTimes(1).Return(nil, sampleResponseError)

	actualOutput, err := mockClient.GetGroups()

	assert.True(t, reflect.DeepEqual(expectedOutput.Error(), err.Error()))
	assert.Nil(t, actualOutput)
}

func Test_GetUsersWithoutPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	userId1 := "user-1-test-id"
	userName1 := "user-1@example.com"
	title1 := "Example title"
	familyName1 := "1"
	givenName1 := "User"
	displayName1 := "User 1"
	addrType1 := "Home"
	country1 := "Canada"
	emailType1 := "work"
	emailValue1 := "user-1@example.com"

	userId2 := "user-2-test-id"
	userName2 := "user-2@example.com"
	familyName2 := "2"
	givenName2 := "User"
	displayName2 := "User 2"
	addrType2a := "Work"
	addrType2b := "Home"
	emailType2a := "work"
	emailValue2a := "user-2@example.com"
	emailType2b := "personal"
	emailValue2b := "user-2-personal@example.com"

	// >1 user response with no pagination (<100 users returned)
	sampleResponseNoPagination := &identitystore.ListUsersOutput{
		Users: []types.User{
			{
				UserId:      &userId1,
				UserName:    &userName1,
				Title:       &title1,
				Name:        &types.Name{FamilyName: &familyName1, GivenName: &givenName1},
				DisplayName: &displayName1,
				Addresses:   []types.Address{{Type: &addrType1, Country: &country1}},
				Emails:      []types.Email{{Primary: true, Type: &emailType1, Value: &emailValue1}},
			},
			{
				UserId:      &userId2,
				UserName:    &userName2,
				Name:        &types.Name{FamilyName: &familyName2, GivenName: &givenName2},
				DisplayName: &displayName2,
				Addresses:   []types.Address{{Type: &addrType2a}, {Type: &addrType2b}},
				Emails: []types.Email{
					{Primary: true, Type: &emailType2a, Value: &emailValue2a},
					{Primary: false, Type: &emailType2b, Value: &emailValue2b},
				},
			},
		},
	}

	expectedOutput := []*aws.User{
		{
			ID:       "user-1-test-id",
			Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			Username: "user-1@example.com",
			Name: struct {
				FamilyName string `json:"familyName"`
				GivenName  string `json:"givenName"`
			}{
				FamilyName: "1",
				GivenName:  "User",
			},
			DisplayName: "User 1",
			Emails: []aws.UserEmail{
				{Primary: true, Type: "work", Value: "user-1@example.com"},
			},
			Addresses: []aws.UserAddress{{Type: "Home"}},
		},
		{
			ID:       "user-2-test-id",
			Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			Username: "user-2@example.com",
			Name: struct {
				FamilyName string `json:"familyName"`
				GivenName  string `json:"givenName"`
			}{
				FamilyName: "2",
				GivenName:  "User",
			},
			DisplayName: "User 2",
			Emails: []aws.UserEmail{
				{Primary: true, Type: "work", Value: "user-2@example.com"},
				{Primary: false, Type: "personal", Value: "user-2-personal@example.com"},
			},
			Addresses: []aws.UserAddress{{Type: "Work"}, {Type: "Home"}},
		},
	}

	mockIdentityStoreClient.EXPECT().ListUsers(gomock.Any(), gomock.Any()).MaxTimes(1).Return(sampleResponseNoPagination, nil)

	actualOutput, err := mockClient.GetUsers()

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
	assert.NoError(t, err)
}

func Test_GetUsersWithPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	// >1 user response with pagination (150 users returned)
	sampleUsersResponseA := make([]types.User, 0, 100)
	sampleUsersResponseB := make([]types.User, 0, 50)
	expectedOutput := make([]*aws.User, 0, 150)

	// Populate responses
	for i := 0; i < 150; i++ {
		id := strconv.Itoa(i)
		displayNameStr := "User " + strconv.Itoa(i)
		emailTypeStr := "work"
		emailValueStr := strconv.Itoa(i) + "@example.com"
		addrTypeStr := "Home"
		givenNameStr := "User"

		usr := types.User{
			UserId:      &id,
			UserName:    &id,
			DisplayName: &displayNameStr,
			Name:        &types.Name{FamilyName: &id, GivenName: &givenNameStr},
			Emails:      []types.Email{{Primary: true, Type: &emailTypeStr, Value: &emailValueStr}},
			Addresses:   []types.Address{{Type: &addrTypeStr}},
		}
		usrNative := aws.User{
			ID:       strconv.Itoa(i),
			Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			Username: strconv.Itoa(i),
			Name: struct {
				FamilyName string `json:"familyName"`
				GivenName  string `json:"givenName"`
			}{
				FamilyName: strconv.Itoa(i),
				GivenName:  "User",
			},
			DisplayName: "User " + strconv.Itoa(i),
			Emails:      []aws.UserEmail{{Primary: true, Type: "work", Value: strconv.Itoa(i) + "@example.com"}},
			Addresses:   []aws.UserAddress{{Type: "Home"}},
		}

		expectedOutput = append(expectedOutput, &usrNative)

		if i < 100 {
			sampleUsersResponseA = append(sampleUsersResponseA, usr)
		} else {
			sampleUsersResponseB = append(sampleUsersResponseB, usr)
		}
	}

	nextToken := "sample NextToken"
	sampleResponsePaginationA := &identitystore.ListUsersOutput{
		Users:     sampleUsersResponseA,
		NextToken: &nextToken,
	}

	sampleResponsePaginationB := &identitystore.ListUsersOutput{
		Users: sampleUsersResponseB,
	}

	gomock.InOrder(
		mockIdentityStoreClient.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(sampleResponsePaginationA, nil),
		mockIdentityStoreClient.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(sampleResponsePaginationB, nil),
	)

	actualOutput, err := mockClient.GetUsers()

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
	assert.NoError(t, err)
}

func Test_GetUsersEmptyResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	// >1 user response with no pagination (<100 users returned)
	sampleResponseNoPagination := &identitystore.ListUsersOutput{}

	expectedOutput := []*aws.User{}

	mockIdentityStoreClient.EXPECT().ListUsers(gomock.Any(), gomock.Any()).MaxTimes(1).Return(sampleResponseNoPagination, nil)

	actualOutput, err := mockClient.GetUsers()

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
	assert.NoError(t, err)
}

func Test_GetUsersErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	sampleResponseError := errors.New("Sample error")

	expectedOutput := errors.New("Sample error")

	mockIdentityStoreClient.EXPECT().ListUsers(gomock.Any(), gomock.Any()).MaxTimes(1).Return(nil, sampleResponseError)

	actualOutput, err := mockClient.GetUsers()

	assert.True(t, reflect.DeepEqual(expectedOutput.Error(), err.Error()))
	assert.Nil(t, actualOutput)
}

func Test_ConvertSdkUserObjToNative(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userId := "user-1-test-id"
	userName := "user-1@example.com"
	title := "Example title"
	familyName := "1"
	givenName := "User"
	displayName := "User 1"
	addrType := "Home"
	country := "Canada"
	emailType := "work"
	emailValue := "user-1@example.com"

	// >1 user response with no pagination (<100 users returned)
	sampleInput := &types.User{
		UserId:      &userId,
		UserName:    &userName,
		Title:       &title,
		Name:        &types.Name{FamilyName: &familyName, GivenName: &givenName},
		DisplayName: &displayName,
		Addresses:   []types.Address{{Type: &addrType, Country: &country}},
		Emails:      []types.Email{{Primary: true, Type: &emailType, Value: &emailValue}},
	}

	expectedOutput := &aws.User{
		ID:       "user-1-test-id",
		Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		Username: "user-1@example.com",
		Name: struct {
			FamilyName string `json:"familyName"`
			GivenName  string `json:"givenName"`
		}{
			FamilyName: "1",
			GivenName:  "User",
		},
		DisplayName: "User 1",
		Emails: []aws.UserEmail{
			{Primary: true, Type: "work", Value: "user-1@example.com"},
		},
		Addresses: []aws.UserAddress{{Type: "Home"}},
	}

	actualOutput := ConvertSdkUserObjToNative(sampleInput)

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
}

func Test_CreateUserIDtoUserObjMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sampleInput := []*aws.User{
		{
			ID:       "user-1-test-id",
			Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			Username: "user-1@example.com",
			Name: struct {
				FamilyName string `json:"familyName"`
				GivenName  string `json:"givenName"`
			}{
				FamilyName: "1",
				GivenName:  "User",
			},
			DisplayName: "User 1",
			Emails: []aws.UserEmail{
				{Primary: true, Type: "work", Value: "user-1@example.com"},
			},
			Addresses: []aws.UserAddress{{Type: "Home"}},
		},
		{
			ID:       "user-2-test-id",
			Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
			Username: "user-2@example.com",
			Name: struct {
				FamilyName string `json:"familyName"`
				GivenName  string `json:"givenName"`
			}{
				FamilyName: "2",
				GivenName:  "User",
			},
			DisplayName: "User 2",
			Emails: []aws.UserEmail{
				{Primary: true, Type: "work", Value: "user-2@example.com"},
				{Primary: false, Type: "personal", Value: "user-2-personal@example.com"},
			},
			Addresses: []aws.UserAddress{{Type: "Work"}, {Type: "Home"}},
		},
	}

	expectedOutput := make(map[string]*aws.User)

	expectedOutput["user-1-test-id"] = &aws.User{
		ID:       "user-1-test-id",
		Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		Username: "user-1@example.com",
		Name: struct {
			FamilyName string `json:"familyName"`
			GivenName  string `json:"givenName"`
		}{
			FamilyName: "1",
			GivenName:  "User",
		},
		DisplayName: "User 1",
		Emails: []aws.UserEmail{
			{Primary: true, Type: "work", Value: "user-1@example.com"},
		},
		Addresses: []aws.UserAddress{{Type: "Home"}},
	}

	expectedOutput["user-2-test-id"] = &aws.User{
		ID:       "user-2-test-id",
		Schemas:  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		Username: "user-2@example.com",
		Name: struct {
			FamilyName string `json:"familyName"`
			GivenName  string `json:"givenName"`
		}{
			FamilyName: "2",
			GivenName:  "User",
		},
		DisplayName: "User 2",
		Emails: []aws.UserEmail{
			{Primary: true, Type: "work", Value: "user-2@example.com"},
			{Primary: false, Type: "personal", Value: "user-2-personal@example.com"},
		},
		Addresses: []aws.UserAddress{{Type: "Work"}, {Type: "Home"}},
	}

	actualOutput := CreateUserIDtoUserObjMap(sampleInput)

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
}

func Test_GetGroupMembershipsLists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	sampleGroupsInput := []*aws.Group{
		{ID: "a", DisplayName: "a"},
		{ID: "b", DisplayName: "b"},
		{ID: "c", DisplayName: "c"},
	}

	sampleUsersMapInput := make(map[string]*aws.User)
	sampleUsersMapInput["1"] = &aws.User{ID: "1"}
	sampleUsersMapInput["2"] = &aws.User{ID: "2"}
	sampleUsersMapInput["3"] = &aws.User{ID: "3"}
	sampleUsersMapInput["4"] = &aws.User{ID: "4"}

	expectedOutput := make(map[string][]*aws.User)
	expectedOutput["a"] = []*aws.User{{ID: "1"}, {ID: "2"}}
	expectedOutput["b"] = []*aws.User{{ID: "2"}, {ID: "3"}, {ID: "4"}}
	expectedOutput["c"] = []*aws.User{}

	groupIdA := "a"
	groupIdB := "b"

	sampleResponseGroupA := &identitystore.ListGroupMembershipsOutput{
		GroupMemberships: []types.GroupMembership{
			{GroupId: &groupIdA, MemberId: &types.MemberIdMemberUserId{Value: "1"}},
			{GroupId: &groupIdA, MemberId: &types.MemberIdMemberUserId{Value: "2"}},
		},
	}

	sampleResponseGroupB := &identitystore.ListGroupMembershipsOutput{
		GroupMemberships: []types.GroupMembership{
			{GroupId: &groupIdB, MemberId: &types.MemberIdMemberUserId{Value: "2"}},
			{GroupId: &groupIdB, MemberId: &types.MemberIdMemberUserId{Value: "3"}},
			{GroupId: &groupIdB, MemberId: &types.MemberIdMemberUserId{Value: "4"}},
		},
	}

	sampleResponseGroupC := &identitystore.ListGroupMembershipsOutput{
		GroupMemberships: []types.GroupMembership{},
	}

	gomock.InOrder(
		mockIdentityStoreClient.EXPECT().ListGroupMemberships(gomock.Any(), gomock.Any()).Return(sampleResponseGroupA, nil),
		mockIdentityStoreClient.EXPECT().ListGroupMemberships(gomock.Any(), gomock.Any()).Return(sampleResponseGroupB, nil),
		mockIdentityStoreClient.EXPECT().ListGroupMemberships(gomock.Any(), gomock.Any()).Return(sampleResponseGroupC, nil),
	)

	actualOutput, err := mockClient.GetGroupMembershipsLists(sampleGroupsInput, sampleUsersMapInput)

	assert.True(t, reflect.DeepEqual(expectedOutput, actualOutput))
	assert.Nil(t, err)
}

func Test_IsUserInGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	sampleUserInput := &aws.User{ID: "test-user-id"}
	sampleGroupInput := &aws.Group{ID: "test-group-id"}

	groupId := "test-group-id"

	// True response
	sampleResponse := &identitystore.IsMemberInGroupsOutput{
		Results: []types.GroupMembershipExistenceResult{
			{
				GroupId:          &groupId,
				MemberId:         &types.MemberIdMemberUserId{Value: "test-user-id"},
				MembershipExists: true,
			},
		},
	}

	mockIdentityStoreClient.EXPECT().IsMemberInGroups(gomock.Any(), gomock.Any()).MaxTimes(1).Return(sampleResponse, nil)

	actualOutput, err := mockClient.IsUserInGroup(sampleUserInput, sampleGroupInput)

	assert.True(t, *actualOutput)
	assert.Nil(t, err)

	// False response
	sampleResponse.Results[0].MembershipExists = false

	mockIdentityStoreClient.EXPECT().IsMemberInGroups(gomock.Any(), gomock.Any()).MaxTimes(1).Return(sampleResponse, nil)

	actualOutput, err = mockClient.IsUserInGroup(sampleUserInput, sampleGroupInput)

	assert.False(t, *actualOutput)
	assert.Nil(t, err)

	// Error response
	sampleResponseErr := errors.New("Sample error")
	expectedResponse := errors.New("Sample error")

	mockIdentityStoreClient.EXPECT().IsMemberInGroups(gomock.Any(), gomock.Any()).MaxTimes(1).Return(nil, sampleResponseErr)

	actualOutput, _ = mockClient.IsUserInGroup(sampleUserInput, sampleGroupInput)

	assert.True(t, reflect.DeepEqual(sampleResponseErr.Error(), expectedResponse.Error()))
	assert.Nil(t, actualOutput)
}

func Test_RemoveUserFromGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityStoreClient := mocks.NewMockIdentityStoreAPI(ctrl)

	mockClient := &syncGSuite{
		aws:                 nil,
		google:              nil,
		cfg:                 &config.Config{IdentityStoreID: "test-identity-store-id"},
		identityStoreClient: mockIdentityStoreClient,
		users:               make(map[string]*aws.User),
	}

	sampleUserInput := "test-user-id"
	sampleGroupInput := "test-group-id"

	membershipId := "test-membership-id"
	sampleResponse := &identitystore.GetGroupMembershipIdOutput{
		MembershipId: &membershipId,
	}

	mockIdentityStoreClient.EXPECT().GetGroupMembershipId(gomock.Any(), gomock.Any()).MaxTimes(1).Return(sampleResponse, nil)
	mockIdentityStoreClient.EXPECT().DeleteGroupMembership(gomock.Any(), gomock.Any()).MaxTimes(1).Return(&identitystore.DeleteGroupMembershipOutput{}, nil)

	err := mockClient.RemoveUserFromGroup(&sampleUserInput, &sampleGroupInput)

	assert.Nil(t, err)
}
