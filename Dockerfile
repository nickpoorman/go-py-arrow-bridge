FROM ubuntu:18.04

# Tools
RUN apt-get update && apt-get install -y \
    g++ \
    gdb \
    git \
    make \
    vim \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Go installation
RUN cd /tmp && \
    wget https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.12.5.linux-amd64.tar.gz && \
    rm go1.12.5.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Python bindings
RUN cd /tmp && \
    wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && \
    bash Miniconda3-latest-Linux-x86_64.sh -b -p /miniconda && \
    rm Miniconda3-latest-Linux-x86_64.sh
ENV PATH="/miniconda/bin:${PATH}"
RUN conda install -c conda-forge -y \
    Cython \
    ipython \
    numpy \
    pkg-config \
    pyarrow=0.13.0

ENV LD_LIBRARY_PATH=/miniconda/lib
ENV CONDA_PREFIX=/miniconda
WORKDIR /src/go-py-arrow-bridge