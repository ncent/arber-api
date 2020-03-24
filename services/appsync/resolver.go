package appsync

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/mitchellh/mapstructure"
	appsync "github.com/rodrigopavezi/appsync-client-go"
	"github.com/rodrigopavezi/appsync-client-go/graphql"
)

var serverURL, _ = os.LookupEnv("AWS_APP_SYNC_URL")

type Resolver struct {
	awsConfig *aws.Config
}

func New() Resolver {
	// get aws credential
	config := aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
	sess := session.Must(session.NewSession(&config))

	r := Resolver{sess.Config}
	return r
}

func (r Resolver) CreateUser(input CreateUserInput) (*User, error) {
	mutation := `mutation CreateUser($input: CreateUserInput!) {
		createUser(input: $input) {
			id
			names
			emails
			phoneNumbers
			pictures
			identity
			token
			etag
			sharedActions {
				items {
					id
				}
				nextToken
			}
			contacts {
				items {
					id
				}
				nextToken
			}
			usersImContactOf {
				items {
					id
				}
				nextToken
			}
			createdAt
			updatedAt
			deletedAt
		}
	}
	`
	inputNewUser := &CreateInput{
		Input: input,
	}
	jsonInputNewUser, err := json.Marshal(inputNewUser)
	variables := json.RawMessage(jsonInputNewUser)
	log.Printf("jsonInputNewUser: %+v", string(jsonInputNewUser))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	appsyncResponse, err := client.Post(graphql.PostRequest{
		Query:     mutation,
		Variables: &variables,
	})
	log.Printf("CreateUser Appsync response status code: %+v", appsyncResponse.StatusCode)
	log.Printf("CreateUser Appsync response Errors: %+v", appsyncResponse.Errors)

	if err != nil {
		log.Printf("Failed to post to appsync: %v", err)
		return nil, err
	}

	var result CreateUserResponse
	err = mapstructure.Decode(appsyncResponse.Data, &result)

	log.Println("CreateUser data: %+v", result.CreateUser)
	return &result.CreateUser, nil
}

func (r Resolver) UpdateUser(input UpdateUserInput) (*User, error) {
	mutation := `mutation UpdateUser($input: UpdateUserInput!) {
		updateUser(input: $input) {
			id
			names
			emails
			phoneNumbers
			pictures
			identity
			token
			etag
			sharedActions {
				items {
					id
				}
				nextToken
			}
			contacts {
				items {
					id
				}
				nextToken
			}
			usersImContactOf {
				items {
					id
				}
				nextToken
			}
			createdAt
			updatedAt
			deletedAt
		}
	}
	`
	inputNewUser := &UpdateInput{
		Input: input,
	}
	jsonInputNewUser, err := json.Marshal(inputNewUser)
	variables := json.RawMessage(jsonInputNewUser)
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	appsyncResponse, err := client.Post(graphql.PostRequest{
		Query:     mutation,
		Variables: &variables,
	})

	log.Printf("UpdateUser Appsync response: %v", appsyncResponse)

	if err != nil {
		log.Printf("Failed to post to appsync: %v", err)
		return nil, err
	}

	var result UpdateUserResponse
	err = mapstructure.Decode(appsyncResponse.Data, &result)

	log.Println("UpdateUser data: %+v", result.UpdateUser)
	return &result.UpdateUser, nil
}

