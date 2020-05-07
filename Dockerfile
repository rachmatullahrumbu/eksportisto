#---------------------------------------------------------------------
# Stage 0: Build Rust library
# Outputs: rust build output @ /bls-zexe/target/
#---------------------------------------------------------------------
FROM ubuntu:16.04 as rustbuilder

RUN apt update && apt install -y curl musl-tools
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH=$PATH:~/.cargo/bin
RUN $HOME/.cargo/bin/rustup install 1.41.0 && $HOME/.cargo/bin/rustup default 1.41.0 && $HOME/.cargo/bin/rustup target add x86_64-unknown-linux-musl
COPY ./external/bls-zexe /bls-zexe
RUN cd /bls-zexe && $HOME/.cargo/bin/cargo build --target x86_64-unknown-linux-musl --release


#---------------------------------------------------------------------
# Stage 1: Build Rosetta
# Outputs: binary @ /rosetta/rosetta
#---------------------------------------------------------------------
FROM golang:1.13-alpine as builder
WORKDIR /rosetta
RUN apk add --no-cache make gcc musl-dev linux-headers git



# Copy BLS library + rust build output
COPY external ./external
RUN mkdir -p external/bls-zexe/target/release/
COPY --from=rustbuilder /bls-zexe/target/x86_64-unknown-linux-musl/release/libepoch_snark.a external/bls-zexe/target/release/

# Downnload dependencies & cache them in docker layer
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build project
#  (this saves to redownload everything when go.mod/sum didn't change)
COPY . .
RUN go build -o eksportisto .

ENTRYPOINT [ "/rosetta/eksportisto" ]