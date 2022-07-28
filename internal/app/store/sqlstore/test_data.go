package sqlstore

import (
	"fmt"

	"github.com/spolyakovs/price-hunter-ITMO/internal/app/model"
)

var (
	TestUsers      []*model.User
	TestPublishers []*model.Publisher
	TestGames      []*model.Game
	TestTags       []*model.Tag
	// TestMarkets            []*model.Market (don't need because markets added already)
	TestUserGameFavourites []*model.UserGameFavourite
	TestGameTags           []*model.GameTag
	TestGameMarketPrices   []*model.GameMarketPrice
)

func FillTestData() {
	fillTestDataUsers()
	fillTestDataPublishers()
	fillTestDataGames()
	fillTestDataTags()
	fillTestDataUserGameFavourites()
	fillTestDataGameTags()
	fillTestDataGameMarketPrices()
}

func fillTestDataUsers() {
	TestUsers = []*model.User{}

	for i := 1; i <= 5; i++ {
		TestUsers = append(TestUsers, &model.User{
			Username: fmt.Sprintf("Test_username_%d", i),
			Email:    fmt.Sprintf("Test_email_%d@example.org", i),
			Password: fmt.Sprintf("Test_password_%d", i),
		})
	}
}

func fillTestDataPublishers() {
	TestPublishers = []*model.Publisher{}

	TestPublishers = append(TestPublishers, &model.Publisher{
		Name: "Valve",
	})

	TestPublishers = append(TestPublishers, &model.Publisher{
		Name: "Hinterland Studio Inc.",
	})

	TestPublishers = append(TestPublishers, &model.Publisher{
		Name: "Wube Software LTD.",
	})

	TestPublishers = append(TestPublishers, &model.Publisher{
		Name: "FromSoftware Inc.",
	})

	TestPublishers = append(TestPublishers, &model.Publisher{
		Name: "Bohemia Interactive",
	})
}

func fillTestDataGames() {
	TestGames = []*model.Game{}

	TestGames = append(TestGames, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/730/header.jpg?t=1641233427",
		Name:           "Counter-Strike: Global Offensive",
		Description:    "Counter-Strike: Global Offensive (CS:GO) расширяет границы ураганной командной игры, представленной ещё 19 лет назад. CS:GO включает в себя новые карты, персонажей, оружие и режимы игры, а также улучшает классическую составляющую CS (de_dust2 и т. п.).",
		ReleaseDate:    "21.08.2012",
		Publisher:      TestPublishers[0],
	})

	TestGames = append(TestGames, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/305620/header.jpg?t=1638931698",
		Name:           "The Long Dark",
		Description:    "The Long Dark is a thoughtful, exploration-survival experience that challenges solo players to think for themselves as they explore an expansive frozen wilderness in the aftermath of a geomagnetic disaster. There are no zombies -- only you, the cold, and all the threats Mother Nature can muster. Welcome to the Quiet Apocalypse.",
		ReleaseDate:    "01.08.2017",
		Publisher:      TestPublishers[1],
	})

	TestGames = append(TestGames, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/427520/header.jpg?t=1620730652",
		Name:           "Factorio",
		Description:    "Factorio is a game about building and creating automated factories to produce items of increasing complexity, within an infinite 2D world. Use your imagination to design your factory, combine simple elements into ingenious structures, and finally protect it from the creatures who don't really like you.",
		ReleaseDate:    "14.08.2020",
		Publisher:      TestPublishers[2],
	})

	TestGames = append(TestGames, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/1245620/header.jpg?t=1649774637",
		Name:           "ELDEN RING",
		Description:    "THE NEW FANTASY ACTION RPG. Rise, Tarnished, and be guided by grace to brandish the power of the Elden Ring and become an Elden Lord in the Lands Between.",
		ReleaseDate:    "25.02.2022",
		Publisher:      TestPublishers[3],
	})

	TestGames = append(TestGames, &model.Game{
		HeaderImageURL: "https://cdn.akamai.steamstatic.com/steam/apps/221100/header.jpg?t=1643209285",
		Name:           "DayZ",
		Description:    "How long can you survive a post-apocalyptic world? A land overrun with an infected &quot;zombie&quot; population, where you compete with other survivors for limited resources. Will you team up with strangers and stay strong together? Or play as a lone wolf to avoid betrayal? This is DayZ – this is your story.",
		ReleaseDate:    "13.12.2018",
		Publisher:      TestPublishers[4],
	})
}