func (r Resolver) GetUser(id string) (*User, error) {
	query := `query GetUser($id: ID!) {
		getUser(id: $id) {
			id
			names
			emails
			phoneNumbers
			pictures
			identity
			token
			etag
			sharedActions {
				items {
					id
				}
				nextToken
			}
			contacts {
				items {
					id
				}
				nextToken
			}
			usersImContactOf {
				items {
					id
				}
				nextToken
			}
			createdAt
			updatedAt
			deletedAt
		}
	}
	`
	log.Printf("ID %v", id)
	log.Printf("serverURL %v", serverURL)
	variables := json.RawMessage(fmt.Sprintf(`{ "id": "%s"	}`, id))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	response, err := client.Post(graphql.PostRequest{
		Query:     query,
		Variables: &variables,
	})
	if err != nil {
		log.Printf("Failed to post to appsync: %+v", err)
		return nil, err
	}

	log.Printf("Graph QL Response: %v", response)
	log.Printf("Graph QL Response status code: %v", response.StatusCode)
	log.Printf("Graph QL Response data: %v", response.Data)
	log.Printf("Graph QL Response errors: %v", response.Errors)

	var result GetUserResponse
	err = mapstructure.Decode(response.Data, &result)

	log.Println("GetUser data: %+v", result.GetUser)
	return &result.GetUser, nil
}

func (r Resolver) MapUsersByEmails(emails []*string) (map[string]User, error) {
	users, err := r.ListUsersByEmails(emails)
	if err != nil {
		return nil, err
	}

	emailsToUserMap := make(map[string]User)
	for _, user := range users {
		for _, email := range user.Emails {
			emailsToUserMap[*email] = user
		}
	}
	return emailsToUserMap, nil
}

func (r Resolver) ListUsersByEmails(emails []*string) ([]User, error) {
	query := `query ListUsers(
		$filter: ModelUserFilterInput
		$limit: Int
		$nextToken: String
	) {
		listUsers(filter: $filter, limit: $limit, nextToken: $nextToken) {
			items {
				id
				names
				emails
				phoneNumbers
				pictures
				identity
				token
				etag
				sharedActions {
					nextToken
				}
				contacts {
					nextToken
				}
				usersImContactOf {
					nextToken
				}
				createdAt
				updatedAt
				deletedAt
			}
			nextToken
		}
	}
	`
	var strs []string

	for _, email := range emails {
		strs = append(strs, fmt.Sprintf(`{"emails": { "contains": "%s" } }`, *email))
	}
	filterJson := fmt.Sprintf(`{"filter": %s, "limit": 1000 }`, strings.Join(strs, `, "or": `))
	log.Printf("filterJson: %v", filterJson)

	variables := json.RawMessage(filterJson)
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	response, err := client.Post(graphql.PostRequest{
		Query:     query,
		Variables: &variables,
	})
	if err != nil {
		log.Printf("Failed to post to appsync: %+v", err)
		return nil, err
	}

	log.Printf("Graph QL Response: %v", response)
	log.Printf("Graph QL Response status code: %v", response.StatusCode)
	log.Printf("Graph QL Response data: %v", response.Data)
	log.Printf("Graph QL Response errors: %v", response.Errors)

	var result ListUsersByEmailsResponse
	err = mapstructure.Decode(response.Data, &result)

	log.Println("ListUsersByEmails data: %+v", result.ListUsers.Items)
	return result.ListUsers.Items, nil
}

func (r Resolver) CreateUserContact(input CreateUserContactInput) (*UserContact, error) {
	query := `mutation CreateUserContact($input: CreateUserContactInput!) {
		createUserContact(input: $input) {
			id
			user {
				id
				names
				emails
				phoneNumbers
				pictures
				identity
				token
				etag
				sharedActions {
					nextToken
				}
				contacts {
					nextToken
				}
				usersImContactOf {
					nextToken
				}
				createdAt
				updatedAt
				deletedAt
			}
			contact {
				id
				names
				emails
				phoneNumbers
				pictures
				identity
				token
				etag
				sharedActions {
					nextToken
				}
				contacts {
					nextToken
				}
				usersImContactOf {
					nextToken
				}
				createdAt
				updatedAt
				deletedAt
			}
		}
	}
	`
	inputNewUserContact := &CreateUserContactInputWrapper{
		Input: input,
	}
	jsonInputNewUserContact, err := json.Marshal(inputNewUserContact)
	variables := json.RawMessage(jsonInputNewUserContact)
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	response, err := client.Post(graphql.PostRequest{
		Query:     query,
		Variables: &variables,
	})
	if err != nil {
		log.Printf("Failed to post to appsync: %+v", err)
		return nil, err
	}

	log.Printf("Graph QL Response: %v", response)
	log.Printf("Graph QL Response status code: %v", response.StatusCode)
	log.Printf("Graph QL Response data: %v", response.Data)
	log.Printf("Graph QL Response errors: %v", response.Errors)

	var result CreateUserContactResponse
	err = mapstructure.Decode(response.Data, &result)

	log.Println("CreateUserContact data: %+v", result.CreateUserContact)
	return &result.CreateUserContact, nil
}

