package sqlstore_test

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store"
	"github.com/spolyakovs/price-hunter-ITMO/internal/app/store/sqlstore"
)

var (
	users              []*model.User
	publishers         []*model.Publisher
	games              []*model.Game
	tags               []*model.Tag
	markets            []*model.Market
	userGameFavourites []*model.UserGameFavourite
	gameTags           []*model.GameTag
	gameMarketPrices   []*model.GameMarketPrice
)

func insertTestData(st *sqlstore.Store) error {
	if err := st.ClearTables(); err != nil {
		return err
	}

	if err := insertTestDataUsers(st); err != nil {
		return err
	}

	if err := insertTestDataPublishers(st); err != nil {
		return err
	}

	if err := insertTestDataGames(st); err != nil {
		return err
	}

	if err := insertTestDataTags(st); err != nil {
		return err
	}

	if err := insertTestDataMarkets(st); err != nil {
		return err
	}

	if err := insertTestDataUserGameFavourites(st); err != nil {
		return err
	}

	if err := insertTestDataGameTags(st); err != nil {
		return err
	}

	if err := insertTestDataGameMarketPrices(st); err != nil {
		return err
	}

	return nil
}

func insertTestDataUsers(st *sqlstore.Store) error {
	tableName := "Users"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for i := 1; i <= 5; i++ {
		users = append(users, &model.User{
			Username: fmt.Sprintf("Test_username_%d", i),
			Email:    fmt.Sprintf("Test_email_%d@example.org", i),
			Password: fmt.Sprintf("Test_password_%d", i),
		})
	}

	for _, user := range users {
		if err := st.Users().Create(user); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func insertTestDataPublishers(st *sqlstore.Store) error {
	tableName := "Publishers"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	publishers = append(publishers, &model.Publisher{
		Name: "Valve",
	})

	publishers = append(publishers, &model.Publisher{
		Name: "Hinterland Studio Inc.",
	})

	publishers = append(publishers, &model.Publisher{
		Name: "Wube Software LTD.",
	})

	publishers = append(publishers, &model.Publisher{
		Name: "FromSoftware Inc.",
	})

	publishers = append(publishers, &model.Publisher{
		Name: "Bohemia Interactive",
	})

	for _, publisher := range publishers {
		if err := st.Publishers().Create(publisher); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func insertTestDataGames(st *sqlstore.Store) error {
	tableName := "Games"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	games = append(games, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/730/header.jpg?t=1641233427",
		Name:           "Counter-Strike: Global Offensive",
		Description:    "Counter-Strike: Global Offensive (CS:GO) расширяет границы ураганной командной игры, представленной ещё 19 лет назад. CS:GO включает в себя новые карты, персонажей, оружие и режимы игры, а также улучшает классическую составляющую CS (de_dust2 и т. п.).",
		ReleaseDate:    "21.08.2012",
		Publisher:      publishers[0],
	})

	games = append(games, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/305620/header.jpg?t=1638931698",
		Name:           "The Long Dark",
		Description:    "The Long Dark is a thoughtful, exploration-survival experience that challenges solo players to think for themselves as they explore an expansive frozen wilderness in the aftermath of a geomagnetic disaster. There are no zombies -- only you, the cold, and all the threats Mother Nature can muster. Welcome to the Quiet Apocalypse.",
		ReleaseDate:    "01.08.2017",
		Publisher:      publishers[1],
	})

	games = append(games, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/427520/header.jpg?t=1620730652",
		Name:           "Factorio",
		Description:    "Factorio is a game about building and creating automated factories to produce items of increasing complexity, within an infinite 2D world. Use your imagination to design your factory, combine simple elements into ingenious structures, and finally protect it from the creatures who don't really like you.",
		ReleaseDate:    "14.08.2020",
		Publisher:      publishers[2],
	})

	games = append(games, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/1245620/header.jpg?t=1649774637",
		Name:           "ELDEN RING",
		Description:    "THE NEW FANTASY ACTION RPG. Rise, Tarnished, and be guided by grace to brandish the power of the Elden Ring and become an Elden Lord in the Lands Between.",
		ReleaseDate:    "25.02.2022",
		Publisher:      publishers[3],
	})

	games = append(games, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/221100/header.jpg?t=1643209285",
		Name:           "DayZ",
		Description:    "How long can you survive a post-apocalyptic world? A land overrun with an infected &quot;zombie&quot; population, where you compete with other survivors for limited resources. Will you team up with strangers and stay strong together? Or play as a lone wolf to avoid betrayal? This is DayZ – this is your story.",
		ReleaseDate:    "13.12.2018",
		Publisher:      publishers[4],
	})

	for _, game := range games {
		if err := st.Games().Create(game); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func insertTestDataTags(st *sqlstore.Store) error {
	tableName := "Tags"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	tags = append(tags, &model.Tag{
		Name: "action",
	})

	tags = append(tags, &model.Tag{
		Name: "free_to_play",
	})

	tags = append(tags, &model.Tag{
		Name: "shooter",
	})

	tags = append(tags, &model.Tag{
		Name: "adventure",
	})

	tags = append(tags, &model.Tag{
		Name: "indie",
	})

	tags = append(tags, &model.Tag{
		Name: "simulation",
	})

	tags = append(tags, &model.Tag{
		Name: "strategy",
	})

	tags = append(tags, &model.Tag{
		Name: "casual",
	})

	tags = append(tags, &model.Tag{
		Name: "rpg",
	})

	tags = append(tags, &model.Tag{
		Name: "mmo",
	})

	for _, tag := range tags {
		if err := st.Tags().Create(tag); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func insertTestDataMarkets(st *sqlstore.Store) error {
	tableName := "Markets"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	markets = append(markets, &model.Market{
		Name: "Steam",
	})

	markets = append(markets, &model.Market{
		Name: "EpicGamesStore",
	})

	markets = append(markets, &model.Market{
		Name: "GOG.com",
	})

	for _, market := range markets {
		if err := st.Markets().Create(market); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func insertTestDataUserGameFavourites(st *sqlstore.Store) error {
	tableName := "UserGameFavourites"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	for i := 0; i <= 4; i++ {
		for j := 0; j <= i; j++ {
			userGameFavourites = append(userGameFavourites, &model.UserGameFavourite{
				User: users[i],
				Game: games[j],
			})
		}
	}

	for _, userGameFavourite := range userGameFavourites {
		if err := st.UserGameFavourites().Create(userGameFavourite); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func insertTestDataGameTags(st *sqlstore.Store) error {
	tableName := "GameTags"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	gameTags = append(gameTags, &model.GameTag{
		Game: games[0],
		Tag:  tags[0],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[0],
		Tag:  tags[1],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[0],
		Tag:  tags[2],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[1],
		Tag:  tags[3],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[1],
		Tag:  tags[4],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[1],
		Tag:  tags[5],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[1],
		Tag:  tags[6],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[2],
		Tag:  tags[4],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[2],
		Tag:  tags[5],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[2],
		Tag:  tags[6],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[2],
		Tag:  tags[7],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[3],
		Tag:  tags[0],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[3],
		Tag:  tags[8],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[4],
		Tag:  tags[0],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[4],
		Tag:  tags[3],
	})

	gameTags = append(gameTags, &model.GameTag{
		Game: games[4],
		Tag:  tags[9],
	})

	for _, gameTag := range gameTags {
		if err := st.GameTags().Create(gameTag); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}

func insertTestDataGameMarketPrices(st *sqlstore.Store) error {
	tableName := "GameMarketPrices"
	errWrapMessage := fmt.Sprintf(store.ErrTestDataInsertionMessageFormat, tableName)

	gameMarketPrices = append(gameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "Free To Play",
		FinalValueFormatted:   "Free To Play",
		DiscountPercent:       0,
		MarketGameURL:         "730",
		Game:                  games[0],
		Market:                markets[0],
	})

	gameMarketPrices = append(gameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "610 руб.",
		FinalValueFormatted:   "247 руб.",
		DiscountPercent:       60,
		MarketGameURL:         "305620",
		Game:                  games[1],
		Market:                markets[0],
	})

	gameMarketPrices = append(gameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "749 руб.",
		FinalValueFormatted:   "749 руб.",
		DiscountPercent:       0,
		MarketGameURL:         "the-long-dark",
		Game:                  games[1],
		Market:                markets[1],
	})

	gameMarketPrices = append(gameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "520 руб.",
		FinalValueFormatted:   "520 руб.",
		DiscountPercent:       0,
		MarketGameURL:         "427520",
		Game:                  games[2],
		Market:                markets[0],
	})

	gameMarketPrices = append(gameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "3 619 руб.",
		FinalValueFormatted:   "3 619 руб.",
		DiscountPercent:       0,
		MarketGameURL:         "factorio",
		Game:                  games[2],
		Market:                markets[2],
	})

	gameMarketPrices = append(gameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "3 999 руб.",
		FinalValueFormatted:   "3 999 руб.",
		DiscountPercent:       0,
		MarketGameURL:         "1245620",
		Game:                  games[3],
		Market:                markets[0],
	})

	gameMarketPrices = append(gameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "1 199 руб.",
		FinalValueFormatted:   "719 руб.",
		DiscountPercent:       40,
		MarketGameURL:         "221100",
		Game:                  games[4],
		Market:                markets[0],
	})

	for _, gameMarketPrice := range gameMarketPrices {
		if err := st.GameMarketPrices().Create(gameMarketPrice); err != nil {
			return errors.Wrap(err, errWrapMessage)
		}
	}

	return nil
}
