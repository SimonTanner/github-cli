# github-cli

### Intro
a cli tool for github allowing you to switch between your work & personal profiles, as well as create a repository locally & on github.com/{user_name}

### Installation
To install the cli tool simply clone the repository enter the folder and in the terminal enter `go build`.
This will install all required dependencies and create a go binary called __github-cli__.
You can either run this by entering `./github-cli`
If you want to be able to run it anywhere then enter `cp github-cli /usr/local/bin` on a mac & linux to enable this for all users. After this you will either need to open a new terminal or use `source ~/.{your_bash_or_profile}` in order to use it anywhere.

## Usage
To use the cli tool simply enter `github-cli` in the terminal which will output the following:

    github-cli is a command line interface for use with multiple github repositories

    Usage:
    github-cli [command]

    Available Commands:
    help        Help about any command
    main        set git user to "main" locally
    set         initialises a github repository locally, using "main" profile and makes it private if [private] flag passed
    which       which config variables
    work        set git user to "work" locally

    Flags:
    -h, --help   help for github-cli

    Use "github-cli [command] --help" for more information about a command.
    
### Commands
#### which
Entering `github-cli which` will show the current user.name & user.email as saved in .git/config. If you're not in a repository it will check the gloabl config and output these.
example:
    `github-cli which`
    No user set in local config, checking global user
    Active user: {SomeUser}, email: {user@some_esp.com}

#### main
Entering `github-cli main` will set the user to the 'main' profile. The 1st time you enter this you will see the following output:
    Current branch: main
    user profile "main" not currently set
    To add a new user use enter -a

Simply enter `github-cli main -a` and you will be prompted to enter your username & email address.

If you're not in a git repository you will see the following output:
    Not a git repository. If you wish to initialise a repository try "github-cli set [private]"

(See `set` command below)

If you run this again iy will show that these details have already been set. To overwrite these simply enter `github-cli main -a -f` and you can re-enter these details.

### work

Is similar to `main` in functionality and the same commands apply. This allows you to store you work github user profile and set it locally to this.

### set

The `set` command initialises a git repository in the current folder with the name as the current folder name.
Adding the `-p` or `--private` flag will create a private repository, the default is public. The 1st time you run it will output:

    No access token for profile "main"
    please enter your github access token:

Here you can enter the access token, which will then attempt to create the repository on github.com and set the remote to this. It will also return the status code if successful or return an error if not. In order to create a token log in to github and navigate to the settings page, then select the "Developer settings" page. Here select "Personal access tokens" and click "Generate new token". Simply add a note e.g. "github-cli" and click the repo tick box to allow "Full control of private repositories". The click "Generate token" and copy this into the github-cli prompt (this will store it within you profile). It will then be able to create your repository on github.

Example of successful creation output:

    github-cli set
    Current Directory: /Users/ITST/hello
    hello repository successfully created, http response status code: 201
    Successfully set remote url

