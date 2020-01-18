FROM golang

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# variables
ARG USERNAME=beanie

# update system
RUN apt-get update

# install desired software
RUN apt-get install -y git
RUN apt-get install -y vim
RUN apt-get install -y zsh
RUN apt-get install -y curl
RUN apt-get install -y wget
RUN apt-get install -y htop
RUN apt-get install -y fonts-powerline
RUN apt-get install -y locales
RUN apt-get install -y sudo

# add non-root user
RUN groupadd ${USERNAME}
RUN useradd -s /bin/zsh -g $USERNAME -m ${USERNAME}

# make user a sudo user
RUN echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME
RUN chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME
WORKDIR /home/$USERNAME

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