func fillTestDataTags() {
	TestTags = []*model.Tag{}

	TestTags = append(TestTags, &model.Tag{
		Name: "action",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "free_to_play",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "shooter",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "adventure",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "indie",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "simulation",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "strategy",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "casual",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "rpg",
	})

	TestTags = append(TestTags, &model.Tag{
		Name: "mmo",
	})
}

func fillTestDataUserGameFavourites() {
	TestUserGameFavourites = []*model.UserGameFavourite{}

	for i := 0; i <= 4; i++ {
		for j := 0; j <= i; j++ {
			TestUserGameFavourites = append(TestUserGameFavourites, &model.UserGameFavourite{
				User: TestUsers[i],
				Game: TestGames[j],
			})
		}
	}
}

func fillTestDataGameTags() {
	TestGameTags = []*model.GameTag{}

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[0],
		Tag:  TestTags[0],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[0],
		Tag:  TestTags[1],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[0],
		Tag:  TestTags[2],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[1],
		Tag:  TestTags[3],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[1],
		Tag:  TestTags[4],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[1],
		Tag:  TestTags[5],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[1],
		Tag:  TestTags[6],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[2],
		Tag:  TestTags[4],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[2],
		Tag:  TestTags[5],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[2],
		Tag:  TestTags[6],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[2],
		Tag:  TestTags[7],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[3],
		Tag:  TestTags[0],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[3],
		Tag:  TestTags[8],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[4],
		Tag:  TestTags[0],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[4],
		Tag:  TestTags[3],
	})

	TestGameTags = append(TestGameTags, &model.GameTag{
		Game: TestGames[4],
		Tag:  TestTags[9],
	})
}

func fillTestDataGameMarketPrices() {
	TestGameMarketPrices = []*model.GameMarketPrice{}

	TestGameMarketPrices = append(TestGameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "Free To Play",
		FinalValueFormatted:   "Free To Play",
		DiscountPercent:       0,
		MarketGameURL:         "730",
		Game:                  TestGames[0],
		Market:                markets[0],
	})

	TestGameMarketPrices = append(TestGameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "610 руб.",
		FinalValueFormatted:   "247 руб.",
		DiscountPercent:       60,
		MarketGameURL:         "305620",
		Game:                  TestGames[1],
		Market:                markets[0],
	})

	TestGameMarketPrices = append(TestGameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "749 руб.",
		FinalValueFormatted:   "749 руб.",
		DiscountPercent:       0,
		MarketGameURL:         "the-long-dark",
		Game:                  TestGames[1],
		Market:                markets[1],
	})

	TestGameMarketPrices = append(TestGameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "520 руб.",
		FinalValueFormatted:   "520 руб.",
		DiscountPercent:       0,
		MarketGameURL:         "427520",
		Game:                  TestGames[2],
		Market:                markets[0],
	})

	TestGameMarketPrices = append(TestGameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "3 619 руб.",
		FinalValueFormatted:   "3 619 руб.",
		DiscountPercent:       0,
		MarketGameURL:         "factorio",
		Game:                  TestGames[2],
		Market:                markets[2],
	})

	TestGameMarketPrices = append(TestGameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "3 999 руб.",
		FinalValueFormatted:   "3 999 руб.",
		DiscountPercent:       0,
		MarketGameURL:         "1245620",
		Game:                  TestGames[3],
		Market:                markets[0],
	})

	TestGameMarketPrices = append(TestGameMarketPrices, &model.GameMarketPrice{
		InitialValueFormatted: "1 199 руб.",
		FinalValueFormatted:   "719 руб.",
		DiscountPercent:       40,
		MarketGameURL:         "221100",
		Game:                  TestGames[4],
		Market:                markets[0],
	})
}