func (r Resolver) CreateChallenge(input CreateChallenge) (*Challenge, error) {
	mutation := `mutation CreateChallenge($input: CreateChallengeInput!) {
		createChallenge(input: $input) {
			id
			name
			sponsorName
		}
	}
	`
	inputNewChallenge := &CreateChallengeInput{
		Input: input,
	}
	jsonInputNewChallenge, err := json.Marshal(inputNewChallenge)
	variables := json.RawMessage(jsonInputNewChallenge)
	log.Printf("jsonInputNewChallenge: %+v", string(jsonInputNewChallenge))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	appsyncResponse, err := client.Post(graphql.PostRequest{
		Query:     mutation,
		Variables: &variables,
	})
	log.Printf("CreateChallenge Appsync response status code: %+v", appsyncResponse.StatusCode)
	log.Printf("CreateChallenge Appsync response Errors: %+v", appsyncResponse.Errors)

	if err != nil {
		log.Printf("Failed to post to appsync: %v", err)
		return nil, err
	}

	var result CreateChallengeResponse
	err = mapstructure.Decode(appsyncResponse.Data, &result)

	log.Println("CreateChallenge data: %+v", result.CreateChallenge)
	return &result.CreateChallenge, nil
}

func (r Resolver) GetChallenge(id string) (*Challenge, error) {
	query := `query GetChallenge($id: ID!) {
		getChallenge(id: $id) {
			id
			name
			description
			sponsorName
		}
	}
	`
	log.Printf("ID %v", id)
	log.Printf("serverURL %v", serverURL)
	variables := json.RawMessage(fmt.Sprintf(`{ "id": "%s"	}`, id))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	response, err := client.Post(graphql.PostRequest{
		Query:     query,
		Variables: &variables,
	})
	if err != nil {
		log.Printf("Failed to post to appsync: %+v", err)
		return nil, err
	}

	log.Printf("Graph QL Response: %v", response)
	log.Printf("Graph QL Response status code: %v", response.StatusCode)
	log.Printf("Graph QL Response data: %v", response.Data)
	log.Printf("Graph QL Response errors: %v", response.Errors)

	var result GetChallengeResponse
	err = mapstructure.Decode(response.Data, &result)

	log.Println("GetChallenge data: %+v", result.GetChallenge)
	return &result.GetChallenge, nil
}

func (r Resolver) CreateShareAction(input CreateShareAction) (*ShareAction, error) {
	mutation := `mutation CreateShareAction($input: CreateShareActionInput!) {
		createShareAction(input: $input) {
			id
		}
	}
	`
	inputNewShareAction := &CreateShareActionInput{
		Input: input,
	}
	jsonInputNewShareAction, err := json.Marshal(inputNewShareAction)
	variables := json.RawMessage(jsonInputNewShareAction)
	log.Printf("jsonInputNewShareAction: %+v", string(jsonInputNewShareAction))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	appsyncResponse, err := client.Post(graphql.PostRequest{
		Query:     mutation,
		Variables: &variables,
	})
	log.Printf("CreateShareAction Appsync response status code: %+v", appsyncResponse.StatusCode)
	log.Printf("CreateShareAction Appsync response Errors: %+v", appsyncResponse.Errors)

	if err != nil {
		log.Printf("Failed to post to appsync: %v", err)
		return nil, err
	}

	var result CreateShareActionResponse
	err = mapstructure.Decode(appsyncResponse.Data, &result)

	log.Println("CreateShareAction data: %+v", result.CreateShareAction)
	return &result.CreateShareAction, nil
}

