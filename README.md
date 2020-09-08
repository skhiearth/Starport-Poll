# Starport-Poll

Tools used:
1. `Starport` -> For generating, serving and scaffolding the application.

`starport app github.com/GITHUB-USERNAME/APP-NAME` to scaffold the application.

In the app directory;
1. `app/app.go` is the main file of our application, ties together all modules and permissions.
2. `cmd` directory with two directories for CLI (Interacting with application) commands and d (Node) commands.
3. `vue` directory for web based view client. We can use any framework we want. The default framework used if `Vue.js`
4. `x` directory containing customer modules.

`starport serve` to run the application.


Generate data type in the blockchain
`startport type NAME PROPERTIES`

starport type poll title options

This generates several files for us in `x/types/poll`