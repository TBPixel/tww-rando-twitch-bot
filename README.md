# TWWR Twitch Bot

The Wind Waker Randomizer community twitch bot. A public chat bot, API and web app for streamers and restreams.

## Legend

- [Roadmap](#roadmap)
- [Development](#development)
- [Specification](#specification)
  - [Chat Commands](#chat-commands)
  - [API](#api)
- [React](#react)

## Roadmap

- [ ] Conceptualize project as a CLI
  - [Refer to the spec](#specification)
- [ ] Document self-hosting and deployment options
- [x] Define data model
- [ ] Support local and remote database drivers
  - [x] Support Badger embedded database
  - [ ] Support remote PostgreSQL database
- [x] Twitch IRC Bot Integration
- [ ] Integrate with Twitch API
  - [x] Oauth2 w/ Twitch API
  - [ ] Integrate with Twitch pubsub APIs (go live events, etc)
- [ ] Integrate with Racetime.gg API
  - [ ] Oauth2 w/ Racetime.gg API
  - [ ] Integrate with Racetime.gg APIs to watch race rooms
  - [x] Integrate with Racetime.gg leaderboard APIs
  - [x] Integrate with Racetime.gg user APIs
  - [x] Integrate with Racetime.gg race room APIs

## Development

Frontend Development requires [Node.js](https://nodejs.org/en/) and [npm](https://www.npmjs.com/). Backend requires [Go](https://golang.org/).

All development environments are welcome; please remember to `.gitignore` any dev specific files.

### Environment Variables

Copy the example env `cp .example.env .env` and edit as needed.

For `TWITCH_` environment variables, generate them by following the official twitch [chatbot/irc documentation on environment variables](https://dev.twitch.tv/docs/irc#get-environment-variables).

### Go Backend

Run the backend with go as such:

```shell
go run ./cmd/cli
```

## Specification

### Chat Commands

- [ ] `!twwr` A help command that lists all the available commands.
- [ ] `!twwr race` Display info about the current race, such as settings and preset info.
- [ ] `!twwr vs` Display (and possibly link to the streams of) the other runners in this race.
- [ ] `!twwr leaderboard` Retrieve the leaderboard position of the current runner, including their RT Bux.
- [ ] `!twwr link` Get a link to the racetime room.
- [ ] `!twwr exampleperma` Get an example permalink for the current settings, if available.
- [ ] `!twwr perma` Get the permalink for the current settings, if available.
- [ ] `!twwr restream` Get a link to the restream, if available.
- [ ] `!twwr multi` Generate a link to a multi-twitch stream view of all the runners in the racetime room.

### API

TODO

## React

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

#### Available Scripts

In the project directory, you can run:

##### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

##### `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

##### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

##### `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.

You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.

#### Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).