func (r Resolver) UpdateShareAction(input UpdateShareAction) (*ShareAction, error) {
	mutation := `mutation UpdateShareAction($input: UpdateShareActionInput!) {
		updateShareAction(input: $input) {
			id
		}
	}
	`
	inputNewShareAction := &UpdateShareActionInput{
		Input: input,
	}
	jsonInputNewShareAction, err := json.Marshal(inputNewShareAction)
	variables := json.RawMessage(jsonInputNewShareAction)
	log.Printf("jsonInputNewShareAction: %+v", string(jsonInputNewShareAction))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	appsyncResponse, err := client.Post(graphql.PostRequest{
		Query:     mutation,
		Variables: &variables,
	})
	log.Printf("UpdateShareAction Appsync response status code: %+v", appsyncResponse.StatusCode)
	log.Printf("UpdateShareAction Appsync response Errors: %+v", appsyncResponse.Errors)

	if err != nil {
		log.Printf("Failed to post to appsync: %v", err)
		return nil, err
	}

	var result UpdateShareActionResponse
	err = mapstructure.Decode(appsyncResponse.Data, &result)

	log.Println("UpdateShareAction data: %+v", result.UpdateShareAction)
	return &result.UpdateShareAction, nil
}

func (r Resolver) CreateShareActionContact(input CreateShareActionContact) (*ShareActionContact, error) {
	mutation := `mutation CreateShareActionContact($input: CreateShareActionContactInput!) {
		createShareActionContact(input: $input) {
			id
		}
	}
	`
	inputNewShareActionContact := &CreateShareActionContactInput{
		Input: input,
	}
	jsonInputNewShareActionContact, err := json.Marshal(inputNewShareActionContact)
	variables := json.RawMessage(jsonInputNewShareActionContact)
	log.Printf("jsonInputNewShareActionContact: %+v", string(jsonInputNewShareActionContact))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	appsyncResponse, err := client.Post(graphql.PostRequest{
		Query:     mutation,
		Variables: &variables,
	})
	log.Printf("CreateShareActionContact Appsync response status code: %+v", appsyncResponse.StatusCode)
	log.Printf("CreateShareActionContact Appsync response Errors: %+v", appsyncResponse.Errors)

	if err != nil {
		log.Printf("Failed to post to appsync: %v", err)
		return nil, err
	}

	var result CreateShareActionContactResponse
	err = mapstructure.Decode(appsyncResponse.Data, &result)

	log.Println("CreateShareActionContact data: %+v", result.CreateShareActionContact)
	return &result.CreateShareActionContact, nil
}

func (r Resolver) CreateTransaction(input CreateTransaction) (*Transaction, error) {
	mutation := `mutation CreateTransaction($input: CreateTransactionInput!) {
		createTransaction(input: $input) {
			id
		}
	}
	`
	inputNewTransaction := &CreateTransactionInput{
		Input: input,
	}
	jsonInputNewTransaction, err := json.Marshal(inputNewTransaction)
	variables := json.RawMessage(jsonInputNewTransaction)
	log.Printf("jsonInputNewTransaction: %+v", string(jsonInputNewTransaction))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	appsyncResponse, err := client.Post(graphql.PostRequest{
		Query:     mutation,
		Variables: &variables,
	})
	log.Printf("CreateTransaction Appsync response status code: %v", *appsyncResponse.StatusCode)
	log.Printf("CreateTransaction Appsync response Errors: %+v", appsyncResponse.Errors)

	if err != nil {
		log.Printf("Failed to post to appsync: %v", err)
		return nil, err
	}

	var result CreateTransactionResponse
	err = mapstructure.Decode(appsyncResponse.Data, &result)

	log.Println("CreateTransaction data: %+v", result.CreateTransaction)
	return &result.CreateTransaction, nil
}

