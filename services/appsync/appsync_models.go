package appsync

import (
	"log"

	"github.com/looplab/fsm"
)

type CreateShareAction struct {
	ID          *string `json:"id,omitempty"`
	ChallengeID *string `json:"challengeId,omitempty"`
	UserID      *string `json:"userId,omitempty"`
}

type CreateShareActionInput struct {
	Input CreateShareAction `json:"input"`
}

type UpdateShareAction struct {
	ID          *string `json:"id,omitempty"`
	ChallengeID *string `json:"challengeId,omitempty"`
	UserID      *string `json:"userId,omitempty"`
}

type UpdateShareActionInput struct {
	Input UpdateShareAction `json:"input"`
}

type ShareActionContact struct {
	ID *string `json:"id,omitempty"`
}

type CreateShareActionContact struct {
	ID                              *string `json:"id,omitempty"`
	ShareActionContactShareActionID *string `json:"shareActionContactShareActionId,omitempty"`
	ShareActionContactContactID     *string `json:"shareActionContactContactId,omitempty"`
}

type CreateShareActionContactInput struct {
	Input CreateShareActionContact `json:"input"`
}

type Challenge struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	SponsorName *string `json:"sponsorName,omitempty"`
}

type CreateChallenge struct {
	ID                         *string `json:"id,omitempty"`
	Name                       *string `json:"name,omitempty"`
	Description                *string `json:"description,omitempty"`
	ImageURL                   *string `json:"imageUrl,omitempty"`
	SponsorName                *string `json:"sponsorName,omitempty"`
	Expiration                 *string `json:"expiration,omitempty"`
	ShareExpiration            *string `json:"shareExpiration,omitempty"`
	MaxShares                  *int    `json:"maxShare,omitempty"`
	MaxRewards                 *int    `json:"maxRewards,omitempty"`
	OffChain                   *bool   `json:"offChain,omitempty"`
	MaxDistributionFeeReward   *int    `json:"maxDistributionFeeReward,omitempty"`
	MaxSharesPerReceivedShare  *int    `json:"maxSharesPerReceivedShare,omitempty"`
	MaxDepth                   *int    `json:"maxDepth,omitempty"`
	MaxNodes                   *int    `json:"maxNodes,omitempty"`
	PublicKey                  *string `json:"publicKey,omitempty"`
	Reward                     *string `json:"reward,omitempty"`
	Active                     *bool   `json:"active,omitempty"`
	ChallengeTemplateID        *string `json:"challengeTemplateId,omitempty"`
	ChallengeParentChallengeID *string `json:"challengeParentChallengeId,omitempty"`
	AttachmentURL              *string `json:"attachmentURL,omitempty"`
}

type CreateChallengeInput struct {
	Input CreateChallenge `json:"input"`
}

type ShareAction struct {
	ID          *string `json:"id,omitempty"`
	ChallengeID *string `json:"challengeId,omitempty"`
}

type ShareActions struct {
	Items     []*ShareAction `json:"items,omitempty"`
	NextToken *string        `json:"nextToken,omitempty"`
}

type Transactions struct {
	Items     []*Transaction `json:"items,omitempty"`
	NextToken *string        `json:"nextToken,omitempty"`
}

type UserContact struct {
	ID      *string `json:"id,omitempty"`
	User    *User   `json:"user,omitempty"`
	Contact *User   `json:"contact,omitempty"`
}

type UserContacts struct {
	Items     []*UserContact `json:"items,omitempty"`
	NextToken *string        `json:"nextToken,omitempty"`
}

type User struct {
	Contacts         *UserContacts `json:"contacts,omitempty"`
	UsersImContactOf *UserContacts `json:"usersImContactOf,omitempty"`
	Emails           []*string     `json:"emails,omitempty"`
	Etag             *string       `json:"etag,omitempty"`
	ID               *string       `json:"id,omitempty"`
	Identity         *string       `json:"identity,omitempty"`
	Names            []*string     `json:"names,omitempty"`
	PhoneNumbers     []*string     `json:"phoneNumbers,omitempty"`
	Pictures         []*string     `json:"pictures,omitempty"`
	Token            *string       `json:"token,omitempty"`
	ShareActions     *ShareActions `json:"shareActions,omitempty"`
}

type CreateUserInput struct {
	Emails       []*string `json:"emails,omitempty"`
	Etag         *string   `json:"etag,omitempty"`
	ID           *string   `json:"id,omitempty"`
	Identity     *string   `json:"identity,omitempty"`
	Names        []*string `json:"names,omitempty"`
	PhoneNumbers []*string `json:"phoneNumbers,omitempty"`
	Pictures     []*string `json:"pictures,omitempty"`
	Token        *string   `json:"token,omitempty"`
}

