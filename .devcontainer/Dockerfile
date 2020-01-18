FROM ubuntu

# set working directory to home
WORKDIR /root

# update system
RUN apt update

# install required software
RUN apt install -y git
RUN apt install -y vim
RUN apt install -y zsh
RUN apt install -y curl
RUN apt install -y wget
RUN apt install -y htop
RUN apt install -y fonts-powerline
RUN apt install -y locales

# install golang
RUN curl -O https://dl.google.com/go/go1.13.6.linux-amd64.tar.gz
RUN tar -xvf go1.13.6.linux-amd64.tar.gz -C /usr/local
RUN rm -rf https://dl.google.com/go/go1.13.6.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# perform setup
RUN locale-gen en_ZA.UTF-8

# terminal colors with xterm
ENV TERM xterm

# set the zsh theme
ENV ZSH_THEME ''

# run the oh-my-zsh installation script  
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true

# set zsh as shell
ENTRYPOINT [ "/bin/zsh" ]