func (r Resolver) GetTransaction(id string) (*Transaction, error) {
	query := `query GetTransaction($id: ID!) {
		getTransaction(id: $id) {
			id
			action {
				id
				challengeId
			}
			parentTransaction {
        id
        parentTransactionId
			}
			parentTransactionId
		}
	}
	`
	log.Printf("ID %v", id)
	variables := json.RawMessage(fmt.Sprintf(`{ "id": "%s"	}`, id))
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	response, err := client.Post(graphql.PostRequest{
		Query:     query,
		Variables: &variables,
	})
	if err != nil {
		log.Printf("Failed to post to appsync: %+v", err)
		return nil, err
	}

	log.Printf("Graph QL Response: %v", response)
	log.Printf("Graph QL Response status code: %v", response.StatusCode)
	log.Printf("Graph QL Response data: %v", response.Data)
	log.Printf("Graph QL Response errors: %v", response.Errors)

	var result GetTransactionResponse
	err = mapstructure.Decode(response.Data, &result)

	log.Println("GetTransaction data: %+v", result.GetTransaction)
	return &result.GetTransaction, nil
}

func (r Resolver) GetShareActionsByChallengeAndUser(challengeID string, userID string) ([]*ShareAction, error) {
	query := `query ListShareActions(
		$filter: ModelShareActionFilterInput
		$limit: Int
		$nextToken: String
	) {
		listShareActions(filter: $filter, limit: $limit, nextToken: $nextToken) {
			items {
				id
			}
			nextToken
		}
	}
	`
	filterJson := fmt.Sprintf(`{"filter": { "challengeId": { "eq": "%s" } }, "and": { "userId" : { "eq": "%s" } }, "limit": 1000 }`, challengeID, userID)
	log.Printf("filterJson: %v", filterJson)

	variables := json.RawMessage(filterJson)
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	response, err := client.Post(graphql.PostRequest{
		Query:     query,
		Variables: &variables,
	})
	if err != nil {
		log.Printf("Failed to post to appsync: %+v", err)
		return nil, err
	}

	log.Printf("Graph QL Response: %v", response)
	log.Printf("Graph QL Response status code: %v", response.StatusCode)
	log.Printf("Graph QL Response data: %v", response.Data)
	log.Printf("Graph QL Response errors: %v", response.Errors)

	var result ListShareActionsByChallengeAndUserResponse
	err = mapstructure.Decode(response.Data, &result)

	log.Println("GetShareActionsByChallengeAndUser data: %+v", result.ListShareActions.Items)

	return result.ListShareActions.Items, nil
}

func (r Resolver) GetTransactionsByShareAction(actionID string) ([]*Transaction, error) {
	query := `query ListTransactions(
		$filter: ModelTransactionFilterInput
		$limit: Int
		$nextToken: String
	) {
		listTransactions(filter: $filter, limit: $limit, nextToken: $nextToken) {
			items {
				id
			}
			nextToken
		}
	}
	`
	filterJson := fmt.Sprintf(`{"filter": { "transactionActionId": { "eq": "%s"} }, "limit": 1000 }`, actionID)
	log.Printf("filterJson: %v", filterJson)

	variables := json.RawMessage(filterJson)
	client := appsync.NewClient(appsync.NewGraphQLClient(graphql.NewClient(serverURL, *r.awsConfig)))
	response, err := client.Post(graphql.PostRequest{
		Query:     query,
		Variables: &variables,
	})
	if err != nil {
		log.Printf("Failed to post to appsync: %+v", err)
		return nil, err
	}

	log.Printf("Graph QL Response: %v", response)
	log.Printf("Graph QL Response status code: %v", response.StatusCode)
	log.Printf("Graph QL Response data: %v", response.Data)
	log.Printf("Graph QL Response errors: %v", response.Errors)

	var result ListTransactionByShareActionResponse
	err = mapstructure.Decode(response.Data, &result)

	log.Println("GetTransactionsByShareAction data: %+v", result.ListTransaction.Items)

	return result.ListTransaction.Items, nil
}
