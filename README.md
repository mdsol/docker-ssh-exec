docker-ssh-exec - Secure SSH key injection for Docker builds
================
Allows commands that require an SSH key to be run from within a `Dockerfile`, without leaving the key in the resulting image.

----------------
Overview
----------------
This program runs in two different modes:

* a server mode, run as the Docker image `mdsol/docker-ssh-exec`, which transmits an SSH key on request to the the client; and
* a client mode, invoked from within the `Dockerfile`, that grabs the key from the server, writes it to the filesystem, runs the desired build command, and then *deletes the key* before the filesystem is snapshotted into the build.

----------------
Installation
----------------
To install the server, just pull it like any other Docker image.

To install the client, just grab it from the [releases page][1], uncompress the archive, and copy the binary to somewhere in your `$PATH`. Remember that the client is run during the `docker build...` process, so either install the client just before invoking it, or make sure it's already present in your source image. Here's an example of the code you might run in your source image, to prepare it for SSH cloning from GitHub:

    # install Medidata docker-ssh-exec build tool from S3 bucket "mybucket"
    curl https://s3.amazonaws.com/mybucket/docker-ssh-exec/\
    docker-ssh-exec_0.3.2_linux_amd64.tar.gz | \
      tar -xz --strip-components=1 -C /usr/local/bin \
      docker-ssh-exec_0.3.2_linux_amd64/docker-ssh-exec
    mkdir -p /root/.ssh && chmod 0700 /root/.ssh
    ssh-keyscan github.com >/root/.ssh/known_hosts


----------------
Usage
----------------
To run the server component, pass it the private half of your SSH key, either as a shared volume:

    docker run -v ~/.ssh/id_rsa:/root/.ssh/id_rsa --name=keyserver -d \
      mdsol/docker-ssh-exec -server

or as an ENV var:

    docker run -e DOCKER-SSH-KEY="$(cat ~/.ssh/id_rsa)" --name=keyserver -d \
      mdsol/docker-ssh-exec -server

Then, run a quick test of the client, to make sure it can get the key:

    docker run --rm -it mdsol/docker-ssh-exec cat /root/.ssh/id_rsa

Finally, as long as the source image is set up to trust (or ignore) GitHub's server key, you can clone private repositories from within the `Dockerfile` like this:

    docker-exec-ssh git clone git@github.com:my_user/my_private_repo.git

The client first transfers the key from the server, writing it to `$HOME/.ssh/id_rsa` (by default), then executes whatever command you supply as arguments. Before exiting, it deletes the key from the filesystem.

Here's the command-line help:

    Usage of docker-ssh-exec:
      -key string
          path to key file (default "~/.ssh/id_rsa")
      -port int
          server receiving port (default 1067)
      -server
          run key server instead of command
      -version
          print version and exit
      -wait int
          client timeout, in seconds (default 3)

The software quits with a non-zero exit code (>100) on any error -- except a timeout from the keyserver, in which case it will just ignore the timeout and try to run the build command anyway. If the build command fails, `docker-ssh-exec` returns the exit code of the failed command.


----------------
Known Limitations / Bugs
----------------
The key data is limited to 4096 bytes.


----------------
Contribution / Development
----------------
This software was created by Benton Roberts _(broberts@mdsol.com)_

To build it yourself, just `go get` and `go install` as usual:

    go get github.com/mdsol/12factor-tools/docker-ssh-exec
    cd $GOPATH/src/github.com/mdsol/12factor-tools/docker-ssh-exec
    go install


--------
[1]: https://github.com/mdsol/12factor-tools/releases
