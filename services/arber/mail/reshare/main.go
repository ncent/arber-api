package mail

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	Resolver "gitlab.com/ncent/arber/api/services/appsync"
	ChallengeController "gitlab.com/ncent/arber/api/services/arber/challenge"
	ShareActionController "gitlab.com/ncent/arber/api/services/arber/share"
	helpers "gitlab.com/ncent/arber/api/services/google/helper"
)

func GenerateReshareBodyByChallenge(resolver Resolver.Resolver, transactionId string, challengeId string) (*string, error) {
	challenge, err := ChallengeController.GetChallenge(
		resolver, challengeId,
	)
	if err != nil {
		return nil, fmt.Errorf("There was a problem in GenerateReshareBodyByChallenge getting challenge: %v", err.Error())
	}
	log.Printf("Found challenge: %+v -- Generating mail body", *challenge)

	transaction, err := ShareActionController.CreateShareActionAndTransactionWithParentTransaction(resolver, transactionId, challengeId)
	if err != nil {
		return nil, fmt.Errorf("There was a problem in Creating new Share Action And Trasaction: %v", err.Error())
	}

	subject := fmt.Sprintf("Love this startup- Can you help us find an %s?", *challenge.Name)
	reshareLink, err := helpers.ShortenUrl(helpers.API_URL + "/reshare?transactionId=" + *transaction.ID + "&challengeId=" + *challenge.ID)
	if err != nil {
		log.Printf("Failed to generate short url for reshare link: %v", err)
		return nil, err
	}
	applyLink, err := helpers.ShortenUrl(helpers.CLIENT_APP_URL + "/apply/" + *transaction.ID)
	if err != nil {
		log.Printf("Failed to generate short url for apply link: %v", err)
		return nil, err
	}

	body := fmt.Sprintf(
		`I immediately thought of you. Please share with your network %s and your contribution will actually be measured and recognized. 
		
		Thanks! (to see more how this works or to apply check out: %s)`, *reshareLink, *applyLink)

	mailto := fmt.Sprintf(
		`mailto:?bcc=%s&subject=%s&body=%s`,
		url.QueryEscape(fmt.Sprintf(`share+%s@redb.ai`, *transaction.ID)),
		(&url.URL{Path: subject}).String(),
		//url.PathEscape(body),
		strings.Replace(url.PathEscape(body), "%0D%0A", "%0A", -1),
	)

	htmlBody := fmt.Sprintf(`
		<html>
			<head>
				<meta http-equiv="refresh" content="0; URL='%s'" />
				<style>
					body {
						margin: 0;
					}
					.background {
						width: 100%%;
						height: 100vh;
						display: flex;
						background-size: contain;
						background-color: #18191B;
						background-image: url(https://images.saatchiart.com/saatchi/841605/art/3529882/2599769-NLGHRQJF-6.jpg);
						background-repeat: no-repeat;
						background-position: center;
						flex-direction: column;
						justify-content: center;
						align-items: center;
					}
					.forwardImage {
						fill: #FFFFFF;
					}
					.forwardButton svg {
						width: 100px;
						height: 100px;
						margin: 0px;
					}
					.forwardButton {
						margin: auto;
						width: 200px;
						height: 200px;
						background-color: #b71c1b;
						background-repeat:no-repeat;
						cursor:pointer;
						overflow: hidden;
						outline:none;
						padding: 0px;
						border: none;
					}
					.forwardButton:hover {
						background-color: #9a1312;   
					}
					.text {
						font-size: 50px;
						color: #FFFFFF;
						margin-bottom: 100px;
					}
				</style>
			</head>
			<body >
				<div class="background">
					<div>
						<h1 class="text">Click to share with your network</h1>
					</div>
					<div>
						<button class="forwardButton" onClick="location.href='%s';">
							<svg version="1.1" id="Capa_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
								class="forwardImage" viewBox="0 0 422.853 422.853" style="enable-background:new 0 0 422.853 422.853;"
								xml:space="preserve">
								<g>
									<path d="M355.142,217.766l33.828,26.008V37.36c0-2.337-1.892-4.23-4.222-4.23H4.224C1.894,33.13,0,35.023,0,37.36v262.125
										c0,2.342,1.894,4.235,4.224,4.235h167.849v-29.462c0-1.488,0.196-2.933,0.437-4.36H33.821V83.417l158.163,115.99
										c1.503,1.086,3.516,1.086,5.005,0l158.167-115.99v134.349H355.142z M194.49,159.286L68.578,66.949h251.817L194.49,159.286z
										M422.853,303.72c0,1.335-0.624,2.615-1.686,3.437l-106.236,81.674c-0.784,0.597-1.727,0.893-2.644,0.893
										c-0.919,0-1.84-0.29-2.617-0.87c-1.552-1.158-2.135-3.239-1.412-5.051l21.353-54.085H204.415c-2.399,0-4.332-1.932-4.332-4.328
										v-43.333c0-2.391,1.937-4.328,4.332-4.328h125.196l-21.353-54.083c-0.723-1.812-0.14-3.882,1.412-5.053
										c1.554-1.171,3.71-1.171,5.261,0.027l106.242,81.675C422.229,301.105,422.853,302.363,422.853,303.72z"/>
								</g>
							</svg>
						</button>
					</div>
				</div>
			</body>
		</html>
	`, mailto, mailto)

	log.Printf("htmlBody: %v", htmlBody)
	return &htmlBody, nil
}