type UpdateUserInput struct {
	Emails       []*string `json:"emails,omitempty"`
	Etag         *string   `json:"etag,omitempty"`
	ID           string    `json:"id,omitempty"`
	Identity     *string   `json:"identity,omitempty"`
	Names        []*string `json:"names,omitempty"`
	PhoneNumbers []*string `json:"phoneNumbers,omitempty"`
	Pictures     []*string `json:"pictures,omitempty"`
	Token        *string   `json:"token,omitempty"`
}

type CreateInput struct {
	Input CreateUserInput `json:"input"`
}

type UpdateInput struct {
	Input UpdateUserInput `json:"input"`
}

type CreateUserContactInput struct {
	UserContactUserId    *string `json:"userContactUserId"`
	UserContactContactId *string `json:"userContactContactId"`
}

type CreateUserContactInputWrapper struct {
	Input CreateUserContactInput `json:"input"`
}

type ListUsers struct {
	Items []User
}

type ListUsersByEmailsResponse struct {
	ListUsers ListUsers
}

type ListShareActionsByChallengeAndUserResponse struct {
	ListShareActions ShareActions
}

type ListTransactionByShareActionResponse struct {
	ListTransaction Transactions
}

type GetTransactionResponse struct {
	GetTransaction Transaction
}

type CreateTransactionResponse struct {
	CreateTransaction Transaction
}

type CreateShareActionContactResponse struct {
	CreateShareActionContact ShareActionContact
}

type CreateShareActionResponse struct {
	CreateShareAction ShareAction
}

type UpdateShareActionResponse struct {
	UpdateShareAction ShareAction
}

type CreateChallengeResponse struct {
	CreateChallenge Challenge
}

type CreateUserContactResponse struct {
	CreateUserContact UserContact
}

type GetUserResponse struct {
	GetUser User
}

type GetChallengeResponse struct {
	GetChallenge Challenge
}

type UpdateUserResponse struct {
	UpdateUser User
}

type CreateUserResponse struct {
	CreateUser User
}

type ActionStatus string

func (as ActionStatus) String() string {
	switch as {
	case ATTEMPED:
		return "ATTEMPED"
	case SCHEDULED:
		return "SCHEDULED"
	case CANCELLED:
		return "CANCELLED"
	case COMPLETED:
		return "COMPLETED"
	case CREATED:
		return "CREATED"
	default:
		panic("Unknown action status")
	}
}

const (
	ATTEMPED  ActionStatus = "ATTEMPED"
	SCHEDULED ActionStatus = "SCHEDULED"
	CANCELLED ActionStatus = "CANCELLED"
	COMPLETED ActionStatus = "COMPLETED"
	CREATED   ActionStatus = "CREATED"
)

type Action struct {
	AttemptCounter int          `json:"attemptCounter,omitempty"`
	Status         ActionStatus `json:status`
	FSM            *fsm.FSM
}

/*
func (a *Transaction) InitStateMachine() *Transaction {
	a.ShareAction.FSM = fsm.NewFSM(
		a.ShareAction.Status.String(),
		fsm.Events{
			{Name: CREATED.String(), Dst: SCHEDULED.String()},
			{Name: SCHEDULED.String(), Src: []string{CREATED.String()}, Dst: COMPLETED.String()},
			{Name: ATTEMPED.String(), Src: []string{CREATED.String(), SCHEDULED.String()}, Dst: COMPLETED.String()},
			{Name: CANCELLED.String(), Src: []string{CREATED.String(), SCHEDULED.String(), ATTEMPED.String()}, Dst: CANCELLED.String()},
			{Name: COMPLETED.String(), Src: []string{CREATED.String(), SCHEDULED.String(), ATTEMPED.String()}, Dst: COMPLETED.String()},
		},
		fsm.Callbacks{
			"status_changed": func(e *fsm.Event) { a.status_changed(e.Src, e.Dst) },
		},
	)
	return a
}*/

func (a *Transaction) status_changed(src string, dst string) error {
	log.Printf("State for action has changed: ")
	log.Printf("Transaction ID: %v", a.ID)
	log.Printf("From: %v to %v", src, dst)
	// TODO: Add switch statement for state changes
	// TODO: Add functionality for each state, starting with CREATE -> COMPLETED
	return nil
}

type Transaction struct {
	ID                *string      `json:"id,omitempty"`
	ParentTransaction *Transaction `json:"parentTransaction"`
	Action            *ShareAction `json:"action"`
}

type CreateTransaction struct {
	ID                  *string `json:"id,omitempty"`
	ParentTransactionID *string `json:"parentTransactionId,omitempty"`
	TransactionActionID *string `json:"transactionActionId,omitempty"`
}

type CreateTransactionInput struct {
	Input CreateTransaction `json:"input"`
}
