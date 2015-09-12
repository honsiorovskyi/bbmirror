# bbmirror
A standalone BitBucket webhooks listener which clones/fetches repositories on changes.

This tool provides a possibility for automatic clonning or fetching
changes from BitBucket when you push a commit there. It uses BitBucket webhooks:
https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html.

[![Build Status](https://travis-ci.org/honsiorovskyi/bbmirror.svg)](https://travis-ci.org/honsiorovskyi/bbmirror)

### Installation

    git clone git@github.com:honsiorovskyi/bbmirror.git
    go build

### Running

    export LISTEN=127.0.0.1:5678
    export REPOSITORY_PATH=/var/lib/bbmirror/repository
    
    ./bbmirror

Both enviromental variables are optional. The default values are provided above.
