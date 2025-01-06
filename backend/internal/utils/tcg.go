package utils

import (
	"fmt"

	tcg "github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg"
	"github.com/PokemonTCG/pokemon-tcg-sdk-go-v2/pkg/request"
)

func GetCardData(tcgID string, client tcg.Client) (*tcg.PokemonCard, error) {
	cards, err := client.GetCards(
		request.Query("id:" + tcgID),
	)

	if err != nil {
		return nil, err
	}

	if len(cards) == 0 {
		return nil, fmt.Errorf("no card found with id %s", tcgID)
	}

	return cards[0], nil
}
