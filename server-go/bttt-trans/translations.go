package bttt_trans

import "github.com/strongo/bots-framework/core"

var TRANS = map[string]map[string]string{
	bots.MessageTextOopsSomethingWentWrong: {
		"ru-RU": "Упс, что-то пошло не так... \xF0\x9F\x98\xB3",
		"en-US": "Oops, something went wrong... \xF0\x9F\x98\xB3",
		"fa-IR": "اوه، یک جای کار مشکل دارد ...  \xF0\x9F\x98\xB3",
		"it-IT": "Ops, qualcosa e' andato storto... \xF0\x9F\x98\xB3",
	},
	MT_START_SELECT_LANG: {
		"en-US": "<b>Please select your language</b>\nПожалуйста выберите язык общения",
		"ru-RU": "<b>Пожалуйста выберите язык общения</b>\nPlease select your language",
	},
	MT_WELCOME: {
		"en-US": `<b>Bidding Tic-Tac-Toe</b>
Welcome!
You can use this bot to play <a href="https://biddingtictactoe.com/#rules">Bidding Tic-Tac-Toe</a> with your friends right in a Telegram chat.`,
		//
		"ru-RU": `<b>Крестики-нолики с аукционом</b>
Приветствуем!
Этот бот позволяет играть в <a href="https://biddingtictactoe.com/ru#rules">аукционные крестики-нолики</a> с вашими друзьями прямо в Телеграмм-чате.`,
	},
	MT_HOW_TO_START_NEW_GAME: {
		"en-US": `<b>To begin new game:</b>
  1. Open Telegram chat with your friend
  2. Type into the text input @BiddingTicTacToeBot and a space
  3. Wait for a popup menu and choose "New game"

<i>First 2 steps can be replaced by clicking the button below!</i>`,
		//
		"ru-RU": `<b>Чтобы начать игру:</b>
  1. Откройте чат с вашим другом
  2. Наберите в поле ввода @BiddingTicTacToeBot и пробел
  3. Дождитесь появлению меню и выберите "Новая игра"

<i>Два первых шага могут быть заменены одним кликом на кнопку ниже!</i>`,
	},
	MT_GAME_NAME: {
		"en-US": "Bidding Tic-Tac-Toe",
		"ru-RU": "Крестики-нолики с аукционом",
	},
	MT_NEW_GAME_WELCOME: {
		"en-US": `To start the game please choose board size.`,
		"ru-RU": `Чтобы начать игру выберите размер доски.`,
	},
	MT_HOW_TO_INLINE: {
		"en-US": `To begin the game and to make first move:
  * choose a cell
  * click Start at bottom of the screen`,
		//
		"ru-RU": `Чтобы начать игру и сделать первый ход:
  * выберите клетку
  * нажмите Start внизу экрана`,
	},
	C_NEW_GAME: {
		"en-US": "New game",
		"ru-RU": "Новая игра",
	},
	C_NEW_GAME_WITH: {
		"en-US": "New game with %v",
		"ru-RU": "Новая игра с %v",
	},
	C_NEW_GAME_THIS_CHAT: {
		"en-US": "New game (this chat)",
		"ru-RU": "Новая игра (этот чат)",
	},
	C_NEW_GAME_OTHER_CHAT: {
		"en-US": "New game with someone else",
		"ru-RU": "Новая игра с кем-то другим",
	},
	MT_PLAYER: {
		"en-US": "Player <b>%v</b>:",
		"ru-RU": "Игрок <b>%v</b>:",
	},
	MT_AWAITING_PLAYER: {
		"en-US": "awaiting...",
		"ru-RU": "ожидается...",
	},
	MT_PLAYER_BALANCE: {
		"en-US": "Balance: %d",
		"ru-RU": "Баланс: %d",
	},
	MT_LAST_BID: {
		"en-US": "Last bid: %d",
		"ru-RU": "Последняя ставка: %d",
	},
	MT_ASK_BID: {
		"en-US": "What is your bid?",
		"ru-RU": "Ваша ставка?",
	},
	MT_ASK_TO_RATE: {
		"en-US": `🙋 <b>Did you like the game?</b> We'll appreciate if you <a href="https://t.me/storebot?start=BiddingTicTacToeBot">rate us</a> with 5⭐s!'`,
		"ru-RU": `🙋 <b>Понравилась игра?</b> Будем признательны если <a href="https://t.me/storebot?start=BiddingTicTacToeBot">оцените нас</a> на 5⭐!`,
	},
	MT_YOUR_BID: {
		"en-US": "your bid: %d",
		"ru-RU": "ваша ставка: %d",
	},
	MT_RIVAL_HAS_TARGET_AND_BID: {
		"en-US": "has target & bid", // TODO: improve English
		"ru-RU": "выбрал цель и ставку",
	},
	MT_RIVAL_HAS_TARGET: {
		"en-US": "has target, making bid...",
		"ru-RU": "выбрал цель, делает ставку...",
	},
	MT_RIVAL_HAS_BID: {
		"en-US": "has bid, aiming...",
		"ru-RU": "сделал ставку, выбирает цель...",
	},
	MT_PLEASE_CHOOSE_YOUR_TARGET: {
		"en-US": "Please choose a cell for next move.",
		"ru-RU": "Выберите клетку для следующего хода.",
	},
	MT_PLEASE_MAKE_A_BID: {
		"en-US": "Please make a bid.",
		"ru-RU": "Сделайте вашу ставку.",
	},
	MT_AWAITING_RIVAL_TURN: {
		"en-US": "Awaiting turn of player %v...",
		"ru-RU": "Ожидается ход игрока %v...",
	},
	MT_PREVIOUS_BID: {
		"en-US": "Previous bid was: %d",
		"ru-RU": "Предыдущая ставка: %d",
	},
	MT_NOT_A_NUMBER: {
		"en-US": "The bid should be an integer positive number.",
		"ru-RU": "Ставка должна быть целым положительным числом.",
	},
	MT_BID_IS_GREATER_THEN_BALANCE: {
		"en-US": "The bid should be an integer positive number.",
		"ru-RU": "Ставка должна быть целым положительным числом.",
	},
	MT_FREE_GAME_SPONSORED_BY: {
		"en-US": "🙏 Free game sponsored by %v",
		"ru-RU": "🙏 Бесплатная игра при поддержке %v - бот для учёта долгов",
	},
	MT_BID_BY: {
		"en-US": "Bid by %v",
		"ru-RU": "Ставка от %v",
	},
	OUR_TWITTER: {
		"en-US": "Our Twitter: %v",
		"ru-RU": "Наш Твиттер %v",
	},
	OUR_FB_PAGE: {
		"en-US": "Our FB page: %v",
		"ru-RU": "Наша FB страница: %v",
	},
	OUR_WEBSITE: {
		"en-US": "Our website: %v",
		"ru-RU": "Наш сайт %v",
	},
	MT_CELL_OCCUPIED: {
		"en-US": "This cell is already occupied.",
		"ru-RU": "Это поле уже занято.",
	},
	MT_TG_BUG_WARNING: {
		"en-US": "There is a bug in Telegram app that prevent updating bid confirmation button. In this case try to remove all text, click cell again and type amount without delay between numbers.",
		"ru-RU": "В Телеграмм есть ошибка когда кнопка подтверждения ставки не обновляется. В этом случае попробуйте удалить весь текст, кликнуть клетку снова и набрать ставку без перерывов между цифрами.",
	},
	MT_BID_TOO_BIG: {
		"en-US": "\xF0\x9F\x9A\xA8 %v, your bid is to big. You have only %d tokens at the moment.",
		"ru-RU": "\xF0\x9F\x9A\xA8 %v, ваша ставка слишком большая. У вас сейчас только %d жетонов.",
	},
	INLINE_BID_TITLE: {
		"en-US": "You are bidding %d tokens",
		"ru-RU": "Ваша ставка: %d жетонов",
	},
	INLINE_BID_DESC: {
		"en-US": "Click here to place the bid",
		"ru-RU": "Кликните тут чтобы сделать ставку",
	},

	MT_RULES: {
		"en-US": "rules",
		"ru-RU": "правила",
	},
	MT_USER_WON_GAME: {
		"en-US": "<b>%v</b> won the game",
		"ru-RU": "<b>%v</b> выиграл(а) игру",
	},
	MT_ENROLLED_TO_TOURNAMENT: {
		"en-US": `<b>You've successfully enrolled to the tournament</b>
		Please read <a href="https://biddingtictactoe.com/#tournament-2017-10">rules & conditions</a> for the tournament prizes giveaway. `,
		"ru-RU": `<b>Вы успешно зарегестрировались на турнир</b>
		Пожалуйста прочитайте <a href="https://biddingtictactoe.com/ru#tournament-2017-10">правила и условия</a> розыгрыша турнирных призов.`,
	},
	MT_TOURNAMENT_201710_SHORT: {
		"en-US": "🎉 <b>We are giving away €120</b> (<i>about $137 USD</i>) on 1st October 2017.",
		"ru-RU": "🎉 <b>Мы разыгрываем €120 евро</b> (<i>примерно $137 долларов</i>) - победители определяются 1го Октября 2017г.",
	},
	MT_TOURNAMENT_201710_LEARN_MORE: {
		"en-US": `To be eligible for the prize & learn rules please <a href="https://t.me/BiddingTicTacToeBot?start=tournament-2017-10-01">enroll to the tournament</a> - <b>it's just 2 clicks and is free</b>!`,
		"ru-RU": `Чтобы претентовать на приз и узнать правила <a href="https://t.me/BiddingTicTacToeBot?start=tournament-2017-10-01">запишитесь в участники турнира</a> - <b>это всего 2 клика и бесплатно</b>!`,
	},
	MT_TOURNAMENT_201710_PRIZES: {
		"en-US": `There would be 3 prizes:

  1. 🏆 €50 to the player who won the most other players.

  2. 🎲 €40 split between a random pair of 2 players who finished at least 1 game.

  3. 📢 Share information about the tournament and be in chance to win another €30.`,
		//
		"ru-RU": `У нас будет 3 приза:

  1. 🏆 €50 получит тот кто обыграет наибольшее количество других игроков.

  2. 🎲 €40 будет поделено между случайной парой игроков которые сыграли хотя бы одну игру.

  3. 📢 Расскажи о турнире другим и ты участник розыгрыша ещё €30.`,
	},

	MT_TOURNAMENT_201710_SPONSOR: {
		"en-US": "💶 This tournament is sponsored by @DebtsTrackerRuBot - debts tracking service. Sends automatic notifications to you and your debtors by Telegram, Email & SMS.",
		"ru-RU": "💶 Данный турнир спонсируетя @DebtsTrackerRuBot - сервис для учёта долгов. Автоматические уведомления вам и вашим должникам по Telegram, Email и SMS.",
	},
}
