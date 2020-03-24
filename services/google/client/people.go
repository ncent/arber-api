package clients

import (
	"context"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	people "google.golang.org/api/people/v1"
)

const maxPageSize = 100

func (gs *GoogleService) GetContacts(config *GoogleConfig, token *oauth2.Token, ctx context.Context) ([]*people.Person, error) {
	connections, err := gs.getConnections(config, token, ctx)
	if err != nil {
		return nil, err
	}

	if len(connections) == 0 {
		return nil, nil
	}

	return connections, err
}

func (gs *GoogleService) GetMyPerson(config *GoogleConfig, token *oauth2.Token, ctx context.Context) (*people.Person, error) {
	return gs.GetPerson(config, token, "people/me", ctx)
}

func (gs *GoogleService) GetPerson(config *GoogleConfig, token *oauth2.Token, personID string, ctx context.Context) (*people.Person, error) {
	service, err := gs.getPeopleService(config, token, ctx)
	if err != nil {
		return nil, err
	}

	r, err := service.People.Get(personID).PersonFields("names,emailAddresses,phoneNumbers").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve person (people/me). %v", err)
		return nil, err
	}

	return r, nil
}

func (gs *GoogleService) getConnections(config *GoogleConfig, token *oauth2.Token, ctx context.Context) ([]*people.Person, error) {
	service, err := gs.getPeopleService(config, token, ctx)
	if err != nil {
		return nil, err
	}

	var connections []*people.Person

	hasMoreConnections := true
	var nextPageToken *string
	for hasMoreConnections {
		query := service.People.Connections.List("people/me").PageSize(maxPageSize).PersonFields("names,emailAddresses,phoneNumbers")
		if nextPageToken != nil {
			query.PageToken(*nextPageToken)
		}
		r, err := query.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve people. %v", err)
			return nil, err
		}

		log.Printf("Connections Response: %+v", r)
		log.Printf("Found %d connections", len(r.Connections))
		nextPageToken = &r.NextPageToken
		if nextPageToken != nil && len(r.Connections) >= maxPageSize {
			log.Printf("Next page found: %v", nextPageToken)
		} else {
			hasMoreConnections = false
		}
		connections = append(connections, r.Connections...)
	}

	log.Printf("Found total contacts: %v", len(connections))
	return connections, nil
}

func (gs *GoogleService) getPeopleService(config *GoogleConfig, token *oauth2.Token, ctx context.Context) (*people.Service, error) {
	client := gs.getPeopleClient(config, token, ctx)
	srv, err := people.New(client)

	if err != nil {
		log.Fatalf("Unable to create service %v", err)
		return nil, err
	}
	return srv, err
}

func (gs *GoogleService) getPeopleClient(config *GoogleConfig, token *oauth2.Token, ctx context.Context) *http.Client {
	var context context.Context
	if ctx != nil {
		context = ctx
	} else {
		context = gs.ctx
	}

	tkr := &TokenRefresher{
		ctx:          context,
		conf:         config,
		refreshToken: token.RefreshToken,
	}

	rtks := &ReuseTokenSource{
		t:   token,
		new: tkr,
	}
	return oauth2.NewClient(context, rtks)
}
