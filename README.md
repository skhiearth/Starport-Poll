# Starport-Poll

Tools used:

## Starport

`Starport` -> For generating, serving and scaffolding the application.

`starport app github.com/GITHUB-USERNAME/APP-NAME` to scaffold the application.

In the app directory;
1. `app/app.go` is the main file of our application, ties together all modules and permissions.
2. `cmd` directory with two directories for CLI (Interacting with application) commands and d (Node) commands.
3. `frontend` directory for web based view client. We can use any framework we want. The default framework used if `Vue.js`
4. `x` directory containing customer modules.

`starport serve` to run the application.


Generate data type in the blockchain
`startport type NAME PROPERTIES`

`starport type poll title options`

This generates several files for us in `x/types/poll`

`x/poll/types/TypePoll.go` -> This file contains definition of the `Poll` type. We can see that a poll has two fields (creator and ID), which will be created automatically, and two fields (title and options) **defined by us**. Since we want Options to be a list of strings, replace `string` with `[]string`.

`x/poll/types/MsgCreatePoll.go` -> This file defines a message that creates a poll.

![](transactionFlow.png)

To write anything to a blockchain or perform any other state transition a client (web app in our case) makes an HTTP POST request with a title and options to `http://localhost:1317/voter/poll` endpoint handler for which is defined in `x/poll/client/rest/txPoll.go`. The handler creates an unsigned transaction which contains an array of messages. The client then signs the transaction and sends it to `http://localhost:1317/txs`. The application then processes the transaction by sending each message to a corresponding handler, in our case `x/poll/handlerMessageCreatePoll.go`. A handler then calls a `CreatePoll` function defined in `x/voter/keeper/poll.go` which writes the poll data into the store.

Going back to `MsgCreatePoll.go`, we need to make options to be stored as a list instead of a string. Replace `Options string` with `Options []string` in `MsgCreatePoll` struct and `options string` with `options []string` in the arguments of `NewMsgCreatePoll` function.

`x/poll/client/rest/txPoll.go` -> Replace `Options string` with `Options []string` in `createPollRequest` struct.

`x/poll/client/cli/txPoll.go` -> A user can also interact with our application through a command line interface.

`pollcli tx poll create-poll "Text editors" "Emacs" "Vim" --from user1`

This command will generate a transaction with "create poll" message, sign it using a private key of `user1` (one of two users created by default) and broadcast it to the blockchain.

The only modification we need to make is to change a line that reads arguments from the console: `argsOptions := args[1:len(args)]`. This will assume that all arguments after the first one represent a list of options.

## Front-end application

Starport has generated a basic front-end for our application. For convenience Vue.js framework is used with Vuex for state management, but since all features of our application are exposed through an HTTP API, clients can be built using any language or framework.

We'll be mostly interested in `vue/src/views` directory, which contains page templates of our app, `vue/src/store/index.js` handles sending transactions and receiving data from our blockchain and `vue/src/components` directory, which contains components, like buttons and forms.

Inside `vue/src/store/index.js` we import CosmJS, a library for handling wallets, creating, signing and broadcasting transactions and define a Vuex store. We'll use `entitySubmit` function for sending data to our blockchain (like a JSON representing a newly created poll), `entityFetch` for requesting a list of polls and `accountUpdate` to fetch information about our token balance.

`frontend/src/view/Index.vue`

Since we don't need the default form component replace `<type-list />` inside of `vue/src/views/Index.vue` with a new component `<poll-form />` that will be created in a new file at `vue/src/components/PollForm.vue`.

Poll form component has an input for title and a list of inputs for options. Clicking "Add option" button adds an empty input and clicking "Create poll" sends, creates and broadcasts a transaction with a "create poll" message.

Refresh the page, sign in with a password and create a new poll. It takes a couple of seconds to process a transaction. Now, if you visit `http://localhost:1317/voter/poll `you should see a list of poll (this endpoint is defined in `x/poll/rest/queryPoll.go`).

## Adding votes

A vote type contains poll ID and a value (string representation of the selected option).

`starport type vote pollID value`

`vue/src/views/Index.vue`:

Add `<poll-list />` into the `vue/src/view/Index.vue` file after the poll form component. Then make a new component at `vue/src/components/PollList.vue`.

The `PollList` component lists for every poll, all the options for that poll, as buttons. Selecting an option triggers a `submit` method that broadcasts a transaction with a "create vote" message and fetches data back from our application.

Try creating polls and casting votes. You may notice that it's possible to cast multiple votes for one poll. This is not what we want, so let's fix this behaviour.

## Casting votes only once

To fix this issue we first have to understand how data is stored in our application.

We can think of our data storage as a lexicographically ordered key value store. You can loop through the entries, filter by key prefix, add, update and delete entries.

Both `poll-` and `vote-` are prefixes. They are added to keys for ease of filtering. By convention, prefixes are defined in `x/poll/key.go`.

Whenever a user casts a vote, a new "create vote" message is handled by a handler and is passed to a keeper. Keeper takes a `vote-` prefix, adds a UUID (unique to every message) and uses this string as a key. `x/poll/keeper/vote.go`:

`key := []byte(types.VotePrefix + vote.ID)`

These strings are unique and we get duplicate votes. To fix that we need to make sure a keeper records every vote only once by choosing the right key. In our case, we can use poll ID and creator address to make sure that one user can cast only one vote per poll.

`key := []byte(types.VotePrefix + vote.PollID + "-" + string(vote.Creator))`

Restart the application and try voting multiple times on a single poll, you'll see you can vote as many times as you want but only your most recent vote is counted.

## Introducing a fee for creating polls

Let's make it so that creating a poll costs 200 tokens.

This feature is very easy to add. We already require users to have accounts registered, and each user has tokens on balance. The only thing we need to do is to send coins from user's account to a module account before we create a poll.

`x/poll/handlerMsgCreatePoll.go`:

```go
moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
payment, _ := sdk.ParseCoins("200token")
if err := k.CoinKeeper.SendCoins(ctx, poll.Creator, moduleAcct, payment); err != nil {
	return nil, err
}
```

Add the code above before `k.CreatePoll(ctx, poll)`. This way, if a user does not have enough tokens, the application will raise an error and will not proceed to creating a poll. Make sure to add `"[github.com/tendermint/tendermint/crypto](http://github.com/tendermint/tendermint/crypto)"` to the import statement (if your text editor didn't do that for you).

Now, restart the app and try creating several polls to see how this affects your token balance.